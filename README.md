# unraidcli
![ChatGPT Image Jan 21, 2026, 10_01_15 PM](https://github.com/user-attachments/assets/7b78d56f-8afb-4116-ad4f-b1492bb18dd7)

> A powerful command-line interface for managing Unraid servers from your terminal

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Manage your Unraid server directly from the command line using the official GraphQL API (Unraid 7.2+). Monitor system health, control Docker containers and VMs, manage the array, and more - all without opening a web browser.

## Why unraidcli?

- 🚀 **Fast & Lightweight** - No browser required, instant access from your terminal
- 🔄 **Real-time Monitoring** - Watch mode for live system stats and container status
- 🎨 **Beautiful Output** - Colorized tables and human-readable formatting
- 🔧 **Automation Ready** - Perfect for scripts, SSH sessions, and remote management
- 📊 **Multi-Format** - JSON/YAML output for integration with other tools
- 🖥️ **Multi-Server** - Manage multiple Unraid servers with named profiles

## Features

- **Server Management**: View system information, status, and health overview
- **Array Control**: Start, stop, and monitor your Unraid storage array
- **Docker Management**: List, start, stop, restart, and view stats/logs for containers
- **VM Management**: Control virtual machines
- **Shares Management**: View and monitor user shares
- **Parity Check**: Monitor and control parity checks
- **Notifications**: View and manage system notifications
- **Metrics**: Real-time CPU and memory monitoring with per-core details
- **Logs**: View system logs
- **Plugin Management**: List, add, and remove plugins
- **Health Check**: Quick system health overview
- **Watch Mode**: Auto-refresh for real-time monitoring
- **Colorized Output**: Easy-to-read colored terminal output
- **Multiple Output Formats**: Table, JSON, and YAML output
- **Multi-Server Support**: Manage multiple Unraid servers with profiles
- **Easy Configuration**: Simple setup with built-in connection testing

## Requirements

- Unraid 7.2+ with GraphQL API enabled
- API key (created via Unraid WebGUI or CLI)

## Installation

### Quick Install (Recommended)

**macOS / Linux - One-liner:**
```bash
curl -sSL https://raw.githubusercontent.com/01dnot/unraidcli/main/install.sh | bash
```

This script automatically:
- Detects your OS and architecture
- Downloads the appropriate binary
- Installs to `/usr/local/bin`
- Verifies the installation

### Alternative Methods

**Using Go:**
```bash
go install github.com/01dnot/unraidcli@latest
```

### Manual Installation

**Download Pre-built Binary:**

Download from the [Releases page](https://github.com/01dnot/unraidcli/releases) or use curl:

```bash
# Linux AMD64
curl -L https://github.com/01dnot/unraidcli/releases/latest/download/unraidcli-linux-amd64 -o unraidcli

# Linux ARM64
curl -L https://github.com/01dnot/unraidcli/releases/latest/download/unraidcli-linux-arm64 -o unraidcli

# macOS Intel
curl -L https://github.com/01dnot/unraidcli/releases/latest/download/unraidcli-darwin-amd64 -o unraidcli

# macOS Apple Silicon
curl -L https://github.com/01dnot/unraidcli/releases/latest/download/unraidcli-darwin-arm64 -o unraidcli

# Make executable and install
chmod +x unraidcli
sudo mv unraidcli /usr/local/bin/
```

**Build from Source:**
```bash
git clone https://github.com/01dnot/unraidcli.git
cd unraidcli
go build -o unraidcli
sudo mv unraidcli /usr/local/bin/
```

## Quick Start

### 1. Create an API Key

Before using unraidcli, you need to create an API key on your Unraid server:

**Via WebGUI:**
1. Go to Settings → Management Access → API Keys
2. Click "Add" to create a new API key
3. Copy the generated key

**Via Unraid CLI:**
```bash
unraid-api apikey --create
```

### 2. Configure unraidcli

```bash
unraidcli config set --url http://192.168.1.100 --apikey YOUR_API_KEY
```

This will:
- Save the configuration to `~/.unraidcli/config.yaml`
- Test the connection to verify it works
- Set it as your default server profile

### 3. Test the Connection

```bash
unraidcli server info
```

## Usage

### Configuration Commands

```bash
# Set up a server configuration
unraidcli config set --url http://192.168.1.100 --apikey YOUR_API_KEY

# Set up a named server profile
unraidcli config set --name remote --url https://unraid.example.com --apikey YOUR_API_KEY

# Show current configuration
unraidcli config show

# List all server profiles
unraidcli config list

# Remove a server profile
unraidcli config remove remote
```

### Server Commands

```bash
# View detailed system information
unraidcli server info

# Check server status and uptime
unraidcli server status
```

### Array Commands

```bash
# View array status and disk information
unraidcli array status

# Start the array
unraidcli array start

# Stop the array
unraidcli array stop
```

### Docker Commands

```bash
# List all containers
unraidcli docker ls

# List with watch mode (auto-refresh every 2 seconds)
unraidcli docker ls --watch

# Filter by state
unraidcli docker ls --state running
unraidcli docker ls --state exited

# List only running containers
unraidcli docker ps

# Start a container
unraidcli docker start plex

# Stop a container
unraidcli docker stop plex

# Restart a container
unraidcli docker restart plex

# Bulk operations
unraidcli docker start-all plex sonarr radarr
unraidcli docker stop-all plex sonarr radarr

# View container stats
unraidcli docker stats
unraidcli docker stats plex sonarr  # Specific containers
unraidcli docker stats --watch      # Real-time monitoring

# View container logs information
unraidcli docker logs plex
```

### VM Commands

```bash
# List all virtual machines
unraidcli vm ls

# Start a VM
unraidcli vm start windows11

# Stop a VM
unraidcli vm stop windows11

# Restart a VM
unraidcli vm restart windows11
```

### Shares Commands

```bash
# List all user shares
unraidcli shares ls

# View detailed share information
unraidcli shares info media
```

### Metrics Commands

```bash
# View system metrics (CPU and memory)
unraidcli metrics

# Show per-core CPU usage
unraidcli metrics --cores

# Watch mode for real-time monitoring
unraidcli metrics --watch
unraidcli metrics --watch --interval 5  # Custom refresh interval
```

### Parity Check Commands

```bash
# View parity check status
unraidcli parity status

# View parity check history
unraidcli parity history

# Start a parity check
unraidcli parity start
unraidcli parity start --correct  # With error correction

# Pause/Resume/Cancel parity check
unraidcli parity pause
unraidcli parity resume
unraidcli parity cancel
```

### Notification Commands

```bash
# List all notifications
unraidcli notif ls

# Filter by type
unraidcli notif ls --type alert
unraidcli notif ls --type warning
unraidcli notif ls --type info

# Filter by importance
unraidcli notif ls --importance normal
unraidcli notif ls --importance urgent

# Archive a notification
unraidcli notif archive <id>

# Overview summary
unraidcli notif overview
```

### Logs Commands

```bash
# List available log files
unraidcli logs ls

# View a specific log file
unraidcli logs view syslog
unraidcli logs view syslog --lines 100    # Limit lines
unraidcli logs view syslog --tail         # Follow mode

# View last N lines of a log
unraidcli logs tail syslog 50
```

### Plugin Commands

```bash
# List all installed plugins
unraidcli plugin ls

# Add one or more plugins
unraidcli plugin add plugin-name
unraidcli plugin add plugin1 plugin2 plugin3

# Remove one or more plugins
unraidcli plugin remove plugin-name
unraidcli plugin rm plugin1 plugin2

# Advanced options
unraidcli plugin add plugin-name --bundled      # Treat as bundled plugin
unraidcli plugin add plugin-name --restart=false # Skip auto-restart
```

### Health Check

```bash
# Quick system health overview
unraidcli health

# Shows status for:
# - Array state and disk health
# - Parity check status and errors
# - Docker containers (running/stopped)
# - System resources (CPU/memory)
# - Notifications (alerts/warnings)
```

### Global Flags

All commands support these global flags:

```bash
# Use a specific server profile
unraidcli docker ls --server remote

# Change output format
unraidcli docker ls --output json
unraidcli docker ls --output yaml
unraidcli docker ls --output table  # default

# Use a custom config file
unraidcli docker ls --config /path/to/config.yaml
```

### Watch Mode

Many commands support watch mode for real-time monitoring:

```bash
# Auto-refresh every 2 seconds (default)
unraidcli docker ls --watch
unraidcli metrics --watch

# Custom refresh interval
unraidcli docker ls --watch --interval 5
unraidcli metrics --watch --interval 10

# Press Ctrl+C to exit watch mode
```

### Colorized Output

Output is automatically colorized for better readability:

- **States**: Green (running/started), Red (stopped/exited), Yellow (paused/starting)
- **Health**: Green (healthy/ok), Red (errors/failed), Yellow (warnings)
- **Percentages**: Green (low), Yellow (medium), Red (high)
- **Temperatures**: Blue (cool), Cyan (normal), Yellow (warm), Red (hot)

Colors are automatically disabled when:
- Output is piped to another command
- `NO_COLOR` environment variable is set
- Output is not a terminal

### Output Format Examples

**Table (default):**
```bash
$ unraidcli docker ls
Name     Image                State    Status              Autostart
plex     plexinc/pms-docker   RUNNING  Up 2 days           ✓
sonarr   linuxserver/sonarr   RUNNING  Up 2 days           ✓
```

**JSON:**
```bash
$ unraidcli docker ls --output json
[
  {
    "id": "abc123",
    "names": ["/plex"],
    "image": "plexinc/pms-docker",
    "state": "running",
    "status": "Up 2 days",
    "autostart": true
  }
]
```

**YAML:**
```bash
$ unraidcli docker ls --output yaml
- id: abc123
  names:
    - /plex
  image: plexinc/pms-docker
  state: running
  status: Up 2 days
  autostart: true
```

## Configuration File

The configuration file is stored at `~/.unraidcli/config.yaml`:

```yaml
default_server: "home"
output_format: "table"

servers:
  home:
    url: "http://192.168.1.100"
    api_key: "your-api-key-here"

  remote:
    url: "https://unraid.example.com"
    api_key: "another-api-key"
```

## Multi-Server Management

You can manage multiple Unraid servers by creating named profiles:

```bash
# Add multiple servers
unraidcli config set --name home --url http://192.168.1.100 --apikey KEY1
unraidcli config set --name remote --url https://unraid.example.com --apikey KEY2

# Use a specific server
unraidcli docker ls --server remote
unraidcli server info --server home

# List all configured servers
unraidcli config list
```

## Security Best Practices

### API Key Security
- **Least Privilege**: Create API keys with only the permissions you need
- **Rotation**: Regularly regenerate API keys
- **Storage**: API keys are stored in `~/.unraidcli/config.yaml` with 0600 permissions
- **Shell History**: After running `config set --apikey`, clear your shell history to avoid exposing the key

### Network Security
- **HTTPS Required**: Always use HTTPS for remote connections to protect your API key in transit
- **Local Networks**: HTTP is acceptable only on trusted local networks
- **VPN Access**: Use Tailscale or WireGuard for secure remote access instead of exposing to the internet

### Clearing Shell History
```bash
# After setting up your API key
history -d $(history 1)  # Removes last command from history
```

## API Permissions

Your API key needs appropriate permissions for the operations you want to perform:

- **Read operations** (ls, status, info): Require read permissions
- **Write operations** (start, stop, restart): Require write permissions

Configure permissions when creating the API key in the Unraid WebGUI.

For more security information, see [SECURITY.md](SECURITY.md).

## Troubleshooting

### Connection Issues

If you get connection errors:

1. Verify your Unraid server URL is correct
2. Ensure you're using Unraid 7.2 or later
3. Check that the GraphQL API is accessible at `http://YOUR_SERVER/graphql`
4. Verify your API key is correct and has necessary permissions

### Testing Connection

```bash
# Test if the server is reachable
curl http://YOUR_SERVER/graphql

# Test connection through unraidcli
unraidcli server info
```

### Debug Mode

For more verbose output, you can use the `--help` flag on any command to see all available options.

## Development

### Building from Source

```bash
git clone https://github.com/01dnot/unraidcli.git
cd unraidcli
go mod download
go build -o unraidcli
```

### Running Tests

```bash
go test ./...
```

### Building for Multiple Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o unraidcli-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o unraidcli-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o unraidcli-darwin-arm64

# Windows
GOOS=windows GOARCH=amd64 go build -o unraidcli-windows-amd64.exe
```

## Project Structure

```
unraidcli/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command and global flags
│   ├── config.go          # Configuration commands
│   ├── server.go          # Server information commands
│   ├── array.go           # Array management commands
│   ├── docker.go          # Docker container commands
│   ├── vm.go              # VM management commands
│   ├── shares.go          # Share management commands
│   ├── metrics.go         # System metrics commands
│   ├── parity.go          # Parity check commands
│   ├── notifications.go   # Notification commands
│   ├── logs.go            # Log viewing commands
│   └── health.go          # Health check command
├── internal/
│   ├── client/            # GraphQL client wrapper
│   │   └── unraid.go
│   ├── config/            # Configuration management
│   │   └── config.go
│   └── output/            # Output formatting
│       ├── formatter.go   # Table, JSON, YAML formatters
│       ├── color.go       # Colorized output
│       └── watch.go       # Watch mode implementation
├── main.go                # Application entry point
├── go.mod                 # Go module definition
├── .gitignore             # Git ignore rules
└── README.md              # This file
```

## Contributing

Contributions are welcome! Here's how you can help:

1. **Report Bugs** - Open an issue with details about the problem
2. **Suggest Features** - Share your ideas for new functionality
3. **Submit PRs** - Fix bugs or add features (please open an issue first to discuss)
4. **Improve Docs** - Help make the documentation clearer

### Development Setup

```bash
# Clone the repository
git clone https://github.com/01dnot/unraidcli.git
cd unraidcli

# Install dependencies
go mod download

# Build
go build -o unraidcli

# Run tests (if available)
go test ./...
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](https://github.com/01dnot/unraidcli/blob/main/README.md)
- 🐛 [Issue Tracker](https://github.com/01dnot/unraidcli/issues)
- 💬 [Discussions](https://github.com/01dnot/unraidcli/discussions)

## Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- Uses [Viper](https://github.com/spf13/viper) for configuration management
- GraphQL client by [machinebox](https://github.com/machinebox/graphql)
- YAML parsing with [go-yaml](https://github.com/go-yaml/yaml)

## Links

- [Unraid API Documentation](https://docs.unraid.net/API/)
- [Unraid API Usage Guide](https://docs.unraid.net/API/how-to-use-the-api/)
- [Issue Tracker](https://github.com/01dnot/unraidcli/issues)
