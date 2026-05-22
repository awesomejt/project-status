package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/awesomejt/project-status/cli/internal/client"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a status record",
	Long:  `Show a status record by ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := GetAPIURL()
		output := GetOutputFormat()
		apiClient := client.NewClient(apiURL)

		id := args[0]

		record, err := apiClient.GetRecord(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if output == "json" {
			b, _ := json.MarshalIndent(record, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Printf("Status Record %s\n", record.ID)
			fmt.Printf("  Project: %s\n", record.ProjectName)
			fmt.Printf("  Short Name: %s\n", record.ShortName)
			fmt.Printf("  Status: %s\n", record.Status)
			if record.Phase != nil {
				fmt.Printf("  Phase: %s\n", *record.Phase)
			}
			fmt.Printf("  Summary: %s\n", record.Summary)
			if record.Reason != nil {
				fmt.Printf("  Reason: %s\n", *record.Reason)
			}
			if record.Details != nil {
				fmt.Printf("  Details: %s\n", *record.Details)
			}
			if len(record.Tags) > 0 {
				tags := ""
				for i, tag := range record.Tags {
					if i > 0 {
						tags += ", "
					}
					tags += tag
				}
				fmt.Printf("  Tags: %s\n", tags)
			}
			if record.Source != nil {
				fmt.Printf("  Source: %s\n", *record.Source)
			}
			fmt.Printf("  Created: %s\n", record.CreatedAt)
			fmt.Printf("  Updated: %s\n", record.UpdatedAt)
		}
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
