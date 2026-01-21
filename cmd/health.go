package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "System health overview",
	Long:  "Display a quick overview of overall system health including array, disks, Docker, VMs, and parity status.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Println("Checking system health...\n")

		var hasErrors bool
		var hasWarnings bool

		// Check Array Status
		fmt.Println("=== Array Status ===")
		arrayInfo, err := apiClient.GetArrayInfo(ctx)
		if err != nil {
			fmt.Printf("%s\n", output.Error("Failed to get array status: "+err.Error()))
			hasErrors = true
		} else {
			if strings.ToUpper(arrayInfo.State) == "STARTED" {
				fmt.Printf("State: %s\n", output.Success("STARTED"))
			} else {
				fmt.Printf("State: %s\n", output.Warning(arrayInfo.State))
				hasWarnings = true
			}

			// Check disk health
			allDisks := arrayInfo.AllDisks()
			healthyDisks := 0
			totalDisks := 0

			for _, disk := range allDisks {
				totalDisks++
				status := strings.ToUpper(disk.Status)
				if status == "DISK_OK" {
					healthyDisks++
				} else if status != "" {
					hasWarnings = true
				}

				// Check temperature
				if disk.Temperature > 0 && disk.Temperature >= 60 {
					hasWarnings = true
				}
			}

			if healthyDisks == totalDisks {
				fmt.Printf("Disks: %s (%d/%d healthy)\n", output.Success(fmt.Sprintf("%d/%d healthy", healthyDisks, totalDisks)), healthyDisks, totalDisks)
			} else {
				fmt.Printf("Disks: %s (%d/%d healthy)\n", output.Warning(fmt.Sprintf("%d/%d healthy", healthyDisks, totalDisks)), healthyDisks, totalDisks)
				hasWarnings = true
			}
		}

		fmt.Println()

		// Check Parity Status
		fmt.Println("=== Parity Status ===")
		parityStatus, err := apiClient.GetParityCheckStatus(ctx)
		if err != nil {
			fmt.Printf("%s\n", output.Error("Failed to get parity status: "+err.Error()))
			hasErrors = true
		} else {
			if parityStatus.Running {
				fmt.Printf("Check Running: %s (Progress: %d%%)\n",
					output.Info("YES"),
					parityStatus.Progress)
			} else {
				fmt.Printf("Check Running: %s\n", output.Success("NO"))
			}

			if parityStatus.Errors > 0 {
				fmt.Printf("Last Check Errors: %s\n", output.Error(fmt.Sprintf("%d errors found", parityStatus.Errors)))
				hasErrors = true
			} else if parityStatus.Date != "" {
				fmt.Printf("Last Check: %s\n", output.Success("No errors"))
			}
		}

		fmt.Println()

		// Check Docker Containers
		fmt.Println("=== Docker Containers ===")
		containers, err := apiClient.GetContainers(ctx)
		if err != nil {
			fmt.Printf("%s\n", output.Error("Failed to get containers: "+err.Error()))
			hasErrors = true
		} else {
			running := 0
			stopped := 0

			for _, container := range containers {
				if strings.ToLower(container.State) == "running" {
					running++
				} else {
					stopped++
				}
			}

			total := running + stopped
			if stopped == 0 {
				fmt.Printf("Status: %s (%d running, %d stopped)\n",
					output.Success(fmt.Sprintf("%d/%d running", running, total)),
					running, stopped)
			} else {
				fmt.Printf("Status: %s (%d running, %d stopped)\n",
					output.Info(fmt.Sprintf("%d/%d running", running, total)),
					running, stopped)
			}
		}

		fmt.Println()

		// Check System Resources
		fmt.Println("=== System Resources ===")
		metrics, err := apiClient.GetMetrics(ctx)
		if err != nil {
			fmt.Printf("%s\n", output.Error("Failed to get metrics: "+err.Error()))
			hasErrors = true
		} else {
			// CPU
			if metrics.CPU.PercentTotal >= 90 {
				fmt.Printf("CPU: %s\n", output.Warning(fmt.Sprintf("%.1f%% (high)", metrics.CPU.PercentTotal)))
				hasWarnings = true
			} else {
				fmt.Printf("CPU: %s\n", output.Success(fmt.Sprintf("%.1f%%", metrics.CPU.PercentTotal)))
			}

			// Memory
			if metrics.Memory.PercentTotal >= 90 {
				fmt.Printf("Memory: %s\n", output.Warning(fmt.Sprintf("%.1f%% (high)", metrics.Memory.PercentTotal)))
				hasWarnings = true
			} else {
				fmt.Printf("Memory: %s\n", output.Success(fmt.Sprintf("%.1f%%", metrics.Memory.PercentTotal)))
			}
		}

		fmt.Println()

		// Check Notifications
		fmt.Println("=== Notifications ===")
		notifOverview, err := apiClient.GetNotificationOverview(ctx)
		if err != nil {
			fmt.Printf("%s\n", output.Error("Failed to get notifications: "+err.Error()))
			hasErrors = true
		} else {
			if notifOverview.Unread.Alert > 0 {
				fmt.Printf("Alerts: %s\n", output.Error(fmt.Sprintf("%d unread", notifOverview.Unread.Alert)))
				hasErrors = true
			} else {
				fmt.Printf("Alerts: %s\n", output.Success("None"))
			}

			if notifOverview.Unread.Warning > 0 {
				fmt.Printf("Warnings: %s\n", output.Warning(fmt.Sprintf("%d unread", notifOverview.Unread.Warning)))
				hasWarnings = true
			} else {
				fmt.Printf("Warnings: %s\n", output.Success("None"))
			}

			if notifOverview.Unread.Info > 0 {
				fmt.Printf("Info: %s unread\n", output.Cyan(fmt.Sprintf("%d", notifOverview.Unread.Info)))
			}
		}

		fmt.Println()
		fmt.Println("===================")

		// Overall summary
		if hasErrors {
			fmt.Printf("\n%s\n", output.Error("System has errors that need attention"))
			return nil
		} else if hasWarnings {
			fmt.Printf("\n%s\n", output.Warning("System has warnings"))
			return nil
		} else {
			fmt.Printf("\n%s\n", output.Success("All systems healthy"))
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
}
