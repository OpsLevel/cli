package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/opslevel/opslevel-go/v2023"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// TODO: examples for create, update

var actionType string

var createActionCmd = &cobra.Command{
	Use:   "action --type=$ACTION_TYPE",
	Short: "Create an action",
	Long:  "Create an action",
	Example: `
cat << EOF | opslevel create action --type=webhook -f -
...
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		switch actionType {
		case "webhook":
			input, err := readWebhookActionInput()
			cobra.CheckErr(err)
			result, err := getClientGQL().CreateWebhookAction(*input)
			cobra.CheckErr(err)
			fmt.Printf("created webhook action: %s\n", result.Id)
		default:
			err := errors.New("unknown action type: " + args[0])
			cobra.CheckErr(err)
		}
	},
}

var getActionCmd = &cobra.Command{
	Use:        "action ID",
	Short:      "Get details about an action",
	Long:       "Get details about an action",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		action, err := getClientGQL().GetCustomAction(*opslevel.NewIdentifier(key))
		cobra.CheckErr(err)
		common.PrettyPrint(action)
	},
}

var listActionCmd = &cobra.Command{
	Use:     "action",
	Aliases: []string{"actions"},
	Short:   "List actions",
	Long:    "List actions",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListCustomActions(nil)
		cobra.CheckErr(err)
		list := resp.Nodes
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("ID", "NAME")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Id, item.Name)
			}
			w.Flush()
		}
	},
}

var updateActionCmd = &cobra.Command{
	Use:   "action ID",
	Short: "Update an action",
	Long:  "Update an action",
	Example: `
cat << EOF | opslevel update action --type=webhook $ACTION_ID -f -
...
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		switch actionType {
		case "webhook":
			input, err := readWebhookActionInput()
			cobra.CheckErr(err)
			updateInput := &opslevel.CustomActionsWebhookActionUpdateInput{
				Id:             *opslevel.NewID(key),
				Name:           &input.Name,
				Description:    input.Description,
				LiquidTemplate: &input.LiquidTemplate,
				WebhookURL:     &input.WebhookURL,
				HTTPMethod:     input.HTTPMethod,
				Headers:        &input.Headers,
			}
			action, err := getClientGQL().UpdateWebhookAction(*updateInput)
			cobra.CheckErr(err)
			common.JsonPrint(json.MarshalIndent(action, "", "    "))
		default:
			err := errors.New("unknown action type: " + actionType)
			cobra.CheckErr(err)
		}
	},
}

var deleteActionCmd = &cobra.Command{
	Use:        "action ID",
	Short:      "Delete an action",
	Long:       "Delete an action",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteWebhookAction(*opslevel.NewIdentifier(key))
		cobra.CheckErr(err)
		fmt.Printf("deleted webhook action: %s\n", key)
	},
}

func init() {
	createCmd.AddCommand(createActionCmd)
	createActionCmd.Flags().StringVar(&actionType, "type", "webhook", "action type, default=webhook")
	updateCmd.AddCommand(updateActionCmd)
	updateActionCmd.Flags().StringVar(&actionType, "type", "webhook", "action type, default=webhook")
	getCmd.AddCommand(getActionCmd)
	listCmd.AddCommand(listActionCmd)
	deleteCmd.AddCommand(deleteActionCmd)
}

func readWebhookActionInput() (*opslevel.CustomActionsWebhookActionCreateInput, error) {
	readInputConfig()
	evt := &opslevel.CustomActionsWebhookActionCreateInput{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
