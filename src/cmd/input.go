package cmd

import (
	"bytes"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var dataFile string

// TODO: deploy.go is still relying on this because it uses viper to bind config
// we need to stop using viper in such a way and just rely on readResourceInput
func readInputConfig() {
	viper.SetConfigType("yaml")
	switch dataFile {
	case ".":
		viper.SetConfigFile("./data.yaml")
	case "-":
		if isStdInFromTerminal() {
			log.Info().Msg("Reading input directly from command line... Press CTRL+D to stop typing")
		}
		viper.ReadConfig(os.Stdin)
	default:
		viper.SetConfigFile(dataFile)
	}
	viper.ReadInConfig()
}

func readResourceInput[T any]() (*T, error) {
	var err error
	var resource T
	var yamlData []byte

	switch dataFile {
	case ".":
		yamlData, err = os.ReadFile("./data.yaml")
	case "-":
		if isStdInFromTerminal() {
			log.Info().Msg("Reading input directly from command line... Press CTRL+D to stop typing")
		}
		buf := bytes.Buffer{}
		_, err = buf.ReadFrom(os.Stdin)
		yamlData = buf.Bytes()
	default:
		yamlData, err = os.ReadFile(dataFile)
	}
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(yamlData, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

func isStdInFromTerminal() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}
