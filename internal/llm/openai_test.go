package llm

import (
	"testing"
)

func TestExtractCodeFilesFromLLMResponse(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		defaultF string
		expect   map[string]string
	}{
		{
			name:     "go code block",
			input:    "```go\npackage main\nfunc main(){}\n```",
			defaultF: "main.go",
			expect:   map[string]string{"main.go": "package main\nfunc main(){}"},
		},
		{
			name:     "python code block",
			input:    "```python\ndef foo():\n    pass\n```",
			defaultF: "main.py",
			expect:   map[string]string{"main.py": "def foo():\n    pass"},
		},
		{
			name:     "no code block",
			input:    "hello world",
			defaultF: "main.go",
			expect:   map[string]string{"main.go": "hello world"},
		},
		{
			name:     "no lang code block",
			input:    "```\nfoo\nbar\n```",
			defaultF: "main.go",
			expect:   map[string]string{"main.go": "foo\nbar"},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := ExtractCodeFilesFromLLMResponse(c.input, c.defaultF)
			if len(got) != len(c.expect) {
				t.Errorf("expect %d files, got %d", len(c.expect), len(got))
			}
			for k, v := range c.expect {
				if got[k] != v {
					t.Errorf("file %s expect %q, got %q", k, v, got[k])
				}
			}
		})
	}
}
