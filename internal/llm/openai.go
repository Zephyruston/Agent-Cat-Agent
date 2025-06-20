package llm

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Zephyruston/Agent-Cat-Agent/internal/config"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type OpenAIClient struct {
	client openai.Client
}

func NewOpenAIClient(cfg config.Config) *OpenAIClient {
	return &OpenAIClient{
		client: openai.NewClient(option.WithAPIKey(cfg.ApiKey), option.WithBaseURL(cfg.BaseUrl)),
	}
}

// ExtractCodeFilesFromLLMResponse 只提取 markdown 代码块内容，依次保存为 main.go/main.py、main1.go 等
func ExtractCodeFilesFromLLMResponse(content, defaultFile string) map[string]string {
	files := make(map[string]string)
	// 优先提取带语言标签的代码块
	reLang := regexp.MustCompile("(?s)```(go|python)\\n(.*?)```")
	matches := reLang.FindAllStringSubmatch(content, -1)
	if len(matches) > 0 {
		for i, m := range matches {
			code := strings.TrimSpace(m[2])
			fname := defaultFile
			if i > 0 {
				dot := strings.LastIndex(defaultFile, ".")
				if dot > 0 {
					fname = defaultFile[:dot] + fmt.Sprintf("%d", i+1) + defaultFile[dot:]
				} else {
					fname = defaultFile + fmt.Sprintf("%d", i+1)
				}
			}
			files[fname] = code
		}
		return files
	}
	// fallback: 提取第一个无标签代码块
	reNoLang := regexp.MustCompile("(?s)```\\n(.*?)```")
	matches2 := reNoLang.FindAllStringSubmatch(content, -1)
	if len(matches2) > 0 {
		files[defaultFile] = strings.TrimSpace(matches2[0][1])
		return files
	}
	// fallback: 没有代码块，全部写入 defaultFile
	files[defaultFile] = strings.TrimSpace(content)
	return files
}

// RawChatCompletion 返回原始 LLM 响应内容
func (c *OpenAIClient) RawChatCompletion(ctx context.Context, prompt, language, model string) (string, error) {
	systemPrompt := "你精通" + language + ", 请你根据用户需求生成代码, 如有多个文件请用注释标明文件名, 每个代码块以markdown单独输出, 例如 // main.go 放在代码块首行"
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(prompt),
		},
		Model: model,
	}
	completion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}
	return completion.Choices[0].Message.Content, nil
}
