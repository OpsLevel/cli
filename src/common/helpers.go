package common

import (
	b64 "encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
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
