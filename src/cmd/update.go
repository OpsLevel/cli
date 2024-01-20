package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"edit"},
	Short:   "Update resources or events from a file or stdin",
	Long:    "Update resources or events from a file or stdin",
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVarP(&dataFile, "file", "f", "-", "File to read update from. If '.' then reads from './data.yaml'. Defaults to reading from stdin.")
	viper.BindPFlags(updateCmd.Flags())
}
