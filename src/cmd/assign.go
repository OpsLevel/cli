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

var unassignCmd = &cobra.Command{
	Use:   "unassign",
	Short: "Unassign properties from resources",
	Long:  "Unassign properties from resources",
}

func init() {
	rootCmd.AddCommand(assignCmd)
	rootCmd.AddCommand(unassignCmd)

	assignCmd.PersistentFlags().StringVarP(&dataFile, "file", "f", "-", "File to read data from. If '.' then reads from './data.yaml'. Defaults to reading from stdin.")
	viper.BindPFlags(assignCmd.Flags())
}
