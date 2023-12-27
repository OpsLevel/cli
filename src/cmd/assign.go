package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var assignCmd = &cobra.Command{
	Use:   "assign",
	Short: "Assign properties to resources",
	Long:  "Assign properties to resources",
}

func init() {
	rootCmd.AddCommand(assignCmd)

	assignCmd.PersistentFlags().StringVarP(&dataFile, "file", "f", "-", "File to read data from. If '.' then reads from './data.yaml'. Defaults to reading from stdin.")
	viper.BindPFlags(assignCmd.Flags())
}
