package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	// ColoredOutput controls whether output should be colored
	ColoredOutput bool
	// ColoredOutputMsg provides help text for colored output
	ColoredOutputMsg = "Enable colored output"
)

// OutputFormat represents different output formats
type OutputFormat string

const (
	// TableOutput formats output as a table
	TableOutput OutputFormat = "table"
	// JSONOutput formats output as JSON
	JSONOutput OutputFormat = "json"
	// YAMLOutput formats output as YAML
	YAMLOutput OutputFormat = "yaml"
	// WideOutput formats output as wide table
	WideOutput OutputFormat = "wide"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

// Printer provides structured output functionality
type Printer struct {
	format     OutputFormat
	colored    bool
	writer     *tabwriter.Writer
}

// NewPrinter creates a new printer with the specified format
func NewPrinter(format OutputFormat) *Printer {
	return &Printer{
		format:  format,
		colored: ColoredOutput,
		writer:  tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0),
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	switch p.format {
	case JSONOutput:
		return p.printJSON(data)
	case YAMLOutput:
		return p.printYAML(data)
	case TableOutput, WideOutput:
		return p.printTable(data)
	default:
		return fmt.Errorf("unsupported output format: %s", p.format)
	}
}

// printJSON outputs data as JSON
func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printYAML outputs data as YAML
func (p *Printer) printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

// printTable outputs data as a formatted table
func (p *Printer) printTable(data interface{}) error {
	defer p.writer.Flush()
	
	switch v := data.(type) {
	case []map[string]interface{}:
		return p.printMapSlice(v)
	case map[string]interface{}:
		return p.printSingleMap(v)
	default:
		return fmt.Errorf("unsupported data type for table output: %T", data)
	}
}

// printMapSlice prints a slice of maps as a table
func (p *Printer) printMapSlice(data []map[string]interface{}) error {
	if len(data) == 0 {
		p.PrintInfo("No resources found")
		return nil
	}

	// Extract headers from first item
	var headers []string
	for key := range data[0] {
		headers = append(headers, strings.ToUpper(key))
	}

	// Print headers
	headerRow := strings.Join(headers, "\t")
	if p.colored {
		headerRow = p.colorize(headerRow, ColorBold)
	}
	fmt.Fprintln(p.writer, headerRow)

	// Print data rows
	for _, item := range data {
		var values []string
		for _, header := range headers {
			key := strings.ToLower(header)
			if val, exists := item[key]; exists {
				values = append(values, fmt.Sprintf("%v", val))
			} else {
				values = append(values, "")
			}
		}
		fmt.Fprintln(p.writer, strings.Join(values, "\t"))
	}

	return nil
}

// printSingleMap prints a single map as key-value pairs
func (p *Printer) printSingleMap(data map[string]interface{}) error {
	for key, value := range data {
		keyStr := strings.ToUpper(key)
		if p.colored {
			keyStr = p.colorize(keyStr, ColorCyan)
		}
		fmt.Fprintf(p.writer, "%s:\t%v\n", keyStr, value)
	}
	return nil
}

// PrintSuccess prints a success message
func PrintSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if ColoredOutput {
		fmt.Printf("%s✓ %s%s\n", ColorGreen, msg, ColorReset)
	} else {
		fmt.Printf("✓ %s\n", msg)
	}
}

// PrintError prints an error message
func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if ColoredOutput {
		fmt.Fprintf(os.Stderr, "%s✗ %s%s\n", ColorRed, msg, ColorReset)
	} else {
		fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
	}
}

// PrintWarning prints a warning message
func PrintWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if ColoredOutput {
		fmt.Printf("%s⚠ %s%s\n", ColorYellow, msg, ColorReset)
	} else {
		fmt.Printf("⚠ %s\n", msg)
	}
}

// PrintInfo prints an info message
func (p *Printer) PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if p.colored {
		fmt.Printf("%sℹ %s%s\n", ColorBlue, msg, ColorReset)
	} else {
		fmt.Printf("ℹ %s\n", msg)
	}
}

// PrintStep prints a step message with timestamp
func PrintStep(step, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05")
	
	if ColoredOutput {
		fmt.Printf("%s[%s]%s %s%s%s %s\n", 
			ColorCyan, timestamp, ColorReset,
			ColorBold, step, ColorReset,
			msg)
	} else {
		fmt.Printf("[%s] %s %s\n", timestamp, step, msg)
	}
}

// colorize applies color codes to text if coloring is enabled
func (p *Printer) colorize(text, color string) string {
	if p.colored {
		return color + text + ColorReset
	}
	return text
}

// ValidateOutputFormat validates the output format string
func ValidateOutputFormat(format string) (OutputFormat, error) {
	switch strings.ToLower(format) {
	case "table", "":
		return TableOutput, nil
	case "json":
		return JSONOutput, nil
	case "yaml", "yml":
		return YAMLOutput, nil
	case "wide":
		return WideOutput, nil
	default:
		return "", fmt.Errorf("invalid output format: %s (supported: table, json, yaml, wide)", format)
	}
}