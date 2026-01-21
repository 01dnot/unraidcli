package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/01dnot/unraidcli/internal/client"
	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

var (
	watchMode      bool
	watchInterval  int
	filterByState  string
)

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Manage Docker containers",
	Long:  "List, start, stop, and restart Docker containers on your Unraid server.",
}

// dockerLsCmd represents the docker ls command
var dockerLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all containers",
	Long:  "List all Docker containers (both running and stopped).",
	RunE: func(cmd *cobra.Command, args []string) error {
		listFunc := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			containers, err := apiClient.GetContainers(ctx)
			if err != nil {
				return fmt.Errorf("failed to get containers: %w", err)
			}

			// Filter by state if specified
			if filterByState != "" {
				var filtered []client.Container
				for _, c := range containers {
					if strings.EqualFold(c.State, filterByState) {
						filtered = append(filtered, c)
					}
				}
				containers = filtered
			}

			if len(containers) == 0 {
				fmt.Println("No containers found.")
				return nil
			}

			if outputFormat == "" || outputFormat == "table" {
				// Show timestamp in watch mode
				if watchMode {
					fmt.Printf("Last updated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
				}

				headers := []string{"Name", "Image", "State", "Status", "Autostart"}
				var rows [][]string

				for _, container := range containers {
					name := container.Names[0]
					if strings.HasPrefix(name, "/") {
						name = name[1:]
					}

					autostart := ""
					if container.Autostart {
						autostart = "✓"
					}

					rows = append(rows, []string{
						name,
						container.Image,
						output.FormatState(container.State),
						container.Status,
						autostart,
					})
				}

				formatter.PrintTable(headers, rows)
			} else {
				formatter.Print(containers)
			}

			return nil
		}

		// Watch mode
		if watchMode {
			// Setup signal handling for graceful exit
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-sigChan
				cancel()
			}()

			interval := time.Duration(watchInterval) * time.Second
			return output.Watch(ctx, interval, listFunc)
		}

		// Normal mode
		return listFunc()
	},
}

// dockerPsCmd represents the docker ps command
var dockerPsCmd = &cobra.Command{
	Use:   "ps",
	Short: "List running containers",
	Long:  "List only running Docker containers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		containers, err := apiClient.GetContainers(ctx)
		if err != nil {
			return fmt.Errorf("failed to get containers: %w", err)
		}

		// Filter for running containers
		var runningContainers []client.Container
		for _, container := range containers {
			if strings.ToLower(container.State) == "running" {
				runningContainers = append(runningContainers, container)
			}
		}

		if len(runningContainers) == 0 {
			fmt.Println("No running containers found.")
			return nil
		}

		if outputFormat == "" || outputFormat == "table" {
			headers := []string{"Name", "Image", "Status", "Autostart"}
			var rows [][]string

			for _, container := range runningContainers {
				name := container.Names[0]
				if strings.HasPrefix(name, "/") {
					name = name[1:]
				}

				autostart := ""
				if container.Autostart {
					autostart = "✓"
				}

				rows = append(rows, []string{
					name,
					container.Image,
					container.Status,
					autostart,
				})
			}

			formatter.PrintTable(headers, rows)
		} else {
			formatter.Print(runningContainers)
		}

		return nil
	},
}

// dockerStartCmd represents the docker start command
var dockerStartCmd = &cobra.Command{
	Use:   "start <container>",
	Short: "Start a container",
	Long:  "Start a Docker container by name or ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		container := args[0]
		fmt.Printf("Starting container '%s'...\n", container)

		if err := apiClient.StartContainer(ctx, container); err != nil {
			return fmt.Errorf("failed to start container: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("✓ Container '%s' started successfully\n", container)
		} else {
			formatter.Print(map[string]string{
				"status":    "success",
				"message":   "Container started successfully",
				"container": container,
			})
		}

		return nil
	},
}

