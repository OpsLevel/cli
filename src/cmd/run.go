package cmd

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Invoked to execute various OpsLevel functions",
	Long:  "Invoked to execute various OpsLevel functions",
}

func init() {
	rootCmd.AddCommand(runCmd)
}
