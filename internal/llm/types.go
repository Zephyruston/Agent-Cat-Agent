package llm

import "github.com/zccccc01/Agent-Cat-Agent/internal/agent"

// LLMService接口
type LLMService interface {
	GenerateCode(task *agent.Task, ctx *ProjectContext) (*LLMResponse, error)
	GenerateTests(task *agent.Task, ctx *ProjectContext) (*LLMResponse, error)
}

type LLMResponse struct {
	Code        string
	Tests       []string
	Explanation string
	Error       error
}

type ProjectContext struct {
	// TODO: 项目结构、依赖、上下文等
}
