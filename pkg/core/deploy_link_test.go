package core

import "testing"

func TestLinkDeploy_PrefersTheIncidentRun(t *testing.T) {
	runs := []PipelineStamp{
		{Project: "api", RunID: "older", Commit: "aaa", StartedAt: "2026-10-21 09:00:00", Outcome: "success"},
		{Project: "api", ProjectID: "p1", RunID: "run-9", Commit: "deadbeef", Branch: "main", Outcome: "failure", StartedAt: "2026-10-21 10:00:00"},
	}
	got, ok := LinkDeploy(runs, "api", "run-9", "2026-10-21 10:06:00")
	if !ok {
		t.Fatal("expected a link")
	}
	if got.RunID != "run-9" || got.Commit != "deadbeef" || got.Branch != "main" {
		t.Fatalf("run: %+v", got)
	}
	if got.MinutesBefore < 5.9 || got.MinutesBefore > 6.1 {
		t.Fatalf("minutes: %v", got.MinutesBefore)
	}
	if !got.Close || got.ProjectID != "p1" {
		t.Fatalf("close/id: %+v", got)
	}
}

func TestLinkDeploy_UsesNearestEarlierRun(t *testing.T) {
	runs := []PipelineStamp{
		{Project: "api", RunID: "morning", StartedAt: "2026-10-21 08:00:00"},
		{Project: "api", RunID: "noon", Commit: "abc1234", Branch: "release", StartedAt: "2026-10-21 12:00:00", Outcome: "success"},
		{Project: "other", RunID: "ignore", StartedAt: "2026-10-21 12:30:00"},
	}
	got, ok := LinkDeploy(runs, "api", "", "2026-10-21 12:20:00")
	if !ok || got.RunID != "noon" || !got.Close {
		t.Fatalf("nearest: %+v ok=%v", got, ok)
	}
}

func TestLinkDeploy_FarRunIsNotClose(t *testing.T) {
	runs := []PipelineStamp{{Project: "api", RunID: "old", StartedAt: "2026-10-20 08:00:00"}}
	got, ok := LinkDeploy(runs, "api", "old", "2026-10-21 10:00:00")
	if !ok || got.Close {
		t.Fatalf("far: %+v ok=%v", got, ok)
	}
}

func TestLinkDeploy_IgnoresALaterRun(t *testing.T) {
	runs := []PipelineStamp{{Project: "api", RunID: "later", StartedAt: "2026-10-21 18:00:00"}}
	if _, ok := LinkDeploy(runs, "api", "", "2026-10-21 10:00:00"); ok {
		t.Fatal("a later run must not attach")
	}
}
