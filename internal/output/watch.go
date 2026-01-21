package output

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// WatchFunc is a function that gets executed repeatedly in watch mode
type WatchFunc func() error

// Watch executes a function repeatedly with a specified interval
func Watch(ctx context.Context, interval time.Duration, fn WatchFunc) error {
	// Clear screen initially
	clearScreen()

	// Run immediately first time
	if err := fn(); err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			clearScreen()
			if err := fn(); err != nil {
				return err
			}
		}
	}
}

// clearScreen clears the terminal screen
func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		// Unix-like systems (Linux, macOS)
		fmt.Print("\033[H\033[2J")
	}
}
