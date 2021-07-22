package cmd

import (
	"fmt"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getServicesCmd = &cobra.Command{
	Use:   "services",
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
	Use:   "service",
	Short: "Delete a service",
	Long:  `Delete a service`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		err := client.DeleteServiceWithAlias(viper.GetString("id"))
		cobra.CheckErr(err)
	},
}

func init() {
	getCmd.AddCommand(getServicesCmd)
	deleteCmd.AddCommand(deleteServiceCmd)

	deleteServiceCmd.Flags().String("id", "", "the id of the service")
	viper.BindPFlags(deleteServiceCmd.Flags())
}
