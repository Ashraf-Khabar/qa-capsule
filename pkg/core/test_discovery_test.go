package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTestNameSearchTokens_RobotNestedSuite(t *testing.T) {
	tokens := TestNameSearchTokens("[RobotFramework] Practice Login > TC-04 Submit Button With Broken Locator")
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	found := false
	for _, tok := range tokens {
		if strings.Contains(tok, "TC-04") || strings.Contains(tok, "Submit Button") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected TC-04 or Submit Button in tokens: %v", tokens)
	}
}

func TestTestNameSearchTokens_Playwright(t *testing.T) {
	tokens := TestNameSearchTokens("[Playwright] checkout.spec.ts > should complete payment")
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
}

func TestFindTestSourceInRepo_WalksTestDirs(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "e2e", "checkout")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := `test('should pay with broken selector', async ({ page }) => {
  await page.locator('button[data-qa="pay-v1"]').click();
});`
	path := filepath.Join(dir, "checkout.spec.ts")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	got := FindTestSourceInRepo(root, "[Playwright] checkout.spec.ts > should pay with broken selector")
	if got == "" {
		t.Fatal("expected to find test source file")
	}
	if !strings.Contains(got, `button[data-qa="pay-v1"]`) {
		t.Errorf("expected broken selector in content, got: %q", got)
	}
}

func TestFindTestSourceInRepo_IgnoresAppSource(t *testing.T) {
	root := t.TempDir()
	web := filepath.Join(root, "web", "js")
	testsDir := filepath.Join(root, "tests", "selenium-pytest")
	if err := os.MkdirAll(web, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	app := "export function deleteHealingFailure() { /* healing selector */ }"
	if err := os.WriteFile(filepath.Join(web, "healing.js"), []byte(app), 0o644); err != nil {
		t.Fatal(err)
	}
	py := "def test_broken_submit_selector_self_healing_demo(self, driver):\n    BROKEN_SUBMIT = '[data-testid=\"sign-in-button\"]'\n"
	if err := os.WriteFile(filepath.Join(testsDir, "test_selenium_sample.py"), []byte(py), 0o644); err != nil {
		t.Fatal(err)
	}
	name := "[Selenium] tests.selenium-pytest.test_selenium_sample.TestFormAuthentication > test_broken_submit_selector_self_healing_demo"
	got := FindTestSourceInRepo(root, name)
	if !strings.Contains(got, "test_broken_submit_selector_self_healing_demo") || strings.Contains(got, "deleteHealingFailure") {
		t.Fatalf("expected the selenium test, got %q", got)
	}
}

func TestReadTestFileAtPaths_RelativePath(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "tests", "login")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	want := "Click    #submit"
	path := filepath.Join(sub, "login.robot")
	if err := os.WriteFile(path, []byte(want), 0o644); err != nil {
		t.Fatal(err)
	}
	got := ReadTestFileAtPaths("tests/login/login.robot", root)
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}

func TestWithinTestTree_AllowsTestsRejectsAppSource(t *testing.T) {
	root := t.TempDir()
	testFile := filepath.Join(root, "tests", "selenium-pytest", "test_selenium_sample.py")
	if err := os.MkdirAll(filepath.Dir(testFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(testFile, []byte("BROKEN_LOGOUT = 'x'\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !WithinTestTree(testFile) {
		t.Fatal("expected a file under tests/ to be allowed")
	}
	appFile := filepath.Join(root, "pkg", "api", "server", "healing_handlers.go")
	if err := os.MkdirAll(filepath.Dir(appFile), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(appFile, []byte("package server\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if WithinTestTree(appFile) {
		t.Fatal("application source must be refused")
	}
	if err := WriteTestFile(appFile, "package nope\n"); err == nil {
		t.Fatal("expected a write outside tests to fail")
	}
	next := "BROKEN_LOGOUT = 'a[href=\"/logout\"]'\n"
	if err := WriteTestFile(testFile, next); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != next {
		t.Fatalf("got %q", got)
	}
}

func TestResolveTestFile_UsesHintUnderTests(t *testing.T) {
	root := t.TempDir()
	rel := filepath.Join("tests", "selenium-pytest", "test_selenium_sample.py")
	abs := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "BROKEN_LOGOUT = '[data-testid=\"logout-button\"]'\n"
	if err := os.WriteFile(abs, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	path, content := ResolveTestFile(root, filepath.ToSlash(rel), "test_broken_logout")
	if path == "" || !strings.Contains(content, "BROKEN_LOGOUT") {
		t.Fatalf("path %q content %q", path, content)
	}
	if !WithinTestTree(path) {
		t.Fatalf("resolved outside the test tree: %s", path)
	}
}
