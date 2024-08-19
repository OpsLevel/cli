package cmd

import (
	"fmt"

	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/spf13/cobra"
)

var createServiceToolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Create service tool",
	Args:  cobra.ExactArgs(1),
	Example: `
cat << EOF | opslevel create service tool my-service-alias -f -
category: deployment
displayName: "fancy tool"
environment: "dev"
url: "https://example.com"
EOF
`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceAlias := args[0]
		toolFromInput, err := readResourceInput[opslevel.ToolCreateInput]()
		cobra.CheckErr(err)

		serviceId, err := getService(serviceAlias)
		cobra.CheckErr(err)

		toolFromInput.ServiceId = &serviceId.Id
		tool, err := getClientGQL().CreateTool(*toolFromInput)
		cobra.CheckErr(err)
		common.PrettyPrint(string(tool.Id))
	},
}

var updateServiceToolCmd = &cobra.Command{
	Use:        "tool TOOL-ID",
	ArgAliases: []string{"TOOL-ID"},
	Short:      "Update service tool",
	Args:       cobra.ExactArgs(1),
	Example: `
cat << EOF | opslevel update service tool tool-ID -f -
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
		toolFromInput, err := readResourceInput[opslevel.ToolUpdateInput]()
		cobra.CheckErr(err)

		toolFromInput.Id = opslevel.ID(id)
		updatedTool, err := getClientGQL().UpdateTool(*toolFromInput)
		cobra.CheckErr(err)
		common.PrettyPrint(string(updatedTool.Id))
	},
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
