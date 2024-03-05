package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var exampleTriggerDefinitionCmd = &cobra.Command{
	Use:     "trigger-definition",
	Aliases: []string{"triggerdefinition", "trigdef", "td"},
	Short:   "Example Scorecard",
	Long:    `Example Scorecard`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.CustomActionsTriggerDefinitionCreateInput]())
	},
}

var createTriggerDefinitionCmd = &cobra.Command{
	Use:     "trigger-definition",
	Aliases: []string{"triggerdefinition", "td"},
	Short:   "Create a trigger definition",
	Long:    "Create a trigger definition",
	Example: `
cat << EOF | opslevel create trigger-definition -f -
name: "Page The On Call"
description: "Pages the On Call"
owner: "some_team"
action: "some_action"
accessControl: "everyone"
manualInputsDefinition: |
  version: 1
  inputs:
    - identifier: IncidentTitle
      displayName: Title
      description: Title of the incident to trigger
      type: text_input
      required: true
      maxLength: 60
      defaultValue: Service Incident Manual Trigger
    - identifier: IncidentDescription
      displayName: Incident Description
      description: The description of the incident
      type: text_area
      required: true
responseTemplate: |
  {% if response.status >= 200 and response.status < 300 %}
  success
  {% else %}
  failure
  {% endif %}
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := ReadResourceInput[opslevel.CustomActionsTriggerDefinitionCreateInput](nil)
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateTriggerDefinition(*input)
		cobra.CheckErr(err)
		fmt.Printf("created trigger definition: %s\n", result.Id)
	},
}

var getTriggerDefinitionCmd = &cobra.Command{
	Use:        "trigger-definition ID|ALIAS",
	Aliases:    []string{"triggerdefinition", "trigdef", "td"},
	Short:      "Get details about a trigger definition",
	Long:       "Get details about a trigger definition",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		triggerDefinition, err := getClientGQL().GetTriggerDefinition(key)
		cobra.CheckErr(err)
		common.PrettyPrint(triggerDefinition)
	},
}

var listTriggerDefinitionCmd = &cobra.Command{
	Use:     "trigger-definition",
	Aliases: []string{"trigger-definitions", "triggerdefinition", "triggerdefinitions", "trigdef", "trigdefs", "td", "tds"},
	Short:   "List trigger definitions",
	Long:    "List trigger definitions",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListTriggerDefinitions(nil)
		cobra.CheckErr(err)
		list := resp.Nodes
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("ID", "NAME", "OWNER")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\n", item.Id, item.Name, item.Owner.Alias)
			}
			w.Flush()
		}
	},
}

var updateTriggerDefinitionCmd = &cobra.Command{
	Use:     "trigger-definition ID",
	Aliases: []string{"triggerdefinition", "trigdef", "td"},
	Short:   "Update a trigger definition",
	Long:    "Update a trigger definition",
	Example: `
cat << EOF | opslevel update trigger-definition $TRIGGER_ID -f -
description: "Pages the On Call via PagerDuty"
accessControl: "service_owners"
extendedTeamAccess:
- alias: "team_alias_1"
- id: "XXX_team_id_XXX"
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := ReadResourceInput[opslevel.CustomActionsTriggerDefinitionUpdateInput](nil)
		input.Id = *opslevel.NewID(key)
		cobra.CheckErr(err)
		triggerDefinition, err := getClientGQL().UpdateTriggerDefinition(*input)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(triggerDefinition, "", "    "))
	},
}

var deleteTriggerDefinitionCmd = &cobra.Command{
	Use:        "trigger-definition ID|ALIAS",
	Aliases:    []string{"triggerdefinition", "trigdef", "td"},
	Short:      "Delete a trigger definition",
	Long:       "Delete a trigger definition",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteTriggerDefinition(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted trigger definition: %s\n", key)
	},
}

func init() {
	exampleCmd.AddCommand(exampleTriggerDefinitionCmd)
	createCmd.AddCommand(createTriggerDefinitionCmd)
	updateCmd.AddCommand(updateTriggerDefinitionCmd)
	getCmd.AddCommand(getTriggerDefinitionCmd)
	listCmd.AddCommand(listTriggerDefinitionCmd)
	deleteCmd.AddCommand(deleteTriggerDefinitionCmd)
}
