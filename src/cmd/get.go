package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getOutputType string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get detailed info about resources in OpsLevel",
	Long:  "Get detailed info about resources in OpsLevel",
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringVarP(&getOutputType, "output", "o", "text", "Output format.  One of: yaml|text [default: text]")
	if err := viper.BindPFlags(getCmd.Flags()); err != nil {
		cobra.CheckErr(err)
	}
}

func isYamlOutput() bool {
	return getOutputType == "yaml"
}
