package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dataFile string

func readInputConfig() {
	if dataFile != "" {
		if dataFile == "-" {
			viper.SetConfigType("yaml")
			viper.ReadConfig(os.Stdin)
		} else if dataFile == "." {
			viper.SetConfigFile("./data.yaml")
		} else {
			viper.SetConfigFile(dataFile)
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
