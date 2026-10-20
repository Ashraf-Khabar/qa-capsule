package healing

import (
	"strings"
	"testing"
)

func TestSpecialistFor_ActivatesOnlyTheMatchingMCP(t *testing.T) {
	locator := SpecialistFor(CategoryLocator, "TC-04", `[data-test="proceed-to-payment"]`)
	if locator.Label != "Locator MCP" || locator.Tool != "record_healing_intervention" {
		t.Fatalf("locator specialist = %+v", locator)
	}
	if locator.Problem == "" || !strings.Contains(locator.Problem, "proceed-to-payment") {
		t.Fatalf("locator problem = %q", locator.Problem)
	}

	assertion := SpecialistFor(CategoryAssertion, "TC-02", "")
	if assertion.Label != "Assertion MCP" || assertion.Tool == "record_healing_intervention" {
		t.Fatalf("assertion specialist = %+v", assertion)
	}
	if strings.Contains(strings.ToLower(assertion.Problem), "locator") {
		t.Fatalf("assertion problem mentions a locator: %q", assertion.Problem)
	}

	timeout := SpecialistFor(CategoryTimeout, "TC-01", "")
	if timeout.Label != "Timeout MCP" {
		t.Fatalf("timeout specialist = %+v", timeout)
	}
}

func TestClassifyError_ScriptAndLocatorStaySeparate(t *testing.T) {
	if got := ClassifyError("SyntaxError: Unexpected token"); got != CategoryScript {
		t.Fatalf("script = %q", got)
	}
	if got := ClassifyError("Timeout 5000ms exceeded waiting for locator('[data-test=\"proceed-to-payment\"]')"); got != CategoryLocator {
		t.Fatalf("locator timeout = %q", got)
	}
}
