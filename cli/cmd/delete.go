package cmd

import (
	"fmt"
	"os"

	"github.com/awesomejt/project-status/cli/internal/client"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a status record",
	Long:  `Delete a status record by ID. Use --force to skip confirmation.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		apiURL := GetAPIURL()
		apiClient := client.NewClient(apiURL)

		id := args[0]

		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Print(fmt.Sprintf("Are you sure you want to delete status record %s? [y/N] ", id))
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				fmt.Println("Cancelled.")
				return
			}
		}

		err := apiClient.DeleteRecord(id)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error deleting record: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Deleted status record %s\n", id)
	},
}

func init() {
	deleteCmd.Flags().BoolP("force", "f", false, "skip confirmation prompt")
	rootCmd.AddCommand(deleteCmd)
}
