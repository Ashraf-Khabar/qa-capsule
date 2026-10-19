package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

const linkTTL = 365 * 24 * time.Hour

type secureLink struct {
	K string `json:"k"`
	I int64  `json:"i"`
	M int64  `json:"m"`
	E int64  `json:"e"`
}

func linkKey() []byte {
	sum := sha256.Sum256(append(append([]byte{}, jwtKey...), []byte("sre-capsule-link-v1")...))
	return sum[:]
}

// mintCommentLink returns an opaque token. The URL does not contain the incident or comment id.
func mintCommentLink(incidentID, commentID int64, now time.Time) (string, error) {
	if len(jwtKey) == 0 || incidentID <= 0 || commentID <= 0 {
		return "", errors.New("invalid link target")
	}
	payload, err := json.Marshal(secureLink{
		K: "comment",
		I: incidentID,
		M: commentID,
		E: now.Add(linkTTL).Unix(),
	})
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(linkKey())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, payload, nil)
	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

func openCommentLink(token string, now time.Time) (int64, int64, error) {
	token = strings.TrimSpace(token)
	if token == "" || len(token) > 512 || strings.ContainsAny(token, "/?#") {
		return 0, 0, errors.New("invalid link")
	}
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return 0, 0, errors.New("invalid link")
	}
	block, err := aes.NewCipher(linkKey())
	if err != nil {
		return 0, 0, errors.New("invalid link")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return 0, 0, errors.New("invalid link")
	}
	nonceSize := gcm.NonceSize()
	if len(raw) <= nonceSize {
		return 0, 0, errors.New("invalid link")
	}
	plain, err := gcm.Open(nil, raw[:nonceSize], raw[nonceSize:], nil)
	if err != nil {
		return 0, 0, errors.New("invalid link")
	}
	var link secureLink
	if err := json.Unmarshal(plain, &link); err != nil || link.K != "comment" || link.I <= 0 || link.M <= 0 {
		return 0, 0, errors.New("invalid link")
	}
	if now.Unix() > link.E {
		return 0, 0, errors.New("expired link")
	}
	return link.I, link.M, nil
}

func handleSecureLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims := claimsFromRequest(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/links/"), "/")
	incidentID, commentID, err := openCommentLink(token, time.Now())
	if err != nil {
		http.Error(w, "This link is not valid", http.StatusNotFound)
		return
	}
	if !core.UserCanAccessIncidentByID(claims.Username, claims.Role, incidentID) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	var project, runID string
	err = core.DB.QueryRow(`SELECT project_name, COALESCE(pipeline_run_id, '') FROM incidents WHERE id = ?`, incidentID).Scan(&project, &runID)
	if err != nil || strings.TrimSpace(runID) == "" {
		http.Error(w, "This link is not valid", http.StatusNotFound)
		return
	}
	var found int64
	if err := core.DB.QueryRow(`SELECT id FROM incident_comments WHERE id = ? AND incident_id = ?`, commentID, incidentID).Scan(&found); err != nil {
		http.Error(w, "This link is not valid", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"project_name":    project,
		"pipeline_run_id": runID,
		"incident_id":     incidentID,
		"comment_id":      commentID,
	})
}

func serveAppIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	w.Header().Set("X-SRE-Capsule-Edition", core.EditionID())
	w.Header().Set("X-SRE-Capsule-Version", core.Version)
	http.ServeFile(w, r, "./web/index.html")
}

func requestPublicBase(r *http.Request) string {
	if base := strings.TrimRight(strings.TrimSpace(os.Getenv("SRE_PUBLIC_URL")), "/"); strings.HasPrefix(base, "http://") || strings.HasPrefix(base, "https://") {
		return base
	}
	proto := "http"
	if r != nil && r.TLS != nil {
		proto = "https"
	}
	if r != nil {
		if fwd := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Proto"), ",")[0]); fwd == "http" || fwd == "https" {
			proto = fwd
		}
	}
	host := ""
	if r != nil {
		host = strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-Host"), ",")[0])
		if host == "" {
			host = r.Host
		}
	}
	if host == "" {
		host = "localhost:9000"
	}
	return proto + "://" + host
}
