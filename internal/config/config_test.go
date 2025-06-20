package config

import "testing"

func TestLoadConfig(t *testing.T) {
	cfg, err := LoadConfig("../../etc/config.yaml")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if cfg.ApiKey == "" {
		t.Errorf("api_key is empty")
	}

	if cfg.BaseUrl == "" {
		t.Errorf("base_url is empty")
	}

	if cfg.Model == "" {
		t.Errorf("model is empty")
	}
	t.Logf("config loaded successfully, %+v", cfg)
}
