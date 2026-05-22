package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/awesomejt/project-status/cli/internal/client"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List status records",
	Long:  `List status records with optional filters.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := GetAPIURL()
		output := GetOutputFormat()
		apiClient := client.NewClient(apiURL)

		page, _ := cmd.Flags().GetInt("page")
		perPage, _ := cmd.Flags().GetInt("per-page")
		status, _ := cmd.Flags().GetString("status")
		phase, _ := cmd.Flags().GetString("phase")

		response, err := apiClient.ListRecords(page, perPage, status, phase)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing records: %v\n", err)
			os.Exit(1)
		}

		if output == "json" {
			b, _ := json.MarshalIndent(response, "", "  ")
			fmt.Println(string(b))
		} else {
			fmt.Printf("%-36s %-20s %-12s %-12s %s\n", "ID", "PROJECT", "STATUS", "PHASE", "SUMMARY")
			fmt.Println("--------------------------------------------------------------------------------------------------------")
			for _, record := range response.Records {
				phaseStr := ""
				if record.Phase != nil {
					phaseStr = *record.Phase
				}
				fmt.Printf("%-36s %-20s %-12s %-12s %s\n", record.ID, record.ProjectName, record.Status, phaseStr, record.Summary)
			}
			fmt.Printf("\nTotal: %d records (page %d of %d)\n", response.Total, response.Page, response.Pages)
		}
	},
}

func init() {
	listCmd.Flags().IntP("page", "p", 1, "page number")
	listCmd.Flags().IntP("per-page", "n", 20, "records per page")
	listCmd.Flags().StringP("status", "s", "", "filter by status")
	listCmd.Flags().StringP("phase", "P", "", "filter by phase")
	rootCmd.AddCommand(listCmd)
}
