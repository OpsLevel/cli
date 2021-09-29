package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateDataFile string

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update resources in OpsLevel",
	Long:  "Update resources in OpsLevel",
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVarP(&updateDataFile, "file", "f", "-", "File to read update from. If '.' then reads from './data.yaml'. Defaults to reading from stdin.")
	viper.BindPFlags(createCmd.Flags())
}

func readUpdateConfigFile() {
	if updateDataFile != "" {
		if updateDataFile == "-" {
			viper.SetConfigType("yaml")
			viper.ReadConfig(os.Stdin)
			return
		} else if updateDataFile == "." {
			viper.SetConfigFile("./data.yaml")
		} else {
			viper.SetConfigFile(updateDataFile)
		}
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.SetConfigName("data")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}
	viper.ReadInConfig()
}
