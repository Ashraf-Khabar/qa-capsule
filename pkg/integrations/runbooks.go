package integrations

import (
	"sort"
	"strings"
	"time"
)

// RunbookTemplate is a curated remediation DAG users can apply to a CI gateway.
type RunbookTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Plugins     []string `json:"required_plugins"`
	Document    WorkflowDocument
}

// RunbookCatalog returns built-in self-healing templates (registry paths only).
func RunbookCatalog() []RunbookTemplate {
	return []RunbookTemplate{
		template502Restart(),
		templateFlakyNotify(),
		templatePerfAlert(),
		templateOOMRestart(),
		templateTimeoutWebhook(),
	}
}

// GetRunbookTemplate finds a template by id.
func GetRunbookTemplate(id string) (*RunbookTemplate, bool) {
	for _, t := range RunbookCatalog() {
		if t.ID == id {
			cp := t
			return &cp, true
		}
	}
	return nil, false
}

// SuggestRunbook picks a catalog template from the incident name and error text.
// An empty result means no template matches; the caller must not invent one.
func SuggestRunbook(name, errText string) *RunbookTemplate {
	blob := strings.ToLower(name + "\n" + errText)
	id := ""
	switch {
	case strings.Contains(blob, "502"):
		id = "502-restart-pod"
	case strings.Contains(blob, "out of memory") || strings.Contains(blob, "oomkilled") || strings.Contains(blob, " oom"):
		id = "oom-restart"
	case strings.Contains(blob, "[flaky]") || strings.Contains(blob, "flaky"):
		id = "flaky-triage"
	case strings.Contains(blob, "timeout") || strings.Contains(blob, "timed out"):
		id = "timeout-cache-flush"
	case strings.Contains(blob, "regression") || strings.Contains(blob, "latency"):
		id = "perf-regression"
	default:
		return nil
	}
	t, ok := GetRunbookTemplate(id)
	if !ok {
		return nil
	}
	return t
}

// PlanRunbook dry-runs the template against one incident and returns the action
// labels the plugin engine would take. An empty list means the conditions do not match.
// No vendor is called.
func PlanRunbook(t *RunbookTemplate, name, errText, status string) []string {
	if t == nil {
		return nil
	}
	doc := CloneRunbookDocument(t)
	plan := NewWorkflowEngine(nil).PlanExecution(doc, WorkflowContext{
		Incident: IncidentContext{Name: name, Error: errText, Status: status},
	})
	labels := map[string]string{}
	for _, n := range doc.Nodes {
		if n.Type != "action" || n.FilePath == "" {
			continue
		}
		label := strings.TrimSpace(n.Label)
		if label == "" {
			label = n.Integration
		}
		if label == "" {
			label = n.FilePath
		}
		labels[n.FilePath] = label
	}
	var steps []string
	seen := map[string]bool{}
	for _, path := range plan.Actions {
		label := labels[path]
		if label == "" {
			label = path
		}
		if seen[label] {
			continue
		}
		seen[label] = true
		steps = append(steps, label)
	}
	return steps
}

// ActionLabels lists the action steps of a template in a stable order.
func ActionLabels(t *RunbookTemplate) []string {
	if t == nil {
		return nil
	}
	var steps []string
	for _, n := range t.Document.Nodes {
		if n.Type == "action" {
			label := strings.TrimSpace(n.Label)
			if label == "" {
				label = n.Integration
			}
			if label != "" {
				steps = append(steps, label)
			}
		}
	}
	sort.Strings(steps)
	return steps
}

// CloneRunbookDocument returns a copy ready to persist (fresh meta timestamp).
func CloneRunbookDocument(t *RunbookTemplate) *WorkflowDocument {
	if t == nil {
		return nil
	}
	doc := t.Document
	doc.Enabled = true
	doc.Version = WorkflowSchemaVersion
	doc.Meta.Name = t.Name
	doc.Meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return &doc
}

func template502Restart() RunbookTemplate {
	return RunbookTemplate{
		ID:          "502-restart-pod",
		Name:        "502 Bad Gateway → K8s restart",
		Description: "When error text mentions HTTP 502, rollout-restart the deployment and notify Slack.",
		Tags:        []string{"502", "gateway", "k8s", "http"},
		Plugins:     []string{"k8s/k8s-restart.json", "slack/slack-notifier.json"},
		Document: WorkflowDocument{
			Version: 1,
			Enabled: true,
			Entry:   "start",
			Nodes: map[string]WorkflowNode{
				"start": {Type: "trigger", Label: "Incident ingested"},
				"cond502": {
					Type:  "condition",
					Label: "502 in error",
					When:  &ConditionExpr{Op: "text", Field: "incident.error", Match: "contains", Value: "502"},
				},
				"k8s": {
					Type:        "action",
					Label:       "Restart deployment",
					FilePath:    "k8s/k8s-restart.json",
					Integration: "k8s",
				},
				"slack": {
					Type:        "action",
					Label:       "Slack notify",
					FilePath:    "slack/slack-notifier.json",
					Integration: "slack",
				},
			},
			Edges: []WorkflowEdge{
				{From: "start", To: "cond502"},
				{From: "cond502", To: "k8s", When: "true"},
				{From: "k8s", To: "slack"},
			},
		},
	}
}

