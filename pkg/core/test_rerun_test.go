package core

import "testing"

func TestPytestNode_ClassMethodFromIncidentTitle(t *testing.T) {
	source := "class TestFormAuthentication:\n    def test_broken_logout_after_login_self_healing_demo(self, driver):\n        pass\n"
	name := "[Selenium] tests.selenium-pytest.test_selenium_sample.TestFormAuthentication > test_broken_logout_after_login_self_healing_demo"
	got, err := PytestNode("tests/selenium-pytest/test_selenium_sample.py", name, source)
	if err != nil {
		t.Fatal(err)
	}
	want := "tests/selenium-pytest/test_selenium_sample.py::TestFormAuthentication::test_broken_logout_after_login_self_healing_demo"
	if got != want {
		t.Fatalf("got %s", got)
	}
}

func TestPytestNode_RejectsOutsidePytest(t *testing.T) {
	if _, err := PytestNode("../pkg/api/server/healing_handlers.go", "test_x", "def test_x():\n"); err == nil {
		t.Fatal("expected a non-pytest path to be refused")
	}
	if _, err := PytestNode("tests/playwright/checkout.spec.js", "should pay", "test('should pay')"); err == nil {
		t.Fatal("expected a javascript file to be refused")
	}
}
