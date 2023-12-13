package cmd

import (
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	return opslevel.JsonOf[T](opslevel.NewExampleOf[T]())
}

func getYaml[T any]() string {
	return opslevel.YamlOf[T](opslevel.NewExampleOf[T]())
}

func init() {
	rootCmd.AddCommand(exampleCmd)

	exampleCmd.PersistentFlags().BoolVar(&exampleIsJson, "json", false, "Output example in JSON format")
	exampleCmd.PersistentFlags().BoolVar(&exampleIsYaml, "yaml", true, "Output example in YAML format")
	exampleCmd.MarkFlagsMutuallyExclusive("json", "yaml")
	viper.BindPFlags(exampleCmd.Flags())
}
