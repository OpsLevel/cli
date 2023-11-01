package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var createGroupCmd = &cobra.Command{
	Use:        "group",
	Short:      "Create a group (deprecated). Use teams instead.",
	Long:       `Create a group (deprecated). Use teams instead.`,
	Deprecated: `Please use Teams instead. See https://docs.opslevel.com/docs/groups`,
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("Groups are deprecated! Please use Teams instead.\nopslevel create team <args>")
		cobra.CheckErr(err)
	},
}

var getGroupCommand = &cobra.Command{
	Use:        "group ID|ALIAS",
	Aliases:    []string{"groups"},
	Short:      "Get details about a group",
	Long:       `Get details about a group`,
	Deprecated: `Please convert this Group into a Team. See https://docs.opslevel.com/docs/groups`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(opslevel.ID(key))
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(group.Id == "", key)
		common.PrettyPrint(group)
	},
}

var getGroupMembersCommand = &cobra.Command{
	Use:        "member ID|ALIAS",
	Aliases:    []string{"members"},
	Short:      "Get members for a group",
	Long:       `The users who are members of the group.`,
	Deprecated: `Please convert this Group into a Team. See https://docs.opslevel.com/docs/groups`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(opslevel.ID(key))
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(group.Id == "", key)
		resp, err := group.Members(getClientGQL(), nil)
		members := resp.Nodes
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
	Use:        "repository ID|ALIAS",
	Aliases:    []string{"repositories"},
	Short:      "Get descendant repositories for a group",
	Long:       `All the repositories that fall under this group - ex. this group's child repositories, all the child repositories of this group's descendants, etc.`,
	Deprecated: `Please convert this Group into a Team. See https://docs.opslevel.com/docs/groups`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(opslevel.ID(key))
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(group.Id == "", key)
		resp, err := group.DescendantRepositories(getClientGQL(), nil)
		descendantRepositories := resp.Nodes
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
	Use:        "service ID|ALIAS",
	Aliases:    []string{"services"},
	Short:      "Get descendant services for a group",
	Long:       `All the services that fall under this group - ex. this group's child services, all the child services of this group's descendants, etc.`,
	Deprecated: `Please convert this Group into a Team. See https://docs.opslevel.com/docs/groups`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(opslevel.ID(key))
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(group.Id == "", key)
		resp, err := group.DescendantServices(getClientGQL(), nil)
		descendantServices := resp.Nodes
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.PrettyPrint(descendantServices)
		} else {
			w := common.NewTabWriter("ALIAS", "ID")
			for _, item := range descendantServices {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Aliases[0], item.Id)
			}
			w.Flush()
		}
	},
}

var getGroupDescendantSubgroupsCommand = &cobra.Command{
	Use:        "subgroup ID|ALIAS",
	Aliases:    []string{"subgroups"},
	Short:      "Get descendant subgroups for a group",
	Long:       `All the groups that fall under this group - ex. this group's child groups, children of those groups, etc.`,
	Deprecated: `Please convert this Group into a Team. See https://docs.opslevel.com/docs/groups`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(opslevel.ID(key))
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(group.Id == "", key)
		resp, err := group.DescendantSubgroups(getClientGQL(), nil)
		descendantSubgroups := resp.Nodes
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
	Use:        "team ID|ALIAS",
	Aliases:    []string{"teams"},
	Short:      "Get descendant teams for a group",
	Long:       `All the teams that fall under this group - ex. this group's child teams, all the child teams of this group's descendants, etc.`,
	Deprecated: `Please convert this parent Group into a Team. See https://docs.opslevel.com/docs/groups`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var group *opslevel.Group
		var err error
		if common.IsID(key) {
			group, err = getClientGQL().GetGroup(opslevel.ID(key))
		} else {
			group, err = getClientGQL().GetGroupWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(group.Id == "", key)
		resp, err := group.DescendantTeams(getClientGQL(), nil)
		descendantTeams := resp.Nodes
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
	Use:        "group",
	Aliases:    []string{"groups"},
	Short:      "Lists the groups",
	Long:       `Lists the groups`,
	Deprecated: `Please convert these Groups into Teams. See https://docs.opslevel.com/docs/groups`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListGroups(nil)
		list := resp.Nodes
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
	Use:        "group ID|ALIAS",
	Short:      "Update a group",
	Long:       `Update a group (deprecated). Use teams instead.`,
	Deprecated: `Please convert this Group into a Team. See https://docs.opslevel.com/docs/groups`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("Groups are deprecated! Please use Teams instead.\nopslevel update team <args>")
		cobra.CheckErr(err)
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
			err := getClientGQL().DeleteGroup(key)
			cobra.CheckErr(err)
		}
		fmt.Printf("deleted '%s' group\n", key)
	},
}

var importGroupsCmd = &cobra.Command{
	Use:        "group",
	Aliases:    []string{"groups"},
	Short:      "(Formerly) Imports groups from a CSV",
	Long:       "(Formerly) Imports groups from a CSV",
	Deprecated: `Groups are deprecated! Please use Teams instead. See https://docs.opslevel.com/docs/groups`,
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("Groups are deprecated! Please use Teams instead.\nopslevel import team <args>")
		cobra.CheckErr(err)
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
	importCmd.AddCommand(importGroupsCmd)
}
