package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources or data from OpsLevel",
	Long:  "Get resources or data from OpsLevel",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
