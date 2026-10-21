package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	pytestFuncRe  = regexp.MustCompile(`^def (test_[A-Za-z0-9_]+)\s*\(`)
	pytestClassRe = regexp.MustCompile(`^class ([A-Za-z_][A-Za-z0-9_]*)\s*[\(:]`)
	pytestNodeRe  = regexp.MustCompile(`^[\w./-]+\.py(?:::[A-Za-z_][A-Za-z0-9_]*){1,3}$`)
)

// PytestNode builds a single pytest id such as tests/demo.py::TestLogin::test_logout.
// The file must be a relative pytest path. The function name comes from the incident title.
func PytestNode(fileHint, testName, source string) (string, error) {
	file := filepath.ToSlash(strings.TrimSpace(fileHint))
	file = strings.TrimPrefix(file, "./")
	file = strings.TrimPrefix(file, "/")
	if file == "" || strings.Contains(file, "..") || !strings.HasSuffix(strings.ToLower(file), ".py") {
		return "", fmt.Errorf("re-run supports a pytest file")
	}
	fn := pytestFunctionName(testName, source)
	if fn == "" {
		return "", fmt.Errorf("could not identify the pytest function")
	}
	node := file + "::" + fn
	if class := classContaining(source, fn); class != "" {
		node = file + "::" + class + "::" + fn
	}
	if !pytestNodeRe.MatchString(node) || strings.Contains(node, "..") {
		return "", fmt.Errorf("refusing this test id")
	}
	return node, nil
}

func pytestFunctionName(testName, source string) string {
	candidate := strings.TrimSpace(testName)
	if i := strings.LastIndex(candidate, ">"); i >= 0 {
		candidate = strings.TrimSpace(candidate[i+1:])
	}
	candidate = strings.Trim(candidate, "`\"' ")
	if strings.HasPrefix(candidate, "test_") && strings.Contains(source, "def "+candidate+"(") {
		return candidate
	}
	best := ""
	for _, line := range strings.Split(source, "\n") {
		m := pytestFuncRe.FindStringSubmatch(strings.TrimSpace(line))
		if m == nil {
			continue
		}
		if strings.Contains(testName, m[1]) && len(m[1]) > len(best) {
			best = m[1]
		}
	}
	return best
}

func classContaining(source, fn string) string {
	current := ""
	for _, line := range strings.Split(source, "\n") {
		if strings.HasPrefix(line, "class ") {
			if m := pytestClassRe.FindStringSubmatch(line); m != nil {
				current = m[1]
			}
			continue
		}
		if strings.HasPrefix(line, "def ") {
			current = ""
		}
		if strings.HasPrefix(strings.TrimSpace(line), "def "+fn+"(") {
			return current
		}
	}
	return ""
}

// PytestRunResult is the outcome of one test id. Passed is true only on exit code 0.
type PytestRunResult struct {
	Passed   bool   `json:"passed"`
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output"`
}

// RunPytestNode asks the configured runner to execute one pytest id.
// runnerURL comes from the server environment, not from the browser.
func RunPytestNode(ctx context.Context, runnerURL, token, node string) (PytestRunResult, error) {
	runnerURL = strings.TrimRight(strings.TrimSpace(runnerURL), "/")
	if runnerURL == "" {
		return PytestRunResult{}, fmt.Errorf("the test runner is not configured")
	}
	if !pytestNodeRe.MatchString(node) || strings.Contains(node, "..") {
		return PytestRunResult{}, fmt.Errorf("refusing this test id")
	}
	payload, err := json.Marshal(map[string]string{"node": node})
	if err != nil {
		return PytestRunResult{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, runnerURL+"/run", bytes.NewReader(payload))
	if err != nil {
		return PytestRunResult{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if strings.TrimSpace(token) != "" {
		req.Header.Set("X-Verify-Token", token)
	}
	client := &http.Client{Timeout: 160 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return PytestRunResult{}, fmt.Errorf("the test runner is unreachable")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if err != nil {
		return PytestRunResult{}, err
	}
	if resp.StatusCode != http.StatusOK {
		var failure struct {
			Error string `json:"error"`
		}
		_ = json.Unmarshal(body, &failure)
		if strings.TrimSpace(failure.Error) == "" {
			failure.Error = "the test runner rejected the request"
		}
		return PytestRunResult{}, fmt.Errorf("%s", failure.Error)
	}
	var result PytestRunResult
	if err := json.Unmarshal(body, &result); err != nil {
		return PytestRunResult{}, fmt.Errorf("the test runner returned an unreadable result")
	}
	return result, nil
}
