package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2025"
	"github.com/spf13/cobra"
)

var getServiceMaturityCmd = &cobra.Command{
	Use:   "maturity ALIAS",
	Short: "Get service maturity data (category and level)",
	Long: `Get service maturity data (category and level)

There are multiple output formats that are useful

	opslevel get service maturity ALIAS
	opslevel get service maturity ALIAS -o yaml | yq
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		client := getClientGQL()

		data, err := client.GetServiceMaturityWithAlias(alias)
		cobra.CheckErr(err)
		if isYamlOutput() {
			common.YamlPrint(*data)
		} else {
			writeOutput(client, []opslevel.ServiceMaturity{*data})
		}
	},
}

var listServiceMaturityCmd = &cobra.Command{
	Use:     "maturity",
	Aliases: []string{"maturities"},
	Short:   "Lists all services maturity data (category and level)",
	Long: `Lists all services maturity data (category and level)

There are multiple output formats that are useful

	opslevel list service maturity
	opslevel list service maturity -o csv > maturity.csv
	opslevel list service maturity -o json | jq
`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClientGQL()
		response, err := client.ListServicesMaturity(nil)
		cobra.CheckErr(err)

		writeOutput(client, response.Nodes)
	},
}

func init() {
	getServiceCmd.AddCommand(getServiceMaturityCmd)
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

func writeOutput(client *opslevel.Client, data []opslevel.ServiceMaturity) {
	headers := getCategoryHeaders(client)

	if isJsonOutput() {
		if len(data) == 1 {
			common.JsonPrint(json.MarshalIndent(data[0], "", "    "))
		} else {
			common.JsonPrint(json.MarshalIndent(data, "", "    "))
		}
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
}

func getCategoryHeaders(client *opslevel.Client) []string {
	headers := []string{"Name", "Overall"}
	categoriesConn, err := client.ListCategories(nil)
	cobra.CheckErr(err)
	categories := categoriesConn.Nodes
	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Name < categories[j].Name
	})
	for _, category := range categories {
		headers = append(headers, category.Name)
	}

	return headers
}
