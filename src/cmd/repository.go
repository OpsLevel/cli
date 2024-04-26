package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var getRepositoryCmd = &cobra.Command{
	Use:        "repository ID|ALIAS",
	Aliases:    []string{"repo"},
	Short:      "Get details about a repository",
	Long:       `Get details about a repository`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		var repository *opslevel.Repository
		var err error
		if opslevel.IsID(key) {
			repository, err = getClientGQL().GetRepository(opslevel.ID(key))
			cobra.CheckErr(err)
		} else {
			repository, err = getClientGQL().GetRepositoryWithAlias(key)
			cobra.CheckErr(err)
		}
		cobra.CheckErr(err)
		common.PrettyPrint(repository)
	},
}

var listRepositoryCmd = &cobra.Command{
	Use:     "repository",
	Aliases: []string{"repositories", "repo", "repos"},
	Short:   "Lists repositories",
	Long:    `Lists repositories`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClientGQL()
		payloadVars := client.InitialPageVariables()

		hiddenOnly, err := cmd.Flags().GetBool("hidden-only")
		cobra.CheckErr(err)
		// visible set to false means "hidden" in the API
		payloadVars["visible"] = !hiddenOnly

		resp, err := client.ListRepositories(&payloadVars)
		list := resp.Nodes
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "ALIAS", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.DefaultAlias, item.Id)
			}
			w.Flush()
		}
	},
}

func init() {
	getCmd.AddCommand(getRepositoryCmd)
	listRepositoryCmd.Flags().Bool("hidden-only", false, "list only hidden repositories when set")
	listCmd.AddCommand(listRepositoryCmd)
}
