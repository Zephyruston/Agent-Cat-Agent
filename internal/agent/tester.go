package agent

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Zephyruston/Agent-Cat-Agent/internal/container"
	"github.com/Zephyruston/Agent-Cat-Agent/internal/llm"
)

type Tester struct {
	OpenAI *llm.OpenAIClient
	Docker *container.DockerClient
}

func NewTester(openai *llm.OpenAIClient, docker *container.DockerClient) *Tester {
	return &Tester{OpenAI: openai, Docker: docker}
}

// GenerateTest 根据 prompt 生成单元测试代码（兼容单文件，返回主文件内容）
func (t *Tester) GenerateTest(ctx context.Context, prompt, language, model string) (string, error) {
	return t.OpenAI.RawChatCompletion(ctx, prompt, language, model)
}

// GenerateTestAndRun 生成测试代码、写入文件并用 Docker 执行
func (t *Tester) GenerateTestAndRun(ctx context.Context, prompt, language, model, workDir string) (string, error) {
	testCode, err := t.OpenAI.RawChatCompletion(ctx, prompt, language, model)
	if err != nil {
		return "", err
	}
	fileName := "main_test.go"
	if language == "python" {
		fileName = "test_main.py"
	}
	filePath := filepath.Join(workDir, fileName)
	if err := os.WriteFile(filePath, []byte(testCode), 0644); err != nil {
		return "", err
	}
	image := "golang:1.24.0"
	cmd := []string{"go", "test", "."}
	if language == "python" {
		image = "python:3.11"
		cmd = []string{"pytest", fileName}
	}
	if t.Docker == nil {
		return testCode, fmt.Errorf("docker client not initialized")
	}
	// 这里应有挂载 workDir 到容器的逻辑，略作伪实现
	_, err = t.Docker.RunContainer(ctx, image, cmd)
	return testCode, err
}

// RunTestInDocker 仅执行 Docker 运行（假定测试代码已写入 workDir/fileName）
func (t *Tester) RunTestInDocker(ctx context.Context, language, workDir, fileName string) (string, error) {
	image := "golang:1.24.0"
	cmd := []string{"go", "test", "."}
	if language == "python" {
		image = "python:3.11"
		cmd = []string{"pytest", fileName}
	}
	if t.Docker == nil {
		return "", fmt.Errorf("docker client not initialized")
	}
	// 这里应有挂载 workDir 到容器的逻辑，略作伪实现
	return t.Docker.RunContainer(ctx, image, cmd)
}
