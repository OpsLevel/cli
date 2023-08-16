package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/creasty/defaults"

	"github.com/opslevel/opslevel-go/v2023"
	"gopkg.in/yaml.v3"

	"github.com/opslevel/cli/common"

	"github.com/spf13/cobra"
)

var secretAlias string

var createSecretCmd = &cobra.Command{
	Use:   "secret",
	Short: "Create a team-owned secret",
	Long:  `Create a team-owned secret`,
	Example: `
cat << EOF | opslevel create secret --alias=my-secret-alias -f -
owner: "devs"
value: "my-really-secure-secret-shhhh"
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readSecretInput()
		cobra.CheckErr(err)
		fmt.Println("creating secret...")
		// TODO: remove headers when API is ready
		headers := map[string]string{"GraphQL-Visibility": "internal"}
		newSecret, err := getClientGQL(opslevel.SetHeaders(headers)).CreateSecret(secretAlias, *input)
		cobra.CheckErr(err)

		fmt.Println(newSecret.ID)
		fmt.Printf("%+v\n", newSecret)
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
		// TODO: remove headers when API is ready
		headers := map[string]string{"GraphQL-Visibility": "internal"}
		newSecret, err := getClientGQL(opslevel.SetHeaders(headers)).GetSecret(identifier)
		cobra.CheckErr(err)

		fmt.Println(newSecret.ID)
		fmt.Printf("%+v\n", newSecret)
	},
}

var listSecretsCmd = &cobra.Command{
	Use:     "secrets",
	Aliases: []string{"secrets"},
	Short:   "List secrets",
	Long:    `List secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listing secrets...")
		// TODO: remove headers when API is ready
		headers := map[string]string{"GraphQL-Visibility": "internal"}
		resp, err := getClientGQL(opslevel.SetHeaders(headers)).ListSecretsVaultsSecret(nil)
		list := resp.Nodes

		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("ALIAS", "ID", "OWNER")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Alias, item.ID, item.Owner.Alias)
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
    owner: "platform"
    value: "09sdf09werlkewlkjs0-9sdf
		EOF
		`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		secretId := args[0]
		input, err := readSecretInput()
		// TODO: remove headers when API is ready
		headers := map[string]string{"GraphQL-Visibility": "internal"}
		secret, err := getClientGQL(opslevel.SetHeaders(headers)).UpdateSecret(secretId, *input)
		cobra.CheckErr(err)
		fmt.Printf("Updated '%s' at %s", secret.Alias, secret.Timestamps.UpdatedAt.Time)
	},
}

var deleteSecretCmd = &cobra.Command{
	Use:        "secret ID|ALIAS",
	Short:      "Delete a system",
	Long:       `Delete a system from OpsLevel`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		secretId := args[0]
		// TODO: remove headers when API is ready
		headers := map[string]string{"GraphQL-Visibility": "internal"}
		err := getClientGQL(opslevel.SetHeaders(headers)).DeleteSecret(secretId)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' secret\n", secretId)
	},
}

func readSecretInput() (*opslevel.SecretInput, error) {
	file, err := io.ReadAll(os.Stdin)
	cobra.CheckErr(err)
	var evt struct {
		Owner string `yaml:"owner"`
		Value string `yaml:"value"`
	}
	cobra.CheckErr(yaml.Unmarshal(file, &evt))
	secretInput := &opslevel.SecretInput{}
	if err := defaults.Set(secretInput); err != nil {
		return nil, err
	}

	secretInput.Value = evt.Value
	secretInput.Owner = *opslevel.NewIdentifier(evt.Owner)
	return secretInput, nil
}

func init() {
	createSecretCmd.Flags().StringVar(&secretAlias, "alias", "", "The alias for the secret")
	createSecretCmd.MarkFlagRequired("alias")

	createCmd.AddCommand(createSecretCmd)
	getCmd.AddCommand(getSecretCmd)
	listCmd.AddCommand(listSecretsCmd)
	updateCmd.AddCommand(updateSecretCmd)
	deleteCmd.AddCommand(deleteSecretCmd)
}
