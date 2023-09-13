package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var dataFile string

func readInputConfig() {
	switch dataFile {
	case ".":
		viper.SetConfigFile("./data.yaml")
	case "-":
		if isStdInFromTerminal() {
			fmt.Println("Reading input directly from command line...")
		}
		viper.SetConfigType("yaml")
		viper.ReadConfig(os.Stdin)
	default:
		viper.SetConfigFile(dataFile)
	}
	viper.ReadInConfig()
}

func isStdInFromTerminal() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}
