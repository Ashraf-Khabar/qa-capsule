package healing

import "strings"

// Specialist is the single MCP that should run for a failure.
// The hub activates one specialist and reports only that problem.
type Specialist struct {
	ID      string
	Label   string
	Tool    string
	Problem string
}

// SpecialistFor picks the MCP that matches the failure category.
func SpecialistFor(category, testName, locator string) Specialist {
	loc := strings.TrimSpace(locator)
	switch category {
	case CategoryLocator:
		problem := "This failure comes from a locator. The element the test targets is not on the page."
		if loc != "" {
			problem = "This failure comes from a locator. The test cannot find " + loc + "."
		}
		return Specialist{ID: CategoryLocator, Label: "Locator MCP", Tool: "record_healing_intervention", Problem: problem}
	case CategoryStaleElement:
		return Specialist{ID: CategoryStaleElement, Label: "Stale element MCP", Tool: "heal_incident", Problem: "This failure comes from a stale element. The test kept a reference after the page changed."}
	case CategoryTimeout:
		return Specialist{ID: CategoryTimeout, Label: "Timeout MCP", Tool: "heal_incident", Problem: "This failure comes from a timeout. The step did not finish before the wait expired."}
	case CategoryAssertion:
		return Specialist{ID: CategoryAssertion, Label: "Assertion MCP", Tool: "propose_healing", Problem: "This failure comes from an assertion. The value on the page does not match what the test expects."}
	case CategoryNetwork:
		return Specialist{ID: CategoryNetwork, Label: "Network MCP", Tool: "heal_incident", Problem: "This failure comes from the network. The test could not reach the service it calls."}
	case CategoryScript:
		return Specialist{ID: CategoryScript, Label: "Script MCP", Tool: "heal_incident", Problem: "This failure comes from the test script. The code raised an error before the check could finish."}
	default:
		problem := "This failure does not match a locator, assertion, timeout, or network error."
		if name := strings.TrimSpace(testName); name != "" {
			problem = "The failure on \"" + name + "\" does not match a locator, assertion, timeout, or network error."
		}
		return Specialist{ID: CategoryUnknown, Label: "Script MCP", Tool: "heal_incident", Problem: problem}
	}
}
