package llm

func CodeGenPrompt(taskDesc string) string {
	return "请根据如下需求生成Go代码：" + taskDesc
}

func TestGenPrompt(code string) string {
	return "请为如下Go代码生成单元测试：" + code
}
