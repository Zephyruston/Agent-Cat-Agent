package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ApiKey  string `yaml:"api_key"`
	BaseUrl string `yaml:"base_url"`
	Model   string `yaml:"model"`
}

func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read config.yaml: %v", err)
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Fatalf("failed to parse config.yaml: %v", err)
		return nil, err
	}
	return cfg, nil
}
