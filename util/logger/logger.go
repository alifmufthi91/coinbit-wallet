package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type logWriter struct {
	LogType string
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(time.Now().Local().Format("2006-01-02 15:04:05") + " [" + writer.LogType + "]: " + string(bytes))
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func Init() {
	InfoLogger = log.New(os.Stderr, "", 0)
	InfoLogger.SetOutput(logWriter{LogType: "INFO"})
	ErrorLogger = log.New(os.Stderr, "", 0)
	ErrorLogger.SetOutput(logWriter{LogType: "ERROR"})
	WarningLogger = log.New(os.Stderr, "", 0)
	WarningLogger.SetOutput(logWriter{LogType: "WARNING"})
}

func Info(format string, a ...interface{}) {
	InfoLogger.Printf(format, a...)
}

func Error(format string, a ...interface{}) {
	ErrorLogger.Printf(format, a...)
}

func Warn(format string, a ...interface{}) {
	WarningLogger.Printf(format, a...)
}
