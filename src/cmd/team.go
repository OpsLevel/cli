package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2022"
	"github.com/spf13/cobra"
)

var createTeamCmd = &cobra.Command{
	Use:        "team NAME",
	Short:      "Create a team",
	Long:       `Create a team`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"NAME"},
	Run: func(cmd *cobra.Command, args []string) {
		team, err := getClientGQL().CreateTeam(opslevel.TeamCreateInput{
			Name: args[0],
		})
		cobra.CheckErr(err)
		fmt.Println(team.Id)
	},
}

var createMemberCmd = &cobra.Command{
	Use:        "member TEAM_ID|TEAM_ALIAS EMAIL",
	Short:      "Add a member to a team",
	Long:       `Add a member to a team`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"TEAM_ID", "TEAM_ALIAS", "EMAIL"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		email := args[1]

		var team *opslevel.Team
		var err error
		if common.IsID(key) {
			team, err = getClientGQL().GetTeam(key)
		} else {
			team, err = getClientGQL().GetTeamWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(team.Id == nil, key)

		_, addErr := getClientGQL().AddMember(&team.TeamId, email)
		cobra.CheckErr(addErr)
		fmt.Printf("add member '%s' on team '%s'\n", email, team.Alias)
	},
}

var contactType string
var createContactCmd = &cobra.Command{
	Use:        "contact TEAM_ID|TEAM_ALIAS ADDRESS DISPLAYNAME",
	Short:      "Add a contact to a team",
	Long:       `Add a contact to a team`,
	Args:       cobra.MinimumNArgs(2),
	ArgAliases: []string{"TEAM_ID", "TEAM_ALIAS", "ADDRESS", "DISPLAYNAME"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		address := args[1]
		displayName := common.GetArg(args, 2, "")

		var team *opslevel.Team
		var err error
		if common.IsID(key) {
			team, err = getClientGQL().GetTeam(key)
		} else {
			team, err = getClientGQL().GetTeamWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(team.Id == nil, key)
		contactInput := opslevel.CreateContactSlack(address, displayName)
		switch contactType {
		case string(opslevel.ContactTypeEmail):
			contactInput = opslevel.CreateContactEmail(address, displayName)
		case string(opslevel.ContactTypeWeb):
			contactInput = opslevel.CreateContactWeb(address, displayName)
		}
		contact, err := getClientGQL().AddContact(team.TeamId.Alias, contactInput)
		cobra.CheckErr(err)
		if contact.Id == nil {
			cobra.CheckErr(fmt.Errorf("unable to create contact '%+v'", contactInput))
		}
		fmt.Printf("create contact '%+v' on team '%s'\n", contactInput, team.Alias)
	},
}

var getTeamCmd = &cobra.Command{
	Use:        "team ID|ALIAS",
	Short:      "Get details about a team",
	Long:       `Get details about a team`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		team, err := GetTeam(key)
		cobra.CheckErr(err)
		common.WasFound(team.Id == nil, key)
		common.PrettyPrint(team)
	},
}

func GetTeam(key string) (*opslevel.Team, error) {
	if common.IsID(key) {
		return getClientGQL().GetTeam(key)
	} else {
		return getClientGQL().GetTeamWithAlias(key)
	}
}

var listTeamCmd = &cobra.Command{
	Use:     "team",
	Aliases: []string{"teams"},
	Short:   "Lists the teams",
	Long:    `Lists the teams`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := getClientGQL().ListTeams()
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

var deleteTeamCmd = &cobra.Command{
	Use:        "team ID|ALIAS",
	Short:      "Delete a team",
	Long:       `Delete a team`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if common.IsID(key) {
			err := getClientGQL().DeleteTeam(key)
			cobra.CheckErr(err)
		} else {
			err := getClientGQL().DeleteTeamWithAlias(key)
			cobra.CheckErr(err)
		}
		fmt.Printf("team '%s' deleted\n", key)
	},
}

var deleteMemberCmd = &cobra.Command{
	Use:        "member TEAM_ID|TEAM_ALIAS EMAIL",
	Short:      "Removes a member from a team",
	Long:       `Removes a member from a team`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"EMAIL"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		email := args[1]

		var team *opslevel.Team
		var err error
		if common.IsID(key) {
			team, err = getClientGQL().GetTeam(key)
		} else {
			team, err = getClientGQL().GetTeamWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(team.Id == nil, key)

		_, removeErr := getClientGQL().RemoveMember(&team.TeamId, email)
		cobra.CheckErr(removeErr)
		fmt.Printf("member '%s' on team '%s' removed\n", email, key)
	},
}

var deleteContactCmd = &cobra.Command{
	Use:        "contact ID",
	Short:      "Removes a contact from a team",
	Long:       `Removes a contact from a team`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().RemoveContact(key)
		cobra.CheckErr(err)
		fmt.Printf("contact '%s' removed\n", key)
	},
}

var importTeamsCmd = &cobra.Command{
	Use:     "team",
	Aliases: []string{"teams"},
	Short:   "Imports teams from a CSV",
	Long: `Imports a list of teams from a CSV file with the column headers:
Name,Manager,Responsibilities,Group

Example:

cat << EOF | opslevel import teams -f -
Name,Manager,Responsibilities,Group
Platform,kyle@opslevel.com,Makes Tools,engineering
Sales,john@opslevel.com,Sells Tools,product
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := readImportFilepathAsCSV()
		cobra.CheckErr(err)
		for reader.Rows() {
			name := reader.Text("Name")
			input := opslevel.TeamCreateInput{
				Name:             name,
				ManagerEmail:     reader.Text("Manager"),
				Responsibilities: reader.Text("Responsibilities"),
			}
			group := reader.Text("Group")
			if group != "" {
				input.Group = opslevel.NewIdentifier(group)
			}
			team, err := getClientGQL().CreateTeam(input)
			if err != nil {
				log.Error().Err(err).Msgf("error creating team '%s'", name)
				continue
			}
			log.Info().Msgf("created team '%s' with id '%s'", team.Name, team.Id)
		}
	},
}

func init() {
	createCmd.AddCommand(createTeamCmd)
	createCmd.AddCommand(createMemberCmd)
	createCmd.AddCommand(createContactCmd)
	getCmd.AddCommand(getTeamCmd)
	listCmd.AddCommand(listTeamCmd)
	deleteCmd.AddCommand(deleteTeamCmd)
	deleteCmd.AddCommand(deleteMemberCmd)
	deleteCmd.AddCommand(deleteContactCmd)
	importCmd.AddCommand(importTeamsCmd)

	createContactCmd.Flags().StringVarP(&contactType, "type", "t", "slack", "The contact type. One of: slack|email|web [default: slack]")
}
