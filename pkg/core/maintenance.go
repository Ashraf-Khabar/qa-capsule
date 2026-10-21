package core

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// MaintenanceWindow is a declared period on one gateway.
// Incidents that start inside it stay visible and are marked, not hidden.
type MaintenanceWindow struct {
	ID        int64  `json:"id"`
	Project   string `json:"project"`
	Reason    string `json:"reason"`
	StartsAt  string `json:"starts_at"`
	EndsAt    string `json:"ends_at"`
	CreatedBy string `json:"created_by"`
}

// ValidateMaintenanceHours accepts 1 to 24 hours.
func ValidateMaintenanceHours(hours int) error {
	if hours < 1 || hours > 24 {
		return fmt.Errorf("duration must be between 1 and 24 hours")
	}
	return nil
}

// MatchMaintenance reports the window that contains created, for one gateway's windows.
// The end instant is excluded.
func MatchMaintenance(created string, windows []MaintenanceWindow) (MaintenanceWindow, bool) {
	at, ok := parseSQLiteTime(created)
	if !ok {
		return MaintenanceWindow{}, false
	}
	for _, window := range windows {
		start, okStart := parseSQLiteTime(window.StartsAt)
		end, okEnd := parseSQLiteTime(window.EndsAt)
		if !okStart || !okEnd {
			continue
		}
		if !at.Before(start) && at.Before(end) {
			return window, true
		}
	}
	return MaintenanceWindow{}, false
}

// OpenMaintenance starts a window now. A gateway can have only one window still open.
func OpenMaintenance(project, reason, actor string, hours int) error {
	if err := ValidateMaintenanceHours(hours); err != nil {
		return err
	}
	project = strings.TrimSpace(project)
	reason = strings.TrimSpace(reason)
	if project == "" || reason == "" {
		return fmt.Errorf("project and reason are required")
	}
	if len(reason) > 160 {
		reason = reason[:160]
	}
	if DB == nil {
		return sql.ErrNoRows
	}
	var n int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM projects WHERE name = ?`, project).Scan(&n); err != nil {
		return err
	}
	if n == 0 {
		return sql.ErrNoRows
	}
	var open int
	if err := DB.QueryRow(`SELECT COUNT(*) FROM gateway_maintenance WHERE project_name = ? AND ends_at > CURRENT_TIMESTAMP`, project).Scan(&open); err != nil {
		return err
	}
	if open > 0 {
		return fmt.Errorf("this gateway already has a maintenance window")
	}
	_, err := DB.Exec(`INSERT INTO gateway_maintenance (project_name, reason, starts_at, ends_at, created_by)
		VALUES (?, ?, CURRENT_TIMESTAMP, datetime('now', ?), ?)`,
		project, reason, fmt.Sprintf("+%d hours", hours), actor)
	return err
}

// EndMaintenance closes the open window for a gateway.
func EndMaintenance(project string) error {
	project = strings.TrimSpace(project)
	if DB == nil || project == "" {
		return sql.ErrNoRows
	}
	res, err := DB.Exec(`UPDATE gateway_maintenance SET ends_at = CURRENT_TIMESTAMP WHERE project_name = ? AND ends_at > CURRENT_TIMESTAMP`, project)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// ListActiveMaintenance returns windows that have not ended.
func ListActiveMaintenance() []MaintenanceWindow {
	if DB == nil {
		return nil
	}
	rows, err := DB.Query(`SELECT id, project_name, reason, starts_at, ends_at, COALESCE(created_by,'')
		FROM gateway_maintenance WHERE ends_at > CURRENT_TIMESTAMP ORDER BY ends_at`)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanMaintenance(rows)
}

// LoadMaintenance returns windows that overlap [from, to].
func LoadMaintenance(from, to time.Time) []MaintenanceWindow {
	if DB == nil || from.IsZero() || to.IsZero() {
		return ListActiveMaintenance()
	}
	rows, err := DB.Query(`SELECT id, project_name, reason, starts_at, ends_at, COALESCE(created_by,'')
		FROM gateway_maintenance
		WHERE ends_at >= ? AND starts_at <= ?
		ORDER BY starts_at`,
		from.UTC().Format("2006-01-02 15:04:05"), to.UTC().Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanMaintenance(rows)
}

func scanMaintenance(rows *sql.Rows) []MaintenanceWindow {
	var out []MaintenanceWindow
	for rows.Next() {
		var window MaintenanceWindow
		if rows.Scan(&window.ID, &window.Project, &window.Reason, &window.StartsAt, &window.EndsAt, &window.CreatedBy) != nil {
			continue
		}
		out = append(out, window)
	}
	return out
}
