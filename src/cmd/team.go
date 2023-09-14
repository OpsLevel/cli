package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var createTeamCmd = &cobra.Command{
	Use:   "team NAME",
	Short: "Create a team",
	Example: `opslevel create team my-team

cat << EOF | opslevel create team my-team" -f -
managerEmail: "manager@example.com""
group:
  alias: "my-group"
responsibilities: "all the things"
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"NAME"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readTeamCreateInput()
		input.Name = key
		cobra.CheckErr(err)
		team, err := getClientGQL().CreateTeam(*input)
		cobra.CheckErr(err)
		fmt.Println(team.Id)
	},
}

var createMemberCmd = &cobra.Command{
	Use:        "member {TEAM_ID|TEAM_ALIAS} EMAIL",
	Short:      "Add a member to a team",
	Example:    `opslevel create member my-team john@example.com`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"TEAM_ID", "TEAM_ALIAS", "EMAIL"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		email := args[1]

		var team *opslevel.Team
		var err error
		if common.IsID(key) {
			team, err = getClientGQL().GetTeam(opslevel.ID(key))
		} else {
			team, err = getClientGQL().GetTeamWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(team.Id == "", key)

		_, addErr := getClientGQL().AddMember(&team.TeamId, email)
		cobra.CheckErr(addErr)
		fmt.Printf("add member '%s' on team '%s'\n", email, team.Alias)
	},
}

var (
	contactType      string
	createContactCmd = &cobra.Command{
		Use:   "contact {TEAM_ID|TEAM_ALIAS} ADDRESS DISPLAYNAME",
		Short: "Add a contact to a team",
		Example: `opslevel create contact --type=slack my-team #general General
opslevel create contact --type=email my-team team@example.com "Mailing List"`,
		Args:       cobra.MinimumNArgs(2),
		ArgAliases: []string{"TEAM_ID", "TEAM_ALIAS", "ADDRESS", "DISPLAYNAME"},
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			address := args[1]
			displayName := common.GetArg(args, 2, "")

			var team *opslevel.Team
			var err error
			if common.IsID(key) {
				team, err = getClientGQL().GetTeam(opslevel.ID(key))
			} else {
				team, err = getClientGQL().GetTeamWithAlias(key)
			}
			cobra.CheckErr(err)
			common.WasFound(team.Id == "", key)
			contactInput := opslevel.CreateContactSlack(address, &displayName)
			switch contactType {
			case string(opslevel.ContactTypeEmail):
				contactInput = opslevel.CreateContactEmail(address, &displayName)
			case string(opslevel.ContactTypeWeb):
				contactInput = opslevel.CreateContactWeb(address, &displayName)
			}
			contact, err := getClientGQL().AddContact(team.TeamId.Alias, contactInput)
			cobra.CheckErr(err)
			if contact.Id == "" {
				cobra.CheckErr(fmt.Errorf("unable to create contact '%+v'", contactInput))
			}
			fmt.Printf("create contact '%+v' on team '%s'\n", contactInput, team.Alias)
		},
	}
)

var createTeamTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Create a team tag",
	Long:  `Create a team tag`,
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("This command is deprecated! Please use \nopslevel create tag <args>")
		cobra.CheckErr(err)
	},
}

var updateTeamCmd = &cobra.Command{
	Use:   "team {ID|ALIAS}",
	Short: "Update a team",
	Example: `
cat << EOF | opslevel update team my-team" -f -
managerEmail: "manager@example.com""
group:
  alias: "my-group"
responsibilities: "all the things"
EOF
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readTeamUpdateInput()
		input.Id = opslevel.ID(key)
		cobra.CheckErr(err)
		team, err := getClientGQL().UpdateTeam(*input)
		cobra.CheckErr(err)
		fmt.Println(team.Id)
	},
}

var getTeamCmd = &cobra.Command{
	Use:        "team {ID|ALIAS}",
	Short:      "Get details about a team",
	Example:    `opslevel get team my-team | jq '.Members.Nodes | map(.Email)'`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		team, err := GetTeam(key)
		cobra.CheckErr(err)
		common.WasFound(team.Id == "", key)
		common.PrettyPrint(team)
	},
}

