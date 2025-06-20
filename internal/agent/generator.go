package agent

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Zephyruston/Agent-Cat-Agent/internal/container"
	"github.com/Zephyruston/Agent-Cat-Agent/internal/llm"
)

type Generator struct {
	OpenAI *llm.OpenAIClient
	Docker *container.DockerClient
}

func NewGenerator(openai *llm.OpenAIClient, docker *container.DockerClient) *Generator {
	return &Generator{OpenAI: openai, Docker: docker}
}

// GenerateCode 根据 prompt 生成代码（兼容单文件，返回主文件内容）
func (g *Generator) GenerateCode(ctx context.Context, prompt, language, model string) (string, error) {
	return g.OpenAI.RawChatCompletion(ctx, prompt, language, model)
}

// RunCodeInDocker 支持传递所有 mainFiles 和 depFiles，go run 时全部传递
func (g *Generator) RunCodeInDocker(ctx context.Context, language, workDir string, mainFiles, depFiles []string, targetDir string) (string, error) {
	image := "golang:1.24.0"
	if language != "go" {
		// 保持原有逻辑
		if language == "python" {
			cmd := []string{"python", mainFiles[0]}
			return g.Docker.RunContainerWithMount(ctx, image, cmd, workDir, targetDir)
		}
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	// Go 多 main 文件智能运行
	var result strings.Builder
	if len(mainFiles) == 1 {
		// 单主程序，全部一起 run
		cmd := []string{"go", "run"}
		cmd = append(cmd, relPaths(mainFiles, workDir)...) // main
		cmd = append(cmd, relPaths(depFiles, workDir)...)  // dep
		out, err := g.Docker.RunContainerWithMount(ctx, image, cmd, workDir, targetDir)
		return out, err
	}
	// 多主程序，先尝试全部 run
	cmdAll := []string{"go", "run"}
	cmdAll = append(cmdAll, relPaths(mainFiles, workDir)...)
	cmdAll = append(cmdAll, relPaths(depFiles, workDir)...)
	out, err := g.Docker.RunContainerWithMount(ctx, image, cmdAll, workDir, targetDir)
	if err == nil || !strings.Contains(out+err.Error(), "main redeclared") {
		// 成功或不是重复 main 错误，直接返回
		return out, err
	}
	// 否则分别运行每个 mainFile
	for _, mf := range mainFiles {
		cmd := []string{"go", "run", relPath(mf, workDir)}
		cmd = append(cmd, relPaths(depFiles, workDir)...)
		out, err := g.Docker.RunContainerWithMount(ctx, image, cmd, workDir, targetDir)
		result.WriteString("==== Output for " + filepath.Base(mf) + " ====" + "\n")
		if err != nil {
			result.WriteString("[ERROR] " + err.Error() + "\n")
		}
		result.WriteString(out + "\n")
	}
	return result.String(), nil
}

// relPaths 批量转相对路径
func relPaths(files []string, workDir string) []string {
	var out []string
	for _, f := range files {
		out = append(out, relPath(f, workDir))
	}
	return out
}

// relPath 单个转相对路径
func relPath(path, workDir string) string {
	if abs, err := filepath.Abs(workDir); err == nil {
		relPath, _ := filepath.Rel(abs, path)
		return relPath
	}
	return path
}

// WriteFilesToDirAndCollectMain 写入多文件并返回所有 main 包和依赖 go 文件路径
func WriteFilesToDirAndCollectMain(files map[string]string, workDir, language string) ([]string, []string, error) {
	var mainFiles []string
	var depFiles []string
	for name, code := range files {
		if language != "go" {
			// 其他语言暂时全部写入 workDir
			path := filepath.Join(workDir, name)
			if err := os.WriteFile(path, []byte(code), 0644); err != nil {
				return nil, nil, err
			}
			mainFiles = append(mainFiles, path)
			continue
		}
		// Go: 检查 package
		pkg := "main"
		for _, line := range strings.Split(code, "\n") {
			if strings.HasPrefix(line, "package ") {
				pkg = strings.TrimSpace(strings.TrimPrefix(line, "package "))
				break
			}
		}
		if pkg == "main" {
			path := filepath.Join(workDir, name)
			if err := os.WriteFile(path, []byte(code), 0644); err != nil {
				return nil, nil, err
			}
			mainFiles = append(mainFiles, path)
		} else {
			dir := filepath.Join(workDir, pkg)
			os.MkdirAll(dir, 0755)
			path := filepath.Join(dir, name)
			if err := os.WriteFile(path, []byte(code), 0644); err != nil {
				return nil, nil, err
			}
			depFiles = append(depFiles, path)
		}
	}
	return mainFiles, depFiles, nil
}

func (g *Generator) GenerateCodeAndWriteFiles(ctx context.Context, prompt, language, model, workDir string, extractFiles func(string, string) map[string]string) (string, string, error) {
	content, err := g.OpenAI.RawChatCompletion(ctx, prompt, language, model)
	if err != nil {
		return "", "", err
	}
	files := extractFiles(content, MainFileName(language))
	mainFiles, _, err := WriteFilesToDirAndCollectMain(files, workDir, language)
	if err != nil {
		return "", "", err
	}
	if language == "go" {
		if err := PostProcessGoFiles(workDir); err != nil {
			return "", "", err
		}
	}
	return content, mainFiles[0], nil
}

// PostProcessGoFiles 自动调用 goimports 修复 import，仅处理 workDir 下所有 go 文件
func PostProcessGoFiles(workDir string) error {
	cmd := exec.Command("goimports", "-w", workDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func MainFileName(language string) string {
	if language == "python" {
		return "main.py"
	}
	return "main.go"
}
