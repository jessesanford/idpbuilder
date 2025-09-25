package kind

import (
	"fmt"
	"io"
	"os"

	"sigs.k8s.io/kind/pkg/log"
)

// KindLogger implements the Kind cluster.Logger interface
type KindLogger struct {
	writer io.Writer
	level  log.Level
}

// NewKindLogger creates a new KindLogger that writes to the specified writer
func NewKindLogger(w io.Writer) *KindLogger {
	if w == nil {
		w = os.Stdout
	}

	return &KindLogger{
		writer: w,
		level:  log.Level(0), // Default to info level
	}
}

// Write implements io.Writer interface
func (l *KindLogger) Write(p []byte) (int, error) {
	if l.writer == nil {
		return len(p), nil // Discard if no writer
	}
	return l.writer.Write(p)
}

// V implements log.Logger interface for verbose logging
func (l *KindLogger) V(level log.Level) log.InfoLogger {
	// Create a copy with the specified verbosity level
	newLogger := &KindLogger{
		writer: l.writer,
		level:  level,
	}
	return newLogger
}

// Info implements cluster.Logger interface
func (l *KindLogger) Info(message string) {
	if l.shouldLog(log.Level(0)) {
		l.Write([]byte("[INFO] " + message + "\n"))
	}
}

// Infof implements cluster.Logger interface
func (l *KindLogger) Infof(format string, args ...interface{}) {
	if l.shouldLog(log.Level(0)) {
		l.Write([]byte(fmt.Sprintf("[INFO] "+format+"\n", args...)))
	}
}

// Warn implements cluster.Logger interface
func (l *KindLogger) Warn(message string) {
	if l.shouldLog(log.Level(1)) {
		l.Write([]byte("[WARN] " + message + "\n"))
	}
}

// Warnf implements cluster.Logger interface
func (l *KindLogger) Warnf(format string, args ...interface{}) {
	if l.shouldLog(log.Level(1)) {
		l.Write([]byte(fmt.Sprintf("[WARN] "+format+"\n", args...)))
	}
}

// Error implements cluster.Logger interface
func (l *KindLogger) Error(message string) {
	l.Write([]byte("[ERROR] " + message + "\n"))
}

// Errorf implements cluster.Logger interface
func (l *KindLogger) Errorf(format string, args ...interface{}) {
	l.Write([]byte(fmt.Sprintf("[ERROR] "+format+"\n", args...)))
}

// Enabled implements log.InfoLogger interface
func (l *KindLogger) Enabled() bool {
	return l.writer != nil
}

// shouldLog determines if a message at the given level should be logged
func (l *KindLogger) shouldLog(messageLevel log.Level) bool {
	return messageLevel <= l.level
}