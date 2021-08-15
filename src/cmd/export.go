package cmd

import (
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export OpsLevel data to other tools or services",
	Long:  "Export OpsLevel data to other tools or services",
}

func init() {
	rootCmd.AddCommand(exportCmd)
}
