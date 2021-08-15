package cmd

import (
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data to OpsLevel.",
	Long:  "Import data to OpsLevel.",
}

func init() {
	rootCmd.AddCommand(importCmd)
}
