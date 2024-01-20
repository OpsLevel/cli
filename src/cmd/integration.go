package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var getIntegrationCmd = &cobra.Command{
	Use:        "integration ID",
	Aliases:    common.GetAliases("Integration"),
	Short:      "Get details about a integration",
	Long:       `Get details about a integration`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		integration, err := getClientGQL().GetIntegration(opslevel.ID(key))
		cobra.CheckErr(err)
		common.PrettyPrint(integration)
	},
}

var listIntegrationCmd = &cobra.Command{
	Use:     "integration",
	Aliases: common.GetAliases("Integration"),
	Short:   "Lists integrations",
	Long:    `Lists integrations`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListIntegrations(nil)
		cobra.CheckErr(err)
		list := resp.Nodes
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "TYPE", "ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", item.Name, item.Type, item.Alias(), item.Id)
			}
			w.Flush()
		}
	},
}

func init() {
	getCmd.AddCommand(getIntegrationCmd)
	listCmd.AddCommand(listIntegrationCmd)
}
