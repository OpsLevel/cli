package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/spf13/cobra"
)

var listServiceMaturityCmd = &cobra.Command{
	Use:   "maturity",
	Short: "Lists all services maturity data (category and level)",
	Long: `Lists all services maturity data (category and level)

There are multiple output formats that are useful

	opslevel list service maturity
	opslevel list service maturity -o csv > maturity.csv
	opslevel list service maturity -o json | jq
`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClientGQL()
		categoriesConn, err := client.ListCategories(nil)
		cobra.CheckErr(err)
		categories := categoriesConn.Nodes
		data, err := client.ListServicesMaturity()
		cobra.CheckErr(err)
		headers := []string{"Name", "Overall"}
		sort.Slice(categories, func(i, j int) bool {
			return categories[i].Name < categories[j].Name
		})
		for _, category := range categories {
			headers = append(headers, category.Name)
		}

		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(data, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			w.Write(headers)
			for _, item := range data {
				w.Write(GetValues(&item, headers...))
			}
			w.Flush()
		} else {
			w := common.NewTabWriter(headers...)
			for _, item := range data {
				fmt.Fprintf(w, "%s\n", strings.Join(GetValues(&item, headers...), "\t"))
			}
			w.Flush()
		}
	},
}

func init() {
	listServiceCmd.AddCommand(listServiceMaturityCmd)
}

func GetValues(s *opslevel.ServiceMaturity, fields ...string) []string {
	var output []string
	for _, field := range fields {
		if field == "Name" {
			output = append(output, s.Name)
			continue
		}
		if field == "Overall" {
			output = append(output, s.MaturityReport.OverallLevel.Name)
			continue
		}
		for _, breakdown := range s.MaturityReport.CategoryBreakdown {
			if field == breakdown.Category.Name {
				output = append(output, breakdown.Level.Name)
				break
			}
		}
	}
	return output
}
