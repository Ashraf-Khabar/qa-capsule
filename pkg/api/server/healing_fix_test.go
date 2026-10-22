package server

import (
	"strings"
	"testing"
)

func TestFinalizeLocatorProposal_CheckoutWithoutCapturedPage(t *testing.T) {
	source := "...    Real checkout button: css=[data-test=\"checkout\"]\n" +
		"...    Broken locator (TC-04): css=[data-test=\"proceed-to-payment\"]\n" +
		"${BROKEN_CHECKOUT}  css=[data-test=\"proceed-to-payment\"]\n"
	loginPage := `<html><body><input id="user-name"><input id="password"><input id="login-button" type="submit"></body></html>`
	// A fetched homepage is not the failure page, so it must not be passed in.
	healed, conf, expl := finalizeLocatorProposal(
		`[data-test="proceed-to-payment"]`, source, "", "Robot", "", 0, "",
	)
	if healed != `css=[data-test="checkout"]` && healed != `[data-test="checkout"]` {
		t.Fatalf("expected the documented checkout selector, got %q (%s)", healed, expl)
	}
	if conf >= 0.85 {
		t.Fatalf("unverified selector must not look confirmed, confidence=%v", conf)
	}
	if !strings.Contains(expl, "not verified") {
		t.Fatalf("explanation should say the page was not stored, got %q", expl)
	}
	// The same selector is confirmed only when the captured cart page contains it.
	cart := `<html><body><button id="checkout" data-test="checkout">Checkout</button></body></html>`
	healed, conf, expl = finalizeLocatorProposal(
		`[data-test="proceed-to-payment"]`, source, cart, "Robot", "", 0, "",
	)
	if healed == "" || conf < 0.9 {
		t.Fatalf("cart page should verify checkout, got %q conf=%v expl=%q", healed, conf, expl)
	}
	_ = loginPage
}

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
