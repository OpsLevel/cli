package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources or events from a file or stdin",
	Long:  "Create resources or events from a file or stdin",
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&dataFile, "file", "f", "", "File to read data from. If '.' then reads from './data.yaml'. Defaults to reading from stdin.")
	viper.BindPFlags(createCmd.Flags())
}
