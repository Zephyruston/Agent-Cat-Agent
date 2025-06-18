package docker

import (
	"os"
	"testing"
)

func TestRenderDockerfile(t *testing.T) {
	tpl := "FROM {{.BaseImage}}\nCMD [\"{{.Cmd}}\"]\n"
	tplPath := "test.Dockerfile.tql"
	_ = os.WriteFile(tplPath, []byte(tpl), 0644)
	defer os.Remove(tplPath)
	data := map[string]any{"BaseImage": "alpine", "Cmd": "echo hello"}
	out, err := RenderDockerfile(tplPath, data)
	if err != nil || out == "" {
		t.Fatalf("render failed: %v, out=%s", err, out)
	}
}
