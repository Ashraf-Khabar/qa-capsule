package core

import "testing"

func TestScoreErrorBudget_QuietWithoutRuns(t *testing.T) {
	got := ScoreErrorBudget("api", 0, 0, 0.15)
	if got.Status != "quiet" || got.RemainingRatio != 1 {
		t.Fatalf("quiet: %+v", got)
	}
}

func TestScoreErrorBudget_HealthyUnderTarget(t *testing.T) {
	got := ScoreErrorBudget("api", 20, 1, 0.15)
	if got.Status != "healthy" {
		t.Fatalf("status: %s burn %v", got.Status, got.BurnRate)
	}
	if got.AllowedFailures < 2.9 || got.AllowedFailures > 3.1 {
		t.Fatalf("allowed: %v", got.AllowedFailures)
	}
	if got.RemainingRatio <= 0.5 {
		t.Fatalf("remaining: %v", got.RemainingRatio)
	}
}

func TestScoreErrorBudget_WatchNearTarget(t *testing.T) {
	got := ScoreErrorBudget("api", 10, 4, 0.5)
	if got.Status != "watch" {
		t.Fatalf("status: %s burn %v", got.Status, got.BurnRate)
	}
}

func TestScoreErrorBudget_FreezeWhenSpent(t *testing.T) {
	got := ScoreErrorBudget("api", 10, 2, 0.15)
	if got.Status != "freeze" || got.RemainingRatio != 0 {
		t.Fatalf("freeze: %+v", got)
	}
}

func TestScoreErrorBudget_DefaultSLO(t *testing.T) {
	got := ScoreErrorBudget("api", 10, 0, 0)
	if got.SLOFailureRate != 0.15 || got.Status != "healthy" {
		t.Fatalf("default slo: %+v", got)
	}
}