// dockerStopCmd represents the docker stop command
var dockerStopCmd = &cobra.Command{
	Use:   "stop <container>",
	Short: "Stop a container",
	Long:  "Stop a Docker container by name or ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		container := args[0]
		fmt.Printf("Stopping container '%s'...\n", container)

		if err := apiClient.StopContainer(ctx, container); err != nil {
			return fmt.Errorf("failed to stop container: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("✓ Container '%s' stopped successfully\n", container)
		} else {
			formatter.Print(map[string]string{
				"status":    "success",
				"message":   "Container stopped successfully",
				"container": container,
			})
		}

		return nil
	},
}

// dockerRestartCmd represents the docker restart command
var dockerRestartCmd = &cobra.Command{
	Use:   "restart <container>",
	Short: "Restart a container",
	Long:  "Restart a Docker container by name or ID.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		container := args[0]
		fmt.Printf("Restarting container '%s'...\n", container)

		if err := apiClient.RestartContainer(ctx, container); err != nil {
			return fmt.Errorf("failed to restart container: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Printf("✓ Container '%s' restarted successfully\n", container)
		} else {
			formatter.Print(map[string]string{
				"status":    "success",
				"message":   "Container restarted successfully",
				"container": container,
			})
		}

		return nil
	},
}

// dockerStartAllCmd represents the docker start-all command
var dockerStartAllCmd = &cobra.Command{
	Use:   "start-all [container1] [container2] ...",
	Short: "Start multiple containers",
	Long:  "Start multiple Docker containers at once.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		fmt.Printf("Starting %d container(s)...\n", len(args))

		var failures []string
		for _, container := range args {
			fmt.Printf("  Starting '%s'... ", container)
			if err := apiClient.StartContainer(ctx, container); err != nil {
				fmt.Printf("✗ Failed: %v\n", err)
				failures = append(failures, container)
			} else {
				fmt.Printf("✓\n")
			}
		}

		if len(failures) > 0 {
			return fmt.Errorf("failed to start %d container(s): %v", len(failures), failures)
		}

		fmt.Printf("\n✓ Successfully started all %d container(s)\n", len(args))
		return nil
	},
}

// dockerStopAllCmd represents the docker stop-all command
var dockerStopAllCmd = &cobra.Command{
	Use:   "stop-all [container1] [container2] ...",
	Short: "Stop multiple containers",
	Long:  "Stop multiple Docker containers at once.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		fmt.Printf("Stopping %d container(s)...\n", len(args))

		var failures []string
		for _, container := range args {
			fmt.Printf("  Stopping '%s'... ", container)
			if err := apiClient.StopContainer(ctx, container); err != nil {
				fmt.Printf("✗ Failed: %v\n", err)
				failures = append(failures, container)
			} else {
				fmt.Printf("✓\n")
			}
		}

		if len(failures) > 0 {
			return fmt.Errorf("failed to stop %d container(s): %v", len(failures), failures)
		}

		fmt.Printf("\n✓ Successfully stopped all %d container(s)\n", len(args))
		return nil
	},
}

// dockerStatsCmd represents the docker stats command
var dockerStatsCmd = &cobra.Command{
	Use:   "stats [container...]",
	Short: "Display container resource usage statistics",
	Long:  "Show CPU, memory, network I/O, and block I/O statistics for containers.",
	RunE: func(cmd *cobra.Command, args []string) error {
		statsFunc := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			containers, err := apiClient.GetContainers(ctx)
			if err != nil {
				return fmt.Errorf("failed to get containers: %w", err)
			}

			// Filter by specific containers if provided
			if len(args) > 0 {
				var filtered []client.Container
				for _, container := range containers {
					name := container.Names[0]
					if strings.HasPrefix(name, "/") {
						name = name[1:]
					}
					for _, arg := range args {
						if name == arg || container.ID == arg || strings.HasPrefix(container.ID, arg) {
							filtered = append(filtered, container)
							break
						}
					}
				}
				containers = filtered
			}

			// Only show running containers for stats
			var running []client.Container
			for _, c := range containers {
				if strings.ToLower(c.State) == "running" {
					running = append(running, c)
				}
			}

			if len(running) == 0 {
				fmt.Println("No running containers found.")
				return nil
			}

			if outputFormat == "" || outputFormat == "table" {
				if watchMode {
					fmt.Printf("Last updated: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
				}

				headers := []string{"Name", "State", "Status", "Image"}
				var rows [][]string

				for _, container := range running {
					name := container.Names[0]
					if strings.HasPrefix(name, "/") {
						name = name[1:]
					}

					rows = append(rows, []string{
						name,
						output.ColorizeState(container.State),
						container.Status,
						container.Image,
					})
				}

				formatter.PrintTable(headers, rows)
				fmt.Println("\nNote: Detailed per-container CPU/Memory stats require additional API support.")
			} else {
				formatter.Print(running)
			}

			return nil
		}

		// Watch mode
		if watchMode {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
			go func() {
				<-sigChan
				cancel()
			}()

			fmt.Println("Press Ctrl+C to exit watch mode\n")
			interval := time.Duration(watchInterval) * time.Second
			return output.Watch(ctx, interval, statsFunc)
		}

		return statsFunc()
	},
}

