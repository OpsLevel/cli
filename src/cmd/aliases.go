package cmd

import (
	"github.com/opslevel/opslevel-go/v2023"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var createAliasCmd = &cobra.Command{
	Use:     "alias ID ALIAS",
	Short:   "Create an alias for a resource",
	Long:    `Create an alias for a resource`,
	Args:    cobra.ExactArgs(2),
	Example: `opslevel create alias XXXX my-new-alias`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		alias := args[1]
		_, err := getClientGQL().CreateAlias(opslevel.AliasCreateInput{
			OwnerId: opslevel.ID(id),
			Alias:   alias,
		})
		if err != nil {
			log.Error().Err(err).Msgf("unable to create alias '%s' for resource with id '%s'", alias, id)
		} else {
			log.Info().Msgf("created '%s' alias", alias)
		}
	},
}

var deleteAliasCmd = &cobra.Command{
	Use:        "alias ALIAS ALIAS_OWNER_TYPE",
	Short:      "Delete an alias on a resource",
	Long:       `Delete an alias on a resource`,
	Example:    `opslevel delete my-new-alias system`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		alias := args[0]
		ownerType := args[1]
		err := getClientGQL().DeleteAlias(opslevel.AliasDeleteInput{
			Alias:     alias,
			OwnerType: opslevel.AliasOwnerTypeEnum(ownerType),
		})
		if err != nil {
			log.Error().Err(err).Msgf("unable to delete alias '%s' for resource type  '%s'", alias, ownerType)
		} else {
			log.Info().Msgf("delete '%s' alias", alias)
		}
	},
}

func init() {
	createCmd.AddCommand(createAliasCmd)
	deleteCmd.AddCommand(deleteAliasCmd)
}
