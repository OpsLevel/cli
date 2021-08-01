package cmd

import (
	"fmt"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var getServiceCmd = &cobra.Command{
	Use:        "service [alias]",
	Short:      "Get details about a service given one of its Aliases",
	Long:       `Get details about a service given one of its Aliases`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"alias"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		service, err := client.GetServiceWithAlias(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(service)
	},
}

var listServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "Lists the services",
	Long:  `Lists the services`,
	Run: func(cmd *cobra.Command, args []string) {

		client := common.NewGraphClient()
		list, err := client.ListServices()
		cobra.CheckErr(err)
		w := common.NewTabWriter("NAME", "ID", "ALIASES")
		if err == nil {
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Id, strings.Join(item.Aliases, ","))
			}
		}
		w.Flush()
	},
}

var deleteServiceCmd = &cobra.Command{
	Use:        "service [id]",
	Short:      "Delete a service",
	Long:       `Delete a service`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"id"},
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		err := client.DeleteServiceWithAlias(args[0])
		cobra.CheckErr(err)
	},
}

func init() {
	getCmd.AddCommand(getServiceCmd)
	listCmd.AddCommand(listServiceCmd)
	deleteCmd.AddCommand(deleteServiceCmd)
}
