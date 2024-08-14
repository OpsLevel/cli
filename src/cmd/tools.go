package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/spf13/cobra"
)

var createServiceToolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Create service tool",
	Example: `
cat << EOF | opslevel create service tool -f -
service: my-service-alias
category: deployment
displayName: "fancy tool"
environment: "dev"
url: "https://example.com"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClientGQL()
		serviceIdentifier := readServiceFieldFromYaml()
		toolFromInput, err := readResourceInput[opslevel.ToolCreateInput]()
		cobra.CheckErr(err)

		err = verifyToolCategory(&toolFromInput.Category)
		cobra.CheckErr(err)

		service, err := getService(serviceIdentifier)
		cobra.CheckErr(err)

		toolFromInput.ServiceId = &service.Id
		tool, err := client.CreateTool(*toolFromInput)
		cobra.CheckErr(err)
		common.PrettyPrint(string(tool.Id))
	},
}

var updateServiceToolCmd = &cobra.Command{
	Use:     "tool TOOL-ID",
	ArgAliases: []string{"TOOL-ID"},
	Short:   "Update service tool",
	Args:    cobra.ExactArgs(1),
	Example: `
cat << EOF | opslevel update service tool tool-ID -f -
service: my-service-alias
category: deployment
displayName: "fancy tool"
environment: "dev"
url: "https://example.com"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if !opslevel.IsID(id) {
			cobra.CheckErr(fmt.Errorf("invalid ID: '%s'", id))
		}
		toolId := opslevel.ID(id)

		serviceIdentifier := readServiceFieldFromYaml()
		toolFromInput, err := readResourceInput[opslevel.ToolUpdateInput]()
		cobra.CheckErr(err)

		err = verifyToolCategory(toolFromInput.Category)
		cobra.CheckErr(err)

		client := getClientGQL()
		service, err := getService(serviceIdentifier)
		cobra.CheckErr(err)

		if !isToolIdInServiceTools(toolId, service.Tools) {
			cobra.CheckErr(fmt.Errorf("no tool with ID '%s' to update on service with identifier '%s'", toolId, serviceIdentifier))
		}
		toolFromInput.Id = toolId
		_, err = client.UpdateTool(*toolFromInput)
		cobra.CheckErr(err)
		common.PrettyPrint(toolId)
	},
}

func isToolIdInServiceTools(toolId opslevel.ID, serviceTools *opslevel.ToolConnection) bool {
	if serviceTools == nil {
		return false
	}
	for _, serviceTool := range serviceTools.Nodes {
		if toolId == serviceTool.Id {
			return true
		}
	}
	return false
}

func verifyToolCategory(toolCategory *opslevel.ToolCategory) error {
	if toolCategory == nil || slices.Contains(opslevel.AllToolCategory, string(*toolCategory)) {
		return nil
	}
	return fmt.Errorf(
		"'category:' field has value '%s' but must be one of '%v'",
		string(*toolCategory),
		strings.Join(opslevel.AllToolCategory, ", "),
	)
}

var deleteServiceToolCmd = &cobra.Command{
	Use:        "tool TOOL-ID",
	ArgAliases: []string{"TOOL-ID"},
	Short:      "Delete a service tool",
	Example:    `opslevel delete service tool <TOOL-ID>`,
	Args:       cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		if !opslevel.IsID(id) {
			cobra.CheckErr(fmt.Errorf("invalid ID: '%s'", id))
		}
		toolId := opslevel.NewID(id)
		err := getClientGQL().DeleteTool(*toolId)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' service tool\n", *toolId)
	},
}

func init() {
	createServiceCmd.AddCommand(createServiceToolCmd)
	updateServiceCmd.AddCommand(updateServiceToolCmd)
	deleteServiceCmd.AddCommand(deleteServiceToolCmd)
}
