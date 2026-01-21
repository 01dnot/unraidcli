package cmd

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

// arrayCmd represents the array command
var arrayCmd = &cobra.Command{
	Use:   "array",
	Short: "Manage the Unraid array",
	Long:  "View status and manage the Unraid storage array.",
}

// arrayStatusCmd represents the array status command
var arrayStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show array status",
	Long:  "Display array state, capacity, and disk information.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		arrayInfo, err := apiClient.GetArrayInfo(ctx)
		if err != nil {
			return fmt.Errorf("failed to get array info: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			// Parse and convert kilobytes to bytes for formatting
			totalKB, _ := strconv.ParseInt(arrayInfo.Capacity.Kilobytes.Total, 10, 64)
			usedKB, _ := strconv.ParseInt(arrayInfo.Capacity.Kilobytes.Used, 10, 64)
			freeKB, _ := strconv.ParseInt(arrayInfo.Capacity.Kilobytes.Free, 10, 64)

			totalBytes := totalKB * 1024
			usedBytes := usedKB * 1024
			freeBytes := freeKB * 1024

			usedPercent := float64(0)
			if totalBytes > 0 {
				usedPercent = float64(usedBytes) / float64(totalBytes) * 100
			}

			// Print array summary
			fmt.Printf("Array State: %s\n", output.FormatState(arrayInfo.State))
			fmt.Printf("Total Capacity: %s\n", output.FormatBytes(totalBytes))
			fmt.Printf("Used: %s (%.1f%%)\n", output.FormatBytes(usedBytes), usedPercent)
			fmt.Printf("Free: %s\n\n", output.FormatBytes(freeBytes))

			// Get all disks (boot, parity, data, cache)
			allDisks := arrayInfo.AllDisks()

			// Print disk table
			if len(allDisks) > 0 {
				headers := []string{"Name", "Device", "Type", "Status", "Size", "Temp", "FS Type"}
				var rows [][]string

				for _, disk := range allDisks {
					temp := float64(disk.Temperature)
					tempStr := "N/A"
					if temp > 0 {
						tempStr = output.ColorizeTemperature(temp)
					}

					rows = append(rows, []string{
						disk.Name,
						disk.Device,
						disk.Type,
						output.ColorizeState(disk.Status),
						output.FormatBytes(disk.Size),
						tempStr,
						disk.FsType,
					})
				}

				formatter.PrintTable(headers, rows)
			}
		} else {
			formatter.Print(arrayInfo)
		}

		return nil
	},
}

// arrayStartCmd represents the array start command
var arrayStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the array",
	Long:  "Start the Unraid storage array.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Println("Starting array...")
		if err := apiClient.StartArray(ctx); err != nil {
			return fmt.Errorf("failed to start array: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Println("✓ Array started successfully")
		} else {
			formatter.Print(map[string]string{
				"status":  "success",
				"message": "Array started successfully",
			})
		}

		return nil
	},
}

// arrayStopCmd represents the array stop command
var arrayStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the array",
	Long:  "Stop the Unraid storage array.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		fmt.Println("Stopping array...")
		if err := apiClient.StopArray(ctx); err != nil {
			return fmt.Errorf("failed to stop array: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Println("✓ Array stopped successfully")
		} else {
			formatter.Print(map[string]string{
				"status":  "success",
				"message": "Array stopped successfully",
			})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(arrayCmd)
	arrayCmd.AddCommand(arrayStatusCmd)
	arrayCmd.AddCommand(arrayStartCmd)
	arrayCmd.AddCommand(arrayStopCmd)
}
