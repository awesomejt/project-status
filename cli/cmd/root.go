package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var version = "dev"

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.project-status")
	viper.AddConfigPath("./")
	viper.SetEnvPrefix("PROJECT_STATUS")
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().String("config", "", "config file path")
	rootCmd.PersistentFlags().StringP("api-url", "u", "http://localhost:5000", "API base URL")
	rootCmd.PersistentFlags().StringP("output", "o", "table", "output format: table or json")

	viper.BindPFlag("api_url", rootCmd.PersistentFlags().Lookup("api-url"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
			os.Exit(1)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:   "status",
	Short: "Project Status CLI",
	Long:  `Project Status CLI - manage project status records via the API.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func SetTestOutput(out, err io.Writer) {
	rootCmd.SetOut(out)
	rootCmd.SetErr(err)
}

func GetAPIURL() string {
	return viper.GetString("api_url")
}

func GetOutputFormat() string {
	return viper.GetString("output")
}

func SetConfig(key string, value interface{}) error {
	viper.Set(key, value)
	return viper.WriteConfig()
}
