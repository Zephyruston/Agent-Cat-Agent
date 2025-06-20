package agent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteFilesToDirAndCollectMain(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"main.go":   "package main\nfunc main(){}",
		"util.go":   "package util\nfunc Foo(){}",
		"helper.go": "package helper\nfunc Bar(){}",
	}
	mainFiles, depFiles, err := WriteFilesToDirAndCollectMain(files, dir, "go")
	if err != nil {
		t.Fatal(err)
	}
	if len(mainFiles) != 1 || len(depFiles) != 2 {
		t.Errorf("expect 1 main, 2 dep, got %d, %d", len(mainFiles), len(depFiles))
	}
	for _, f := range mainFiles {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("main file not exist: %s", f)
		}
	}
	for _, f := range depFiles {
		if _, err := os.Stat(f); err != nil {
			t.Errorf("dep file not exist: %s", f)
		}
		utilDir := filepath.Join(dir, "util")
		helperDir := filepath.Join(dir, "helper")
		fClean := filepath.Clean(f)
		if !strings.HasPrefix(fClean, utilDir+string(os.PathSeparator)) && !strings.HasPrefix(fClean, helperDir+string(os.PathSeparator)) {
			t.Errorf("dep file not in correct dir: %s", f)
		}
	}
}
