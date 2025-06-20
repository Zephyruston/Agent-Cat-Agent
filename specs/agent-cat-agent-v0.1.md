# Agent-Cat-Agent

**自动化系统**, “协作式 AI 工具链”, ai 可阅读项目代码然后根据用户需求编写单元测试代码/仅根据用户 prompt 生成代码, 然后将代码放入 docker 里面运行, 运行成功给用户反馈, 并提供完整代码

## apis

### aca test <prompt> <language>

根据用户的`test`prompt, 生成 unit test, 进行测试

```bash
aca test "给我<path>里面的<file>里面所有的函数生成单元测试" <language> # 系统应该先调用openai-go sdk 与大模型对话, 传递上下文, 等待响应的代码, 然后准备docker镜像, 用go-docker-sdk 将所有依赖的代码和单测放入docker里面运行, 用户需要指定编程语言
```

### aca gen <prompt> <language>

根据用户的`gen`prompt, 生成用户需要的代码, 放入 docker 运行

```bash
aca gen "给我生成计算斐波那契的函数" <language> # 系统应该先调用openai-go sdk 与大模型对话, 等待响应的代码, 然后准备docker镜像, 用go-docker-sdk 将所有依赖的代码放入docker里面运行, 用户需要指定编程语言
```

## Dependencies

### go-docker version: v28.2.2

examples:

build

```go
cmdBuild := exec.Command("docker", "build", "-t", imageName, "-f", dockerfilePath, ".")
```

run

```go
package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", image.PullOptions{})
	if err != nil {
		panic(err)
	}

	defer reader.Close()
	// cli.ImagePull is asynchronous.
	// The reader needs to be read completely for the pull operation to complete.
	// If stdout is not required, consider using io.Discard instead of os.Stdout.
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}
```

print log

```go
package main

import (
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	options := container.LogsOptions{ShowStdout: true}
	// Replace this ID with a container that really exists
	out, err := cli.ContainerLogs(ctx, "f1064a8a4c82", options)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)
}
```

stop container

```go
package main

import (
	"context"
	"fmt"

	containertypes "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, containertypes.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Print("Stopping container ", container.ID[:10], "... ")
		noWaitTimeout := 0 // to not wait for the container to exit gracefully
		if err := cli.ContainerStop(ctx, container.ID, containertypes.StopOptions{Timeout: &noWaitTimeout}); err != nil {
			panic(err)
		}
		fmt.Println("Success")
	}
}
```

### openai-go version: 1.6.0

chat-completion

```go
package main
import (
	"context"
	"github.com/openai/openai-go"
)
func main() {
	client := openai.NewClient()
	ctx := context.Background()
	question := "Write me a haiku"
	print("> ")
	println(question)
	println()
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT4o,
	}
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}
	println(completion.Choices[0].Message.Content)
}
```

chat-completion-accumulating

