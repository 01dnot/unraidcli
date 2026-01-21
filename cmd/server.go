package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server information and status",
	Long:  "View Unraid server information, status, and health.",
}

// serverInfoCmd represents the server info command
var serverInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show server information",
	Long:  "Display detailed system information including CPU, memory, platform, and version.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		info, err := apiClient.GetSystemInfo(ctx)
		if err != nil {
			return fmt.Errorf("failed to get system info: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			// Calculate total memory from layout
			var totalMem int64
			for _, module := range info.Memory.Layout {
				totalMem += module.Size
			}

			cpuBrand := info.CPU.Brand
			if cpuBrand == "" {
				cpuBrand = info.CPU.Manufacturer
			}

			data := map[string]interface{}{
				"Hostname": info.OS.Hostname,
				"Platform": info.OS.Platform,
				"Version":  info.Versions.Core.Unraid,
				"Uptime":   info.OS.Uptime,
				"CPU":      fmt.Sprintf("%s (%d cores, %d threads)", cpuBrand, info.CPU.Cores, info.CPU.Threads),
				"CPU Speed": fmt.Sprintf("%.2f GHz", info.CPU.Speed),
				"Total Memory": output.FormatBytes(totalMem),
			}
			formatter.PrintKeyValue(data)
		} else {
			formatter.Print(info)
		}

		return nil
	},
}

// serverStatusCmd represents the server status command
var serverStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show server status",
	Long:  "Display overall server health status and uptime.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		info, err := apiClient.GetSystemInfo(ctx)
		if err != nil {
			return fmt.Errorf("failed to get system status: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("Server: %s\n", info.OS.Hostname)
			fmt.Printf("Status: ✓ Online\n")
			fmt.Printf("Uptime: %s\n", info.OS.Uptime)
			fmt.Printf("Version: %s\n", info.Versions.Core.Unraid)
			fmt.Printf("Platform: %s\n", info.OS.Platform)
		} else {
			status := map[string]interface{}{
				"hostname": info.OS.Hostname,
				"status":   "online",
				"uptime":   info.OS.Uptime,
				"version":  info.Versions.Core.Unraid,
				"platform": info.OS.Platform,
			}
			formatter.Print(status)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverInfoCmd)
	serverCmd.AddCommand(serverStatusCmd)
}
