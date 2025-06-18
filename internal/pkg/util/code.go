package util

import (
	"regexp"
)

// ExtractGoCode 提取markdown中的Go代码块
func ExtractGoCode(md string) string {
	re := regexp.MustCompile("(?s)```go\\s*(.*?)```")
	matches := re.FindStringSubmatch(md)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}
