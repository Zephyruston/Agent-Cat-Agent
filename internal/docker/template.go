package docker

import (
	"bytes"
	"os"
	"text/template"
)

// RenderDockerfile 渲染用户自定义.tql模板，data为上下文参数
func RenderDockerfile(tplPath string, data any) (string, error) {
	tplBytes, err := os.ReadFile(tplPath)
	if err != nil {
		return "", err
	}
	tpl, err := template.New("dockerfile").Parse(string(tplBytes))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
