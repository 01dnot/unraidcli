package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/01dnot/unraidcli/internal/client"
	"github.com/01dnot/unraidcli/internal/output"
	"github.com/spf13/cobra"
)

// sharesCmd represents the shares command
var sharesCmd = &cobra.Command{
	Use:   "shares",
	Short: "Manage user shares",
	Long:  "List and view information about user shares on your Unraid server.",
}

// sharesLsCmd represents the shares ls command
var sharesLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all user shares",
	Long:  "List all user shares with size and usage information.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shares, err := apiClient.GetShares(ctx)
		if err != nil {
			return fmt.Errorf("failed to get shares: %w", err)
		}

		if len(shares) == 0 {
			fmt.Println("No shares found.")
			return nil
		}

		if outputFormat == "" || outputFormat == "table" {
			headers := []string{"Name", "Total", "Used", "Free", "% Used", "Cache", "Comment"}
			var rows [][]string

			for _, share := range shares {
				usedPercent := "0%"
				if share.Size > 0 {
					usedPercent = fmt.Sprintf("%.1f%%", float64(share.Used)/float64(share.Size)*100)
				}

				cache := ""
				if share.Cache {
					cache = "✓"
				}

				comment := share.Comment
				if len(comment) > 30 {
					comment = comment[:27] + "..."
				}

				rows = append(rows, []string{
					share.Name,
					output.FormatBytes(share.Size),
					output.FormatBytes(share.Used),
					output.FormatBytes(share.Free),
					usedPercent,
					cache,
					comment,
				})
			}

			formatter.PrintTable(headers, rows)
		} else {
			formatter.Print(shares)
		}

		return nil
	},
}

// sharesInfoCmd represents the shares info command
var sharesInfoCmd = &cobra.Command{
	Use:   "info <share-name>",
	Short: "Show detailed share information",
	Long:  "Display detailed information about a specific user share.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		shareName := args[0]

		shares, err := apiClient.GetShares(ctx)
		if err != nil {
			return fmt.Errorf("failed to get shares: %w", err)
		}

		// Find the share
		var found *client.Share

		for i, share := range shares {
			if share.Name == shareName {
				found = &shares[i]
				break
			}
		}

		if found == nil {
			return fmt.Errorf("share '%s' not found", shareName)
		}

		if outputFormat == "" || outputFormat == "table" {
			usedPercent := float64(0)
			if found.Size > 0 {
				usedPercent = float64(found.Used) / float64(found.Size) * 100
			}

			data := map[string]interface{}{
				"Name":          found.Name,
				"Total Size":    output.FormatBytes(found.Size),
				"Used":          fmt.Sprintf("%s (%.1f%%)", output.FormatBytes(found.Used), usedPercent),
				"Free":          output.FormatBytes(found.Free),
				"Cache":         found.Cache,
				"Comment":       found.Comment,
				"Include Disks": fmt.Sprintf("%v", found.Include),
				"Exclude Disks": fmt.Sprintf("%v", found.Exclude),
			}
			formatter.PrintKeyValue(data)
		} else {
			formatter.Print(found)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(sharesCmd)
	sharesCmd.AddCommand(sharesLsCmd)
	sharesCmd.AddCommand(sharesInfoCmd)
}
