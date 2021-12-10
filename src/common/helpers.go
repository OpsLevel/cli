package common

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"strings"

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
		cobra.CheckErr(fmt.Errorf("Not found - '%s'", key))
	}
}

func IsID(value string) bool {
	decoded, err := b64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return false
	}
	return strings.HasPrefix(string(decoded), "gid://")
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
	fmt.Printf("%s\n", string(b.Bytes()))
}
