package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listOutputType string

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all resources in OpsLevel",
	Long:    "List all resources in OpsLevel",
}

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().StringVarP(&listOutputType, "output", "o", "text", "Output format.  One of: json|csv|text [default: text]")
	viper.BindPFlags(listCmd.Flags())
}

func isJsonOutput() bool {
	return listOutputType == "json" || getOutputType == "json"
}

func isCsvOutput() bool {
	return listOutputType == "csv"
}
