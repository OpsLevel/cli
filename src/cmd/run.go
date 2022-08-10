package cmd

import cobra "github.com/spf13/cobra"

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Invoked to execute various OpsLevel functions",
	Long:  "Invoked to execute various OpsLevel functions",
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().String("github-token", "", "Github API Token. Overrides environment variable 'GITHUB_API_TOKEN'")
}
