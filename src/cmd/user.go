package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2022"
	"github.com/rs/zerolog/log"
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

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

var importUsersCmd = &cobra.Command{
	Use:     "user",
	Aliases: []string{"users"},
	Short:   "Imports users from a CSV",
	Long: `Imports a list of users from a CSV file with the column headers:
Name,Email,Role,Team

Example:

cat << EOF | opslevel import user -f -
Name,Email,Role,Team
Kyle Rockman,kyle@opslevel.com,Admin,platform
Edger Ochoa,edgar@opslevel.com,Admin,platform
Adam Del Gobbo,adam@opslevel.com,User,sales
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		reader, err := readImportFilepathAsCSV()
		cobra.CheckErr(err)
		for reader.Rows() {
			name := reader.Text("Name")
			email := reader.Text("Email")
			role := reader.Text("Role")
			if email == "" {
				log.Error().Msgf("user '%s' has invalid email '%s'", name, email)
				continue
			}
			userRole := opslevel.UserRoleUser
			if Contains(opslevel.AllUserRole(), role) {
				userRole = opslevel.UserRole(role)
			}
			input := opslevel.UserInput{
				Name: name,
				Role: userRole,
			}
			user, err := getClientGQL().InviteUser(email, input)
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
				_, err = getClientGQL().AddMember(&t.TeamId, user.Email)
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
	createCmd.AddCommand(createUserCmd)
	listCmd.AddCommand(listUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)
	importCmd.AddCommand(importUsersCmd)
}
