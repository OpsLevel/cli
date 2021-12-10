package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/gosimple/slug"
	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var getIntegrationCmd = &cobra.Command{
	Use:        "integration ID",
	Short:      "Get details about a integration",
	Long:       `Get details about a integration`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		integration, err := graphqlClient.GetIntegration(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(integration)
	},
}

var listIntegrationCmd = &cobra.Command{
	Use:     "integration",
	Aliases: []string{"integrations"},
	Short:   "Lists integrations",
	Long:    `Lists integrations`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := graphqlClient.ListIntegrations()
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "TYPE", "ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", item.Name, item.Type, fmt.Sprintf("%s-%s", slug.Make(item.Type), slug.Make(item.Name)), item.Id)
			}
			w.Flush()
		}
	},
}

func init() {
	getCmd.AddCommand(getIntegrationCmd)
	listCmd.AddCommand(listIntegrationCmd)
}
