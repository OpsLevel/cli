package cmd

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources or data from OpsLevel",
	Long:  "List resources or data from OpsLevel",
}

func init() {
	rootCmd.AddCommand(listCmd)
}
