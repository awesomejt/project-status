package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/awesomejt/project-status/cli/internal/client"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a status record",
	Long:  `Update a status record by ID. At least one field must be provided.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := GetAPIURL()
		output := GetOutputFormat()
		apiClient := client.NewClient(apiURL)

		id := args[0]

		projectName, _ := cmd.Flags().GetString("project-name")
		shortName, _ := cmd.Flags().GetString("short-name")
		status, _ := cmd.Flags().GetString("status")
		phase, _ := cmd.Flags().GetString("phase")
		summary, _ := cmd.Flags().GetString("summary")
		reason, _ := cmd.Flags().GetString("reason")
		details, _ := cmd.Flags().GetString("details")
		tagsStr, _ := cmd.Flags().GetString("tags")

		record := client.StatusRecordUpdate{}
		if cmd.Flags().Changed("project-name") {
			record.ProjectName = &projectName
		}
		if cmd.Flags().Changed("short-name") {
			record.ShortName = &shortName
		}
		if cmd.Flags().Changed("status") {
			record.Status = &status
		}
		if cmd.Flags().Changed("phase") {
			record.Phase = &phase
		}
		if cmd.Flags().Changed("summary") {
			record.Summary = &summary
		}
		if cmd.Flags().Changed("reason") {
			record.Reason = &reason
		}
		if cmd.Flags().Changed("details") {
			record.Details = &details
		}
		if cmd.Flags().Changed("tags") {
			tags := splitTags(tagsStr)
			record.Tags = &tags
		}

		updated, err := apiClient.UpdateRecord(id, record)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating record: %v\n", err)
			os.Exit(1)
		}

		if output == "json" {
			b, _ := json.MarshalIndent(updated, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Printf("Updated status record %s\n", updated.ID)
			fmt.Printf("  Project: %s (%s)\n", updated.ProjectName, updated.ShortName)
			fmt.Printf("  Status: %s\n", updated.Status)
			if updated.Phase != nil {
				fmt.Printf("  Phase: %s\n", *updated.Phase)
			}
			fmt.Printf("  Summary: %s\n", updated.Summary)
		}
	},
}

func init() {
	updateCmd.Flags().StringP("project-name", "n", "", "project name")
	updateCmd.Flags().StringP("short-name", "s", "", "short project identifier")
	updateCmd.Flags().StringP("status", "t", "", "status: active, paused, blocked, working, error, stopped, completed")
	updateCmd.Flags().StringP("phase", "p", "", "workflow phase: planning, implementation, validation, release")
	updateCmd.Flags().StringP("summary", "m", "", "short status summary")
	updateCmd.Flags().String("reason", "", "explanation for paused, blocked, error, or stopped states")
	updateCmd.Flags().String("details", "", "longer notes")
	updateCmd.Flags().String("tags", "", "comma-separated tags")
	rootCmd.AddCommand(updateCmd)
}
