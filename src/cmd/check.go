package cmd

import (
	"fmt"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var getCheckCmd = &cobra.Command{
	Use:        "check [id]",
	Short:      "Get details about a rubic check given its ID",
	Long:       `Get details about a rubic check given its ID`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"id"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		category, err := client.GetCheck(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(category)
	},
}

var listCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Lists the rubric checks",
	Long:  `Lists the rubric checks`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		list, err := client.ListChecks()
		cobra.CheckErr(err)
		w := common.NewTabWriter("NAME", "ID")
		if err == nil {
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Name, item.Id)
			}
		}
		w.Flush()
	},
}

func init() {
	getCmd.AddCommand(getCheckCmd)
	listCmd.AddCommand(listCheckCmd)
}
