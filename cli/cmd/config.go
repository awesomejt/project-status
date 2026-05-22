package cmd

import (
	"fmt"

	"github.com/awesomejt/project-status/cli/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show or set configuration",
	Long:  `Show current configuration or set a config value.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Current Configuration:")
		cmd.Printf("  API URL: %s\n", viper.GetString("api_url"))
		cmd.Printf("  Output:  %s\n", viper.GetString("output"))

		if viper.ConfigFileUsed() != "" {
			cmd.Printf("  Config File: %s\n", viper.ConfigFileUsed())
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long:  `Set a configuration value. Available keys: api_url, output`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		value := args[1]

		switch key {
		case "api_url":
			if err := client.ValidateURL(value); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: invalid URL: %v\n", err)
				return
			}
			viper.Set("api_url", value)
		case "output":
			if value != "table" && value != "json" {
				fmt.Fprint(cmd.ErrOrStderr(), "Error: output must be 'table' or 'json'\n")
				return
			}
			viper.Set("output", value)
		default:
			fmt.Fprintf(cmd.ErrOrStderr(), "Error: unknown config key: %s\n", key)
			return
		}

		if err := viper.WriteConfig(); err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "Error writing config: %v\n", err)
			return
		}

		fmt.Printf("Set %s = %s\n", key, value)
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}
