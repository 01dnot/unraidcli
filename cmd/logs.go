package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

var (
	logLines int
	logTail  bool
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View system logs",
	Long:  "List and view system log files from your Unraid server.",
}

// logsListCmd represents the logs ls command
var logsListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List available log files",
	Long:    "Display all available log files on the system.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		logFiles, err := apiClient.GetLogFiles(ctx)
		if err != nil {
			return fmt.Errorf("failed to get log files: %w", err)
		}

		if len(logFiles) == 0 {
			fmt.Println("No log files found.")
			return nil
		}

		if outputFormat == "" || outputFormat == "table" {
			headers := []string{"Name", "Size", "Modified"}
			var rows [][]string

			for _, logFile := range logFiles {
				rows = append(rows, []string{
					logFile.Name,
					output.FormatBytes(int64(logFile.Size)),
					logFile.ModifiedAt,
				})
			}

			formatter.PrintTable(headers, rows)
		} else {
			formatter.Print(logFiles)
		}

		return nil
	},
}

// logsViewCmd represents the logs view command
var logsViewCmd = &cobra.Command{
	Use:   "view <log-file-path>",
	Short: "View a log file",
	Long:  "Display the content of a specific log file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		logPath := args[0]

		startLine := 0
		if logTail {
			// Get file info first to calculate start line
			logContent, err := apiClient.GetLogFile(ctx, logPath, 0, 0)
			if err != nil {
				return fmt.Errorf("failed to get log file: %w", err)
			}

			if logContent.TotalLines > logLines {
				startLine = logContent.TotalLines - logLines
			}
		}

		logContent, err := apiClient.GetLogFile(ctx, logPath, logLines, startLine)
		if err != nil {
			return fmt.Errorf("failed to get log file: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("Log file: %s\n", logContent.Path)
			fmt.Printf("Total lines: %d\n", logContent.TotalLines)
			if logTail {
				fmt.Printf("Showing last %d lines:\n\n", logLines)
			} else if logLines > 0 {
				fmt.Printf("Showing first %d lines:\n\n", logLines)
			} else {
				fmt.Printf("\n")
			}
			fmt.Println(logContent.Content)
		} else {
			formatter.Print(logContent)
		}

		return nil
	},
}

// logsTailCmd represents the logs tail command
var logsTailCmd = &cobra.Command{
	Use:   "tail <log-file-path>",
	Short: "Show the last lines of a log file",
	Long:  "Display the last N lines of a log file (similar to tail command).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		logPath := args[0]

		// Get file info first to calculate start line
		logContent, err := apiClient.GetLogFile(ctx, logPath, 0, 0)
		if err != nil {
			return fmt.Errorf("failed to get log file: %w", err)
		}

		startLine := 0
		if logContent.TotalLines > logLines {
			startLine = logContent.TotalLines - logLines
		}

		logContent, err = apiClient.GetLogFile(ctx, logPath, logLines, startLine)
		if err != nil {
			return fmt.Errorf("failed to get log file: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Println(logContent.Content)
		} else {
			formatter.Print(logContent)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.AddCommand(logsListCmd)
	logsCmd.AddCommand(logsViewCmd)
	logsCmd.AddCommand(logsTailCmd)

	// Add flags for logs view
	logsViewCmd.Flags().IntVarP(&logLines, "lines", "n", 100, "Number of lines to display (0 for all)")
	logsViewCmd.Flags().BoolVarP(&logTail, "tail", "t", false, "Show last N lines instead of first N lines")

	// Add flags for logs tail
	logsTailCmd.Flags().IntVarP(&logLines, "lines", "n", 50, "Number of lines to display")
}
