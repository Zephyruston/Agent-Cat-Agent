package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	cfg := Load()
	if cfg == nil {
		t.Fatal("config should not be nil")
	}
}
