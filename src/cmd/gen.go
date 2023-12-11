package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate example yaml files for OpsLevel resources",
	Long:  "Generate example yaml files for OpsLevel resources",
}

func init() {
	rootCmd.AddCommand(genCmd)

	viper.BindPFlags(genCmd.Flags())
}
