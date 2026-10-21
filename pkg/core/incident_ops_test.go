package core

import (
	"strings"
	"testing"
)

func TestComposePostmortemUsesStoredFacts(t *testing.T) {
	text := ComposePostmortem(PostmortemFacts{
		ID: 5, Project: "cupress", Name: "logout",
		Error: "locator missing", Severity: "sev2",
		CreatedAt: "2026-10-21 10:00:00", AckAt: "2026-10-21 10:04:00", AckBy: "akr",
		Mitigated: "2026-10-21 10:12:00", ResolvedAt: "2026-10-21 10:20:00", ResolvedBy: "akr",
		Deploy:   "abc1234 main failed, 6 min before",
		Runbook:  "Flaky test — Slack, Jira",
		Comments: []string{"<p>checked the page</p>"},
	})
	for _, want := range []string{"#5", "Detected: 2026-10-21 10:00:00", "sev2", "abc1234", "checked the page", "Slack, Jira"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in %s", want, text)
		}
	}
}

func TestNormalizeSeverity(t *testing.T) {
	if normalizeSeverity("SEV1") != "sev1" {
		t.Fatal(normalizeSeverity("SEV1"))
	}
	if normalizeSeverity("urgent") != "sev3" {
		t.Fatal("unknown severity must fall back to sev3")
	}
}
