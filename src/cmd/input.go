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
// we need to stop using viper in such a way and just rely on ReadResourceInput
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

func readInput() ([]byte, error) {
	var err error
	var input []byte

	switch dataFile {
	case ".":
		input, err = os.ReadFile("./data.yaml")
	case "-":
		if isStdInFromTerminal() {
			log.Info().Msg("Reading input directly from command line... Press CTRL+D to stop typing")
		}
		buf := bytes.Buffer{}
		_, err = buf.ReadFrom(os.Stdin)
		input = buf.Bytes()
	default:
		input, err = os.ReadFile(dataFile)
	}
	if err != nil {
		return input, err
	}
	return input, nil
}

func ReadResource[T any](input []byte) (*T, error) {
	var resource T
	if err := yaml.Unmarshal(input, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

func ReadResourceInput[T any](mockInputs ...[]byte) (*T, error) {
	if len(mockInputs) > 0 {
		return ReadResource[T](mockInputs[0])
	}

	newInput, err := readInput()
	if err != nil {
		return nil, err
	}
	return ReadResource[T](newInput)
}

func isStdInFromTerminal() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}
