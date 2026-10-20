package server

import (
	"strings"
	"testing"
)

func TestApplyLocatorFixToSource_RobotBrokenVariable(t *testing.T) {
	source := `*** Variables ***
${BROKEN_SUBMIT}    button[data-qa="submit-v1"]

*** Test Cases ***
TC-04 Broken
    Click        ${BROKEN_SUBMIT}
`
	got := applyLocatorFixToSource(source, `button[data-qa="submit-v1"]`, "#submit")
	if !strings.Contains(got, "Click        #submit") {
		t.Fatalf("expected fixed click line, got:\n%s", got)
	}
}

func TestApplyLocatorFixToSource_PythonBrokenConstant(t *testing.T) {
	source := "BROKEN_SUBMIT = '[data-testid=\"sign-in-button\"]'\ndriver.find_element(By.CSS_SELECTOR, BROKEN_SUBMIT).click()\n"
	got := applyLocatorFixToSource(source, `[data-testid="sign-in-button"]`, `button[type='submit']`)
	if !strings.Contains(got, `button[type='submit']`) || strings.Contains(got, `sign-in-button`) {
		t.Fatalf("expected the broken selector replaced, got:\n%s", got)
	}
	if !strings.Contains(got, `driver.find_element(By.CSS_SELECTOR, "button[type='submit']").click()`) {
		t.Fatalf("expected the call site rewritten, got:\n%s", got)
	}
}

func TestApplyLocatorFixToSource_PythonLogoutShowsTheCallLine(t *testing.T) {
	source := "BROKEN_LOGOUT = '[data-testid=\"logout-button\"]'\n" +
		"def test_broken_logout_after_login_self_healing_demo(self, driver):\n" +
		"    driver.find_element(By.CSS_SELECTOR, BROKEN_LOGOUT).click()\n"
	original := `[data-testid="logout-button"]`
	healed := `a[href="/logout"]`
	got := applyLocatorFixToSource(source, original, healed)
	if !strings.Contains(got, `BROKEN_LOGOUT = 'a[href="/logout"]'`) {
		t.Fatalf("constant not updated:\n%s", got)
	}
	if !strings.Contains(got, `driver.find_element(By.CSS_SELECTOR, 'a[href="/logout"]').click()`) {
		t.Fatalf("call site not updated:\n%s", got)
	}
	lines := buildLocatorChangeLines(source, got, original, healed)
	var sawCall bool
	for _, row := range lines {
		after, _ := row["after"].(string)
		if strings.Contains(after, `find_element`) && strings.Contains(after, healed) {
			sawCall = true
			if row["line"] != 3 {
				t.Fatalf("expected the click on line 3, got %#v", row)
			}
		}
	}
	if !sawCall {
		t.Fatalf("change lines missing the find_element edit: %#v", lines)
	}
}

func TestBuildLocatorPrompt_IncludesTestIntentAndPage(t *testing.T) {
	html := `<html><head><title>Login</title></head><body><h1>Sign in</h1><button data-test="login-button">Login</button></body></html>`
	src := "test('user can login', async ({ page }) => {\n  await page.locator('button[data-qa=\"submit-v1\"]').click();\n});"
	prompt := buildLocatorPrompt(`button[data-qa="submit-v1"]`, "Playwright", html, "waiting for locator", "user can login", src)
	for _, want := range []string{"TEST SOURCE", "user can login", "login-button", "PAGE CONTEXT", "intent"} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt missing %q", want)
		}
	}
}

func TestIntentFromExplanation(t *testing.T) {
	got := intentFromExplanation("Intent: submit the login form. The login button is data-test=login-button.")
	if got != "submit the login form" {
		t.Fatalf("got %q", got)
	}
}

func TestBuildFixHighlights_RobotFix(t *testing.T) {
	before := "    Click        ${BROKEN_SUBMIT}\n"
	after := "    Click        #submit\n"
	hl := buildFixHighlights(before, after, "#submit")
	if len(hl) != 1 || hl[0]["line"] != 1 || hl[0]["token"] != "#submit" {
		t.Fatalf("unexpected highlights: %#v", hl)
	}
}
