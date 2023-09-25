package cmd

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

var dataFile string

func readInputConfig() {
	fmt.Printf("dataFile is %s\n", dataFile)
	viper.SetConfigType("yaml")

	if dataFile == "." {
		// TODO: validate file exists, if doesn't exist, check data.yml
		fmt.Printf("using data.yaml\n")
		viper.SetConfigFile("./data.yaml")
	} else if dataFile == "-" || dataFile == "" {
		fmt.Printf("using stdin\n")
		if isStdInFromTerminal() {
			log.Info().Msg("Reading input directly from command line...")
		}
		viper.ReadConfig(os.Stdin)
	} else {
		// TODO: validate file exists
		fmt.Printf("using %s\n", dataFile)
		viper.SetConfigFile(dataFile)
	}

	viper.ReadInConfig()
}

func isStdInFromTerminal() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}