func templateFlakyNotify() RunbookTemplate {
	return RunbookTemplate{
		ID:          "flaky-triage",
		Name:        "Flaky test → Slack + Jira",
		Description: "Routes [FLAKY] incidents to Slack and opens a Jira ticket.",
		Tags:        []string{"flaky", "quality"},
		Plugins:     []string{"slack/slack-notifier.json", "jira/jira-ticket.json"},
		Document: WorkflowDocument{
			Version: 1,
			Enabled: true,
			Entry:   "start",
			Nodes: map[string]WorkflowNode{
				"start": {Type: "trigger"},
				"cond": {
					Type:  "condition",
					Label: "Flaky prefix",
					When:  &ConditionExpr{Op: "tag", Match: "prefix", Value: "[FLAKY]", Field: "incident.name"},
				},
				"slack": {Type: "action", FilePath: "slack/slack-notifier.json", Integration: "slack", Label: "Slack"},
				"jira":  {Type: "action", FilePath: "jira/jira-ticket.json", Integration: "jira", Label: "Jira"},
			},
			Edges: []WorkflowEdge{
				{From: "start", To: "cond"},
				{From: "cond", To: "slack", When: "true"},
				{From: "slack", To: "jira"},
			},
		},
	}
}

func templatePerfAlert() RunbookTemplate {
	return RunbookTemplate{
		ID:          "perf-regression",
		Name:        "Performance regression → Datadog + Slack",
		Description: "Notifies on [PERF] degradation incidents.",
		Tags:        []string{"perf", "latency"},
		Plugins:     []string{"datadog/datadog-event.json", "slack/slack-notifier.json"},
		Document: WorkflowDocument{
			Version: 1,
			Enabled: true,
			Entry:   "start",
			Nodes: map[string]WorkflowNode{
				"start": {Type: "trigger"},
				"cond": {
					Type:  "condition",
					When:  &ConditionExpr{Op: "tag", Match: "prefix", Value: "[PERF]", Field: "incident.name"},
					Label: "Perf prefix",
				},
				"dd":    {Type: "action", FilePath: "datadog/datadog-event.json", Integration: "datadog", Label: "Datadog event"},
				"slack": {Type: "action", FilePath: "slack/slack-notifier.json", Integration: "slack", Label: "Slack"},
			},
			Edges: []WorkflowEdge{
				{From: "start", To: "cond"},
				{From: "cond", To: "dd", When: "true"},
				{From: "dd", To: "slack"},
			},
		},
	}
}

func templateOOMRestart() RunbookTemplate {
	return RunbookTemplate{
		ID:          "oom-restart",
		Name:        "OOM / CrashLoop → K8s restart",
		Description: "Restarts the deployment when logs mention OOM or CrashLoopBackOff.",
		Tags:        []string{"oom", "k8s", "crashloop"},
		Plugins:     []string{"k8s/k8s-restart.json"},
		Document: WorkflowDocument{
			Version: 1,
			Enabled: true,
			Entry:   "start",
			Nodes: map[string]WorkflowNode{
				"start": {Type: "trigger"},
				"cond": {
					Type:  "condition",
					Label: "OOM or crash loop",
					When: &ConditionExpr{
						Op: "or",
						Or: []ConditionExpr{
							{Op: "text", Field: "incident.error", Match: "contains", Value: "oom"},
							{Op: "text", Field: "incident.error", Match: "contains", Value: "crashloop"},
						},
					},
				},
				"k8s": {Type: "action", FilePath: "k8s/k8s-restart.json", Integration: "k8s", Label: "K8s restart"},
			},
			Edges: []WorkflowEdge{
				{From: "start", To: "cond"},
				{From: "cond", To: "k8s", When: "true"},
			},
		},
	}
}

func templateTimeoutWebhook() RunbookTemplate {
	return RunbookTemplate{
		ID:          "timeout-cache-flush",
		Name:        "Timeout → custom webhook (cache flush)",
		Description: "Calls the custom webhook integration when errors mention timeout (wire your cache purge URL in plugin config).",
		Tags:        []string{"timeout", "cache", "webhook"},
		Plugins:     []string{"webhook/custom-webhook.json"},
		Document: WorkflowDocument{
			Version: 1,
			Enabled: true,
			Entry:   "start",
			Nodes: map[string]WorkflowNode{
				"start": {Type: "trigger"},
				"cond": {
					Type:  "condition",
					Label: "Timeout in error",
					When:  &ConditionExpr{Op: "text", Field: "incident.error", Match: "contains", Value: "timeout"},
				},
				"hook": {Type: "action", FilePath: "webhook/custom-webhook.json", Integration: "webhook", Label: "Cache flush webhook"},
			},
			Edges: []WorkflowEdge{
				{From: "start", To: "cond"},
				{From: "cond", To: "hook", When: "true"},
			},
		},
	}
}
