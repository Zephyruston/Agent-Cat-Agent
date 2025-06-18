package llm

import (
	"testing"
)

func TestOpenAIClient_Configurable(t *testing.T) {
	cfg := OpenAIConfig{APIKey: "test", Model: "gpt-4", BaseURL: "https://api.example.com"}
	client := NewOpenAIClient(cfg)
	if client.model != "gpt-4" {
		t.Errorf("model not set")
	}
	// 只测试配置和结构体可用性，不实际请求
}
