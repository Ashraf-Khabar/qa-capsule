package server

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

var (
	presetAvatars = map[string]bool{
		"a01": true, "a02": true, "a03": true, "a04": true,
		"a05": true, "a06": true, "a07": true, "a08": true,
		"a09": true, "a10": true, "a11": true, "a12": true,
	}
	photoDataRe = regexp.MustCompile(`^data:image/(jpeg|png|webp);base64,[A-Za-z0-9+/=]+$`)
)

func registerCollaborationRoutes(config *core.Config) {
	http.HandleFunc("/api/users/directory", jwtAuthMiddleware(config, "", handleUserDirectory))
	http.HandleFunc("/api/me/avatar", jwtAuthMiddleware(config, "", handleMyAvatar))
	http.HandleFunc("/api/incident-comments", jwtAuthMiddleware(config, "", handleIncidentCommentBatch))
}

func handleUserDirectory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	rows, err := core.DB.Query(`SELECT username, COALESCE(fullname, ''), COALESCE(avatar_kind, ''), COALESCE(avatar_value, '')
		FROM users WHERE is_active = 1 ORDER BY username LIMIT 200`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	users := make([]map[string]string, 0)
	for rows.Next() {
		var username, fullname, kind, value string
		if err := rows.Scan(&username, &fullname, &kind, &value); err != nil {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(username), q) && !strings.Contains(strings.ToLower(fullname), q) {
			continue
		}
		users = append(users, map[string]string{
			"username":     username,
			"fullname":     fullname,
			"avatar_kind":  kind,
			"avatar_value": value,
		})
		if q != "" && len(users) >= 8 {
			break
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"users": users})
}

func handleMyAvatar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims := claimsFromRequest(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var body struct {
		Kind  string `json:"kind"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	kind, value, err := normalizeAvatar(body.Kind, body.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := core.DB.Exec(`UPDATE users SET avatar_kind = ?, avatar_value = ? WHERE username = ?`, kind, value, claims.Username); err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "saved", "avatar_kind": kind, "avatar_value": value})
}

func normalizeAvatar(kind, value string) (string, string, error) {
	kind = strings.TrimSpace(kind)
	value = strings.TrimSpace(value)
	switch kind {
	case "", "initials":
		return "", "", nil
	case "preset":
		if !presetAvatars[value] {
			return "", "", errString("Unknown avatar")
		}
		return "preset", value, nil
	case "photo":
		if len(value) > 200000 || !photoDataRe.MatchString(value) {
			return "", "", errString("Photo must be a small PNG, JPEG, or WebP")
		}
		return "photo", value, nil
	default:
		return "", "", errString("Unknown avatar kind")
	}
}

type errString string

func (e errString) Error() string { return string(e) }

func handleIncidentComments(w http.ResponseWriter, r *http.Request, idRaw string) {
	claims := claimsFromRequest(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	incidentID, err := strconv.ParseInt(idRaw, 10, 64)
	if err != nil || incidentID <= 0 {
		http.Error(w, "Invalid incident id", http.StatusBadRequest)
		return
	}
	if !core.UserCanAccessIncidentByID(claims.Username, claims.Role, incidentID) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := loadComments([]int64{incidentID})
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"comments": items[incidentID]})
	case http.MethodPost:
		var body struct {
			Blocks []CommentBlock `json:"blocks"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		rendered, ok := RenderCommentBlocks(body.Blocks)
		if !ok {
			http.Error(w, "Comment is empty", http.StatusBadRequest)
			return
		}
		res, err := core.DB.Exec(`INSERT INTO incident_comments (incident_id, author_username, body_html) VALUES (?, ?, ?)`,
			incidentID, claims.Username, rendered)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		id, _ := res.LastInsertId()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "body_html": rendered})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleIncidentCommentBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims := claimsFromRequest(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	ids := make([]int64, 0, 16)
	for _, part := range strings.Split(r.URL.Query().Get("ids"), ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil || id <= 0 {
			continue
		}
		if !core.UserCanAccessIncidentByID(claims.Username, claims.Role, id) {
			continue
		}
		ids = append(ids, id)
		if len(ids) >= 100 {
			break
		}
	}
	grouped, err := loadComments(ids)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"comments": grouped})
}

type commentRow struct {
	ID           int64  `json:"id"`
	Author       string `json:"author"`
	Fullname     string `json:"fullname"`
	AvatarKind   string `json:"avatar_kind"`
	AvatarValue  string `json:"avatar_value"`
	BodyHTML     string `json:"body_html"`
	CreatedAt    string `json:"created_at"`
	IncidentID   int64  `json:"-"`
}

func loadComments(ids []int64) (map[string][]commentRow, error) {
	out := map[string][]commentRow{}
	for _, id := range ids {
		out[strconv.FormatInt(id, 10)] = []commentRow{}
	}
	if len(ids) == 0 {
		return out, nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `SELECT c.id, c.incident_id, c.author_username, COALESCE(u.fullname, ''), COALESCE(u.avatar_kind, ''), COALESCE(u.avatar_value, ''), c.body_html, c.created_at
		FROM incident_comments c
		LEFT JOIN users u ON u.username = c.author_username
		WHERE c.incident_id IN (` + strings.Join(placeholders, ",") + `)
		ORDER BY c.id ASC`
	rows, err := core.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var row commentRow
		var created sql.NullString
		if err := rows.Scan(&row.ID, &row.IncidentID, &row.Author, &row.Fullname, &row.AvatarKind, &row.AvatarValue, &row.BodyHTML, &created); err != nil {
			return nil, err
		}
		row.CreatedAt = created.String
		key := strconv.FormatInt(row.IncidentID, 10)
		out[key] = append(out[key], row)
	}
	return out, rows.Err()
}
