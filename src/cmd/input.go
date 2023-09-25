package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

var dataFile string

func readInputConfig() {
	viper.SetConfigType("yaml")
	fmt.Printf("dataFile is %s\n", dataFile)
	switch dataFile {
	case ".":
		viper.SetConfigFile("./data.yaml")
	case "-":
		if isStdInFromTerminal() {
			log.Info().Msg("Reading input directly from command line...")
		}
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
