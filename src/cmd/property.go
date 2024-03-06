package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/opslevel/cli/common"

	"github.com/spf13/cobra"
)

var examplePropertyCmd = &cobra.Command{
	Use:     "property",
	Aliases: []string{"prop"},
	Short:   "Example Property",
	Long:    `Example Property`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.PropertyInput]())
	},
}

var getPropertyCmd = &cobra.Command{
	Use:        "property",
	Aliases:    []string{"prop"},
	Short:      "Get details about an assigned property",
	Long:       `Get details about an assigned property`,
	Example:    `opslevel get property owner-alias property-id`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		ownerId := args[0]
		propertyId := args[1]

		result, err := getClientGQL().GetProperty(ownerId, propertyId)
		cobra.CheckErr(err)
		if result.Definition.Id == "" && result.Owner.Id() == "" {
			err = fmt.Errorf("property '%s' on entity '%s' not found\n", propertyId, ownerId)
			cobra.CheckErr(err)
		}

		if isYamlOutput() {
			common.YamlPrint(result)
		} else {
			common.PrettyPrint(result)
		}
	},
}

var listPropertyCmd = &cobra.Command{
	Use:        "property",
	Short:      "List properties on a Service",
	Aliases:    []string{"properties", "prop", "props"},
	Long:       "List properties on a Service identified by ID or Alias",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"SERVICE_ID", "SERVICE_ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		var service *opslevel.Service
		var err error
		if opslevel.IsID(args[0]) {
			service, err = getClientGQL().GetService(*opslevel.NewID(args[0]))
		} else {
			service, err = getClientGQL().GetServiceWithAlias(args[0])
		}
		cobra.CheckErr(err)
		properties, err := service.GetProperties(getClientGQL(), nil)
		cobra.CheckErr(err)

		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(properties.Nodes, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			w.Write([]string{"ID", "LOCKED", "ALIASES"})
			for _, p := range properties.Nodes {
				aliases := strings.Join(p.Definition.Aliases, "/")
				w.Write([]string{string(p.Definition.Id), fmt.Sprintf("%t", p.Locked), aliases})
			}
			w.Flush()
		} else {
			w := common.NewTabWriter("ID", "LOCKED", "VALUE", "ALIASES")
			for _, p := range properties.Nodes {
				var valueOutput string
				if p.Value != nil {
					valueOutput = string(*p.Value)
				}
				aliases := strings.Join(p.Definition.Aliases, ",")
				fmt.Fprintf(w, "%s\t%t\t%s\t%s\n", string(p.Definition.Id), p.Locked, valueOutput, aliases)
			}
			w.Flush()
		}
	},
}

var assignPropertyCmd = &cobra.Command{
	Use:     "property",
	Aliases: []string{"prop"},
	Short:   "Assign a Property",
	Long:    `Assign a Property to an Entity by Id or Alias`,
	Example: fmt.Sprintf(`
cat << EOF | opslevel assign property -f -
%s
EOF`, getYaml[opslevel.PropertyInput]()),
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readResourceInput[opslevel.PropertyInput]()
		cobra.CheckErr(err)
		newProperty, err := getClientGQL().PropertyAssign(*input)
		cobra.CheckErr(err)

		fmt.Printf("assigned property '%s' on '%s'\n", newProperty.Definition.Id, newProperty.Owner.Id())
	},
}

var unassignPropertyCmd = &cobra.Command{
	Use:        "property",
	Aliases:    []string{"prop"},
	Short:      "Unassign a Property",
	Long:       `Unassign a Property from an Entity by Id or Alias`,
	Example:    `opslevel unassign property owner-alias property-id`,
	Args:       cobra.ExactArgs(2),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		ownerId := args[0]
		propertyId := args[1]

		err := getClientGQL().PropertyUnassign(ownerId, propertyId)
		cobra.CheckErr(err)

		fmt.Printf("unassigned property '%s' from '%s'\n", propertyId, ownerId)
	},
}

var examplePropertyDefinitionCmd = &cobra.Command{
	Use:     "property-definition",
	Aliases: []string{"propertydefinition", "propdef", "pd"},
	Short:   "Example Property Definition",
	Long:    `Example Property Definition`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.PropertyDefinitionInput]())
	},
}

var createPropertyDefinitionCmd = &cobra.Command{
	Use:     "property-definition",
	Aliases: []string{"propertydefinition", "propdef", "pd"},
	Short:   "Create a property-definition",
	Long:    `Create a property-definition`,
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

var updatePropertyDefinitionCmd = &cobra.Command{
	Use:     "property-definition",
	Aliases: []string{"propertydefinition", "propdef", "pd"},
	Short:   "Update a property-definition",
	Long:    `Update a property-definition`,
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
	d, err := readResourceInput[opslevel.JSONSchema]()
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
	jsonSchema := opslevel.JSONSchema(schema)
	propDefInput := opslevel.PropertyDefinitionInput{
		Name:   opslevel.RefOf(name),
		Schema: &jsonSchema,
	}

	if description, ok := data["description"].(string); ok {
		propDefInput.Description = opslevel.RefOf(description)
	}
	if propertyDisplayStatus, ok := data["propertyDisplayStatus"].(string); ok {
		propDefInput.PropertyDisplayStatus = opslevel.RefOf(opslevel.PropertyDisplayStatusEnum(propertyDisplayStatus))
	}

	return &propDefInput, nil
}

var getPropertyDefinitionCmd = &cobra.Command{
	Use:        "property-definition",
	Aliases:    []string{"propertydefinition", "propdef", "pd"},
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
	Aliases: []string{"property-definitions", "propertydefinition", "propertydefinitions", "propdef", "propdefs", "pd", "pds"},
	Short:   "List property definitions",
	Long:    "List property definitions",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListPropertyDefinitions(nil)
		list := resp.Nodes

		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else if isCsvOutput() {
			w := csv.NewWriter(os.Stdout)
			w.Write([]string{"ALIASES", "ID", "NAME", "DISPLAY_STATUS", "ALLOWED_IN_CONFIG_FILES"})
			for _, d := range list {
				aliases := strings.Join(d.Aliases, "/")
				w.Write([]string{aliases, string(d.Id), d.Name, string(d.PropertyDisplayStatus), fmt.Sprintf("%t", d.AllowedInConfigFiles)})
			}
			w.Flush()
		} else {
			w := common.NewTabWriter("ALIASES", "ID", "NAME", "DISPLAY_STATUS", "ALLOWED_IN_CONFIG_FILES")
			for _, d := range list {
				aliases := strings.Join(d.Aliases, "/")
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\n", aliases, d.Id, d.Name, d.PropertyDisplayStatus, d.AllowedInConfigFiles)
			}
			w.Flush()
		}
	},
}

var deletePropertyDefinitionCmd = &cobra.Command{
	Use:        "property-definition ID",
	Aliases:    []string{"propertydefinition", "propdef", "pd"},
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
	getCmd.AddCommand(getPropertyCmd)
	listCmd.AddCommand(listPropertyCmd)

	// Property Definition Commands
	exampleCmd.AddCommand(examplePropertyDefinitionCmd)
	createCmd.AddCommand(createPropertyDefinitionCmd)
	updateCmd.AddCommand(updatePropertyDefinitionCmd)
	getCmd.AddCommand(getPropertyDefinitionCmd)
	listCmd.AddCommand(listPropertyDefinitionsCmd)
	deleteCmd.AddCommand(deletePropertyDefinitionCmd)
}
