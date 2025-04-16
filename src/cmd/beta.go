package cmd

import "github.com/spf13/cobra"

var betaCmd = &cobra.Command{
	Use:   "beta",
	Short: "Beta commands that are subject to removal",
	Long:  "Beta commands that are subject to removal",
}

func init() {
	rootCmd.AddCommand(betaCmd)
}
