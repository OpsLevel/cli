package cmd

import (
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "Delete or remove data in OpsLevel",
	Long:    "Delete or remove data in OpsLevel",
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
