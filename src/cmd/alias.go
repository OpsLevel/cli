package cmd

import (
	"fmt"
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var createAliasCommand = &cobra.Command{
	Use:     "alias ID ALIAS",
	Short:   "Create an alias",
	Long:    "Create an alias",
	Args:    cobra.MinimumNArgs(2),
	Example: `opslevel create alias XXXXX my-alias`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		alias := args[1]
		client := getClientGQL()
		_, err := client.CreateAlias(opslevel.AliasCreateInput{
			OwnerId: opslevel.ID(id),
			Alias:   alias,
		})
		if err != nil {
			log.Error().Err(err).Msg("failed to create alias")
		} else {
			log.Info().Msg("alias created")
		}
	},
}

var deleteAliasCommand = &cobra.Command{
	Use:   "alias ALIAS",
	Short: "Delete an alias",
	Long:  "Delete an alias",
	Args:  cobra.MinimumNArgs(1),
	Example: `
# Delete alias on a service
opslevel delete alias my-service-alias
# Or specify -t to target a different resource type
opslevel delete alias -t group my-group-alias
opslevel delete alias -t infrastructure-resource my-infra-alias`,
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		aliasType := cmd.Flags().Lookup("type").Value.String()
		if !Contains(opslevel.AllAliasOwnerTypeEnum, aliasType) {
			log.Error().Msgf("invalid alias type '%s'", aliasType)
			os.Exit(1)
		}
		client := getClientGQL()
		err := client.DeleteAlias(opslevel.AliasDeleteInput{
			Alias:     alias,
			OwnerType: opslevel.AliasOwnerTypeEnum(aliasType),
		})
		if err != nil {
			log.Error().Err(err).Msgf("failed to delete alias '%s'", alias)
		} else {
			log.Info().Msgf("alias '%s' deleted", alias)
		}
	},
}

func init() {
	createCmd.AddCommand(createAliasCommand)
	deleteCmd.AddCommand(deleteAliasCommand)

	deleteAliasCommand.Flags().StringP("type", "t", "service", fmt.Sprintf("the resource type that the alias is on.  Can be one of [%s]", strings.Join(opslevel.AllAliasOwnerTypeEnum, ", ")))
}
