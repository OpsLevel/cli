package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createFilterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Create a filter",
	Long: `Create a filter

cat << EOF | opslevel create filter -f -
name: "Tier 1 Rails apps using RDS"
connective: "and"
predicates:
  - key: "tier_index"
    type: "equals"
    value: "1"
  - key: "filter_id"
    type: "matches"
    value: "Z2lkOi8vb3BzbGV2ZWwvRmlsdGVyLzEyNjQ"
  - key: "tags"
    keyData: "db"
    type: "equals"
    value: "rds"
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readFilterCreateInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateFilter(*input)
		cobra.CheckErr(err)
		fmt.Println(result.Id)
	},
}

var getFilterCmd = &cobra.Command{
	Use:        "filter ID",
	Short:      "Get details about a filter",
	Long:       `Get details about a filter`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		filter, err := getClientGQL().GetFilter(opslevel.ID(key))
		cobra.CheckErr(err)
		common.PrettyPrint(filter)
	},
}

var listFilterCmd = &cobra.Command{
	Use:     "filter",
	Aliases: []string{"filters"},
	Short:   "Lists filters",
	Long:    `Lists filters`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListFilters(nil)
		cobra.CheckErr(err)
		list := resp.Nodes
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Alias(), item.Id)
			}
			w.Flush()
		}
	},
}

var updateFilterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Update a filter",
	Long: `Update a filter

cat << EOF | opslevel update filter -f -
id: "Z2lkOi8vb3BzbGV2ZWwvRmlsdGVyLzIzNTk"
name: "Apps using RDS or DynamoDB"
connective: "or"
predicates:
  - key: "tags"
    keyData: "db"
    type: "equals"
    value: "rds"
  - key: "tags"
    keyData: "db"
    type: "equals"
    value: "dynamo"
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readFilterUpdateInput()
		cobra.CheckErr(err)
		filter, err := getClientGQL().UpdateFilter(*input)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(filter, "", "    "))
	},
}

var deleteFilterCmd = &cobra.Command{
	Use:        "filter ID",
	Short:      "Delete a filter",
	Long:       `Delete a filter`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteFilter(opslevel.ID(key))
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' filter\n", key)
	},
}

func init() {
	createCmd.AddCommand(createFilterCmd)
	updateCmd.AddCommand(updateFilterCmd)
	getCmd.AddCommand(getFilterCmd)
	listCmd.AddCommand(listFilterCmd)
	deleteCmd.AddCommand(deleteFilterCmd)
}

func readFilterCreateInput() (*opslevel.FilterCreateInput, error) {
	readCreateConfigFile()
	evt := &opslevel.FilterCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}

func readFilterUpdateInput() (*opslevel.FilterUpdateInput, error) {
	readUpdateConfigFile()
	evt := &opslevel.FilterUpdateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
