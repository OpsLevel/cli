package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"

	"github.com/spf13/cobra"
)

var createPropertyDefinitonCmd = &cobra.Command{
	Use:   "property-definition",
	Short: "Create a property-definition",
	Long:  `Create a property-definition`,
	Example: `
cat << EOF | opslevel create property-definition  -f -
name: "Is Beta Feature"
schema: {"$schema":"https://json-schema.org/draft/2020-12/schema","type":"boolean"}
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		// input, err := readResourceInput[opslevel.PropertyDefinitionInput]()
		// cobra.CheckErr(err)
		input := opslevel.PropertyDefinitionInput{
			Name:   "Is Beta Feature",
			Schema: opslevel.JSONString(`{"$schema":"https://json-schema.org/draft/2020-12/schema","type":"boolean"}`),
		}
		newPropertyDefinition, err := getClientGQL().CreatePropertyDefinition(input)
		cobra.CheckErr(err)

		fmt.Println(newPropertyDefinition.Id)
	},
}

var getPropertyDefinition = &cobra.Command{
	Use:        "property-definition",
	Short:      "Get details about a property definition",
	Long:       `Get details about a property definition`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		result, err := getClientGQL().GetPropertyDefinition(identifier)
		cobra.CheckErr(err)

		common.WasFound(result.Id == "", identifier)
		if isYamlOutput() {
			common.YamlPrint(result)
		} else {
			common.PrettyPrint(result)
		}
	},
}

var listPropertyDefinitionsCmd = &cobra.Command{
	Use:   "property-definitions",
	Short: "List property definitions",
	Long:  "List property definitions",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListPropertyDefinitions(nil)
		list := resp.Nodes

		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("ALIASES", "ID", "NAME", "SCHEMA")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", item.Aliases, item.Id, item.Name, item.Schema.ToJSON())
			}
			w.Flush()
		}
	},
}

// var updateSecretCmd = &cobra.Command{
// 	Use:   "secret",
// 	Short: "Update an OpsLevel secret",
// 	Long:  `Update an OpsLevel secret`,
// 	Example: `
// cat << EOF | opslevel update secret XXX_secret_id_XXX -f -
// owner:
//   alias: "platform"
// value: "09sdf09werlkewlkjs0-9sdf
// EOF`,
// 	Args:       cobra.ExactArgs(1),
// 	ArgAliases: []string{"ID"},
// 	Run: func(cmd *cobra.Command, args []string) {
// 		secretId := args[0]
// 		input, err := readSecretInput()
// 		cobra.CheckErr(err)
// 		secret, err := getClientGQL().UpdateSecret(secretId, *input)
// 		cobra.CheckErr(err)
// 		fmt.Println(secret.ID)
// 	},
// }

var deletePropertyDefinitonCmd = &cobra.Command{
	Use:        "property-definition ID",
	Short:      "Delete a property definitions",
	Long:       "Delete a property definitions",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		propertyDefinitionId := args[0]
		err := getClientGQL().DeletePropertyDefinition(propertyDefinitionId)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' property definition\n", propertyDefinitionId)
	},
}

func init() {
	createCmd.AddCommand(createPropertyDefinitonCmd)
	getCmd.AddCommand(getPropertyDefinition)
	listCmd.AddCommand(listPropertyDefinitionsCmd)
	deleteCmd.AddCommand(deletePropertyDefinitonCmd)
}
