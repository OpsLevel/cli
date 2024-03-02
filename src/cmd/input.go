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
func readInputConfig() error {
	viper.SetConfigType("yaml")
	switch dataFile {
	case ".":
		viper.SetConfigFile("./data.yaml")
	case "-":
		b, err := isStdInFromTerminal()
		if err != nil {
			return err
		}
		if b {
			log.Info().Msg("Reading input directly from command line... Press CTRL+D to stop typing")
		}
		viper.ReadConfig(os.Stdin)
	default:
		viper.SetConfigFile(dataFile)
	}
	viper.ReadInConfig()
	return nil
}

func readResourceInput[T any]() (*T, error) {
	var err error
	var resource T
	var yamlData []byte

	switch dataFile {
	case ".":
		yamlData, err = os.ReadFile("./data.yaml")
	case "-":
		b, err := isStdInFromTerminal()
		if err != nil {
			return nil, err
		}
		if b {
			log.Info().Msg("Reading input directly from command line... Press CTRL+D to stop typing")
		}
		buf := bytes.Buffer{}
		_, err = buf.ReadFrom(os.Stdin)
		if err != nil {
			return nil, err
		}
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

func isStdInFromTerminal() (bool, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false, err
	}
	return fi.Mode()&os.ModeCharDevice != 0, nil
}
