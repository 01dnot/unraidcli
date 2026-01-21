package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	notifImportance string
	notifLimit      int
)

// notificationsCmd represents the notifications command
var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "View system notifications",
	Long:  "View and manage system notifications from your Unraid server.",
	Aliases: []string{"notif", "alerts"},
}

// notificationsListCmd represents the notifications ls command
var notificationsListCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List unread notifications",
	Long:    "Display unread notifications from your Unraid server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		notifications, err := apiClient.GetNotifications(ctx, "UNREAD", notifImportance, 0, notifLimit)
		if err != nil {
			return fmt.Errorf("failed to get notifications: %w", err)
		}

		if len(notifications) == 0 {
			fmt.Println("No unread notifications.")
			return nil
		}

		if outputFormat == "" || outputFormat == "table" {
			for i, notif := range notifications {
				if i > 0 {
					fmt.Println()
				}

				// Format importance with indicator
				importance := notif.Importance
				switch notif.Importance {
				case "ALERT":
					importance = "🔴 " + importance
				case "WARNING":
					importance = "🟡 " + importance
				case "INFO":
					importance = "🔵 " + importance
				}

				fmt.Printf("[%s] %s\n", importance, notif.Title)
				fmt.Printf("Subject: %s\n", notif.Subject)
				if notif.Description != "" {
					fmt.Printf("Description: %s\n", notif.Description)
				}
				if notif.Timestamp != "" {
					fmt.Printf("Time: %s\n", notif.Timestamp)
				}
				if notif.Link != "" {
					fmt.Printf("Link: %s\n", notif.Link)
				}
			}
		} else {
			formatter.Print(notifications)
		}

		return nil
	},
}

// notificationsArchiveCmd represents the notifications archive command
var notificationsArchiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "List archived notifications",
	Long:  "Display archived notifications from your Unraid server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		notifications, err := apiClient.GetNotifications(ctx, "ARCHIVE", notifImportance, 0, notifLimit)
		if err != nil {
			return fmt.Errorf("failed to get notifications: %w", err)
		}

		if len(notifications) == 0 {
			fmt.Println("No archived notifications.")
			return nil
		}

		if outputFormat == "" || outputFormat == "table" {
			headers := []string{"Importance", "Title", "Subject", "Time"}
			var rows [][]string

			for _, notif := range notifications {
				rows = append(rows, []string{
					notif.Importance,
					notif.Title,
					notif.Subject,
					notif.Timestamp,
				})
			}

			formatter.PrintTable(headers, rows)
		} else {
			formatter.Print(notifications)
		}

		return nil
	},
}

// notificationsOverviewCmd represents the notifications overview command
var notificationsOverviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "Show notification overview",
	Long:  "Display a summary of all notifications by type and importance.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		overview, err := apiClient.GetNotificationOverview(ctx)
		if err != nil {
			return fmt.Errorf("failed to get notification overview: %w", err)
		}

		if outputFormat == "" || outputFormat == "table" {
			fmt.Println("Unread Notifications:")
			fmt.Printf("  Alerts: %d\n", overview.Unread.Alert)
			fmt.Printf("  Warnings: %d\n", overview.Unread.Warning)
			fmt.Printf("  Info: %d\n", overview.Unread.Info)
			fmt.Printf("  Total: %d\n\n", overview.Unread.Total)

			fmt.Println("Archived Notifications:")
			fmt.Printf("  Alerts: %d\n", overview.Archive.Alert)
			fmt.Printf("  Warnings: %d\n", overview.Archive.Warning)
			fmt.Printf("  Info: %d\n", overview.Archive.Info)
			fmt.Printf("  Total: %d\n", overview.Archive.Total)
		} else {
			formatter.Print(overview)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(notificationsCmd)
	notificationsCmd.AddCommand(notificationsListCmd)
	notificationsCmd.AddCommand(notificationsArchiveCmd)
	notificationsCmd.AddCommand(notificationsOverviewCmd)

	// Add flags
	notificationsListCmd.Flags().StringVar(&notifImportance, "importance", "", "Filter by importance: ALERT, WARNING, INFO")
	notificationsListCmd.Flags().IntVar(&notifLimit, "limit", 20, "Maximum number of notifications to display")

	notificationsArchiveCmd.Flags().StringVar(&notifImportance, "importance", "", "Filter by importance: ALERT, WARNING, INFO")
	notificationsArchiveCmd.Flags().IntVar(&notifLimit, "limit", 20, "Maximum number of notifications to display")
}
