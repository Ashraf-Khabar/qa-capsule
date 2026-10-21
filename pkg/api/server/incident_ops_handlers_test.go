package server

import (
	"testing"

	"github.com/SRE-Capsule/sre-capsule-community/pkg/core"
)

func TestPolishPostmortemKeepsDraftWithoutAI(t *testing.T) {
	prev := core.AIService
	core.AIService = nil
	t.Cleanup(func() { core.AIService = prev })
	draft := "Postmortem draft — cupress / logout (#5)\n\nTimeline\n- Detected: 2026-10-21 10:00:00"
	if got := polishPostmortem(draft); got != draft {
		t.Fatalf("draft changed: %q", got)
	}
}
