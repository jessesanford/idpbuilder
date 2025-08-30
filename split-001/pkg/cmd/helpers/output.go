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
	ColoredOutput    bool
	ColoredOutputMsg = "Enable colored output"
)

type OutputFormat string

const (
	TableOutput OutputFormat = "table"
	JSONOutput  OutputFormat = "json"
	YAMLOutput  OutputFormat = "yaml"
	WideOutput  OutputFormat = "wide"
)

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

type Printer struct {
	format  OutputFormat
	colored bool
	writer  *tabwriter.Writer
}

func NewPrinter(format OutputFormat) *Printer {
	return &Printer{
		format:  format,
		colored: ColoredOutput,
		writer:  tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0),
	}
}

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

func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (p *Printer) printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	encoder.SetIndent(2)
	return encoder.Encode(data)
}

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

func (p *Printer) printMapSlice(data []map[string]interface{}) error {
	if len(data) == 0 {
		fmt.Println("No resources found")
		return nil
	}

	var headers []string
	for key := range data[0] {
		headers = append(headers, strings.ToUpper(key))
	}

	fmt.Fprintln(p.writer, strings.Join(headers, "\t"))

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

func (p *Printer) printSingleMap(data map[string]interface{}) error {
	for key, value := range data {
		fmt.Fprintf(p.writer, "%s:\t%v\n", strings.ToUpper(key), value)
	}
	return nil
}

func PrintSuccess(format string, args ...interface{}) {
	fmt.Printf("✓ %s\n", fmt.Sprintf(format, args...))
}

func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "✗ %s\n", fmt.Sprintf(format, args...))
}

func PrintWarning(format string, args ...interface{}) {
	fmt.Printf("⚠ %s\n", fmt.Sprintf(format, args...))
}

func (p *Printer) PrintInfo(format string, args ...interface{}) {
	fmt.Printf("ℹ %s\n", fmt.Sprintf(format, args...))
}

func PrintStep(step, format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s %s\n", timestamp, step, fmt.Sprintf(format, args...))
}

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
		return "", fmt.Errorf("invalid output format: %s", format)
	}
}