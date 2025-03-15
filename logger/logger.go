package logger

import (
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
	logger.logFile.WriteString(msg)
}