func GetTeam(key string) (*opslevel.Team, error) {
	if common.IsID(key) {
		return getClientGQL().GetTeam(opslevel.ID(key))
	} else {
		return getClientGQL().GetTeamWithAlias(key)
	}
}

var listTeamCmd = &cobra.Command{
	Use:     "team",
	Aliases: []string{"teams"},
	Short:   "Lists the teams",
	Example: `
opslevel list team
opslevel list team -o json | jq 'map((.Members.Nodes | map(.Email)))'
`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListTeams(nil)
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

var getTeamTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Get a team's tag",
	Long:  `Get a team's tag`,
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("This command is deprecated! Please use \nopslevel get tag <args>")
		cobra.CheckErr(err)
	},
}

var deleteTeamCmd = &cobra.Command{
	Use:   "team {ID|ALIAS}",
	Short: "Delete a team",
	Example: `
opslevel delete team my-team
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		if common.IsID(key) {
			err := getClientGQL().DeleteTeam(opslevel.ID(key))
			cobra.CheckErr(err)
		} else {
			err := getClientGQL().DeleteTeamWithAlias(key)
			cobra.CheckErr(err)
		}
		fmt.Printf("team '%s' deleted\n", key)
	},
}

var deleteMemberCmd = &cobra.Command{
	Use:        "member {TEAM_ID|TEAM_ALIAS} EMAIL",
	Short:      "Removes a member from a team",
	Example:    `opslevel delete member my-team john@example.com`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"EMAIL"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		email := args[1]

		var team *opslevel.Team
		var err error
		if common.IsID(key) {
			team, err = getClientGQL().GetTeam(opslevel.ID(key))
		} else {
			team, err = getClientGQL().GetTeamWithAlias(key)
		}
		cobra.CheckErr(err)
		common.WasFound(team.Id == "", key)

		_, removeErr := getClientGQL().RemoveMember(&team.TeamId, email)
		cobra.CheckErr(removeErr)
		fmt.Printf("member '%s' on team '%s' removed\n", email, key)
	},
}

var deleteContactCmd = &cobra.Command{
	Use:        "contact ID",
	Short:      "Removes a contact from a team",
	Example:    `opslevel delete contact XXXXXXXXXXXXXXXX`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().RemoveContact(opslevel.ID(key))
		cobra.CheckErr(err)
		fmt.Printf("contact '%s' removed\n", key)
	},
}

var deleteTeamTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Delete a team's tag",
	Long:  `Delete a team's tag`,
	Run: func(cmd *cobra.Command, args []string) {
		err := errors.New("This command is deprecated! Please use \nopslevel delete tag <args>")
		cobra.CheckErr(err)
	},
}

var importTeamsCmd = &cobra.Command{
	Use:     "team",
	Aliases: []string{"teams"},
	Short:   "Imports teams from a CSV",
	Long: `Imports a list of teams from a CSV file with the column headers:
Name,Manager,Responsibilities,Group`,
	Example: `
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
	createTeamCmd.AddCommand(createTeamTagCmd)
	updateCmd.AddCommand(updateTeamCmd)
	getCmd.AddCommand(getTeamCmd)
	getTeamCmd.AddCommand(getTeamTagCmd)
	listCmd.AddCommand(listTeamCmd)
	deleteCmd.AddCommand(deleteTeamCmd)
	deleteCmd.AddCommand(deleteMemberCmd)
	deleteCmd.AddCommand(deleteContactCmd)
	deleteTeamCmd.AddCommand(deleteTeamTagCmd)
	importCmd.AddCommand(importTeamsCmd)

	createContactCmd.Flags().StringVarP(&contactType, "type", "t", "slack", "The contact type. One of: slack|email|web [default: slack]")
}

func readTeamCreateInput() (*opslevel.TeamCreateInput, error) {
	readInputConfig()
	evt := &opslevel.TeamCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}

func readTeamUpdateInput() (*opslevel.TeamUpdateInput, error) {
	readInputConfig()
	evt := &opslevel.TeamUpdateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
