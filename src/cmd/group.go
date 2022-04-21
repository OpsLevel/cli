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
		common.WasFound(group.Id == nil, key)
		common.PrettyPrint(group)
	},
}

var getGroupMembersCommand = &cobra.Command{
	Use:        "members ID|ALIAS",
	Short:      "Get members for a group",
	Long:       `The users who are members of the group.`,
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
		common.WasFound(group.Id == nil, key)
		members, err := group.Members(getClientGQL())
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.PrettyPrint(members)
		} else {
			w := common.NewTabWriter("EMAIL", "ID")
			for _, item := range members {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Email, item.Id)
			}
			w.Flush()
		}
	},
}

var getGroupDescendantRepositoriesCommand = &cobra.Command{
	Use:        "repositories ID|ALIAS",
	Short:      "Get descendant repositories for a group",
	Long:       `All the repositories that fall under this group - ex. this group's child repositories, all the child repositories of this group's descendants, etc.`,
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
		common.WasFound(group.Id == nil, key)
		descendantRepositories, err := group.DescendantRepositories(getClientGQL())
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.PrettyPrint(descendantRepositories)
		} else {
			w := common.NewTabWriter("ALIAS", "ID")
			for _, item := range descendantRepositories {
				fmt.Fprintf(w, "%s\t%s\t\n", item.DefaultAlias, item.Id)
			}
			w.Flush()
		}
	},
}

var getGroupDescendantServicesCommand = &cobra.Command{
	Use:        "services ID|ALIAS",
	Short:      "Get descendant services for a group",
	Long:       `All the services that fall under this group - ex. this group's child services, all the child services of this group's descendants, etc.`,
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
		common.WasFound(group.Id == nil, key)
		descendantServices, err := group.DescendantServices(getClientGQL())
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.PrettyPrint(descendantServices)
		} else {
			w := common.NewTabWriter("ID")
			for _, item := range descendantServices {
				fmt.Fprintf(w, "%s\t\n", item.Id)
			}
			w.Flush()
		}
	},
}

var getGroupDescendantSubgroupsCommand = &cobra.Command{
	Use:        "subgroups ID|ALIAS",
	Short:      "Get descendant subgroups for a group",
	Long:       `All the groups that fall under this group - ex. this group's child groups, children of those groups, etc.`,
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
		common.WasFound(group.Id == nil, key)
		descendantSubgroups, err := group.DescendantSubgroups(getClientGQL())
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.PrettyPrint(descendantSubgroups)
		} else {
			w := common.NewTabWriter("ALIAS", "ID")
			for _, item := range descendantSubgroups {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Alias, item.Id)
			}
			w.Flush()
		}
	},
}

var getGroupDescendantTeamsCommand = &cobra.Command{
	Use:        "teams ID|ALIAS",
	Short:      "Get descendant teams for a group",
	Long:       `All the teams that fall under this group - ex. this group's child teams, all the child teams of this group's descendants, etc.`,
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
		common.WasFound(group.Id == nil, key)
		descendantTeams, err := group.DescendantTeams(getClientGQL())
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.PrettyPrint(descendantTeams)
		} else {
			w := common.NewTabWriter("ALIAS", "ID")
			for _, item := range descendantTeams {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Alias, item.Id)
			}
			w.Flush()
		}
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
	getGroupCommand.AddCommand(getGroupMembersCommand)
	getGroupCommand.AddCommand(getGroupDescendantRepositoriesCommand)
	getGroupCommand.AddCommand(getGroupDescendantServicesCommand)
	getGroupCommand.AddCommand(getGroupDescendantSubgroupsCommand)
	getGroupCommand.AddCommand(getGroupDescendantTeamsCommand)
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
