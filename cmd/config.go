package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/01dnot/unraidcli/internal/client"
	"github.com/01dnot/unraidcli/internal/config"
	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

var (
	configURL    string
	configAPIKey string
	configName   string
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Manage unraidcli configuration including server profiles and settings.",
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set server configuration",
	Long: `Set or update server configuration. Creates a new server profile or updates an existing one.

Examples:
  unraidcli config set --url http://192.168.1.100 --apikey YOUR_API_KEY
  unraidcli config set --name remote --url https://unraid.example.com --apikey YOUR_API_KEY`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load or create config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Use default name if not specified
		name := configName
		if name == "" {
			name = "default"
		}

		// Validate required fields
		if configURL == "" {
			return fmt.Errorf("--url is required")
		}
		if configAPIKey == "" {
			return fmt.Errorf("--apikey is required")
		}

		// Test connection
		fmt.Printf("Testing connection to %s...\n", configURL)
		testClient := client.New(configURL, configAPIKey)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := testClient.TestConnection(ctx); err != nil {
			return fmt.Errorf("connection test failed: %w", err)
		}

		// Save configuration
		cfg.SetServer(name, configURL, configAPIKey)
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✓ Server '%s' configured successfully\n", name)
		if cfg.DefaultServer == name {
			fmt.Printf("✓ Set as default server\n")
		}

		return nil
	},
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  "Display the current configuration including all server profiles.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		configPath, _ := config.GetConfigPath()
		fmt.Printf("Config file: %s\n\n", configPath)

		if len(cfg.Servers) == 0 {
			fmt.Println("No servers configured.")
			fmt.Println("\nRun 'unraidcli config set --url <url> --apikey <key>' to add a server.")
			return nil
		}

		fmt.Printf("Default server: %s\n", cfg.DefaultServer)
		fmt.Printf("Output format: %s\n\n", cfg.OutputFormat)

		fmt.Println("Servers:")
		for name, server := range cfg.Servers {
			defaultMarker := ""
			if name == cfg.DefaultServer {
				defaultMarker = " (default)"
			}
			fmt.Printf("  %s%s:\n", name, defaultMarker)
			fmt.Printf("    URL: %s\n", server.URL)
			// Mask API key for security
			maskedKey := "***" + server.APIKey[len(server.APIKey)-4:]
			if len(server.APIKey) < 8 {
				maskedKey = "***"
			}
			fmt.Printf("    API Key: %s\n", maskedKey)
		}

		return nil
	},
}

// configListCmd represents the config ls command
var configListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all server profiles",
	Long:    "List all configured server profiles.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if len(cfg.Servers) == 0 {
			fmt.Println("No servers configured.")
			return nil
		}

		formatter := output.New(outputFormat)
		if outputFormat == "" || outputFormat == "table" {
			headers := []string{"Name", "URL", "Default"}
			var rows [][]string

			for name, server := range cfg.Servers {
				isDefault := ""
				if name == cfg.DefaultServer {
					isDefault = "✓"
				}
				rows = append(rows, []string{name, server.URL, isDefault})
			}

			formatter.PrintTable(headers, rows)
		} else {
			type ServerInfo struct {
				Name      string `json:"name" yaml:"name"`
				URL       string `json:"url" yaml:"url"`
				IsDefault bool   `json:"is_default" yaml:"is_default"`
			}

			var servers []ServerInfo
			for name, server := range cfg.Servers {
				servers = append(servers, ServerInfo{
					Name:      name,
					URL:       server.URL,
					IsDefault: name == cfg.DefaultServer,
				})
			}

			formatter.Print(servers)
		}

		return nil
	},
}

// configRemoveCmd represents the config remove command
var configRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a server profile",
	Long:  "Remove a server profile from the configuration.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		name := args[0]
		if err := cfg.RemoveServer(name); err != nil {
			return err
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Printf("✓ Server '%s' removed\n", name)
		if cfg.DefaultServer != "" {
			fmt.Printf("✓ Default server is now '%s'\n", cfg.DefaultServer)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configRemoveCmd)

	// Flags for config set
	configSetCmd.Flags().StringVar(&configURL, "url", "", "Unraid server URL (e.g., http://192.168.1.100)")
	configSetCmd.Flags().StringVar(&configAPIKey, "apikey", "", "API key for authentication")
	configSetCmd.Flags().StringVar(&configName, "name", "", "Server profile name (default: 'default')")
}
