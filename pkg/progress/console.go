package progress

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ConsoleReporter implements ProgressReporter for console output.
// It provides human-readable progress updates suitable for command-line tools.
type ConsoleReporter struct {
	output io.Writer
	quiet  bool
	active map[string]*OperationProgress
}

// NewConsoleReporter creates a new console progress reporter.
// If output is nil, it defaults to os.Stdout.
func NewConsoleReporter(output io.Writer, quiet bool) *ConsoleReporter {
	if output == nil {
		output = os.Stdout
	}
	return &ConsoleReporter{
		output: output,
		quiet:  quiet,
		active: make(map[string]*OperationProgress),
	}
}

// StartOperation begins tracking a new operation with console output.
func (c *ConsoleReporter) StartOperation(ctx context.Context, opType OperationType, reference string) context.Context {
	if c.quiet {
		return ctx
	}

	progress := &OperationProgress{
		OperationType: opType,
		Reference:     reference,
		State:         StateStarted,
		StartTime:     time.Now(),
		LastUpdate:    time.Now(),
	}
	c.active[reference] = progress

	fmt.Fprintf(c.output, "Starting %s: %s\n", opType, reference)
	return ctx
}

// UpdateProgress reports current progress with a console progress bar.
func (c *ConsoleReporter) UpdateProgress(ctx context.Context, progress OperationProgress) {
	if c.quiet {
		return
	}

	c.active[progress.Reference] = &progress

	// Calculate percentage if total is known
	var percentage float64
	if progress.Total > 0 {
		percentage = float64(progress.Current) / float64(progress.Total) * 100
	}

	// Create simple progress bar
	var bar string
	if progress.Total > 0 {
		barWidth := 30
		filled := int(percentage * float64(barWidth) / 100)
		bar = "[" + strings.Repeat("=", filled) + strings.Repeat(" ", barWidth-filled) + "]"
	} else {
		bar = "[???]" // Indeterminate progress
	}

	if progress.Total > 0 {
		fmt.Fprintf(c.output, "\r%s %s %.1f%% (%d/%d) %s",
			progress.OperationType, bar, percentage, progress.Current, progress.Total, progress.Message)
	} else {
		fmt.Fprintf(c.output, "\r%s %s %s", progress.OperationType, bar, progress.Message)
	}
}

// CompleteOperation marks an operation as completed with final status.
func (c *ConsoleReporter) CompleteOperation(ctx context.Context, reference string, message string) {
	if c.quiet {
		return
	}

	if progress, exists := c.active[reference]; exists {
		duration := time.Since(progress.StartTime)
		fmt.Fprintf(c.output, "\nCompleted %s: %s (%v)\n", progress.OperationType, reference, duration)
		delete(c.active, reference)
	} else {
		fmt.Fprintf(c.output, "Completed: %s - %s\n", reference, message)
	}
}

// ErrorOperation reports that an operation failed.
func (c *ConsoleReporter) ErrorOperation(ctx context.Context, reference string, err error) {
	fmt.Fprintf(c.output, "\nError: %s - %v\n", reference, err)
	delete(c.active, reference)
}

// CancelOperation handles operation cancellation.
func (c *ConsoleReporter) CancelOperation(ctx context.Context, reference string, reason string) {
	if c.quiet {
		return
	}

	fmt.Fprintf(c.output, "\nCancelled: %s - %s\n", reference, reason)
	delete(c.active, reference)
}