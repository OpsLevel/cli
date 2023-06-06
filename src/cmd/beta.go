package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var betaCmd = &cobra.Command{
	Use:   "beta",
	Short: "Beta features",
}

var syncOpsLevelCmd = &cobra.Command{
	Use:   "sync opslevel",
	Short: "Trigger a sync of opslevel.yml for all services",
	Run: func(cmd *cobra.Command, args []string) {
		client := getClientGQL()
		resp, err := client.ListServices(nil)
		cobra.CheckErr(err)
		services := resp.Nodes
		for _, service := range services {
			edgesCount := len(service.Repositories.Edges)
			if edgesCount == 0 {
				continue
			}
			if edgesCount == 1 {
				client.SyncOpsLevelYml(string(service.Id), service.Repositories.Edges[0].Node.DefaultAlias)
				continue
			}
			// Prompt User
			templates := &promptui.SelectTemplates{
				Label:    "{{ .Node.DefaultAlias }}?",
				Active:   fmt.Sprintf("%s {{ .Node.DefaultAlias | cyan }}", promptui.IconSelect),
				Inactive: "    {{ .Node.DefaultAlias | cyan }}",
				Selected: fmt.Sprintf("%s {{ .Node.DefaultAlias | faint }}", promptui.IconGood),
			}

			prompt := promptui.Select{
				Label:     "Select Repository",
				Items:     service.Repositories.Edges,
				Templates: templates,
				Size:      common.MinInt(6, edgesCount),
			}

			index, _, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				continue
			}
			client.SyncOpsLevelYml(string(service.Id), service.Repositories.Edges[index].Node.DefaultAlias)
		}
	},
}

func init() {
	rootCmd.AddCommand(betaCmd)
	betaCmd.AddCommand(syncOpsLevelCmd)
}
