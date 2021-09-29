package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var getCheckCmd = &cobra.Command{
	Use:        "check ID",
	Short:      "Get details about a rubic check",
	Long:       `Get details about a rubic check`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		category, err := client.GetCheck(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(category)
	},
}

var listCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"checks"},
	Short:   "Lists the rubric checks",
	Long:    `Lists the rubric checks`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := graphqlClient.ListChecks()
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Name, item.Id)
			}
			w.Flush()
		}
	},
}

var deleteCheckCmd = &cobra.Command{
	Use:        "check ID",
	Short:      "Delete a rubric check",
	Long:       `Delete a rubric check`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := graphqlClient.DeleteCheck(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' check\n", key)
	},
}

func init() {
	getCmd.AddCommand(getCheckCmd)
	listCmd.AddCommand(listCheckCmd)
	deleteCmd.AddCommand(deleteCheckCmd)
}
