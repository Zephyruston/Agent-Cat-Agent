package llm

import (
	"context"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIConfig struct {
	APIKey  string
	Model   string
	BaseURL string // 可选
}

type OpenAIClient struct {
	cli   openai.Client
	model string
}

func NewOpenAIClient(cfg OpenAIConfig) *OpenAIClient {
	opts := []option.RequestOption{option.WithAPIKey(cfg.APIKey)}
	if cfg.BaseURL != "" {
		opts = append(opts, option.WithBaseURL(cfg.BaseURL))
	}
	return &OpenAIClient{
		cli:   openai.NewClient(opts...),
		model: cfg.Model,
	}
}

func (c *OpenAIClient) Completion(ctx context.Context, prompt string) (string, error) {
	resp, err := c.cli.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("你是golang开发专家,请根据用户需求,生成符合golang最佳实践的代码"),
			openai.UserMessage(prompt),
		},
		Model: c.model,
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}
