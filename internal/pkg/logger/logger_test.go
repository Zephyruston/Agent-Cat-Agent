package logger

import "testing"

func TestInfo(t *testing.T) {
	Info("test info log", "key", "value")
}

func TestError(t *testing.T) {
	Error("test error log", "err", "something wrong")
}
