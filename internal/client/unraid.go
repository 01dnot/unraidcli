package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/machinebox/graphql"
)

// Client represents an Unraid GraphQL API client
type Client struct {
	graphql *graphql.Client
	apiKey  string
	url     string
}

// New creates a new Unraid API client
func New(serverURL, apiKey string) *Client {
	// Ensure URL ends with /graphql
	if serverURL[len(serverURL)-8:] != "/graphql" {
		serverURL = serverURL + "/graphql"
	}

	client := graphql.NewClient(serverURL)

	return &Client{
		graphql: client,
		apiKey:  apiKey,
		url:     serverURL,
	}
}

// Query executes a GraphQL query
func (c *Client) Query(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	req := graphql.NewRequest(query)

	// Add API key header
	req.Header.Set("x-api-key", c.apiKey)

	// Add variables if provided
	if variables != nil {
		for key, value := range variables {
			req.Var(key, value)
		}
	}

	if err := c.graphql.Run(ctx, req, response); err != nil {
		return fmt.Errorf("GraphQL query failed: %w", err)
	}

	return nil
}

// Mutate executes a GraphQL mutation
func (c *Client) Mutate(ctx context.Context, mutation string, variables map[string]interface{}, response interface{}) error {
	req := graphql.NewRequest(mutation)

	// Add API key header
	req.Header.Set("x-api-key", c.apiKey)

	// Add variables if provided
	if variables != nil {
		for key, value := range variables {
			req.Var(key, value)
		}
	}

	if err := c.graphql.Run(ctx, req, response); err != nil {
		return fmt.Errorf("GraphQL mutation failed: %w", err)
	}

	return nil
}

// TestConnection verifies the connection to the Unraid server
func (c *Client) TestConnection(ctx context.Context) error {
	query := `
		query {
			info {
				id
				os {
					platform
				}
			}
		}
	`

	var response struct {
		Info struct {
			ID string `json:"id"`
			OS struct {
				Platform string `json:"platform"`
			} `json:"os"`
		} `json:"info"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}

	if response.Info.ID == "" {
		return fmt.Errorf("received empty response from server")
	}

	return nil
}

// SystemInfo contains system information
type SystemInfo struct {
	CPU struct {
		Manufacturer string `json:"manufacturer"`
		Brand        string `json:"brand"`
		Cores        int    `json:"cores"`
		Threads      int    `json:"threads"`
		Speed        float64 `json:"speed"`
	} `json:"cpu"`
	Memory struct {
		Layout []struct {
			Size int64 `json:"size"`
		} `json:"layout"`
	} `json:"memory"`
	OS struct {
		Platform string `json:"platform"`
		Hostname string `json:"hostname"`
		Uptime   string `json:"uptime"`
	} `json:"os"`
	Versions struct {
		Core struct {
			Unraid string `json:"unraid"`
		} `json:"core"`
	} `json:"versions"`
}

// GetSystemInfo retrieves system information
func (c *Client) GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	query := `
		query {
			info {
				cpu {
					manufacturer
					brand
					cores
					threads
					speed
				}
				memory {
					layout {
						size
					}
				}
				os {
					platform
					hostname
					uptime
				}
				versions {
					core {
						unraid
					}
				}
			}
			metrics {
				memory {
					total
					used
					available
				}
			}
		}
	`

	var response struct {
		Info SystemInfo `json:"info"`
		Metrics struct {
			Memory struct {
				Total     int64 `json:"total"`
				Used      int64 `json:"used"`
				Available int64 `json:"available"`
			} `json:"memory"`
		} `json:"metrics"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	// Return the info, we'll handle metrics separately if needed
	return &response.Info, nil
}

// ArrayInfo contains array information
type ArrayInfo struct {
	State    string `json:"state"`
	Capacity struct {
		Kilobytes struct {
			Total string `json:"total"`
			Used  string `json:"used"`
			Free  string `json:"free"`
		} `json:"kilobytes"`
	} `json:"capacity"`
	Boot     *ArrayDisk   `json:"boot"`
	Parities []ArrayDisk  `json:"parities"`
	Disks    []ArrayDisk  `json:"disks"`
	Caches   []ArrayDisk  `json:"caches"`
}

// AllDisks returns all disks (boot, parities, data disks, and caches) in order
func (a *ArrayInfo) AllDisks() []ArrayDisk {
	var allDisks []ArrayDisk

	// Add boot disk
	if a.Boot != nil {
		allDisks = append(allDisks, *a.Boot)
	}

	// Add parity disks
	allDisks = append(allDisks, a.Parities...)

	// Add data disks
	allDisks = append(allDisks, a.Disks...)

	// Add cache disks
	allDisks = append(allDisks, a.Caches...)

	return allDisks
}

// ArrayDisk represents a disk in the array
type ArrayDisk struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Device      string  `json:"device"`
	Status      string  `json:"status"`
	Size        int64   `json:"size"`
	Temperature int     `json:"temp"`
	Type        string  `json:"type"`
	FsType      string  `json:"fsType"`
}

