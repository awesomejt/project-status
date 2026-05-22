package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/awesomejt/project-status/cli/internal/client"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new status record",
	Long:  `Add a new status record via the API.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
	apiURL := GetAPIURL()
	output := GetOutputFormat()
	apiClient := client.NewClient(apiURL)

	projectName, _ := cmd.Flags().GetString("project-name")
		shortName, _ := cmd.Flags().GetString("short-name")
		status, _ := cmd.Flags().GetString("status")
		phase, _ := cmd.Flags().GetString("phase")
		summary, _ := cmd.Flags().GetString("summary")
		reason, _ := cmd.Flags().GetString("reason")
		details, _ := cmd.Flags().GetString("details")
		tagsStr, _ := cmd.Flags().GetString("tags")

		if projectName == "" {
			projectName, _ = cmd.Flags().GetString("name")
		}

		if projectName == "" || shortName == "" || status == "" || summary == "" {
			fmt.Fprint(os.Stderr, "Error: --project-name, --short-name, --status, and --summary are required\n")
			cmd.Help()
			os.Exit(1)
		}

		var tags []string
		if tagsStr != "" {
			for _, tag := range splitTags(tagsStr) {
				if tag != "" {
					tags = append(tags, tag)
				}
			}
		}

	record := client.StatusRecordCreate{
		ProjectName: projectName,
		ShortName:   shortName,
		Status:      status,
		Phase:       &phase,
		Summary:     summary,
		Reason:      &reason,
		Details:     &details,
		Tags:        tags,
		Source:      stringPtr("cli"),
	}

	created, err := apiClient.CreateRecord(record)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating record: %v\n", err)
			os.Exit(1)
		}

		if output == "json" {
			b, _ := json.MarshalIndent(created, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Printf("Created status record #%d\n", created.ID)
			fmt.Printf("  Project: %s (%s)\n", created.ProjectName, created.ShortName)
			fmt.Printf("  Status: %s\n", created.Status)
			if created.Phase != nil {
				fmt.Printf("  Phase: %s\n", *created.Phase)
			}
			fmt.Printf("  Summary: %s\n", created.Summary)
		}
	},
}

func init() {
	addCmd.Flags().StringP("project-name", "n", "", "project name (required)")
	addCmd.Flags().StringP("short-name", "s", "", "short project identifier (required)")
	addCmd.Flags().StringP("status", "t", "", "status: active, paused, blocked, working, error, stopped, completed (required)")
	addCmd.Flags().StringP("phase", "p", "", "workflow phase: planning, implementation, validation, release")
	addCmd.Flags().StringP("summary", "m", "", "short status summary (required)")
	addCmd.Flags().String("reason", "", "explanation for paused, blocked, error, or stopped states")
	addCmd.Flags().String("details", "", "longer notes")
	addCmd.Flags().String("tags", "", "comma-separated tags")
	rootCmd.AddCommand(addCmd)
}

func splitTags(s string) []string {
	result := []string{}
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func stringPtr(s string) *string {
	return &s
}
