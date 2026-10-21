package core

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
)

// IncidentOps is the lifecycle, last runbook, and postmortem attached to one incident.
type IncidentOps struct {
	Fingerprint    string
	Acknowledged   string
	AcknowledgedBy string
	Severity       string
	MitigatedAt    string
	Runbook        string
	Postmortem     string
}

var tagStrip = regexp.MustCompile(`<[^>]*>`)

// LoadIncidentOps reads lifecycle fields for the given incident ids.
func LoadIncidentOps(ids []int) map[int]IncidentOps {
	out := map[int]IncidentOps{}
	if DB == nil || len(ids) == 0 {
		return out
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	q := `SELECT i.id,
		COALESCE(i.fingerprint,''),
		COALESCE(i.acknowledged_at,''),
		COALESCE(i.acknowledged_by,''),
		COALESCE(i.severity,''),
		COALESCE(i.mitigated_at,''),
		COALESCE((SELECT r.template_name || ' — ' || r.summary FROM incident_runbooks r WHERE r.incident_id = i.id ORDER BY r.id DESC LIMIT 1), ''),
		COALESCE((SELECT p.body FROM incident_postmortems p WHERE p.incident_id = i.id), '')
		FROM incidents i WHERE i.id IN (` + strings.Join(placeholders, ",") + `)`
	rows, err := DB.Query(q, args...)
	if err != nil {
		return out
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var op IncidentOps
		if err := rows.Scan(&id, &op.Fingerprint, &op.Acknowledged, &op.AcknowledgedBy, &op.Severity, &op.MitigatedAt, &op.Runbook, &op.Postmortem); err != nil {
			continue
		}
		out[id] = op
	}
	return out
}

// AcknowledgeIncident records the owner and severity. A second call does not replace the first owner.
func AcknowledgeIncident(id int64, actor, severity string) error {
	severity = normalizeSeverity(severity)
	res, err := DB.Exec(`UPDATE incidents
		SET acknowledged_at = COALESCE(acknowledged_at, CURRENT_TIMESTAMP),
		    acknowledged_by = CASE WHEN COALESCE(acknowledged_by,'') = '' THEN ? ELSE acknowledged_by END,
		    severity = CASE WHEN COALESCE(severity,'') = '' THEN ? ELSE severity END
		WHERE id = ?`, actor, severity, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// MitigateIncident marks the incident mitigated and acknowledges it when nobody has yet.
func MitigateIncident(id int64, actor string) error {
	res, err := DB.Exec(`UPDATE incidents
		SET mitigated_at = COALESCE(mitigated_at, CURRENT_TIMESTAMP),
		    acknowledged_at = COALESCE(acknowledged_at, CURRENT_TIMESTAMP),
		    acknowledged_by = CASE WHEN COALESCE(acknowledged_by,'') = '' THEN ? ELSE acknowledged_by END
		WHERE id = ?`, actor, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// RecordRunbook stores the confirmed runbook on the incident. It does not call a vendor.
func RecordRunbook(id int64, templateID, templateName, summary, actor string) error {
	_, err := DB.Exec(`INSERT INTO incident_runbooks (incident_id, template_id, template_name, summary, actor)
		VALUES (?, ?, ?, ?, ?)`, id, templateID, templateName, summary, actor)
	return err
}

// SavePostmortem replaces the draft for one incident.
func SavePostmortem(id int64, body, author string) error {
	body = strings.TrimSpace(body)
	if body == "" {
		return fmt.Errorf("empty postmortem")
	}
	_, err := DB.Exec(`INSERT INTO incident_postmortems (incident_id, body, author, updated_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(incident_id) DO UPDATE SET body = excluded.body, author = excluded.author, updated_at = CURRENT_TIMESTAMP`,
		id, body, author)
	return err
}

// PostmortemFacts is the incident record a draft is written from.
type PostmortemFacts struct {
	ID         int64
	Project    string
	Name       string
	Error      string
	Severity   string
	CreatedAt  string
	AckAt      string
	AckBy      string
	Mitigated  string
	ResolvedAt string
	ResolvedBy string
	Deploy     string
	Runbook    string
	Comments   []string
}

// ComposePostmortem writes a chronology from stored facts. It does not invent a root cause.
func ComposePostmortem(f PostmortemFacts) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Postmortem draft — %s / %s (#%d)\n\n", strings.TrimSpace(f.Project), strings.TrimSpace(f.Name), f.ID)
	b.WriteString("Timeline\n")
	writeLine(&b, "Detected", f.CreatedAt)
	if f.AckAt != "" {
		who := strings.TrimSpace(f.AckBy)
		if f.Severity != "" {
			who = strings.TrimSpace(who + " " + f.Severity)
		}
		writeLine(&b, "Acknowledged", strings.TrimSpace(f.AckAt+" "+who))
	}
	writeLine(&b, "Mitigated", f.Mitigated)
	if f.ResolvedAt != "" {
		writeLine(&b, "Resolved", strings.TrimSpace(f.ResolvedAt+" "+f.ResolvedBy))
	}
	if strings.TrimSpace(f.Deploy) != "" {
		fmt.Fprintf(&b, "\nNearest deploy\n%s\n", strings.TrimSpace(f.Deploy))
	}
	if strings.TrimSpace(f.Runbook) != "" {
		fmt.Fprintf(&b, "\nRunbook recorded\n%s\n", strings.TrimSpace(f.Runbook))
	}
	errText := strings.TrimSpace(f.Error)
	if len(errText) > 500 {
		errText = errText[:500]
	}
	if errText != "" {
		fmt.Fprintf(&b, "\nError\n%s\n", errText)
	}
	if len(f.Comments) > 0 {
		b.WriteString("\nComments\n")
		for _, c := range f.Comments {
			c = strings.TrimSpace(tagStrip.ReplaceAllString(c, " "))
			c = strings.Join(strings.Fields(c), " ")
			if c == "" {
				continue
			}
			if len(c) > 240 {
				c = c[:240]
			}
			fmt.Fprintf(&b, "- %s\n", c)
		}
	}
	b.WriteString("\nWhat to fill in: cause, what fixed it, and what to change next time.\n")
	return b.String()
}

func writeLine(b *strings.Builder, label, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	fmt.Fprintf(b, "- %s: %s\n", label, value)
}

func normalizeSeverity(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "sev1", "sev2", "sev3", "sev4":
		return strings.ToLower(strings.TrimSpace(s))
	default:
		return "sev3"
	}
}

// LoadPostmortemFacts reads the rows a draft is built from.
func LoadPostmortemFacts(id int64, deploy string) (PostmortemFacts, error) {
	var f PostmortemFacts
	f.ID = id
	f.Deploy = deploy
	err := DB.QueryRow(`SELECT project_name, name,
		COALESCE(error_message,''), COALESCE(severity,''),
		COALESCE(created_at,''), COALESCE(acknowledged_at,''), COALESCE(acknowledged_by,''),
		COALESCE(mitigated_at,''), COALESCE(resolved_at,''), COALESCE(resolved_by,''),
		COALESCE((SELECT r.template_name || ' — ' || r.summary FROM incident_runbooks r WHERE r.incident_id = incidents.id ORDER BY r.id DESC LIMIT 1), '')
		FROM incidents WHERE id = ?`, id).Scan(
		&f.Project, &f.Name, &f.Error, &f.Severity, &f.CreatedAt, &f.AckAt, &f.AckBy, &f.Mitigated, &f.ResolvedAt, &f.ResolvedBy, &f.Runbook)
	if err != nil {
		return f, err
	}
	rows, err := DB.Query(`SELECT body_html FROM incident_comments WHERE incident_id = ? AND parent_id IS NULL ORDER BY id LIMIT 8`, id)
	if err != nil {
		return f, nil
	}
	defer rows.Close()
	for rows.Next() {
		var body string
		if rows.Scan(&body) == nil {
			f.Comments = append(f.Comments, body)
		}
	}
	return f, nil
}
