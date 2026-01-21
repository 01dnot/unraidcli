package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ServerConfig holds configuration for a single Unraid server
type ServerConfig struct {
	URL    string `yaml:"url"`
	APIKey string `yaml:"api_key"`
}

// Config represents the application configuration
type Config struct {
	DefaultServer string                  `yaml:"default_server"`
	OutputFormat  string                  `yaml:"output_format"`
	Servers       map[string]ServerConfig `yaml:"servers"`
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, ".unraidcli", "config.yaml"), nil
}

// Load reads the configuration from disk
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// Return default config if file doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			OutputFormat: "table",
			Servers:      make(map[string]ServerConfig),
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults if not specified
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = "table"
	}
	if cfg.Servers == nil {
		cfg.Servers = make(map[string]ServerConfig)
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetServer returns the configuration for a specific server
// If serverName is empty, returns the default server
func (c *Config) GetServer(serverName string) (*ServerConfig, error) {
	if serverName == "" {
		serverName = c.DefaultServer
	}

	if serverName == "" {
		return nil, fmt.Errorf("no server specified and no default server configured")
	}

	server, ok := c.Servers[serverName]
	if !ok {
		return nil, fmt.Errorf("server '%s' not found in configuration", serverName)
	}

	return &server, nil
}

// SetServer adds or updates a server configuration
func (c *Config) SetServer(name string, url, apiKey string) {
	if c.Servers == nil {
		c.Servers = make(map[string]ServerConfig)
	}

	c.Servers[name] = ServerConfig{
		URL:    url,
		APIKey: apiKey,
	}

	// Set as default if it's the first server
	if c.DefaultServer == "" {
		c.DefaultServer = name
	}
}

// RemoveServer removes a server from the configuration
func (c *Config) RemoveServer(name string) error {
	if _, ok := c.Servers[name]; !ok {
		return fmt.Errorf("server '%s' not found", name)
	}

	delete(c.Servers, name)

	// Clear default if we just removed it
	if c.DefaultServer == name {
		c.DefaultServer = ""
		// Set a new default if there are other servers
		for serverName := range c.Servers {
			c.DefaultServer = serverName
			break
		}
	}

	return nil
}
