package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	oldOut := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = oldOut }()
	f()
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func captureError(f func()) string {
	r, w, _ := os.Pipe()
	oldErr := os.Stderr
	os.Stderr = w
	defer func() { os.Stderr = oldErr }()
	f()
	w.Close()
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestInfoAndWarningAndError(t *testing.T) {
	out := captureOutput(func() {
		Info("hello", 123)
		Warning("warn", "msg")
	})
	err := captureError(func() {
		Error("err", 456)
	})
	if !strings.Contains(out, "[INFO] hello123") {
		t.Errorf("Info log not found: %q", out)
	}
	if !strings.Contains(out, "[WARN] warnmsg") {
		t.Errorf("Warning log not found: %q", out)
	}
	if !strings.Contains(err, "[ERROR] err456") {
		t.Errorf("Error log not found: %q", err)
	}
}
