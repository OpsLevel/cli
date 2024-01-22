package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var actionType string

var exampleActionCmd = &cobra.Command{
	Use:   "action",
	Short: "Example action",
	Long:  `Example action`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.CustomActionsWebhookActionCreateInput]())
	},
}

var createActionCmd = &cobra.Command{
	Use:   "action --type=$ACTION_TYPE",
	Short: "Create an action",
	Long:  "Create an action",
	Example: `
cat << EOF | opslevel create action --type=webhook -f -
name: "Page The On Call"
description: "Pages the On Call"
webhookUrl: "https://api.pagerduty.com/incidents"
httpMethod: "POST"
headers:
  accept: "application/vnd.pagerduty+json;version=2"
  authorization: "Token token=XXXXXXXXXXXXX"
  from: "someone@example.com"
liquidTemplate: |
  {
      "incident": {
          "type": "incident",
          "title": "{{manualInputs.IncidentTitle}}",
          "service": {
            "id": "{{ service | tag_value: 'pd_id' }}",
            "type": "service_reference"
          },
          "body": {
            "type": "incident_body",
            "details": "Incident triggered from OpsLevel by {{user.name}} with the email {{user.email}}. {{manualInputs.IncidentDescription}}"
          }
      }
  }
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		switch actionType {
		case "webhook":
			input, err := readResourceInput[opslevel.CustomActionsWebhookActionCreateInput]()
			cobra.CheckErr(err)
			result, err := getClientGQL().CreateWebhookAction(*input)
			cobra.CheckErr(err)
			fmt.Printf("created webhook action: %s\n", result.Id)
		default:
			err := errors.New("unknown action type: " + actionType)
			cobra.CheckErr(err)
		}
	},
}

var getActionCmd = &cobra.Command{
	Use:        "action ID|ALIAS",
	Short:      "Get details about an action",
	Long:       "Get details about an action",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		action, err := getClientGQL().GetCustomAction(key)
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
			w := common.NewTabWriter("ID", "NAME", "HTTP_METHOD", "WEBHOOK_URL")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", item.Id, item.Name, item.HTTPMethod, item.WebhookURL)
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
description: "Pages the oncall and creates an incident"
headers:
  accept: "application/vnd.pagerduty+json;version=2"
  authorization: "Token token=XXXXXXXXXXXXX"
  from: "someone-else@example.com"
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		switch actionType {
		case "webhook":
			input, err := readResourceInput[opslevel.CustomActionsWebhookActionUpdateInput]()
			input.Id = *opslevel.NewID(key)
			cobra.CheckErr(err)
			action, err := getClientGQL().UpdateWebhookAction(*input)
			cobra.CheckErr(err)
			common.JsonPrint(json.MarshalIndent(action, "", "    "))
		default:
			err := errors.New("unknown action type: " + actionType)
			cobra.CheckErr(err)
		}
	},
}

var deleteActionCmd = &cobra.Command{
	Use:        "action ID|ALIAS",
	Short:      "Delete an action",
	Long:       "Delete an action",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := getClientGQL().DeleteWebhookAction(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted webhook action: %s\n", key)
	},
}

func init() {
	exampleCmd.AddCommand(exampleActionCmd)
	createCmd.AddCommand(createActionCmd)
	createActionCmd.Flags().StringVar(&actionType, "type", "webhook", "action type, default=webhook")
	updateCmd.AddCommand(updateActionCmd)
	updateActionCmd.Flags().StringVar(&actionType, "type", "webhook", "action type, default=webhook")
	getCmd.AddCommand(getActionCmd)
	listCmd.AddCommand(listActionCmd)
	deleteCmd.AddCommand(deleteActionCmd)
}
