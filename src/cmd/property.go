package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"

	"github.com/spf13/cobra"
)

var examplePropertyDefinitionCmd = &cobra.Command{
	Use:   "property-definition",
	Short: "Example Property Definition",
	Long:  `Example Property Definition`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.PropertyDefinitionInput]())
	},
}

var createPropertyDefinitonCmd = &cobra.Command{
	Use:   "property-definition",
	Short: "Create a property-definition",
	Long:  `Create a property-definition`,
	Example: fmt.Sprintf(`
cat << EOF | opslevel create property-definition -f -
%s
EOF`, getYaml[opslevel.PropertyDefinitionInput]()),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readPropertyDefinitionInput()
		cobra.CheckErr(err)
		newPropertyDefinition, err := getClientGQL().CreatePropertyDefinition(*input)
		cobra.CheckErr(err)

		fmt.Println(newPropertyDefinition.Id)
	},
}

// The schema in PropertyDefinitionInput can be a nested map[string]any and needs to be handled separately
func readPropertyDefinitionInput() (*opslevel.PropertyDefinitionInput, error) {
	d, err := readResourceInput[map[string]any]()
	if err != nil {
		return nil, err
	}
	data := *d
	name, ok := data["name"].(string)
	if !ok {
		return nil, fmt.Errorf("name is required and must be a string")
	}
	schema, ok := data["schema"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("schema is required and must be a JSON object")
	}
	propDefInput := opslevel.PropertyDefinitionInput{
		Name:   name,
		Schema: opslevel.JSON(schema),
	}

	if description, ok := data["description"].(string); !ok {
		propDefInput.Description = description
	}
	if propertyDisplayStatus, ok := data["propertyDisplayStatus"].(string); !ok {
		propDefInput.PropertyDisplayStatus = opslevel.PropertyDisplayStatusEnum(propertyDisplayStatus)
	}

	return &propDefInput, nil
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

		if isYamlOutput() {
			common.YamlPrint(result)
		} else {
			common.PrettyPrint(result)
		}
	},
}

var listPropertyDefinitionsCmd = &cobra.Command{
	Use:     "property-definition",
	Short:   "List property definitions",
	Aliases: []string{"property-definitions"},
	Long:    "List property definitions",
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
		fmt.Printf("deleted property definition '%s'\n", propertyDefinitionId)
	},
}

func init() {
	exampleCmd.AddCommand(examplePropertyDefinitionCmd)
	createCmd.AddCommand(createPropertyDefinitonCmd)
	getCmd.AddCommand(getPropertyDefinition)
	listCmd.AddCommand(listPropertyDefinitionsCmd)
	deleteCmd.AddCommand(deletePropertyDefinitonCmd)
}
