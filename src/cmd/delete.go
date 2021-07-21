package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources or data from OpsLevel",
	Long:  "Delete resources or data from OpsLevel",
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
