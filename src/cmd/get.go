package cmd

import (
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get detailed info about resources in OpsLevel",
	Long:  "Get detailed info about resources in OpsLevel",
}

func init() {
	rootCmd.AddCommand(getCmd)
}
