package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/opslevel/cli/common"
	"github.com/spf13/cobra"
)

var exampleScorecardCmd = &cobra.Command{
	Use:     "scorecard",
	Aliases: []string{"sc"},
	Short:   "Example Scorecard",
	Long:    `Example Scorecard`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(getExample[opslevel.ScorecardInput]())
	},
}

var createScorecardCmd = &cobra.Command{
	Use:     "scorecard",
	Aliases: []string{"sc"},
	Short:   "Create a scorecard",
	Long:    "Create a scorecard",
	Example: `
cat << EOF | opslevel create scorecard -f -
name: "new scorecard"
description: "a newly created scorecard"
ownerId: "XXX_team_id_XXX"
affectsOverallServiceLevels: false
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readResourceInput[opslevel.ScorecardInput]()
		cobra.CheckErr(err)
		result, err := getClientGQL().CreateScorecard(*input)
		cobra.CheckErr(err)
		fmt.Printf("created scorecard: %s\n", result.Id)
	},
}

var getScorecardCmd = &cobra.Command{
	Use:        "scorecard ID|ALIAS",
	Aliases:    []string{"sc"},
	Short:      "Get details about a scorecard",
	Long:       "Get details about a scorecard",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		scorecard, err := getClientGQL().GetScorecard(key)
		cobra.CheckErr(err)
		common.PrettyPrint(scorecard)
	},
}

var listScorecardCmd = &cobra.Command{
	Use:     "scorecard",
	Aliases: []string{"scorecards", "sc", "scs"},
	Short:   "List scorecards",
	Long:    "List scorecards",
	Example: `
opslevel list scorecards -o json | jq 'map( {(.Name): (.ServiceCount)} )'
`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListScorecards(nil)
		cobra.CheckErr(err)
		list := resp.Nodes
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("ID", "NAME", "PASSING_CHECKS", "CHECKS_COUNT", "SERVICE_COUNT")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\n", item.Id, item.Name, item.PassingChecks, item.ChecksCount, item.ServiceCount)
			}
			w.Flush()
		}
	},
}

var updateScorecardCmd = &cobra.Command{
	Use:     "scorecard ID|ALIAS",
	Aliases: []string{"sc"},
	Short:   "Update a scorecard",
	Long:    "Update a scorecard",
	Example: `
cat << EOF | opslevel update scorecard $ID -f -
name: "updated scorecard"
description: "an updated scorecard"
ownerId: "XXX_team_id_XXX"
affectsOverallServiceLevels: true
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		input, err := readResourceInput[opslevel.ScorecardInput]()
		cobra.CheckErr(err)
		scorecard, err := getClientGQL().UpdateScorecard(key, *input)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(scorecard, "", "    "))
	},
}

var deleteScorecardCmd = &cobra.Command{
	Use:        "scorecard ID|ALIAS",
	Aliases:    []string{"sc"},
	Short:      "Delete a scorecard",
	Long:       "Delete a scorecard",
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID", "ALIAS"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		deletedScorecardId, err := getClientGQL().DeleteScorecard(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted scorecard: %s\n", *deletedScorecardId)
	},
}

func init() {
	exampleCmd.AddCommand(exampleScorecardCmd)
	createCmd.AddCommand(createScorecardCmd)
	updateCmd.AddCommand(updateScorecardCmd)
	getCmd.AddCommand(getScorecardCmd)
	listCmd.AddCommand(listScorecardCmd)
	deleteCmd.AddCommand(deleteScorecardCmd)
}
