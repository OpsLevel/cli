package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dataFile string

func hasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return stat.Size() > 0
}

func readInputFile() ([]byte, error) {
	if hasStdin() {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		file, err := os.Open(dataFile)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

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
