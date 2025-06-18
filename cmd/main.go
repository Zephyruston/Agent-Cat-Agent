package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/zccccc01/Agent-Cat-Agent/internal/agent"
	"github.com/zccccc01/Agent-Cat-Agent/internal/docker"
	"github.com/zccccc01/Agent-Cat-Agent/internal/llm"
	"github.com/zccccc01/Agent-Cat-Agent/internal/notify"
	"github.com/zccccc01/Agent-Cat-Agent/internal/pkg/util"
)

func main() {
	// 初始化持久化
	err := agent.InitDB("agent.db")
	if err != nil {
		fmt.Println("DB init error:", err)
		os.Exit(1)
	}
	defer os.Remove("agent.db")

	// 初始化各核心模块
	queue := agent.NewTaskQueue(10)
	notifier := notify.NewConsoleNotifier()
	workAgent := agent.NewWorkAgent(10)
	worker := agent.NewWorker(queue)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	worker.Start(ctx)

	llmClient := llm.NewOpenAIClient(llm.OpenAIConfig{
		APIKey:  os.Getenv("OPENAI_API_KEY"),
		Model:   os.Getenv("OPENAI_MODEL"),
		BaseURL: os.Getenv("OPENAI_BASE_URL"),
	})
	dockerCli, _ := docker.GetDockerClient()
	containerMgr := docker.NewContainerManager(dockerCli)

	result, err := workAgent.ProcessTask("task-001")
	if err != nil {
		fmt.Println("ProcessTask error:", err)
		return
	}
	notifier.Send(&notify.Notification{Type: "info", TaskID: "task-001", Message: result.Output, Timestamp: time.Now()})

	time.Sleep(50 * time.Millisecond)
	taskStatus, _ := agent.GetTaskStatus("task-001")
	notifier.Send(&notify.Notification{Type: "status", TaskID: "task-001", Message: "status=" + taskStatus, Timestamp: time.Now()})

	// LLM调用示例
	llmResp, _ := llmClient.Completion(ctx, llm.CodeGenPrompt("实现斐波那契函数, 然后必须放入main函数中调用"))
	fmt.Println("LLM生成代码示例:", llmResp)
	code := util.ExtractGoCode(llmResp)
	if code == "" {
		fmt.Println("未提取到Go代码段")
		return
	}
	fmt.Println("提取到Go代码:\n", code)

	// 保存代码到main.go
	os.MkdirAll("agentcat/cmd", 0755)
	codePath := "agentcat/cmd/main.go"
	err = os.WriteFile(codePath, []byte(code), 0644)
	if err != nil {
		fmt.Println("保存Go代码失败:", err)
		return
	}

	// Dockerfile模板渲染
	dockerfileData := map[string]any{
		"BaseImage":    "golang:1.24.0",
		"ServicePath":  "./agentcat/cmd",
		"Output":       "agentcat_bin",
		"RuntimeImage": "debian@sha256:5484adc33b4c352c5a9f4c4ae295fc49aed1cb54a7a0712a1b29674fb6f4f10f",
		"CmdArray":     `"./agentcat_bin"`,
	}
	dockerfile, err := docker.RenderDockerfile("cmd/default.tpl", dockerfileData)
	if err != nil {
		fmt.Println("Dockerfile渲染失败:", err)
		return
	}
	fmt.Println("渲染后的Dockerfile:\n", dockerfile)

	// 保存Dockerfile
	dockerfilePath := "example.Dockerfile"
	err = os.WriteFile(dockerfilePath, []byte(dockerfile), 0644)
	if err != nil {
		fmt.Println("保存Dockerfile失败:", err)
		return
	}

	// 构建镜像并运行（需本地有docker环境）
	imageName := "agentcat-test:latest"
	cmdBuild := exec.Command("docker", "build", "-t", imageName, "-f", dockerfilePath, ".")
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	fmt.Println("开始构建docker镜像...")
	if err := cmdBuild.Run(); err != nil {
		fmt.Println("docker build失败:", err)
		return
	}
	fmt.Println("镜像构建完成，开始运行...")
	cmdRun := exec.Command("docker", "run", "--rm", imageName)
	output, err := cmdRun.CombinedOutput()
	if err != nil {
		fmt.Println("docker run失败:", err)
	}
	fmt.Println("容器输出:\n", string(output))

	_ = containerMgr // 仅演示集成
}
