package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	pluginBundled bool
	pluginRestart bool
)

// pluginCmd represents the plugin command
var pluginCmd = &cobra.Command{
	Use:     "plugin",
	Aliases: []string{"plugins"},
	Short:   "Manage plugins",
	Long:    "List, add, and remove plugins on your Unraid server.",
}

// pluginLsCmd represents the plugin ls command
var pluginLsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all plugins",
	Long:    "List all installed plugins.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		plugins, err := apiClient.GetPlugins(ctx)
		if err != nil {
			return fmt.Errorf("failed to get plugins: %w", err)
		}

		if len(plugins) == 0 {
			fmt.Println("No plugins found.")
			return nil
		}

		if outputFormat == "" || outputFormat == "table" {
			headers := []string{"Name", "Version", "API Module", "CLI Module"}
			var rows [][]string

			for _, plugin := range plugins {
				apiModule := "No"
				if plugin.HasApiModule != nil && *plugin.HasApiModule {
					apiModule = "Yes"
				}

				cliModule := "No"
				if plugin.HasCliModule != nil && *plugin.HasCliModule {
					cliModule = "Yes"
				}

				rows = append(rows, []string{
					plugin.Name,
					plugin.Version,
					apiModule,
					cliModule,
				})
			}

			formatter.PrintTable(headers, rows)
			fmt.Printf("\nTotal: %d plugin(s)\n", len(plugins))
		} else {
			formatter.Print(plugins)
		}

		return nil
	},
}

// pluginAddCmd represents the plugin add command
var pluginAddCmd = &cobra.Command{
	Use:   "add <plugin> [plugin2] [plugin3]...",
	Short: "Add one or more plugins",
	Long:  "Add one or more plugins to the system.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()

		fmt.Printf("Adding %d plugin(s)...\n", len(args))

		if err := apiClient.AddPlugin(ctx, args, pluginBundled, pluginRestart); err != nil {
			return fmt.Errorf("failed to add plugins: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("✓ Successfully added %d plugin(s)\n", len(args))
			if pluginRestart {
				fmt.Println("\nNote: API restart may be required for changes to take effect.")
			}
		} else {
			formatter.Print(map[string]interface{}{
				"status":  "success",
				"message": "Plugins added successfully",
				"plugins": args,
			})
		}

		return nil
	},
}

// pluginRemoveCmd represents the plugin remove command
var pluginRemoveCmd = &cobra.Command{
	Use:     "remove <plugin> [plugin2] [plugin3]...",
	Aliases: []string{"rm", "uninstall"},
	Short:   "Remove one or more plugins",
	Long:    "Remove/uninstall one or more plugins from the system.",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		fmt.Printf("Removing %d plugin(s)...\n", len(args))

		if err := apiClient.RemovePlugin(ctx, args, pluginBundled, pluginRestart); err != nil {
			return fmt.Errorf("failed to remove plugins: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("✓ Successfully removed %d plugin(s)\n", len(args))
			if pluginRestart {
				fmt.Println("\nNote: API restart may be required for changes to take effect.")
			}
		} else {
			formatter.Print(map[string]interface{}{
				"status":  "success",
				"message": "Plugins removed successfully",
				"plugins": args,
			})
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
	pluginCmd.AddCommand(pluginLsCmd)
	pluginCmd.AddCommand(pluginAddCmd)
	pluginCmd.AddCommand(pluginRemoveCmd)

	// Add flags for add and remove commands
	pluginAddCmd.Flags().BoolVar(&pluginBundled, "bundled", false, "Treat plugins as bundled plugins")
	pluginAddCmd.Flags().BoolVar(&pluginRestart, "restart", true, "Restart the API after the operation")

	pluginRemoveCmd.Flags().BoolVar(&pluginBundled, "bundled", false, "Treat plugins as bundled plugins")
	pluginRemoveCmd.Flags().BoolVar(&pluginRestart, "restart", true, "Restart the API after the operation")
}
