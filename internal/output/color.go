package output

import (
	"fmt"
	"os"
	"strings"
)

// ANSI color codes
const (
	ColorReset   = "\033[0m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"
	ColorGray    = "\033[90m"

	ColorBoldRed    = "\033[1;31m"
	ColorBoldGreen  = "\033[1;32m"
	ColorBoldYellow = "\033[1;33m"
	ColorBoldBlue   = "\033[1;34m"
)

var colorsEnabled = true

// DisableColors disables colored output
func DisableColors() {
	colorsEnabled = false
}

// EnableColors enables colored output
func EnableColors() {
	colorsEnabled = true
}

// ColorsEnabled returns whether colors are enabled
func ColorsEnabled() bool {
	return colorsEnabled
}

// init checks if we should disable colors based on environment
func init() {
	// Disable colors if NO_COLOR is set or output is not a terminal
	if os.Getenv("NO_COLOR") != "" {
		colorsEnabled = false
	}

	// Check if stdout is a terminal
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		colorsEnabled = false
	}
}

// Colorize wraps text with color codes if colors are enabled
func Colorize(text, color string) string {
	if !colorsEnabled || color == "" {
		return text
	}
	return color + text + ColorReset
}

// Red returns red colored text
func Red(text string) string {
	return Colorize(text, ColorRed)
}

// Green returns green colored text
func Green(text string) string {
	return Colorize(text, ColorGreen)
}

// Yellow returns yellow colored text
func Yellow(text string) string {
	return Colorize(text, ColorYellow)
}

// Blue returns blue colored text
func Blue(text string) string {
	return Colorize(text, ColorBlue)
}

// Cyan returns cyan colored text
func Cyan(text string) string {
	return Colorize(text, ColorCyan)
}

// Gray returns gray colored text
func Gray(text string) string {
	return Colorize(text, ColorGray)
}

// BoldRed returns bold red colored text
func BoldRed(text string) string {
	return Colorize(text, ColorBoldRed)
}

// BoldGreen returns bold green colored text
func BoldGreen(text string) string {
	return Colorize(text, ColorBoldGreen)
}

// BoldYellow returns bold yellow colored text
func BoldYellow(text string) string {
	return Colorize(text, ColorBoldYellow)
}

// Success returns a green checkmark with text
func Success(text string) string {
	return BoldGreen("✓") + " " + text
}

// Error returns a red X with text
func Error(text string) string {
	return BoldRed("✗") + " " + text
}

// Warning returns a yellow warning symbol with text
func Warning(text string) string {
	return BoldYellow("⚠") + " " + text
}

// Info returns a blue info symbol with text
func Info(text string) string {
	return Colorize("ℹ", ColorBoldBlue) + " " + text
}

// ColorizeState returns colored state text based on the state
func ColorizeState(state string) string {
	stateUpper := strings.ToUpper(state)

	switch stateUpper {
	case "STARTED", "RUNNING", "UP", "ONLINE":
		return Green(stateUpper)
	case "STOPPED", "DOWN", "OFFLINE", "EXITED":
		return Red(stateUpper)
	case "PAUSED", "STOPPING", "STARTING":
		return Yellow(stateUpper)
	case "DISK_OK", "GOOD", "HEALTHY":
		return Green(stateUpper)
	case "DISK_DSBL", "DISK_INVALID", "DISK_WRONG", "ERROR", "FAILED":
		return Red(stateUpper)
	case "WARNING", "DISK_NP":
		return Yellow(stateUpper)
	default:
		return stateUpper
	}
}

// ColorizePercentage returns colored percentage based on value
// High percentages are red, medium are yellow, low are green
func ColorizePercentage(percent float64, reverse bool) string {
	text := fmt.Sprintf("%.1f%%", percent)

	// For things like disk usage, high is bad
	// For things like CPU idle, high is good (reverse=true)

	if !reverse {
		// High is bad (disk usage, CPU usage, memory usage)
		if percent >= 90 {
			return Red(text)
		} else if percent >= 75 {
			return Yellow(text)
		}
		return Green(text)
	} else {
		// High is good (battery level, free space)
		if percent >= 75 {
			return Green(text)
		} else if percent >= 50 {
			return Yellow(text)
		}
		return Red(text)
	}
}

// ColorizeTemperature returns colored temperature based on value
func ColorizeTemperature(temp float64) string {
	text := fmt.Sprintf("%.1f°C", temp)

	if temp >= 60 {
		return Red(text)
	} else if temp >= 50 {
		return Yellow(text)
	} else if temp >= 40 {
		return Cyan(text)
	}
	return Blue(text)
}
