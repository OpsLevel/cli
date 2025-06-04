package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2025"

	"github.com/opslevel/cli/common"

	"github.com/spf13/cobra"
)

var secretAlias string

var exampleSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Example Secret",
	Long:  `Example Secret`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.SecretInput]())
	},
}

var createSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Create a team-owned secret",
	Long:  `Create a team-owned secret`,
	Example: `
cat << EOF | opslevel create secret --alias=my-secret-alias -f -
owner:
  alias: "devs"
value: "my-really-secure-secret-shhhh"
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readResourceInput[opslevel.SecretInput]()
		cobra.CheckErr(err)
		newSecret, err := getClientGQL().CreateSecret(secretAlias, *input)
		cobra.CheckErr(err)

		fmt.Printf("%s", string(newSecret.Id))
	},
}

var getSecretCmd = &cobra.Command{
	Use:        "secret",
	Short:      "Get details about a secret",
	Long:       `Get details about a secret`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		result, err := getClientGQL().GetSecret(identifier)
		cobra.CheckErr(err)

		common.WasFound(result.Id == "", identifier)
		if isYamlOutput() {
			common.YamlPrint(result)
		} else {
			common.PrettyPrint(result)
		}
	},
}

var listSecretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secrets"},
	Short:   "List secrets",
	Long:    `List secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListSecretsVaultsSecret(nil)
		list := resp.Nodes

		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("ALIAS", "ID", "OWNER", "UPDATED_AT")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", item.Alias, item.Id, item.Owner.Alias, item.Timestamps.UpdatedAt.Time)
			}
			w.Flush()
		}
	},
}

var updateSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Update an OpsLevel secret",
	Long:  `Update an OpsLevel secret`,
	Example: `
cat << EOF | opslevel update secret XXX_secret_id_XXX -f -
owner:
  alias: "platform"
value: "09sdf09werlkewlkjs0-9sdf
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		secretId := args[0]
		input, err := readResourceInput[opslevel.SecretInput]()
		cobra.CheckErr(err)
		secret, err := getClientGQL().UpdateSecret(secretId, *input)
		cobra.CheckErr(err)
		fmt.Printf("%s", string(secret.Id))
	},
}

var deleteSecretCmd = &cobra.Command{
	Use:        "secret ID|ALIAS",
	Short:      "Delete a secret",
	Long:       `Delete a secret from OpsLevel`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		secretId := args[0]
		err := getClientGQL().DeleteSecret(secretId)
		cobra.CheckErr(err)
		fmt.Printf("deleted secret: %s\n", secretId)
	},
}

func init() {
	createSecretCmd.Flags().StringVar(&secretAlias, "alias", "", "The alias for the secret")
	createSecretCmd.MarkFlagRequired("alias")

	exampleCmd.AddCommand(exampleSecretCmd)
	createCmd.AddCommand(createSecretCmd)
	getCmd.AddCommand(getSecretCmd)
	listCmd.AddCommand(listSecretsCmd)
	updateCmd.AddCommand(updateSecretCmd)
	deleteCmd.AddCommand(deleteSecretCmd)
}
