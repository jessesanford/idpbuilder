package helpers

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

var (
	LogLevel    string = "info"
	LogLevelMsg        = "Set log level (debug, info, warn, error)"
	Verbose     bool
	VerboseMsg         = "Enable verbose output"
	logger      *slog.Logger
)

type LogLevelValue int

const (
	DebugLevel LogLevelValue = -4
	InfoLevel  LogLevelValue = 0
	WarnLevel  LogLevelValue = 4
	ErrorLevel LogLevelValue = 8
)

func InitializeLogging() {
	level := parseLogLevel(LogLevel)
	opts := &slog.HandlerOptions{
		Level:     slog.Level(level),
		AddSource: Verbose,
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, opts))
	slog.SetDefault(logger)
}

func LogDebug(format string, args ...interface{}) {
	if logger != nil {
		logger.Debug(fmt.Sprintf(format, args...))
	}
}

func LogInfo(format string, args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprintf(format, args...))
	}
}

func LogWarn(format string, args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprintf(format, args...))
	}
}

func LogError(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf(format, args...))
	}
}

func LogWarning(format string, args ...interface{}) {
	LogWarn(format, args...)
}

func parseLogLevel(level string) LogLevelValue {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info", "":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

func SetLogLevel(level string) {
	LogLevel = level
	InitializeLogging()
}

func GetLogger() *slog.Logger {
	return logger
}

