package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents the output format type
type Format string

const (
	// FormatTable represents table output format
	FormatTable Format = "table"
	// FormatJSON represents JSON output format
	FormatJSON Format = "json"
	// FormatYAML represents YAML output format
	FormatYAML Format = "yaml"
)

// Formatter handles output formatting
type Formatter struct {
	format Format
	writer io.Writer
}

// New creates a new formatter
func New(format string) *Formatter {
	f := Format(strings.ToLower(format))

	// Default to table if invalid format
	if f != FormatTable && f != FormatJSON && f != FormatYAML {
		f = FormatTable
	}

	return &Formatter{
		format: f,
		writer: os.Stdout,
	}
}

// Print outputs data in the configured format
func (f *Formatter) Print(data interface{}) error {
	switch f.format {
	case FormatJSON:
		return f.printJSON(data)
	case FormatYAML:
		return f.printYAML(data)
	default:
		// Table format is handled by specific methods
		return fmt.Errorf("table format requires using PrintTable method")
	}
}

// PrintTable outputs data as a table
func (f *Formatter) PrintTable(headers []string, rows [][]string) {
	if f.format != FormatTable {
		// If not table format, convert to map and print
		data := make([]map[string]string, len(rows))
		for i, row := range rows {
			rowMap := make(map[string]string)
			for j, header := range headers {
				if j < len(row) {
					rowMap[header] = row[j]
				}
			}
			data[i] = rowMap
		}
		f.Print(data)
		return
	}

	// Calculate column widths
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}

	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Print headers
	for i, header := range headers {
		if i > 0 {
			fmt.Fprint(f.writer, "  ")
		}
		fmt.Fprintf(f.writer, "%-*s", colWidths[i], header)
	}
	fmt.Fprintln(f.writer)

	// Print separator line
	for i := range headers {
		if i > 0 {
			fmt.Fprint(f.writer, "  ")
		}
		fmt.Fprint(f.writer, strings.Repeat("-", colWidths[i]))
	}
	fmt.Fprintln(f.writer)

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			if i > 0 {
				fmt.Fprint(f.writer, "  ")
			}
			if i < len(colWidths) {
				fmt.Fprintf(f.writer, "%-*s", colWidths[i], cell)
			} else {
				fmt.Fprint(f.writer, cell)
			}
		}
		fmt.Fprintln(f.writer)
	}
}

// PrintKeyValue outputs key-value pairs
func (f *Formatter) PrintKeyValue(data map[string]interface{}) error {
	if f.format == FormatTable {
		for key, value := range data {
			fmt.Fprintf(f.writer, "%s:\t%v\n", key, value)
		}
		return nil
	}

	return f.Print(data)
}

// printJSON outputs data as JSON
func (f *Formatter) printJSON(data interface{}) error {
	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// printYAML outputs data as YAML
func (f *Formatter) printYAML(data interface{}) error {
	encoder := yaml.NewEncoder(f.writer)
	defer encoder.Close()
	return encoder.Encode(data)
}

// PrintSuccess prints a success message
func (f *Formatter) PrintSuccess(message string) {
	if f.format == FormatTable {
		fmt.Fprintf(f.writer, "✓ %s\n", message)
	} else {
		f.Print(map[string]string{
			"status":  "success",
			"message": message,
		})
	}
}

// PrintError prints an error message
func (f *Formatter) PrintError(message string) {
	if f.format == FormatTable {
		fmt.Fprintf(f.writer, "✗ %s\n", message)
	} else {
		f.Print(map[string]string{
			"status":  "error",
			"message": message,
		})
	}
}

// FormatBytes formats bytes into human-readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatUptime formats uptime in seconds to human-readable format
func FormatUptime(seconds int64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	} else {
		return fmt.Sprintf("%dm", minutes)
	}
}

// FormatTemperature formats temperature value
func FormatTemperature(temp float64) string {
	if temp == 0 {
		return "N/A"
	}
	return fmt.Sprintf("%.1f°C", temp)
}

// FormatState formats state with color indicators (for table format)
func FormatState(state string) string {
	return ColorizeState(state)
}

// FormatBool formats a boolean value
func FormatBool(value bool) string {
	if value {
		return Green("Yes")
	}
	return Red("No")
}
