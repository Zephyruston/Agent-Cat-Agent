package logger

import (
	"fmt"
	"os"
)

func Info(args ...interface{}) {
	fmt.Fprintln(os.Stdout, "[INFO]", fmt.Sprint(args...))
}

func Warning(args ...interface{}) {
	fmt.Fprintln(os.Stdout, "[WARN]", fmt.Sprint(args...))
}

func Error(args ...interface{}) {
	fmt.Fprintln(os.Stderr, "[ERROR]", fmt.Sprint(args...))
}
