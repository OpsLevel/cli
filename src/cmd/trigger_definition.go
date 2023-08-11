package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: examples for create, update

var createTriggerDefinitionCmd = &cobra.Command{
	Use:   "trigger-definition",
	Short: "Create a trigger definition",
	Long:  "Create a trigger definition",
	Example: `
cat << EOF | opslevel create trigger-definition -f -
...
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readTriggerDefinitionCreateInput()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateTriggerDefinition(*input)
		cobra.CheckErr(err)
		fmt.Printf("created trigger definition: %s\n", result.Id)
	},
}

var getTriggerDefinitionCmd = &cobra.Command{
	Use:        "trigger-definition ID",
	Short:      "Get details about a trigger definition",
	Long:       "Get details about a trigger definition",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		triggerDefinition, err := getClientGQL().GetTriggerDefinition(*opslevel.NewIdentifier(key))
		cobra.CheckErr(err)
		common.PrettyPrint(triggerDefinition)
	},
}

var listTriggerDefinitionCmd = &cobra.Command{
	Use:     "trigger-definition",
	Aliases: []string{"trigger-definitions"},
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
	Use:   "trigger-definition ID",
	Short: "Update a trigger definition",
	Long:  "Update a trigger definition",
	Example: `
cat << EOF | opslevel update trigger-definition $TRIGGER_ID -f -
...
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readTriggerDefinitionUpdateInput()
		input.Id = *opslevel.NewID(key)
		cobra.CheckErr(err)
		triggerDefinition, err := getClientGQL().UpdateTriggerDefinition(*input)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(triggerDefinition, "", "    "))
	},
}

var deleteTriggerDefinitionCmd = &cobra.Command{
	Use:        "trigger-definition ID",
	Short:      "Delete a trigger definition",
	Long:       "Delete a trigger definition",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteTriggerDefinition(*opslevel.NewIdentifier(key))
		cobra.CheckErr(err)
		fmt.Printf("deleted trigger definition: %s\n", key)
	},
}

func init() {
	createCmd.AddCommand(createTriggerDefinitionCmd)
	updateCmd.AddCommand(updateTriggerDefinitionCmd)
	getCmd.AddCommand(getTriggerDefinitionCmd)
	listCmd.AddCommand(listTriggerDefinitionCmd)
	deleteCmd.AddCommand(deleteTriggerDefinitionCmd)
}

func readTriggerDefinitionCreateInput() (*opslevel.CustomActionsTriggerDefinitionCreateInput, error) {
	readInputConfig()
	evt := &opslevel.CustomActionsTriggerDefinitionCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}

func readTriggerDefinitionUpdateInput() (*opslevel.CustomActionsTriggerDefinitionUpdateInput, error) {
	readInputConfig()
	evt := &opslevel.CustomActionsTriggerDefinitionUpdateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
