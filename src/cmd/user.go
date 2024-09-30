package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var exampleUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Example User",
	Long:  `Example User`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.UserInput]())
	},
}

var createUserCmd = &cobra.Command{
	Use:   "user EMAIL NAME [ROLE]",
	Short: "Create a User",
	Long:  "Create a User and optionally define the role (options `User`|`Admin`).",
	Example: `
opslevel create user "john@example.com" "John Doe"
opslevel create user "jane@example.com" "Jane Doe" Admin --skip-send-invite
opslevel create user "jane@example.com" "Jane Doe" Admin --skip-welcome-email
opslevel create user "jane@example.com" "Jane Doe" Admin --skip-send-invite --skip-welcome-email
`,
	Args: cobra.RangeArgs(2, 4),
	Run: func(cmd *cobra.Command, args []string) {
		email := args[0]
		name := args[1]
		role := opslevel.UserRoleUser
		if len(args) > 2 {
			desiredRole := strings.ToLower(args[2])
			if slices.Contains(opslevel.AllUserRole, desiredRole) {
				role = opslevel.UserRole(desiredRole)
			}
		}

		sendInvite, err := cmd.Flags().GetBool("skip-send-invite")
		cobra.CheckErr(err)
		skipEmail, err := cmd.Flags().GetBool("skip-welcome-email")
		cobra.CheckErr(err)

		userInput := opslevel.UserInput{
			Name:             opslevel.RefOf(name),
			Role:             opslevel.RefOf(role),
			SkipWelcomeEmail: opslevel.RefOf(skipEmail),
		}
		resource, err := getClientGQL().InviteUser(email, userInput, sendInvite)
		cobra.CheckErr(err)
		fmt.Println(resource.Id)
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "user {ID|EMAIL}",
	Short: "Update a user",
	Long:  `Update a group`,
	Example: `
cat << EOF | opslevel update user "john@example.com" -f -
name: John Foobar Doe
role: Admin
EOF
`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readResourceInput[opslevel.UserInput]()
		cobra.CheckErr(err)
		filter, err := getClientGQL().UpdateUser(key, *input)
		cobra.CheckErr(err)
		fmt.Println(filter.Id)
	},
}

var getUserCmd = &cobra.Command{
	Use:        "user {ID|EMAIL}",
	Short:      "Get details about a user",
	Example:    `opslevel get user john@example.com`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		filter, err := getClientGQL().GetUser(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(filter)
	},
}

var listUserCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Lists the users",
	Example: `
opslevel list user
opslevel list user --ignore-deactivated
opslevel list user -o json | jq 'map({"key": .Name, "value": .Role}) | from_entries'
`,
	Run: func(cmd *cobra.Command, args []string) {
		// payloadVars should remain nil if '--ignore-deactivated' not set
		var payloadVars *opslevel.PayloadVariables

		ignoreDeactivated, err := cmd.Flags().GetBool("ignore-deactivated")
		cobra.CheckErr(err)

		client := getClientGQL()
		if ignoreDeactivated {
			payloadVars = client.InitialPageVariablesPointer().WithoutDeactivedUsers()
		}

		resp, err := getClientGQL().ListUsers(payloadVars)
		cobra.CheckErr(err)
		list := resp.Nodes
		sort.Slice(list, func(i, j int) bool {
			return list[i].Email < list[j].Email
		})
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			w.Write([]string{"ID", "EMAIL", "NAME", "ROLE", "URL"})
			for _, item := range list {
				w.Write([]string{string(item.Id), item.Email, item.Name, string(item.Role), item.HTMLUrl})
			}
			w.Flush()
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
	Use:     "user {ID|EMAIL}",
	Short:   "Delete a User",
	Example: `opslevel delete user john@example.com`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]

		err := getClientGQL().DeleteUser(args[0])
		cobra.CheckErr(err)

		fmt.Printf("user '%s' deleted\n", key)
	},
}

var importUsersCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Imports users from a CSV",
	Long: `Imports a list of users from a CSV file with the column headers:
Name,Email,Role,Team`,
	Example: `
cat << EOF | opslevel import user -f -
Name,Email,Role,Team
Kyle Rockman,kyle@opslevel.com,Admin,platform
Edgar Ochoa,edgar@opslevel.com,Admin,platform
Adam Del Gobbo,adam@opslevel.com,User,sales
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := readImportFilepathAsCSV()
		cobra.CheckErr(err)
		for reader.Rows() {
			name := reader.Text("Name")
			email := reader.Text("Email")
			role := strings.ToLower(reader.Text("Role"))
			if email == "" {
				log.Error().Msgf("user '%s' has invalid email '%s'", name, email)
				continue
			}
			userRole := opslevel.UserRoleUser
			if slices.Contains(opslevel.AllUserRole, role) {
				userRole = opslevel.UserRole(role)
			}
			input := opslevel.UserInput{
				Name: opslevel.RefOf(name),
				Role: opslevel.RefOf(userRole),
			}
			user, err := getClientGQL().InviteUser(email, input, false)
			if err != nil {
				log.Error().Err(err).Msgf("error inviting user '%s' with email '%s'", name, email)
				continue
			}
			log.Info().Msgf("invited user '%s' with email '%s'", user.Name, user.Email)
			team := reader.Text("Team")
			if team != "" {
				t, err := GetTeam(team)
				if err != nil {
					log.Error().Err(err).Msgf("error finding team '%s' for user '%s'", team, user.Email)
					continue
				}
				newMembership := opslevel.TeamMembershipUserInput{
					User: opslevel.NewUserIdentifier(email),
					Role: opslevel.RefOf(string(user.Role)),
				}
				_, err = getClientGQL().AddMemberships(&t.TeamId, newMembership)
				if err != nil {
					log.Error().Err(err).Msgf("error adding user '%s' to team '%s'", user.Email, t.Name)
					continue
				}
				log.Info().Msgf("added user '%s' to team '%s'", user.Email, t.Name)
			}

		}
	},
}

func init() {
	createUserCmd.Flags().Bool("skip-send-invite", false, "If this flag is set the welcome e-mail will be not be sent")
	createUserCmd.Flags().Bool("skip-welcome-email", false, "If this flag is set send an invite email even if notifications are disabled for the account")
	listUserCmd.Flags().Bool("ignore-deactivated", false, "If this flag is set only return active users")

	exampleCmd.AddCommand(exampleUserCmd)
	createCmd.AddCommand(createUserCmd)
	updateCmd.AddCommand(updateUserCmd)
	getCmd.AddCommand(getUserCmd)
	listCmd.AddCommand(listUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)
	importCmd.AddCommand(importUsersCmd)
}
