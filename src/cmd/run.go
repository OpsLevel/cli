package cmd

import (
	cobra "github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Invoked to execute various OpsLevel functions",
	Long:  "Invoked to execute various OpsLevel functions",
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().String("github-token", "", "Github API Token to get repo data. Overrides environment variable 'GITHUB_API_TOKEN'")

	viper.BindPFlags(runCmd.PersistentFlags())
	viper.BindEnv("github-token", "GITHUB_API_TOKEN")
}
