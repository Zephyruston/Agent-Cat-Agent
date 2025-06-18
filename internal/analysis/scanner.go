package analysis

import (
	"go/parser"
	"go/token"
	"os"
)

type FileInfo struct {
	Path  string
	Funcs []string
	Types []string
}

type ProjectScanner struct{}

func NewProjectScanner() *ProjectScanner {
	return &ProjectScanner{}
}

func (s *ProjectScanner) ScanGoFile(path string) (*FileInfo, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	info := &FileInfo{Path: path}
	for _, decl := range f.Decls {
		// TODO: 收集函数/类型名
		_ = decl
	}
	return info, nil
}

func (s *ProjectScanner) ScanPythonFile(path string) (*FileInfo, error) {
	// TODO: Python AST解析，可用第三方库
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	return &FileInfo{Path: path}, nil
}
