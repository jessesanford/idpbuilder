// Package printer provides output formatting utilities for the image builder
package printer

import (
	"encoding/json"
	"fmt"
	"io"
)

// Format represents different output formats
type Format string

const (
	JSONFormat  Format = "json"
	PlainFormat Format = "plain"
)

// Printer interface defines methods for printing data
type Printer interface {
	Print(data interface{}) error
	SetFormat(format Format)
}

// DefaultPrinter is the default implementation of the Printer interface
type DefaultPrinter struct {
	format Format
	writer io.Writer
}

// NewPrinter creates a new printer with the specified format and writer
func NewPrinter(format Format, writer io.Writer) *DefaultPrinter {
	return &DefaultPrinter{
		format: format,
		writer: writer,
	}
}

// SetFormat sets the output format
func (p *DefaultPrinter) SetFormat(format Format) {
	p.format = format
}

// Print prints data using the configured format
func (p *DefaultPrinter) Print(data interface{}) error {
	switch p.format {
	case JSONFormat:
		encoder := json.NewEncoder(p.writer)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	case PlainFormat:
		_, err := fmt.Fprintf(p.writer, "%v\n", data)
		return err
	default:
		return fmt.Errorf("unsupported format: %s", p.format)
	}
}

// ProgressPrinter provides simple progress indication
type ProgressPrinter struct {
	writer  io.Writer
	total   int
	current int
}

// NewProgressPrinter creates a new progress printer
func NewProgressPrinter(writer io.Writer, total int) *ProgressPrinter {
	return &ProgressPrinter{
		writer: writer,
		total:  total,
	}
}

// Update updates the progress
func (p *ProgressPrinter) Update(current int) {
	p.current = current
	percentage := float64(current) / float64(p.total) * 100
	fmt.Fprintf(p.writer, "Progress: %d/%d (%.1f%%)\n", current, p.total, percentage)
}

// Finish completes the progress
func (p *ProgressPrinter) Finish() {
	fmt.Fprintf(p.writer, "Complete: %d/%d\n", p.total, p.total)
}