```go
package main
import (
	"context"
	"github.com/openai/openai-go"
)
func main() {
	client := openai.NewClient()
	ctx := context.Background()
	sysprompt := "Share only a brief description of the place in 50 words. Then immediately make some tool calls and announce them."
	question := "Tell me about Greece's largest city."
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(sysprompt),
		openai.UserMessage(question),
	}
	print("> ")
	println(question)
	println()
	params := openai.ChatCompletionNewParams{
		Messages: messages,
		Seed:     openai.Int(0),
		Model:    openai.ChatModelGPT4o,
		Tools:    tools,
	}
	stream := client.Chat.Completions.NewStreaming(ctx, params)
	acc := openai.ChatCompletionAccumulator{}
	for stream.Next() {
		chunk := stream.Current()
		acc.AddChunk(chunk)
		// When this fires, the current chunk value will not contain content data
		if _, ok := acc.JustFinishedContent(); ok {
			println()
			println("finish-event: Content stream finished")
		}
		if refusal, ok := acc.JustFinishedRefusal(); ok {
			println()
			println("finish-event: refusal stream finished:", refusal)
			println()
		}
		if tool, ok := acc.JustFinishedToolCall(); ok {
			println("finish-event: tool call stream finished:", tool.Index, tool.Name, tool.Arguments)
		}
		// It's best to use chunks after handling JustFinished events.
		// Here we print the delta of the content, if it exists.
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			print(chunk.Choices[0].Delta.Content)
		}
	}
	if err := stream.Err(); err != nil {
		panic(err)
	}
	if acc.Usage.TotalTokens > 0 {
		println("Total Tokens:", acc.Usage.TotalTokens)
	}
}
var tools = []openai.ChatCompletionToolParam{
	{
		Function: openai.FunctionDefinitionParam{
			Name:        "get_live_weather",
			Description: openai.String("Get weather at the given location"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]string{
						"type": "string",
					},
				},
				"required": []string{"location"},
			},
		},
	},
	{
		Function: openai.FunctionDefinitionParam{
			Name:        "get_population",
			Description: openai.String("Get population of a given town"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"town": map[string]string{
						"type": "string",
					},
					"nation": map[string]string{
						"type": "string",
					},
					"rounding": map[string]string{
						"type":        "integer",
						"description": "Nearest base 10 to round to, e.g. 1000 or 1000000",
					},
				},
				"required": []string{"town", "nation"},
			},
		},
	},
}
// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}
// Mock function to simulate population data retrieval
func getPopulation(town, nation string, rounding int) string {
	// In a real implementation, this function would call a population API
	return "Athens, Greece: 664,046"
}
```

chat-completion-tool-calling

```go
package main
import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/openai/openai-go"
)
func main() {
	client := openai.NewClient()
	ctx := context.Background()
	question := "What is the weather in New York City?"
	print("> ")
	println(question)
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		Tools: []openai.ChatCompletionToolParam{
			{
				Function: openai.FunctionDefinitionParam{
					Name:        "get_weather",
					Description: openai.String("Get weather at the given location"),
					Parameters: openai.FunctionParameters{
						"type": "object",
						"properties": map[string]interface{}{
							"location": map[string]string{
								"type": "string",
							},
						},
						"required": []string{"location"},
					},
				},
			},
		},
		Seed:  openai.Int(0),
		Model: openai.ChatModelGPT4o,
	}
	// Make initial chat completion request
	completion, err := client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}
	toolCalls := completion.Choices[0].Message.ToolCalls
	// Return early if there are no tool calls
	if len(toolCalls) == 0 {
		fmt.Printf("No function call")
		return
	}
	// If there is a was a function call, continue the conversation
	params.Messages = append(params.Messages, completion.Choices[0].Message.ToParam())
	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_weather" {
			// Extract the location from the function call arguments
			var args map[string]interface{}
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
			if err != nil {
				panic(err)
			}
			location := args["location"].(string)
			// Simulate getting weather data
			weatherData := getWeather(location)
			// Print the weather data
			fmt.Printf("Weather in %s: %s\n", location, weatherData)
			params.Messages = append(params.Messages, openai.ToolMessage(weatherData, toolCall.ID))
		}
	}
	completion, err = client.Chat.Completions.New(ctx, params)
	if err != nil {
		panic(err)
	}
	println(completion.Choices[0].Message.Content)
}
// Mock function to simulate weather data retrieval
func getWeather(location string) string {
	// In a real implementation, this function would call a weather API
	return "Sunny, 25°C"
}
```

responses

```go
package main
import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/responses"
)
func main() {
	client := openai.NewClient()
	ctx := context.Background()
	question := "Write me a haiku about computers"
	resp, err := client.Responses.New(ctx, responses.ResponseNewParams{
		Input: responses.ResponseNewParamsInputUnion{OfString: openai.String(question)},
		Model: openai.ChatModelGPT4,
	})
	if err != nil {
		panic(err)
	}
	println(resp.OutputText())
}
```
