package cmd

import (
	"encoding/json"

	"github.com/opslevel/opslevel-go/v2025"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var exampleIsJson, exampleIsYaml bool

var exampleCmd = &cobra.Command{
	Use:   "example",
	Short: "Examples of OpsLevel resources",
	Long:  "Examples of OpsLevel resources in different formats",
}

func getExample[T any]() string {
	if exampleIsJson {
		return getJson[T]()
	}
	return getYaml[T]()
}

func getJson[T any]() string {
	var (
		out []byte
		err error
	)
	t := opslevel.NewExampleOf[T]()
	out, err = json.Marshal(t)
	if err != nil {
		panic("unexpected error getting example json")
	}
	return string(out)
}

func getYaml[T any]() string {
	var (
		out []byte
		err error
	)
	t := opslevel.NewExampleOf[T]()
	out, err = yaml.Marshal(t)
	if err != nil {
		panic("unexpected error getting example yaml")
	}
	return string(out)
}

func init() {
	rootCmd.AddCommand(exampleCmd)

	exampleCmd.PersistentFlags().BoolVar(&exampleIsJson, "json", false, "Output example in JSON format")
	exampleCmd.PersistentFlags().BoolVar(&exampleIsYaml, "yaml", true, "Output example in YAML format")
	exampleCmd.MarkFlagsMutuallyExclusive("json", "yaml")
	viper.BindPFlags(exampleCmd.Flags())
}
