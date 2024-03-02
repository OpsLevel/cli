package common

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func GetArg(args []string, index int, defaultValue string) string {
	if len(args) > index {
		return args[index]
	}
	return defaultValue
}

func WasFound(condition bool, key string) {
	if condition {
		cobra.CheckErr(fmt.Errorf("not found - '%s'", key))
	}
}

func JsonPrint(jsonBytes []byte, err error) {
	cobra.CheckErr(err)
	fmt.Printf("%s\n", string(jsonBytes))
}

func YamlPrint(object interface{}) {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(4) // this is what you're looking for
	err := yamlEncoder.Encode(&object)
	cobra.CheckErr(err)
	fmt.Printf("---\n%s\n", b.String())
}

func MinInt(values ...int) int {
	if len(values) < 1 {
		panic("MinInt: unexpected received no values")
	}

	minValue := values[0]

	for _, val := range values {
		if val < minValue {
			minValue = val
		}
	}

	return minValue
}
