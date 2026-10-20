package healing

import "strings"

const (
	CategoryTimeout      = "timeout"
	CategoryLocator      = "locator"
	CategoryAssertion    = "assertion"
	CategoryNetwork      = "network"
	CategoryStaleElement = "stale_element"
	CategoryScript       = "script_failure"
	CategoryUnknown      = "unknown"
)

// ClassifyError maps failure text to a framework-agnostic category.
func ClassifyError(err string) string {
	lower := strings.ToLower(err)
	switch {
	// Locator before timeout — "waiting for locator(...)" must not become timeout.
	case strings.Contains(lower, "locator"), strings.Contains(lower, "selector"), strings.Contains(lower, "getbyrole"),
		strings.Contains(lower, "getbytext"), strings.Contains(lower, "element not found"), strings.Contains(lower, "no such element"),
		strings.Contains(lower, "data-qa"), strings.Contains(lower, "data-testid"), strings.Contains(lower, "unable to find element"):
		return CategoryLocator
	case strings.Contains(lower, "stale element"), strings.Contains(lower, "staleelementreference"):
		return CategoryStaleElement
	case strings.Contains(lower, "timeout"), strings.Contains(lower, "timed out"):
		return CategoryTimeout
	case strings.Contains(lower, "waiting for"):
		return CategoryTimeout
	case strings.Contains(lower, "assert"), strings.Contains(lower, " should be"),
		strings.Contains(lower, "should equal"), strings.Contains(lower, "does not match"),
		strings.Contains(lower, " expected "), strings.HasPrefix(lower, "expected "):
		return CategoryAssertion
	case strings.Contains(lower, "econnrefused"), strings.Contains(lower, "network"), strings.Contains(lower, "fetch failed"),
		strings.Contains(lower, "connection refused"), strings.Contains(lower, "502"), strings.Contains(lower, "503"):
		return CategoryNetwork
	case strings.Contains(lower, "syntaxerror"), strings.Contains(lower, "referenceerror"),
		strings.Contains(lower, "typeerror"), strings.Contains(lower, "indentationerror"),
		strings.Contains(lower, "traceback (most recent call last)"):
		return CategoryScript
	default:
		return CategoryUnknown
	}
}

// SuggestedActions returns actionable hints for MCP agents (any test framework).
func SuggestedActions(category string) []string {
	switch category {
	case CategoryTimeout:
		return []string{
			"Increase wait timeout or use explicit wait-for-stable-state before assertion.",
			"Verify the target service responds within SLA in CI (not just locally).",
			"Check for race conditions: element visible but not yet interactive.",
		}
	case CategoryLocator:
		return []string{
			"Replace brittle CSS/XPath with role- or text-based selectors where possible.",
			"Confirm the element exists in the DOM at assertion time (screenshot/trace if available).",
			"Check for dynamic IDs, iframe nesting, or shadow DOM boundaries.",
		}
	case CategoryStaleElement:
		return []string{
			"Re-query the element after navigation or DOM refresh before interacting.",
			"Avoid holding element references across page transitions or AJAX updates.",
		}
	case CategoryAssertion:
		return []string{
			"Compare expected vs actual from the failure message; fix test data or application bug.",
			"Ensure test isolation: no dependency on execution order or shared state.",
		}
	case CategoryNetwork:
		return []string{
			"Verify API/base URL env vars match the CI environment.",
			"Add retry for transient upstream failures or mock external dependencies in CI.",
		}
	default:
		return []string{
			"Open Propose Fix on this incident to see the cause and the suggested change.",
			"Inspect the stack trace and console logs to locate the failing step.",
			"Apply the smallest fix, then re-run the failing test.",
		}
	}
}

func extractSelectorHint(err string) string {
	lower := strings.ToLower(err)
	for _, kw := range []string{"locator", "selector", "getbyrole", "getbytext", "xpath", "css", "data-testid"} {
		if idx := strings.Index(lower, kw); idx >= 0 {
			return truncate(err[idx:], 160)
		}
	}
	return ""
}

func buildMCPPrompt(incidentID int64, category string) string {
	spec := SpecialistFor(category, "", "")
	id := formatInt(incidentID)
	return strings.TrimSpace(strings.Join([]string{
		spec.Problem,
		spec.Label + " is active. Call only " + spec.Tool + " with incident_id=" + id + ".",
		"Fix only that problem. Leave every other step unchanged.",
	}, "\n"))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func formatInt(n int64) string {
	if n <= 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
