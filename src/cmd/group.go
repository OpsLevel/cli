package cmd

import (
	"fmt"

	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createGroupCmd = &cobra.Command{
	Use:   "group",
	Short: "Create a group",
	Long: `Create a group

cat << EOF | opslevel create group -f -
name: "My Group"
description: "Hello World Group"
parent:
  alias: "my-other-group-alias"
members:
  - email: info@opslevel.com
  - email: support@opslevel.com
team:
  - alias: "my-team-alias"
  - id: "s90s90ewr0fgd09sdf"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readGroupInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateGroup(*input)
		cobra.CheckErr(err)
		fmt.Println(result.Id)
	},
}

// TODO: Get Group

// TODO: List Groups

var updateGroupCmd = &cobra.Command{
	Use:   "group ID|ALIAS",
	Short: "Update a group",
	Long: `Update a group

cat << EOF | opslevel update group "my-group-alias" -f -
description: My updated group description
parent:
  alias: "next-group-alias"
EOF
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readGroupInput()
		cobra.CheckErr(err)
		filter, err := getClientGQL().UpdateGroup(key, *input)
		cobra.CheckErr(err)
		fmt.Println(filter.Id)
	},
}

var deleteGroupCmd = &cobra.Command{
	Use:        "group ID|ALIAS",
	Short:      "Delete a group",
	Long:       `Delete a group`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if common.IsID(key) {
			err := getClientGQL().DeleteGroup(key)
			cobra.CheckErr(err)
		} else {
			err := getClientGQL().DeleteGroupWithAlias(key)
			cobra.CheckErr(err)
		}
		fmt.Printf("deleted '%s' group\n", key)
	},
}

func init() {
	createCmd.AddCommand(createGroupCmd)
	updateCmd.AddCommand(updateGroupCmd)
	deleteCmd.AddCommand(deleteGroupCmd)
}

func readGroupInput() (*opslevel.GroupInput, error) {
	readCreateConfigFile()
	evt := &opslevel.GroupInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
