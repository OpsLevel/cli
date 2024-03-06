package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/opslevel/opslevel-go/v2024"

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

func ReadResourceHandleJSONSchema[T any](input []byte) (*T, error) {
	var err error
	if input == nil {
		input, err = readInput()
		if err != nil {
			return nil, fmt.Errorf("error reading from input: %w", err)
		}
	}

	m, err := ReadResource[map[string]any](input)
	if err != nil {
		return nil, fmt.Errorf("error creating map from input: %w", err)
	}
	toMap := *m

	if _, ok := toMap["schema"]; !ok {
		return nil, errors.New("required field 'schema' not found")
	}

	// if this is a string, the user should have provided JSON so parse the key value pairs
	if schemaString, ok := toMap["schema"].(string); ok {
		jsonSchema, err := opslevel.NewJSONSchema(schemaString)
		if err != nil {
			return nil, fmt.Errorf("error creating JSONSchema from field 'schema': %w", err)
		}
		toMap["schema"] = jsonSchema
	}

	var finalInput T
	err = mapstructure.Decode(toMap, &finalInput)
	if err != nil {
		return nil, fmt.Errorf("error decoding map as type %T: %w", finalInput, err)
	}
	return &finalInput, nil
}

func ReadResourceInput[T any](input []byte) (*T, error) {
	var err error
	if input == nil {
		input, err = readInput()
		if err != nil {
			return nil, err
		}
	}
	return ReadResource[T](input)
}

func isStdInFromTerminal() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}
