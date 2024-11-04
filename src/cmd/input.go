package cmd

import (
	"bytes"
	"errors"
	"os"
	"reflect"

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
	yamlData, err := getYamlData()
	if err != nil {
		return nil, err
	}
	return yamlUnmarshalInto[T](yamlData)
}

func yamlUnmarshalInto[T any](yamlData []byte) (*T, error) {
	var resource T
	if err := yaml.Unmarshal(yamlData, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

// for yaml unmarshaling into a generic struct field - a struct tag name workaound
func yamlUnmarshalIntoStructField[T any](yamlData []byte, structField reflect.StructField) (*T, error) {
	withExtraFields := reflect.StructOf([]reflect.StructField{structField})
	v := reflect.New(withExtraFields).Elem()
	s := v.Addr().Interface()

	r := bytes.NewReader(yamlData)
	if err := yaml.NewDecoder(r).Decode(s); err != nil {
		return nil, err
	}

	thisOne := v.FieldByName(structField.Name).Addr().Interface()
	thing, ok := thisOne.(*T)
	if !ok {
		return nil, errors.New("could not get extra data yaml data")
	}
	return thing, nil
}

func getYamlData() ([]byte, error) {
	var yamlData []byte
	var err error
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
	return yamlData, err
}

func isStdInFromTerminal() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice != 0
}