// GetArrayInfo retrieves array information
func (c *Client) GetArrayInfo(ctx context.Context) (*ArrayInfo, error) {
	query := `
		query {
			array {
				state
				capacity {
					kilobytes {
						total
						used
						free
					}
				}
				boot {
					id
					name
					device
					status
					size
					temp
					type
					fsType
				}
				parities {
					id
					name
					device
					status
					size
					temp
					type
					fsType
				}
				disks {
					id
					name
					device
					status
					size
					temp
					type
					fsType
				}
				caches {
					id
					name
					device
					status
					size
					temp
					type
					fsType
				}
			}
		}
	`

	var response struct {
		Array ArrayInfo `json:"array"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return &response.Array, nil
}

// StartArray starts the Unraid array
func (c *Client) StartArray(ctx context.Context) error {
	mutation := `
		mutation {
			array {
				setState(input: {desiredState: START}) {
					state
				}
			}
		}
	`

	var response struct {
		Array struct {
			SetState struct {
				State string `json:"state"`
			} `json:"setState"`
		} `json:"array"`
	}

	if err := c.Mutate(ctx, mutation, nil, &response); err != nil {
		return err
	}

	return nil
}

// StopArray stops the Unraid array
func (c *Client) StopArray(ctx context.Context) error {
	mutation := `
		mutation {
			array {
				setState(input: {desiredState: STOP}) {
					state
				}
			}
		}
	`

	var response struct {
		Array struct {
			SetState struct {
				State string `json:"state"`
			} `json:"setState"`
		} `json:"array"`
	}

	if err := c.Mutate(ctx, mutation, nil, &response); err != nil {
		return err
	}

	return nil
}

// Container represents a Docker container
type Container struct {
	ID        string   `json:"id"`
	Names     []string `json:"names"`
	Image     string   `json:"image"`
	State     string   `json:"state"`
	Status    string   `json:"status"`
	Autostart bool     `json:"autoStart"`
}

// GetContainers retrieves all Docker containers
func (c *Client) GetContainers(ctx context.Context) ([]Container, error) {
	query := `
		query {
			docker {
				containers {
					id
					names
					image
					state
					status
					autoStart
				}
			}
		}
	`

	var response struct {
		Docker struct {
			Containers []Container `json:"containers"`
		} `json:"docker"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return response.Docker.Containers, nil
}

// FindContainerID finds a container ID by name or partial ID
func (c *Client) FindContainerID(ctx context.Context, nameOrID string) (string, error) {
	containers, err := c.GetContainers(ctx)
	if err != nil {
		return "", err
	}

	// Try exact ID match first
	for _, container := range containers {
		if container.ID == nameOrID {
			return container.ID, nil
		}
	}

	// Try name match
	for _, container := range containers {
		for _, name := range container.Names {
			// Names might have leading slash
			cleanName := strings.TrimPrefix(name, "/")
			if cleanName == nameOrID {
				return container.ID, nil
			}
		}
	}

	// Try partial ID match
	for _, container := range containers {
		if strings.HasPrefix(container.ID, nameOrID) {
			return container.ID, nil
		}
	}

	return "", fmt.Errorf("container not found: %s", nameOrID)
}

// StartContainer starts a Docker container
func (c *Client) StartContainer(ctx context.Context, nameOrID string) error {
	// Find the container ID
	id, err := c.FindContainerID(ctx, nameOrID)
	if err != nil {
		return err
	}

	mutation := `
		mutation($id: PrefixedID!) {
			docker {
				start(id: $id) {
					id
					state
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Docker struct {
			Start struct {
				ID    string `json:"id"`
				State string `json:"state"`
			} `json:"start"`
		} `json:"docker"`
	}

	if err := c.Mutate(ctx, mutation, variables, &response); err != nil {
		return err
	}

	return nil
}

// StopContainer stops a Docker container
func (c *Client) StopContainer(ctx context.Context, nameOrID string) error {
	// Find the container ID
	id, err := c.FindContainerID(ctx, nameOrID)
	if err != nil {
		return err
	}

	mutation := `
		mutation($id: PrefixedID!) {
			docker {
				stop(id: $id) {
					id
					state
				}
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		Docker struct {
			Stop struct {
				ID    string `json:"id"`
				State string `json:"state"`
			} `json:"stop"`
		} `json:"docker"`
	}

	if err := c.Mutate(ctx, mutation, variables, &response); err != nil {
		return err
	}

	return nil
}

// RestartContainer restarts a Docker container (stop then start)
func (c *Client) RestartContainer(ctx context.Context, nameOrID string) error {
	// Stop the container
	if err := c.StopContainer(ctx, nameOrID); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	// Wait a moment
	time.Sleep(2 * time.Second)

	// Start the container
	if err := c.StartContainer(ctx, nameOrID); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	return nil
}

// VM represents a virtual machine
type VM struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
}

// GetVMs retrieves all virtual machines
func (c *Client) GetVMs(ctx context.Context) ([]VM, error) {
	query := `
		query {
			vms {
				domains {
					id
					name
					state
				}
			}
		}
	`

	var response struct {
		VMs struct {
			Domains []VM `json:"domains"`
		} `json:"vms"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return response.VMs.Domains, nil
}

// FindVMID finds a VM ID by name or ID
func (c *Client) FindVMID(ctx context.Context, nameOrID string) (string, error) {
	vms, err := c.GetVMs(ctx)
	if err != nil {
		return "", err
	}

	// Try exact ID match first
	for _, vm := range vms {
		if vm.ID == nameOrID {
			return vm.ID, nil
		}
	}

	// Try name match
	for _, vm := range vms {
		if vm.Name == nameOrID {
			return vm.ID, nil
		}
	}

	// Try partial ID match
	for _, vm := range vms {
		if strings.HasPrefix(vm.ID, nameOrID) {
			return vm.ID, nil
		}
	}

	return "", fmt.Errorf("VM not found: %s", nameOrID)
}

// StartVM starts a virtual machine
func (c *Client) StartVM(ctx context.Context, nameOrID string) error {
	// Find the VM ID
	id, err := c.FindVMID(ctx, nameOrID)
	if err != nil {
		return err
	}

	mutation := `
		mutation($id: PrefixedID!) {
			vm {
				start(id: $id)
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		VM struct {
			Start bool `json:"start"`
		} `json:"vm"`
	}

	if err := c.Mutate(ctx, mutation, variables, &response); err != nil {
		return err
	}

	return nil
}

// StopVM stops a virtual machine
func (c *Client) StopVM(ctx context.Context, nameOrID string) error {
	// Find the VM ID
	id, err := c.FindVMID(ctx, nameOrID)
	if err != nil {
		return err
	}

	mutation := `
		mutation($id: PrefixedID!) {
			vm {
				stop(id: $id)
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		VM struct {
			Stop bool `json:"stop"`
		} `json:"vm"`
	}

	if err := c.Mutate(ctx, mutation, variables, &response); err != nil {
		return err
	}

	return nil
}

// RestartVM restarts a virtual machine using reboot
func (c *Client) RestartVM(ctx context.Context, nameOrID string) error {
	// Find the VM ID
	id, err := c.FindVMID(ctx, nameOrID)
	if err != nil {
		return err
	}

	mutation := `
		mutation($id: PrefixedID!) {
			vm {
				reboot(id: $id)
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response struct {
		VM struct {
			Reboot bool `json:"reboot"`
		} `json:"vm"`
	}

	if err := c.Mutate(ctx, mutation, variables, &response); err != nil {
		return err
	}

	return nil
}

// Share represents a user share
type Share struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Free    int64    `json:"free"`
	Used    int64    `json:"used"`
	Size    int64    `json:"size"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
	Cache   bool     `json:"cache"`
	Comment string   `json:"comment"`
}

// GetShares retrieves all user shares
func (c *Client) GetShares(ctx context.Context) ([]Share, error) {
	query := `
		query {
			shares {
				id
				name
				free
				used
				size
				include
				exclude
				cache
				comment
			}
		}
	`

	var response struct {
		Shares []Share `json:"shares"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return response.Shares, nil
}

// Metrics contains system metrics
type Metrics struct {
	CPU struct {
		PercentTotal float64 `json:"percentTotal"`
		CPUs         []struct {
			PercentTotal  float64 `json:"percentTotal"`
			PercentUser   float64 `json:"percentUser"`
			PercentSystem float64 `json:"percentSystem"`
			PercentIdle   float64 `json:"percentIdle"`
		} `json:"cpus"`
	} `json:"cpu"`
	Memory struct {
		Total            int64   `json:"total"`
		Used             int64   `json:"used"`
		Free             int64   `json:"free"`
		Available        int64   `json:"available"`
		PercentTotal     float64 `json:"percentTotal"`
		SwapTotal        int64   `json:"swapTotal"`
		SwapUsed         int64   `json:"swapUsed"`
		SwapFree         int64   `json:"swapFree"`
		PercentSwapTotal float64 `json:"percentSwapTotal"`
	} `json:"memory"`
}

// GetMetrics retrieves current system metrics
func (c *Client) GetMetrics(ctx context.Context) (*Metrics, error) {
	query := `
		query {
			metrics {
				cpu {
					percentTotal
					cpus {
						percentTotal
						percentUser
						percentSystem
						percentIdle
					}
				}
				memory {
					total
					used
					free
					available
					percentTotal
					swapTotal
					swapUsed
					swapFree
					percentSwapTotal
				}
			}
		}
	`

	var response struct {
		Metrics Metrics `json:"metrics"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return &response.Metrics, nil
}

// ParityCheck represents a parity check operation
type ParityCheck struct {
	Date      string `json:"date"`
	Duration  int    `json:"duration"`
	Speed     string `json:"speed"`
	Status    string `json:"status"`
	Errors    int    `json:"errors"`
	Progress  int    `json:"progress"`
	Correcting bool  `json:"correcting"`
	Paused    bool   `json:"paused"`
	Running   bool   `json:"running"`
}

// GetParityCheckStatus retrieves current parity check status
func (c *Client) GetParityCheckStatus(ctx context.Context) (*ParityCheck, error) {
	query := `
		query {
			array {
				parityCheckStatus {
					date
					duration
					speed
					status
					errors
					progress
					correcting
					paused
					running
				}
			}
		}
	`

	var response struct {
		Array struct {
			ParityCheckStatus ParityCheck `json:"parityCheckStatus"`
		} `json:"array"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return &response.Array.ParityCheckStatus, nil
}

// GetParityHistory retrieves parity check history
func (c *Client) GetParityHistory(ctx context.Context) ([]ParityCheck, error) {
	query := `
		query {
			parityHistory {
				date
				duration
				speed
				status
				errors
			}
		}
	`

	var response struct {
		ParityHistory []ParityCheck `json:"parityHistory"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return response.ParityHistory, nil
}

// StartParityCheck starts a parity check
func (c *Client) StartParityCheck(ctx context.Context, correct bool) error {
	mutation := `
		mutation($correct: Boolean!) {
			parityCheck {
				start(correct: $correct)
			}
		}
	`

	variables := map[string]interface{}{
		"correct": correct,
	}

	var response map[string]interface{}

	if err := c.Mutate(ctx, mutation, variables, &response); err != nil {
		return err
	}

	return nil
}

// PauseParityCheck pauses a running parity check
func (c *Client) PauseParityCheck(ctx context.Context) error {
	mutation := `
		mutation {
			parityCheck {
				pause
			}
		}
	`

	var response map[string]interface{}

	if err := c.Mutate(ctx, mutation, nil, &response); err != nil {
		return err
	}

	return nil
}

// ResumeParityCheck resumes a paused parity check
func (c *Client) ResumeParityCheck(ctx context.Context) error {
	mutation := `
		mutation {
			parityCheck {
				resume
			}
		}
	`

	var response map[string]interface{}

	if err := c.Mutate(ctx, mutation, nil, &response); err != nil {
		return err
	}

	return nil
}

// CancelParityCheck cancels a running parity check
func (c *Client) CancelParityCheck(ctx context.Context) error {
	mutation := `
		mutation {
			parityCheck {
				cancel
			}
		}
	`

	var response map[string]interface{}

	if err := c.Mutate(ctx, mutation, nil, &response); err != nil {
		return err
	}

	return nil
}

// Notification represents a system notification
type Notification struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Importance  string `json:"importance"`
	Link        string `json:"link"`
	Type        string `json:"type"`
	Timestamp   string `json:"timestamp"`
}

// NotificationCounts contains notification counts by importance
type NotificationCounts struct {
	Info    int `json:"info"`
	Warning int `json:"warning"`
	Alert   int `json:"alert"`
	Total   int `json:"total"`
}

// NotificationOverview contains overview of all notifications
type NotificationOverview struct {
	Unread  NotificationCounts `json:"unread"`
	Archive NotificationCounts `json:"archive"`
}

// GetNotifications retrieves notifications
func (c *Client) GetNotifications(ctx context.Context, notifType string, importance string, offset int, limit int) ([]Notification, error) {
	query := `
		query($type: NotificationType!, $importance: NotificationImportance, $offset: Int!, $limit: Int!) {
			notifications {
				list(filter: {type: $type, importance: $importance, offset: $offset, limit: $limit}) {
					id
					title
					subject
					description
					importance
					link
					type
					timestamp
				}
			}
		}
	`

	variables := map[string]interface{}{
		"type":   notifType,
		"offset": offset,
		"limit":  limit,
	}

	if importance != "" {
		variables["importance"] = importance
	}

	var response struct {
		Notifications struct {
			List []Notification `json:"list"`
		} `json:"notifications"`
	}

	if err := c.Query(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return response.Notifications.List, nil
}

// GetNotificationOverview retrieves notification overview
func (c *Client) GetNotificationOverview(ctx context.Context) (*NotificationOverview, error) {
	query := `
		query {
			notifications {
				overview {
					unread {
						info
						warning
						alert
						total
					}
					archive {
						info
						warning
						alert
						total
					}
				}
			}
		}
	`

	var response struct {
		Notifications struct {
			Overview NotificationOverview `json:"overview"`
		} `json:"notifications"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return &response.Notifications.Overview, nil
}

// ArchiveNotification archives a notification
func (c *Client) ArchiveNotification(ctx context.Context, id string) error {
	mutation := `
		mutation($id: PrefixedID!) {
			archiveNotification(id: $id) {
				id
			}
		}
	`

	variables := map[string]interface{}{
		"id": id,
	}

	var response map[string]interface{}

	if err := c.Mutate(ctx, mutation, variables, &response); err != nil {
		return err
	}

	return nil
}

// LogFile represents a log file on the system
type LogFile struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Size       int    `json:"size"`
	ModifiedAt string `json:"modifiedAt"`
}

// LogFileContent represents the content of a log file
type LogFileContent struct {
	Path       string `json:"path"`
	Content    string `json:"content"`
	TotalLines int    `json:"totalLines"`
	StartLine  int    `json:"startLine"`
}

// GetLogFiles retrieves the list of available log files
func (c *Client) GetLogFiles(ctx context.Context) ([]LogFile, error) {
	query := `
		query {
			logFiles {
				name
				path
				size
				modifiedAt
			}
		}
	`

	var response struct {
		LogFiles []LogFile `json:"logFiles"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return response.LogFiles, nil
}

// GetLogFile retrieves the content of a specific log file
func (c *Client) GetLogFile(ctx context.Context, path string, lines int, startLine int) (*LogFileContent, error) {
	query := `
		query($path: String!, $lines: Int, $startLine: Int) {
			logFile(path: $path, lines: $lines, startLine: $startLine) {
				path
				content
				totalLines
				startLine
			}
		}
	`

	variables := map[string]interface{}{
		"path": path,
	}

	if lines > 0 {
		variables["lines"] = lines
	}

	if startLine > 0 {
		variables["startLine"] = startLine
	}

	var response struct {
		LogFile LogFileContent `json:"logFile"`
	}

	if err := c.Query(ctx, query, variables, &response); err != nil {
		return nil, err
	}

	return &response.LogFile, nil
}

// Plugin represents a plugin on the Unraid system
type Plugin struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	HasApiModule *bool  `json:"hasApiModule"`
	HasCliModule *bool  `json:"hasCliModule"`
}

// GetPlugins retrieves the list of installed plugins
func (c *Client) GetPlugins(ctx context.Context) ([]Plugin, error) {
	query := `
		query {
			plugins {
				name
				version
				hasApiModule
				hasCliModule
			}
		}
	`

	var response struct {
		Plugins []Plugin `json:"plugins"`
	}

	if err := c.Query(ctx, query, nil, &response); err != nil {
		return nil, err
	}

	return response.Plugins, nil
}

// AddPlugin adds one or more plugins
func (c *Client) AddPlugin(ctx context.Context, names []string, bundled bool, restart bool) error {
	query := `
		mutation($input: PluginManagementInput!) {
			addPlugin(input: $input)
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"names":   names,
			"bundled": bundled,
			"restart": restart,
		},
	}

	var response struct {
		AddPlugin bool `json:"addPlugin"`
	}

	if err := c.Query(ctx, query, variables, &response); err != nil {
		return err
	}

	return nil
}

// RemovePlugin removes one or more plugins
func (c *Client) RemovePlugin(ctx context.Context, names []string, bundled bool, restart bool) error {
	query := `
		mutation($input: PluginManagementInput!) {
			removePlugin(input: $input)
		}
	`

	variables := map[string]interface{}{
		"input": map[string]interface{}{
			"names":   names,
			"bundled": bundled,
			"restart": restart,
		},
	}

	var response struct {
		RemovePlugin bool `json:"removePlugin"`
	}

	if err := c.Query(ctx, query, variables, &response); err != nil {
		return err
	}

	return nil
}
