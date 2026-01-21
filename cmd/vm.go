package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

// vmCmd represents the vm command
var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Manage virtual machines",
	Long:  "List, start, stop, and restart virtual machines on your Unraid server.",
}

// vmLsCmd represents the vm ls command
var vmLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all VMs",
	Long:  "List all virtual machines.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		vms, err := apiClient.GetVMs(ctx)
		if err != nil {
			return fmt.Errorf("failed to get VMs: %w", err)
		}

		if len(vms) == 0 {
			fmt.Println("No virtual machines found.")
			return nil
		}

		if outputFormat == "" || outputFormat == "table" {
			headers := []string{"Name", "State"}
			var rows [][]string

			for _, vm := range vms {
				rows = append(rows, []string{
					vm.Name,
					output.FormatState(vm.State),
				})
			}

			formatter.PrintTable(headers, rows)
		} else {
			formatter.Print(vms)
		}

		return nil
	},
}

// vmStartCmd represents the vm start command
var vmStartCmd = &cobra.Command{
	Use:   "start <vm>",
	Short: "Start a VM",
	Long:  "Start a virtual machine by name or UUID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		vm := args[0]
		fmt.Printf("Starting VM '%s'...\n", vm)

		if err := apiClient.StartVM(ctx, vm); err != nil {
			return fmt.Errorf("failed to start VM: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("✓ VM '%s' started successfully\n", vm)
		} else {
			formatter.Print(map[string]string{
				"status":  "success",
				"message": "VM started successfully",
				"vm":      vm,
			})
		}

		return nil
	},
}

// vmStopCmd represents the vm stop command
var vmStopCmd = &cobra.Command{
	Use:   "stop <vm>",
	Short: "Stop a VM",
	Long:  "Stop a virtual machine by name or UUID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		vm := args[0]
		fmt.Printf("Stopping VM '%s'...\n", vm)

		if err := apiClient.StopVM(ctx, vm); err != nil {
			return fmt.Errorf("failed to stop VM: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("✓ VM '%s' stopped successfully\n", vm)
		} else {
			formatter.Print(map[string]string{
				"status":  "success",
				"message": "VM stopped successfully",
				"vm":      vm,
			})
		}

		return nil
	},
}

// vmRestartCmd represents the vm restart command
var vmRestartCmd = &cobra.Command{
	Use:   "restart <vm>",
	Short: "Restart a VM",
	Long:  "Restart a virtual machine by name or UUID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		vm := args[0]
		fmt.Printf("Restarting VM '%s'...\n", vm)

		if err := apiClient.RestartVM(ctx, vm); err != nil {
			return fmt.Errorf("failed to restart VM: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("✓ VM '%s' restarted successfully\n", vm)
		} else {
			formatter.Print(map[string]string{
				"status":  "success",
				"message": "VM restarted successfully",
				"vm":      vm,
			})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(vmCmd)
	vmCmd.AddCommand(vmLsCmd)
	vmCmd.AddCommand(vmStartCmd)
	vmCmd.AddCommand(vmStopCmd)
	vmCmd.AddCommand(vmRestartCmd)
}
