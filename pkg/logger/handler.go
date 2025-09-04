// Package logger provides basic logging capabilities for the image builder
package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Level represents the logging level
type Level int

const (
	InfoLevel Level = iota
	WarnLevel
	ErrorLevel
)

// String returns the string representation of the logging level
func (l Level) String() string {
	switch l {
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides basic logging capabilities
type Logger struct {
	writer io.Writer
	level  Level
}

// New creates a new Logger with the specified writer
func New(writer io.Writer) *Logger {
	return &Logger{
		writer: writer,
		level:  InfoLevel,
	}
}

// NewDefault creates a new Logger with stdout as output
func NewDefault() *Logger {
	return New(os.Stdout)
}

// SetLevel sets the minimum logging level
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Debug logs a debug-level message
func (l *Logger) Debug(msg string, args ...interface{}) {
	// Debug not implemented in minimal version
}

// Info logs an info-level message
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, msg, args...)
	}
}

// Warn logs a warning-level message
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, msg, args...)
	}
}

// Error logs an error-level message
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, msg, args...)
	}
}

// log formats and writes a log entry
func (l *Logger) log(level Level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(msg, args...)
	fmt.Fprintf(l.writer, "%s [%s] %s\n", timestamp, level.String(), message)
}

// Global logger instance
var defaultLogger = NewDefault()

// Package-level logging functions

// SetDefaultLevel sets the level for the default logger
func SetDefaultLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// Info logs an info message using the default logger
func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

// Warn logs a warning message using the default logger
func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

// Error logs an error message using the default logger
func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}