package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

var (
	metricsWatch    bool
	metricsInterval int
	showCores       bool
)

// metricsCmd represents the metrics command
var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show system metrics",
	Long:  "Display real-time system metrics including CPU and memory usage.",
	RunE: func(cmd *cobra.Command, args []string) error {
		metricsFunc := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			metrics, err := apiClient.GetMetrics(ctx)
			if err != nil {
				return fmt.Errorf("failed to get metrics: %w", err)
			}

			if outputFormat == "" || outputFormat == "table" {
				// Show timestamp in watch mode
				if metricsWatch {
					fmt.Printf("Last updated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
				}

				// CPU Usage
				fmt.Printf("CPU Usage: %s\n", output.ColorizePercentage(metrics.CPU.PercentTotal, false))

				// Show per-core usage if requested
				if showCores && len(metrics.CPU.CPUs) > 0 {
					fmt.Printf("\nPer-Core Usage:\n")
					headers := []string{"Core", "Total", "User", "System", "Idle"}
					var rows [][]string

					for i, cpu := range metrics.CPU.CPUs {
						rows = append(rows, []string{
							fmt.Sprintf("Core %d", i),
							fmt.Sprintf("%.1f%%", cpu.PercentTotal),
							fmt.Sprintf("%.1f%%", cpu.PercentUser),
							fmt.Sprintf("%.1f%%", cpu.PercentSystem),
							fmt.Sprintf("%.1f%%", cpu.PercentIdle),
						})
					}
					formatter.PrintTable(headers, rows)
					fmt.Println()
				}

				// Memory Usage
				fmt.Printf("Memory Usage: %s (%s / %s)\n",
					output.ColorizePercentage(metrics.Memory.PercentTotal, false),
					output.FormatBytes(metrics.Memory.Used),
					output.FormatBytes(metrics.Memory.Total))
				fmt.Printf("  Used: %s\n", output.FormatBytes(metrics.Memory.Used))
				fmt.Printf("  Available: %s\n", output.FormatBytes(metrics.Memory.Available))
				fmt.Printf("  Free: %s\n", output.FormatBytes(metrics.Memory.Free))

				// Swap Usage
				if metrics.Memory.SwapTotal > 0 {
					fmt.Printf("\nSwap Usage: %.1f%% (%s / %s)\n",
						metrics.Memory.PercentSwapTotal,
						output.FormatBytes(metrics.Memory.SwapUsed),
						output.FormatBytes(metrics.Memory.SwapTotal))
				}
			} else {
				formatter.Print(metrics)
			}

			return nil
		}

		// Watch mode
		if metricsWatch {
			// Setup signal handling for graceful exit
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-sigChan
				cancel()
			}()

			fmt.Println("Press Ctrl+C to exit watch mode\n")
			interval := time.Duration(metricsInterval) * time.Second
			return output.Watch(ctx, interval, metricsFunc)
		}

		// Normal mode
		return metricsFunc()
	},
}

func init() {
	rootCmd.AddCommand(metricsCmd)

	metricsCmd.Flags().BoolVarP(&metricsWatch, "watch", "w", false, "Watch mode - auto-refresh every N seconds")
	metricsCmd.Flags().IntVarP(&metricsInterval, "interval", "i", 2, "Refresh interval in seconds for watch mode")
	metricsCmd.Flags().BoolVarP(&showCores, "cores", "c", false, "Show per-core CPU usage")
}
