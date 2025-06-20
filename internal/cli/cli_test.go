package cli

import (
	"os"
	"testing"
)

func TestCliHelp(t *testing.T) {
	os.Args = []string{"aca", "--help"}
	defer func() { recover() }() // 捕获 os.Exit
	Execute()
}

func TestCliGenArgs(t *testing.T) {
	os.Args = []string{"aca", "--mode=gen", "--prompt=hello", "--language=go", "--config=internal/config/config.yaml", "--workdir=./tmp", "--mount=/app"}
	defer func() { recover() }() // 捕获 os.Exit
	Execute()
}

func TestCliTestArgs(t *testing.T) {
	os.Args = []string{"aca", "--mode=test", "--prompt=hello", "--language=go", "--config=internal/config/config.yaml", "--workdir=./tmp", "--mount=/app"}
	defer func() { recover() }() // 捕获 os.Exit
	Execute()
}

func TestCliMissingPrompt(t *testing.T) {
	os.Args = []string{"aca", "--mode=gen", "--language=go", "--config=etc/config.yaml"}
	defer func() { recover() }()
	Execute()
}

func TestCliInvalidMode(t *testing.T) {
	os.Args = []string{"aca", "--mode=foo", "--prompt=hello", "--language=go", "--config=etc/config.yaml"}
	defer func() { recover() }()
	Execute()
}

func TestCliUnknownFlag(t *testing.T) {
	os.Args = []string{"aca", "--unknown", "--prompt=hello", "--mode=gen", "--language=go", "--config=etc/config.yaml"}
	defer func() { recover() }()
	Execute()
}

func TestCliEmptyPrompt(t *testing.T) {
	os.Args = []string{"aca", "--mode=gen", "--prompt=", "--language=go", "--config=etc/config.yaml"}
	defer func() { recover() }()
	Execute()
}
