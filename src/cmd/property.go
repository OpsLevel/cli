package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"

	"github.com/spf13/cobra"
)

var examplePropertyCmd = &cobra.Command{
	Use:   "property",
	Short: "Example Property",
	Long:  `Example Property`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.PropertyInput]())
	},
}

var assignPropertyCmd = &cobra.Command{
	Use:   "property",
	Short: "Assign a Property",
	Long:  `Assign a Property`,
	Example: fmt.Sprintf(`
cat << EOF | opslevel assign property-definition -f -
%s
EOF`, getYaml[opslevel.PropertyInput]()),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readResourceInput[opslevel.PropertyInput]()
		cobra.CheckErr(err)
		newProperty, err := getClientGQL().PropertyAssign(*input)
		cobra.CheckErr(err)

		fmt.Println(newProperty.Definition.Id)
	},
}

var unassignPropertyCmd = &cobra.Command{
	Use:        "property",
	Short:      "Unassign a Property",
	Long:       `Unassign a Property from an Owner by Id or Alias`,
	Example:    `opslevel unassign property owner-alias property-id`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"OWNER_ID", "PROPERTY_ID"},
	Run: func(cmd *cobra.Command, args []string) {
		ownerId := args[0]
		propertyId := args[1]

		err := getClientGQL().PropertyUnassign(ownerId, propertyId)
		cobra.CheckErr(err)

		fmt.Printf("unassigned property '%s' from '%s'\n", propertyId, ownerId)
	},
}

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

var updatePropertyDefinitonCmd = &cobra.Command{
	Use:   "property-definition",
	Short: "Update a property-definition",
	Long:  `Update a property-definition`,
	Example: fmt.Sprintf(`
cat << EOF | opslevel update property-definition propdef3 -f -
%s
EOF`, getYaml[opslevel.PropertyDefinitionInput]()),
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		identifier := args[0]
		input, err := readPropertyDefinitionInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().UpdatePropertyDefinition(identifier, *input)
		cobra.CheckErr(err)

		if isYamlOutput() {
			common.YamlPrint(result)
		} else {
			common.PrettyPrint(result)
		}
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

	if description, ok := data["description"].(string); ok {
		propDefInput.Description = description
	}
	if propertyDisplayStatus, ok := data["propertyDisplayStatus"].(string); ok {
		propDefInput.PropertyDisplayStatus = opslevel.PropertyDisplayStatusEnum(propertyDisplayStatus)
	}

	return &propDefInput, nil
}

var getPropertyDefinitionCmd = &cobra.Command{
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
	// Property Commands
	exampleCmd.AddCommand(examplePropertyCmd)
	assignCmd.AddCommand(assignPropertyCmd)
	unassignCmd.AddCommand(unassignPropertyCmd)

	// Property Definition Commands
	exampleCmd.AddCommand(examplePropertyDefinitionCmd)
	createCmd.AddCommand(createPropertyDefinitonCmd)
	updateCmd.AddCommand(updatePropertyDefinitonCmd)
	getCmd.AddCommand(getPropertyDefinitionCmd)
	listCmd.AddCommand(listPropertyDefinitionsCmd)
	deleteCmd.AddCommand(deletePropertyDefinitonCmd)
}
