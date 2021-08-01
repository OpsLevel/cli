package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createDataFile string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create resources or events from a file or stdin",
	Long:  "Create resources or events from a file or stdin",
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&createDataFile, "file", "f", "./data.yaml", "")
	viper.BindPFlags(createCmd.Flags())
}

func readCreateConfigFile() {
	if createDataFile != "" {
		if createDataFile == "-" {
			viper.SetConfigType("yaml")
			viper.ReadConfig(os.Stdin)
			return
		} else if createDataFile == "." {
			viper.SetConfigFile("./data.yaml")
		} else {
			viper.SetConfigFile(createDataFile)
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
