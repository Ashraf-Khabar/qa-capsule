package core

import "strings"

// DOMCaptureHint returns setup guidance when no page HTML was stored for an incident.
func DOMCaptureHint(framework string) string {
	switch strings.ToLower(strings.TrimSpace(framework)) {
	case "robotframework", "robot":
		return "Robot: add --listener tests/robotframework/resources/dom_capture_listener.py to your robot command (already wired in e2e-tests-robot.yml)."
	case "cypress":
		return "Cypress: keep supportFile cypress/support/qa-capsule.js in cypress.config.js — it embeds the DOM snapshot in the JUnit failure output on each failed test."
	case "playwright":
		return "Playwright: import { test, expect } from fixtures/qa-capsule-fixture.js (not @playwright/test) so the DOM is captured on failure."
	case "selenium", "pytest":
		return "Selenium: tests/selenium-pytest/conftest.py captures driver.page_source on failure and writes it into the JUnit report."
	case "postman", "newman", "junit", "java":
		return "API suites (Newman, JUnit, Pytest API) do not run a browser — locator healing uses live URL fetch or error text only."
	default:
		return "E2E frameworks in this repo (Cypress, Playwright, Robot, Selenium) ship DOM capture hooks under tests/. Re-run the pipeline after verifying the hook is enabled."
	}
}

// IsBrowserFramework reports whether DOM capture at failure time is expected.
func IsBrowserFramework(framework string) bool {
	switch strings.ToLower(strings.TrimSpace(framework)) {
	case "robotframework", "robot", "cypress", "playwright", "selenium":
		return true
	default:
		return false
	}
}
