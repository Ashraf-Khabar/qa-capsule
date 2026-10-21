package core

import (
	"database/sql"
	"fmt"
	"math"
	"time"
)

const (
	defaultSLOFailureRate = 0.15
	budgetWatchBurn       = 0.8
	budgetFreezeBurn      = 1
)

// ErrorBudget is the change-failure budget for one gateway over a window.
type ErrorBudget struct {
	Project         string  `json:"project"`
	SLOFailureRate  float64 `json:"slo_failure_rate"`
	Deployments     int     `json:"deployments"`
	Failed          int     `json:"failed"`
	AllowedFailures float64 `json:"allowed_failures"`
	RemainingRatio  float64 `json:"remaining_ratio"`
	BurnRate        float64 `json:"burn_rate"`
	Status          string  `json:"status"`
}

// ScoreErrorBudget turns deployment counts and an allowed failure rate into a decision.
// Status is quiet, healthy, watch, or freeze.
func ScoreErrorBudget(project string, deployments, failed int, slo float64) ErrorBudget {
	if slo <= 0 || slo > 1 || math.IsNaN(slo) {
		slo = defaultSLOFailureRate
	}
	out := ErrorBudget{
		Project:        project,
		SLOFailureRate: slo,
		Deployments:    deployments,
		Failed:         failed,
		Status:         "quiet",
		RemainingRatio: 1,
	}
	if deployments <= 0 {
		return out
	}
	allowed := float64(deployments) * slo
	out.AllowedFailures = allowed
	if allowed <= 0 {
		if failed > 0 {
			out.BurnRate = float64(failed)
			out.RemainingRatio = 0
			out.Status = "freeze"
		} else {
			out.Status = "healthy"
		}
		return out
	}
	consumed := float64(failed)
	remaining := allowed - consumed
	if remaining < 0 {
		remaining = 0
	}
	out.BurnRate = consumed / allowed
	out.RemainingRatio = remaining / allowed
	switch {
	case out.BurnRate >= budgetFreezeBurn:
		out.Status = "freeze"
	case out.BurnRate >= budgetWatchBurn:
		out.Status = "watch"
	default:
		out.Status = "healthy"
	}
	return out
}

// ErrorBudgets reads pipeline outcomes per gateway for the window.
func ErrorBudgets(project string, from, to time.Time) []ErrorBudget {
	if DB == nil {
		return nil
	}
	fromS := from.UTC().Format("2006-01-02 15:04:05")
	toS := to.UTC().Format("2006-01-02 15:04:05")
	rows, err := DB.Query(`
		SELECT r.project_name,
			COALESCE(p.slo_failure_rate, 0.15),
			COUNT(*),
			SUM(CASE WHEN r.outcome = 'failure' THEN 1 ELSE 0 END)
		FROM pipeline_runs r
		LEFT JOIN projects p ON p.name = r.project_name
		WHERE r.started_at >= ? AND r.started_at <= ?
		  AND (? = '' OR r.project_name = ?)
		GROUP BY r.project_name
		ORDER BY r.project_name`, fromS, toS, project, project)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []ErrorBudget
	for rows.Next() {
		var name string
		var slo sql.NullFloat64
		var deploys, failed int
		if rows.Scan(&name, &slo, &deploys, &failed) != nil {
			continue
		}
		rate := defaultSLOFailureRate
		if slo.Valid && slo.Float64 > 0 && slo.Float64 <= 1 {
			rate = slo.Float64
		}
		out = append(out, ScoreErrorBudget(name, deploys, failed, rate))
	}
	if project != "" && len(out) == 0 {
		out = append(out, ScoreErrorBudget(project, 0, 0, defaultSLOFailureRate))
	}
	return out
}

// BudgetDecision is the headline for a set of gateway budgets.
func BudgetDecision(rows []ErrorBudget) string {
	freeze, watch := 0, 0
	for _, row := range rows {
		switch row.Status {
		case "freeze":
			freeze++
		case "watch":
			watch++
		}
	}
	if freeze > 0 {
		return fmt.Sprintf("%d gateway(s) spent the error budget", freeze)
	}
	if watch > 0 {
		return fmt.Sprintf("%d gateway(s) are close to the error budget", watch)
	}
	return ""
}
