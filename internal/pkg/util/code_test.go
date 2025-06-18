package util

import "testing"

func TestExtractGoCode(t *testing.T) {
	md := "这是描述\n```go\npackage main\nfunc main(){}\n```\n后面还有内容"

	code := ExtractGoCode(md)
	if code == "" || code[:7] != "package" {
		t.Errorf("提取失败，got: %s", code)
	}
}
