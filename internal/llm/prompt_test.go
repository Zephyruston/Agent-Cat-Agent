package llm

import "testing"

func TestCodeGenPrompt(t *testing.T) {
	if got := CodeGenPrompt("test"); got == "" {
		t.Error("should not be empty")
	}
}

func TestTestGenPrompt(t *testing.T) {
	if got := TestGenPrompt("code"); got == "" {
		t.Error("should not be empty")
	}
}