// dockerLogsCmd represents the docker logs command
var dockerLogsCmd = &cobra.Command{
	Use:   "logs <container>",
	Short: "Fetch container logs",
	Long:  "Display logs from a specific container.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		containerName := args[0]

		// Find the container to verify it exists
		containers, err := apiClient.GetContainers(ctx)
		if err != nil {
			return fmt.Errorf("failed to get containers: %w", err)
		}

		var found *client.Container
		for _, container := range containers {
			name := container.Names[0]
			if strings.HasPrefix(name, "/") {
				name = name[1:]
			}
			if name == containerName || container.ID == containerName || strings.HasPrefix(container.ID, containerName) {
				found = &container
				break
			}
		}

		if found == nil {
			return fmt.Errorf("container '%s' not found", containerName)
		}

		name := found.Names[0]
		if strings.HasPrefix(name, "/") {
			name = name[1:]
		}

		fmt.Printf("Container: %s\n", name)
		fmt.Printf("Image: %s\n", found.Image)
		fmt.Printf("State: %s\n", output.ColorizeState(found.State))
		fmt.Printf("Status: %s\n\n", found.Status)

		fmt.Println("Note: Direct container log streaming requires Docker socket access.")
		fmt.Println("You can view logs on your Unraid server using:")
		fmt.Printf("  docker logs %s\n", name)
		fmt.Printf("  docker logs -f %s  (follow mode)\n", name)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.AddCommand(dockerLsCmd)
	dockerCmd.AddCommand(dockerPsCmd)
	dockerCmd.AddCommand(dockerStartCmd)
	dockerCmd.AddCommand(dockerStopCmd)
	dockerCmd.AddCommand(dockerRestartCmd)
	dockerCmd.AddCommand(dockerStartAllCmd)
	dockerCmd.AddCommand(dockerStopAllCmd)
	dockerCmd.AddCommand(dockerStatsCmd)
	dockerCmd.AddCommand(dockerLogsCmd)

	// Add flags for docker ls
	dockerLsCmd.Flags().BoolVarP(&watchMode, "watch", "w", false, "Watch mode - auto-refresh every N seconds")
	dockerLsCmd.Flags().IntVarP(&watchInterval, "interval", "i", 2, "Refresh interval in seconds for watch mode")
	dockerLsCmd.Flags().StringVar(&filterByState, "state", "", "Filter by state (running, exited)")

	// Add flags for docker ps
	dockerPsCmd.Flags().BoolVarP(&watchMode, "watch", "w", false, "Watch mode - auto-refresh every N seconds")
	dockerPsCmd.Flags().IntVarP(&watchInterval, "interval", "i", 2, "Refresh interval in seconds for watch mode")

	// Add flags for docker stats
	dockerStatsCmd.Flags().BoolVarP(&watchMode, "watch", "w", false, "Watch mode - auto-refresh every N seconds")
	dockerStatsCmd.Flags().IntVarP(&watchInterval, "interval", "i", 2, "Refresh interval in seconds for watch mode")
}
