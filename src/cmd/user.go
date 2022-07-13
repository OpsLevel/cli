package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2022"
	"github.com/spf13/cobra"
	"sort"
)

var createUserCmd = &cobra.Command{
	Use:   "user EMAIL NAME",
	Short: "Create a User",
	Long:  `Create a User`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		name := args[1]

		resource, err := getClientGQL().InviteUser(email, opslevel.UserInput{
			Name: name,
		})
		cobra.CheckErr(err)
		fmt.Println(resource.Id)
	},
}

var listUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Lists the users",
	Long:    `Lists the users`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := getClientGQL().ListUsers()
		sort.Slice(list, func(i, j int) bool {
			return list[i].Email < list[j].Email
		})
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "EMAIL", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Email, item.Id)
			}
			w.Flush()
		}
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user ID",
	Short: "Delete a User",
	Long:  `Delete a User`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		err := getClientGQL().DeleteUser(args[0])
		cobra.CheckErr(err)

		fmt.Printf("user '%s' deleted\n", key)
	},
}

func init() {
	createCmd.AddCommand(createUserCmd)
	listCmd.AddCommand(listUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)
}
