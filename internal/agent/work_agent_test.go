package agent

import "testing"

func TestWorkAgent_ProcessTask(t *testing.T) {
	agent := NewWorkAgent(2)
	result, err := agent.ProcessTask("task-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil || result.Output == "" {
		t.Fatal("expected non-empty result output")
	}
}
