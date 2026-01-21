package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var correctingErrors bool

// parityCmd represents the parity command
var parityCmd = &cobra.Command{
	Use:   "parity",
	Short: "Manage parity checks",
	Long:  "Start, stop, pause, and monitor parity check operations.",
}

// parityStatusCmd represents the parity status command
var parityStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show parity check status",
	Long:  "Display current parity check status and progress.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		status, err := apiClient.GetParityCheckStatus(ctx)
		if err != nil {
			return fmt.Errorf("failed to get parity check status: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("Status: %s\n", status.Status)

			if status.Running {
				fmt.Printf("Running: ✓ (Progress: %d%%)\n", status.Progress)
			} else {
				fmt.Printf("Running: ✗\n")
			}

			if status.Paused {
				fmt.Printf("Paused: ✓\n")
			}

			if status.Correcting {
				fmt.Printf("Correcting: ✓\n")
			}

			if status.Date != "" {
				fmt.Printf("Last Check: %s\n", status.Date)
			}

			if status.Duration > 0 {
				duration := time.Duration(status.Duration) * time.Second
				fmt.Printf("Duration: %s\n", duration.String())
			}

			if status.Speed != "" {
				fmt.Printf("Speed: %s\n", status.Speed)
			}

			if status.Errors > 0 {
				fmt.Printf("Errors: %d\n", status.Errors)
			} else if status.Date != "" {
				fmt.Printf("Errors: 0\n")
			}
		} else {
			formatter.Print(status)
		}

		return nil
	},
}

// parityHistoryCmd represents the parity history command
var parityHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show parity check history",
	Long:  "Display history of previous parity checks.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		history, err := apiClient.GetParityHistory(ctx)
		if err != nil {
			return fmt.Errorf("failed to get parity history: %w", err)
		}

		if len(history) == 0 {
			fmt.Println("No parity check history found.")
			return nil
		}

		if outputFormat == "" || outputFormat == "table" {
			headers := []string{"Date", "Status", "Duration", "Speed", "Errors"}
			var rows [][]string

			for _, check := range history {
				duration := ""
				if check.Duration > 0 {
					d := time.Duration(check.Duration) * time.Second
					duration = d.String()
				}

				rows = append(rows, []string{
					check.Date,
					check.Status,
					duration,
					check.Speed,
					fmt.Sprintf("%d", check.Errors),
				})
			}

			formatter.PrintTable(headers, rows)
		} else {
			formatter.Print(history)
		}

		return nil
	},
}

// parityStartCmd represents the parity start command
var parityStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a parity check",
	Long:  "Start a parity check operation. Use --correct to enable error correction.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if correctingErrors {
			fmt.Println("Starting parity check with error correction...")
		} else {
			fmt.Println("Starting parity check (read-only)...")
		}

		if err := apiClient.StartParityCheck(ctx, correctingErrors); err != nil {
			return fmt.Errorf("failed to start parity check: %w", err)
		}

		fmt.Println("✓ Parity check started successfully")
		return nil
	},
}

// parityPauseCmd represents the parity pause command
var parityPauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause a running parity check",
	Long:  "Pause the currently running parity check.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Println("Pausing parity check...")

		if err := apiClient.PauseParityCheck(ctx); err != nil {
			return fmt.Errorf("failed to pause parity check: %w", err)
		}

		fmt.Println("✓ Parity check paused successfully")
		return nil
	},
}

// parityResumeCmd represents the parity resume command
var parityResumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused parity check",
	Long:  "Resume a previously paused parity check.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Println("Resuming parity check...")

		if err := apiClient.ResumeParityCheck(ctx); err != nil {
			return fmt.Errorf("failed to resume parity check: %w", err)
		}

		fmt.Println("✓ Parity check resumed successfully")
		return nil
	},
}

// parityCancelCmd represents the parity cancel command
var parityCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel a running parity check",
	Long:  "Cancel the currently running parity check.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Println("Canceling parity check...")

		if err := apiClient.CancelParityCheck(ctx); err != nil {
			return fmt.Errorf("failed to cancel parity check: %w", err)
		}

		fmt.Println("✓ Parity check canceled successfully")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(parityCmd)
	parityCmd.AddCommand(parityStatusCmd)
	parityCmd.AddCommand(parityHistoryCmd)
	parityCmd.AddCommand(parityStartCmd)
	parityCmd.AddCommand(parityPauseCmd)
	parityCmd.AddCommand(parityResumeCmd)
	parityCmd.AddCommand(parityCancelCmd)

	// Add flags
	parityStartCmd.Flags().BoolVarP(&correctingErrors, "correct", "c", false, "Enable error correction (write mode)")
}
