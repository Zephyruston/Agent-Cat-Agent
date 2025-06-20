package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Zephyruston/Agent-Cat-Agent/internal/agent"
	"github.com/Zephyruston/Agent-Cat-Agent/internal/config"
	"github.com/Zephyruston/Agent-Cat-Agent/internal/container"
	"github.com/Zephyruston/Agent-Cat-Agent/internal/llm"
	"github.com/Zephyruston/Agent-Cat-Agent/internal/logger"
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "aca",
		Short: "Agent-Cat-Agent CLI",
		Run:   runRoot,
	}

	rootCmd.Flags().StringP("mode", "m", "gen", "gen/test")
	rootCmd.Flags().StringP("prompt", "p", "", "prompt for code or test generation")
	rootCmd.Flags().StringP("language", "l", "go", "programming language (default: go)")
	rootCmd.Flags().StringP("config", "c", "etc/config.yaml", "config file path")
	rootCmd.Flags().String("workdir", "./tmp", "working directory")
	rootCmd.Flags().String("mount", "/app", "container mount dir")
	rootCmd.MarkFlagRequired("prompt")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) {
	mode, _ := cmd.Flags().GetString("mode")
	prompt, _ := cmd.Flags().GetString("prompt")
	language, _ := cmd.Flags().GetString("language")
	configPath, _ := cmd.Flags().GetString("config")
	workDir, _ := cmd.Flags().GetString("workdir")
	mountDir, _ := cmd.Flags().GetString("mount")

	absWorkDir, err := filepath.Abs(workDir)
	if err != nil {
		fmt.Println("[ERROR] 解析工作目录绝对路径失败:", err)
		os.Exit(1)
	}
	os.MkdirAll(absWorkDir, 0755)
	logger.Info("创建/检查工作目录:", absWorkDir)
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Error("配置文件加载失败:", err)
		os.Exit(1)
	}
	ctx := context.Background()
	openaiClient := llm.NewOpenAIClient(*cfg)
	dockerClient, err := container.NewDockerClient()
	if err != nil {
		fmt.Println("[ERROR] Docker 初始化失败:", err)
		os.Exit(1)
	}
	defer dockerClient.Close()

	switch mode {
	case "gen":
		runGen(ctx, openaiClient, dockerClient, prompt, language, cfg.Model, absWorkDir, mountDir)
	case "test":
		runTest(ctx, openaiClient, dockerClient, prompt, language, cfg.Model, absWorkDir, mountDir)
	default:
		fmt.Println("[ERROR] --mode 仅支持 gen/test")
		os.Exit(1)
	}
}

func runGen(ctx context.Context, openaiClient *llm.OpenAIClient, dockerClient *container.DockerClient, prompt, language, model, absWorkDir, mountDir string) {
	generator := agent.NewGenerator(openaiClient, dockerClient)
	logger.Info("请求 LLM 生成代码...")
	content, _, err := generator.GenerateCodeAndWriteFiles(ctx, prompt, language, model, absWorkDir, llm.ExtractCodeFilesFromLLMResponse)
	mainFiles, depFiles, _ := agent.WriteFilesToDirAndCollectMain(llm.ExtractCodeFilesFromLLMResponse(content, agent.MainFileName(language)), absWorkDir, language)
	if err != nil {
		fmt.Println("[ERROR] 代码生成或写入失败:", err)
		os.Exit(1)
	}
	logger.Info("LLM 响应内容如下:\n====================\n", content, "\n====================")
	logger.Info("用 Docker 执行代码...")
	_, err = generator.RunCodeInDocker(ctx, language, absWorkDir, mainFiles, depFiles, mountDir)
	if err != nil {
		fmt.Println("[ERROR] Docker 执行失败:", err)
		os.Exit(1)
	}
	logger.Info("Docker 执行完成，详见上方输出")
}

func runTest(ctx context.Context, openaiClient *llm.OpenAIClient, dockerClient *container.DockerClient, prompt, language, model, absWorkDir, mountDir string) {
	tester := agent.NewTester(openaiClient, dockerClient)
	logger.Info("请求 LLM 生成单元测试代码...")
	testCode, err := tester.GenerateTest(ctx, prompt, language, model)
	if err != nil || testCode == "" {
		fmt.Println("[ERROR] 测试代码生成失败:", err)
		os.Exit(1)
	}
	logger.Info("LLM 响应的测试代码如下:\n====================\n", testCode, "\n====================")
	fileName := agent.MainFileName(language)
	filePath := filepath.Join(absWorkDir, fileName)
	if err := os.WriteFile(filePath, []byte(testCode), 0644); err != nil {
		fmt.Println("[ERROR] 写入文件失败:", err)
		os.Exit(1)
	}
	logger.Info("用 Docker 执行测试代码...")
	_, err = tester.RunTestInDocker(ctx, language, absWorkDir, fileName)
	if err != nil {
		fmt.Println("[ERROR] Docker 执行失败:", err)
		os.Exit(1)
	}
	logger.Info("Docker 执行完成，详见上方输出")
}
