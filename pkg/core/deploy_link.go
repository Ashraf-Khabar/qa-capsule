package core

import (
	"database/sql"
	"strings"
	"time"
)

const deployCloseMinutes = 30.0

// PipelineStamp is one recorded pipeline run.
type PipelineStamp struct {
	Project   string
	ProjectID string
	RunID     string
	Commit    string
	Branch    string
	Outcome   string
	StartedAt string
}

// DeployLink is the pipeline run that sits next to an incident.
type DeployLink struct {
	RunID         string  `json:"run_id"`
	Commit        string  `json:"commit"`
	Branch        string  `json:"branch"`
	Outcome       string  `json:"outcome"`
	StartedAt     string  `json:"started_at"`
	MinutesBefore float64 `json:"minutes_before"`
	Close         bool    `json:"close"`
	ProjectID     string  `json:"project_id"`
}

// LinkDeploy picks the incident's own run, or the latest run that started before it.
func LinkDeploy(runs []PipelineStamp, project, runID, createdAt string) (DeployLink, bool) {
	var candidates []PipelineStamp
	for _, run := range runs {
		if run.Project == project {
			candidates = append(candidates, run)
		}
	}
	if len(candidates) == 0 {
		return DeployLink{}, false
	}
	created, hasCreated := parseSQLiteTime(createdAt)
	chosen, found := PipelineStamp{}, false
	plainRun := strings.TrimSpace(runID)
	if plainRun != "" && !strings.HasPrefix(plainRun, "legacy-") {
		for _, run := range candidates {
			if run.RunID == plainRun {
				chosen = run
				found = true
				break
			}
		}
	}
	if !found && hasCreated {
		var best time.Time
		for _, run := range candidates {
			started, ok := parseSQLiteTime(run.StartedAt)
			if !ok || started.After(created) {
				continue
			}
			if !found || started.After(best) {
				chosen = run
				best = started
				found = true
			}
		}
	}
	if !found {
		return DeployLink{}, false
	}
	link := DeployLink{
		RunID:     chosen.RunID,
		Commit:    chosen.Commit,
		Branch:    chosen.Branch,
		Outcome:   chosen.Outcome,
		StartedAt: chosen.StartedAt,
		ProjectID: chosen.ProjectID,
	}
	if started, ok := parseSQLiteTime(chosen.StartedAt); ok && hasCreated {
		link.MinutesBefore = created.Sub(started).Minutes()
		link.Close = link.MinutesBefore >= -2 && link.MinutesBefore <= deployCloseMinutes
	}
	return link, true
}

// LoadPipelineStamps reads runs for the given projects inside the window.
func LoadPipelineStamps(projects []string, from, to time.Time) []PipelineStamp {
	if DB == nil || len(projects) == 0 {
		return nil
	}
	placeholders := make([]string, len(projects))
	args := make([]interface{}, 0, len(projects)+2)
	for i, name := range projects {
		placeholders[i] = "?"
		args = append(args, name)
	}
	args = append(args, from.UTC().Format("2006-01-02 15:04:05"), to.UTC().Format("2006-01-02 15:04:05"))
	query := `
		SELECT r.project_name, COALESCE(p.id, ''), r.pipeline_run_id,
			COALESCE(r.commit_sha, ''), COALESCE(r.branch, ''), COALESCE(r.outcome, ''),
			COALESCE(r.started_at, '')
		FROM pipeline_runs r
		LEFT JOIN projects p ON p.name = r.project_name
		WHERE r.project_name IN (` + strings.Join(placeholders, ",") + `)
		  AND r.started_at >= ? AND r.started_at <= ?`
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []PipelineStamp
	for rows.Next() {
		var stamp PipelineStamp
		var projectID, commit, branch, outcome, started sql.NullString
		if rows.Scan(&stamp.Project, &projectID, &stamp.RunID, &commit, &branch, &outcome, &started) != nil {
			continue
		}
		stamp.ProjectID = projectID.String
		stamp.Commit = commit.String
		stamp.Branch = branch.String
		stamp.Outcome = outcome.String
		stamp.StartedAt = started.String
		out = append(out, stamp)
	}
	return out
}

func parseSQLiteTime(value string) (time.Time, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, false
	}
	layouts := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		time.RFC3339,
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}
