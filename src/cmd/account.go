package cmd

import (
	"fmt"

	"github.com/opslevel/cli/common"

	"github.com/opslevel/opslevel-go"
	"github.com/spf13/cobra"
)

var lifecycleCmd = &cobra.Command{
	Use:   "lifecycle",
	Short: "Lists the valid alias for lifecycles in your account",
	Long:  `Lists the valid alias for lifecycles in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		list, err := client.ListLifecycles()
		if err == nil {
			for _, item := range list {
				fmt.Println(item.Alias)
			}
		}
	},
}

var tierCmd = &cobra.Command{
	Use:   "tier",
	Short: "Lists the valid alias for tiers in your account",
	Long:  `Lists the valid alias for tiers in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		list, err := client.ListTiers()
		if err == nil {
			for _, item := range list {
				fmt.Println(item.Alias)
			}
		}
	},
}

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Lists teams in your account",
	Long:  `Lists teams in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		client := common.NewGraphClient()
		list, err := client.ListTeams()
		cobra.CheckErr(err)
		w := common.NewTabWriter("Name", "ID", "Alias")
		if err == nil {
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Id, item.Alias)
			}
		}
		w.Flush()
	},
}

var toolsCmd = &cobra.Command{
	Use:   "tool",
	Short: "Lists the valid alias for tools in your account",
	Long:  `Lists the valid alias for tools in your account`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(opslevel.ToolCategoryAdmin)
		fmt.Println(opslevel.ToolCategoryCode)
		fmt.Println(opslevel.ToolCategoryContinuousIntegration)
		fmt.Println(opslevel.ToolCategoryDeployment)
		fmt.Println(opslevel.ToolCategoryErrors)
		fmt.Println(opslevel.ToolCategoryFeatureFlag)
		fmt.Println(opslevel.ToolCategoryHealthChecks)
		fmt.Println(opslevel.ToolCategoryIncidents)
		fmt.Println(opslevel.ToolCategoryLogs)
		fmt.Println(opslevel.ToolCategoryMetrics)
		fmt.Println(opslevel.ToolCategoryOrchestrator)
		fmt.Println(opslevel.ToolCategoryRunbooks)
		fmt.Println(opslevel.ToolCategoryStatusPage)
		fmt.Println(opslevel.ToolCategoryWiki)
		fmt.Println(opslevel.ToolCategoryOther)
	},
}

func init() {
	listCmd.AddCommand(lifecycleCmd)
	listCmd.AddCommand(tierCmd)
	listCmd.AddCommand(teamCmd)
	listCmd.AddCommand(toolsCmd)
}
