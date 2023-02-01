package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var createFilterCmd = &cobra.Command{
	Use:        "filter NAME",
	Short:      "Create a filter",
	Long:       `Create a filter`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"NAME"},
	Run: func(cmd *cobra.Command, args []string) {
		filter, err := getClientGQL().CreateFilter(opslevel.FilterCreateInput{
			Name: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(filter.Id)
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
	getCmd.AddCommand(getFilterCmd)
	listCmd.AddCommand(listFilterCmd)
	deleteCmd.AddCommand(deleteFilterCmd)
}
