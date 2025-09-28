package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type consoleReporter struct {
	writer    io.Writer
	operation *Operation
	verbose   bool
	mutex     sync.Mutex
}

// NewConsoleReporter creates a progress reporter for terminal output
func NewConsoleReporter() ProgressReporter {
	return &consoleReporter{
		writer: os.Stdout,
	}
}

// NewConsoleReporterWithWriter creates a console reporter with custom writer for testing
func NewConsoleReporterWithWriter(writer io.Writer) ProgressReporter {
	return &consoleReporter{
		writer: writer,
	}
}

func (c *consoleReporter) Start(message string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.operation = &Operation{
		Name:      message,
		StartTime: time.Now(),
	}
	fmt.Fprintf(c.writer, "Starting: %s\n", message)
}

func (c *consoleReporter) ReportProgress(current, total int64, message string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.operation != nil {
		c.operation.Current = current
		c.operation.Total = total
	}

	if total > 0 {
		percentage := float64(current) / float64(total) * 100
		fmt.Fprintf(c.writer, "[%.1f%%] %s\n", percentage, message)
	} else {
		fmt.Fprintf(c.writer, "[%d] %s\n", current, message)
	}
}

func (c *consoleReporter) Complete(message string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.operation != nil {
		duration := time.Since(c.operation.StartTime)
		fmt.Fprintf(c.writer, "Completed: %s (took %s)\n", message, formatDuration(duration))
	} else {
		fmt.Fprintf(c.writer, "Completed: %s\n", message)
	}
}

func (c *consoleReporter) Error(err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if err != nil {
		fmt.Fprintf(c.writer, "Error: %v\n", err)
	}
}