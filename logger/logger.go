package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	logFile *os.File
}

var logger Logger

func SetLogFile(f *os.File) {
	logger.logFile = f
}

func Debug(msg string) {
	fmt.Println(logger)
	logger.logFile.WriteString(msg)
}