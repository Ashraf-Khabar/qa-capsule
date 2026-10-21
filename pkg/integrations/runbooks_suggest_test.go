package integrations

import "testing"

func TestSuggestRunbook(t *testing.T) {
	if got := SuggestRunbook("login", "HTTP 502 from gateway"); got == nil || got.ID != "502-restart-pod" {
		t.Fatalf("502: %#v", got)
	}
	if got := SuggestRunbook("[FLAKY] checkout", "intermittent"); got == nil || got.ID != "flaky-triage" {
		t.Fatalf("flaky: %#v", got)
	}
	if SuggestRunbook("checkout", "element not found") != nil {
		t.Fatal("a locator miss is not a catalog runbook")
	}
	labels := ActionLabels(SuggestRunbook("api", "context deadline exceeded timeout"))
	if len(labels) == 0 {
		t.Fatal("timeout template should list action steps")
	}
}

func TestPlanRunbookFollowsTheMatchingBranch(t *testing.T) {
	tmpl, ok := GetRunbookTemplate("502-restart-pod")
	if !ok {
		t.Fatal("missing 502 template")
	}
	steps := PlanRunbook(tmpl, "gateway", "HTTP 502 from upstream", "FAILED")
	if len(steps) == 0 {
		t.Fatal("a 502 incident should walk the action branch")
	}
	flaky, _ := GetRunbookTemplate("flaky-triage")
	if got := PlanRunbook(flaky, "checkout", "element not found", "FAILED"); len(got) != 0 {
		t.Fatalf("a locator miss must not take the flaky branch: %v", got)
	}
}
