package docker

import "github.com/zccccc01/Agent-Cat-Agent/internal/agent"

// DockerEngine接口
type DockerEngine interface {
	ExecuteCode(code, lang string) (*agent.TaskResult, error)
	ExecuteTests(code string, tests []string, lang string) (*agent.TaskResult, error)
}
