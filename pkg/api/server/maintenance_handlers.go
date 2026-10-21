package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

func registerMaintenanceRoutes(config *core.Config) {
	http.HandleFunc("/api/maintenance", jwtAuthMiddleware(config, "", handleMaintenance))
}

func handleMaintenance(w http.ResponseWriter, r *http.Request) {
	claims := parseClaims(r)
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"windows": visibleMaintenance(claims.Username, claims.Role)})
	case http.MethodPost:
		if !core.CanAccessFinOps(claims.Role) && !core.IsAdmin(claims.Role) {
			writeJSONError(w, "Manager access required", http.StatusForbidden)
			return
		}
		var body struct {
			Project string `json:"project"`
			Reason  string `json:"reason"`
			Hours   int    `json:"hours"`
			Action  string `json:"action"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeJSONError(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		var err error
		if body.Action == "end" {
			err = core.EndMaintenance(body.Project)
		} else {
			err = core.OpenMaintenance(strings.TrimSpace(body.Project), body.Reason, claims.Username, body.Hours)
		}
		if err != nil {
			msg := err.Error()
			if msg == "sql: no rows in result set" {
				msg = "Unknown gateway or no open window"
			}
			writeJSONError(w, msg, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func visibleMaintenance(username, role string) []core.MaintenanceWindow {
	all := core.ListActiveMaintenance()
	if core.IsAdmin(role) || core.IsManager(role) {
		return all
	}
	var out []core.MaintenanceWindow
	for _, window := range all {
		var projectID string
		if err := core.DB.QueryRow(`SELECT id FROM projects WHERE name = ?`, window.Project).Scan(&projectID); err != nil {
			continue
		}
		if core.UserCanAccessProject(username, role, projectID) {
			out = append(out, window)
		}
	}
	return out
}

func attachMaintenance(incidents []map[string]interface{}, from, to string) {
	fromT, _ := time.Parse("2006-01-02 15:04:05", from)
	toT, _ := time.Parse("2006-01-02 15:04:05", to)
	windows := core.LoadMaintenance(fromT, toT)
	if len(windows) == 0 {
		return
	}
	byProject := map[string][]core.MaintenanceWindow{}
	for _, window := range windows {
		byProject[window.Project] = append(byProject[window.Project], window)
	}
	for _, inc := range incidents {
		project, _ := inc["project_name"].(string)
		created, _ := inc["created_at"].(string)
		hit, ok := core.MatchMaintenance(created, byProject[project])
		if !ok {
			continue
		}
		inc["maintenance"] = map[string]string{
			"reason":  hit.Reason,
			"ends_at": hit.EndsAt,
		}
	}
}
