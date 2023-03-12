package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"

	"github.com/spf13/cobra"
)

func NewLifecycleCmd(client *opslevel.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lifecycle",
		Aliases: []string{"lifecycles"},
		Short:   "Lists lifecycles",
		Long:    `Lists lifecycles`,
		Run: func(cmd *cobra.Command, args []string) {
			list, err := client.ListLifecycles()
			cobra.CheckErr(err)
			if isJsonOutput() {
				data, err := json.MarshalIndent(list, "", "  ")
				cobra.CheckErr(err)
				cmd.Print(string(data))
			} else {
				w := common.NewTabWriter("ALIAS", "ID")
				for _, item := range list {
					fmt.Fprintf(w, "%s\t%s\t\n", item.Alias, item.Id)
				}
				w.Flush()
			}
		},
	}
	cmd.Flags().StringVarP(&listOutputType, "output", "o", "text", "Output format.  One of: json|csv|text [default: text]")
	return cmd
}

var tierCmd = &cobra.Command{
	Use:     "tier",
	Aliases: []string{"tiers"},
	Short:   "Lists tiers",
	Long:    `Lists tiers`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := getClientGQL().ListTiers()
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Alias, item.Id)
			}
			w.Flush()
		}
	},
}

var toolsCmd = &cobra.Command{
	Use:     "tool",
	Aliases: []string{"tools"},
	Short:   "Lists the valid alias for tools",
	Long:    `Lists the valid alias for tools`,
	Run: func(cmd *cobra.Command, args []string) {
		list := opslevel.AllToolCategory()
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t\n", item)
			}
			w.Flush()
		}
	},
}

func init() {
	//client := getClientGQL()
	//listCmd.AddCommand(NewLifecycleCmd(client))
	listCmd.AddCommand(tierCmd)
	listCmd.AddCommand(toolsCmd)
}
