package cmd

import (
	"encoding/json"
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

var getGroupCommand = &cobra.Command{
	Use:        "group ID|ALIAS",
	Short:      "Get details about a group",
	Long:       `Get details about a group`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(key)
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		common.PrettyPrint(group)
		common.WasFound(group.Id == nil, key)
	},
}

var getMembersCommand = &cobra.Command{
	Use:        "members ID|ALIAS",
	Short:      "Get members for a group",
	Long:       `Get members for a group`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(key)
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		members, err := group.Members(getClientGQL())
		cobra.CheckErr(err)
		common.PrettyPrint(members)
		common.WasFound(group.Id == nil, key)
	},
}

var getDescendantRepositoriesCommand = &cobra.Command{
	Use:        "repositories ID|ALIAS",
	Short:      "Get descendant repositories for a group",
	Long:       `Get descendant repositories for a group`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(key)
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		descendantRepositories, err := group.DescendantRepositories(getClientGQL())
		cobra.CheckErr(err)
		common.PrettyPrint(descendantRepositories)
		common.WasFound(group.Id == nil, key)
	},
}

var getDescendantServicesCommand = &cobra.Command{
	Use:        "services ID|ALIAS",
	Short:      "Get descendant services for a group",
	Long:       `Get descendant services for a group`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(key)
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		descendantServices, err := group.DescendantServices(getClientGQL())
		cobra.CheckErr(err)
		common.PrettyPrint(descendantServices)
		common.WasFound(group.Id == nil, key)
	},
}

var getDescendantSubgroupsCommand = &cobra.Command{
	Use:        "subgroups ID|ALIAS",
	Short:      "Get descendant subgroups for a group",
	Long:       `Get descendant subgroups for a group`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(key)
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		descendantSubgroups, err := group.DescendantSubgroups(getClientGQL())
		cobra.CheckErr(err)
		common.PrettyPrint(descendantSubgroups)
		common.WasFound(group.Id == nil, key)
	},
}

var getDescendantTeamsCommand = &cobra.Command{
	Use:        "teams ID|ALIAS",
	Short:      "Get descendant teams for a group",
	Long:       `Get descendant teams for a group`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(key)
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		descendantTeams, err := group.DescendantTeams(getClientGQL())
		cobra.CheckErr(err)
		common.PrettyPrint(descendantTeams)
		common.WasFound(group.Id == nil, key)
	},
}

var listGroupCmd = &cobra.Command{
	Use:     "group",
	Aliases: []string{"groups"},
	Short:   "Lists the groups",
	Long:    `Lists the groups`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := getClientGQL().ListGroups()
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Alias, item.Id)
			}
			w.Flush()
		}
	},
}

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
	getCmd.AddCommand(getGroupCommand)
	getGroupCommand.AddCommand(getMembersCommand)
	getGroupCommand.AddCommand(getDescendantRepositoriesCommand)
	getGroupCommand.AddCommand(getDescendantServicesCommand)
	getGroupCommand.AddCommand(getDescendantSubgroupsCommand)
	getGroupCommand.AddCommand(getDescendantTeamsCommand)
	listCmd.AddCommand(listGroupCmd)
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
