package common

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func PrettyPrint(object interface{}) {
	bytes, err := json.MarshalIndent(object, "", "  ")
	cobra.CheckErr(err)
	fmt.Println(string(bytes))
}

func NewTabWriter(headers ...string) *tabwriter.Writer {
	longestHeader := 0
	var headerFormat strings.Builder
	headersCasted := make([]interface{}, len(headers))
	for i, header := range headers {
		headerFormat.WriteString("%s\t")
		headersCasted[i] = header
		headerLength := len(header)
		if headerLength > longestHeader {
			longestHeader = headerLength
		}
	}
	headerFormat.WriteString("\n")
	w := tabwriter.NewWriter(os.Stdout, longestHeader, longestHeader, 2, ' ', 0)
	fmt.Fprintf(w, headerFormat.String(), headersCasted...)
	return w
}
