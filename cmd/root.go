package cmd

import (
	"fmt"
	"os"

	"github.com/01dnot/unraidcli/internal/client"
	"github.com/01dnot/unraidcli/internal/config"
	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfgFile      string
	outputFormat string
	serverName   string
	cfg          *config.Config
	apiClient    *client.Client
	formatter    *output.Formatter
)

// Version information (set during build)
var (
	Version   = "dev"
	BuildDate = "unknown"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "unraidcli",
	Short: "A CLI tool for managing Unraid servers",
	Long: `unraidcli is a command-line interface for interacting with Unraid servers
using the official GraphQL API (available in Unraid 7.2+).

It allows you to manage Docker containers, VMs, the array, and view system
information directly from your terminal.`,
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config loading for config commands
		if cmd.Name() == "config" || cmd.Parent().Name() == "config" {
			return nil
		}

		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Get output format (command flag overrides config)
		format := outputFormat
		if format == "" {
			format = cfg.OutputFormat
		}
		formatter = output.New(format)

		// Get server configuration
		server, err := cfg.GetServer(serverName)
		if err != nil {
			return fmt.Errorf("failed to get server config: %w\n\nRun 'unraidcli config set' to configure a server", err)
		}

		// Create API client
		apiClient = client.New(server.URL, server.APIKey)

		return nil
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.unraidcli/config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format: table, json, yaml (default from config or 'table')")
	rootCmd.PersistentFlags().StringVarP(&serverName, "server", "s", "", "server profile name (default from config)")

	// Set version template
	rootCmd.SetVersionTemplate(fmt.Sprintf("unraidcli version %s (built %s)\n", Version, BuildDate))

	// Hide completion command from help (still accessible via 'unraidcli completion')
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
