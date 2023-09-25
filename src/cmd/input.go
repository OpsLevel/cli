package cmd

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
)

var dataFile string

func readInputConfig() {
	viper.SetConfigType("yaml")
	switch dataFile {
	case ".":
		// TODO: does this block ever actually ever run?
		viper.SetConfigFile("./data.yaml")
	case "-":
		if isStdInFromTerminal() {
			// TODO: this can take up to half a second to output which interrupts the user's experience
			log.Info().Msg("Reading input directly from command line... Press CTRL+D to stop typing")
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
