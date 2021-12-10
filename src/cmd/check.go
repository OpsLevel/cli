package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var getCheckCmd = &cobra.Command{
	Use:        "check ID",
	Short:      "Get details about a rubic check",
	Long:       `Get details about a rubic check`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		check, err := graphqlClient.GetCheck(args[0])
		cobra.CheckErr(err)
		common.PrettyPrint(check)
	},
}

var listCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"checks"},
	Short:   "Lists the rubric checks",
	Long:    `Lists the rubric checks`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := graphqlClient.ListChecks()
		cobra.CheckErr(err)
		if isJsonOutput() {
			common.JsonPrint(json.MarshalIndent(list, "", "    "))
		} else {
			w := common.NewTabWriter("NAME", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t\n", item.Name, item.Id)
			}
			w.Flush()
		}
	},
}

var checkCreateCmd = &cobra.Command{
	Use:   "check",
	Short: "Create a rubric check",
	Long: `Create a rubric check

Examples:

	opslevel create check -f my_cec.yaml
`,
	Run: func(cmd *cobra.Command, args []string) {
		common.AliasCache.CacheCategories(graphqlClient)
		common.AliasCache.CacheLevels(graphqlClient)
		common.AliasCache.CacheTeams(graphqlClient)
		common.AliasCache.CacheFilters(graphqlClient)
		common.AliasCache.CacheIntegrations(graphqlClient)
		input, err := readCheckCreateInput()
		cobra.CheckErr(err)
		switch input.Kind {
		case opslevel.CheckTypeCustom:
			check, err := graphqlClient.CreateCheckCustomEvent(*input.AsCustomEventCreate())
			cobra.CheckErr(err)
			fmt.Printf("Created: %s - %s\n", check.Name, check.Id)
		}
	},
}

var deleteCheckCmd = &cobra.Command{
	Use:        "check ID",
	Short:      "Delete a rubric check",
	Long:       `Delete a rubric check`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		err := graphqlClient.DeleteCheck(key)
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' check\n", key)
	},
}

func init() {
	createCmd.AddCommand(checkCreateCmd)
	getCmd.AddCommand(getCheckCmd)
	listCmd.AddCommand(listCheckCmd)
	deleteCmd.AddCommand(deleteCheckCmd)
}

type CheckCreateType struct {
	Kind opslevel.CheckType
	Spec map[string]interface{}
}

func (self *CheckCreateType) resolveAliases() {
	if item, ok := self.Spec["category"]; ok {
		if value, ok := common.AliasCache.TryGetCategory(item.(string)); ok {
			delete(self.Spec, "category")
			self.Spec["categoryId"] = value.Id.(interface{})
		}
	}
	if item, ok := self.Spec["level"]; ok {
		if value, ok := common.AliasCache.TryGetLevel(item.(string)); ok {
			delete(self.Spec, "level")
			self.Spec["levelId"] = value.Id.(interface{})
		}
	}
	if item, ok := self.Spec["owner"]; ok {
		if value, ok := common.AliasCache.TryGetTeam(item.(string)); ok {
			delete(self.Spec, "owner")
			self.Spec["ownerId"] = value.Id.(interface{})
		}
	}
	if item, ok := self.Spec["filter"]; ok {
		if value, ok := common.AliasCache.TryGetFilter(item.(string)); ok {
			delete(self.Spec, "filter")
			self.Spec["filterId"] = value.Id.(interface{})
		}
	}
}

func (self *CheckCreateType) AsCustomEventCreate() *opslevel.CheckCustomEventCreateInput {
	if item, ok := self.Spec["integration"]; ok {
		if value, ok := common.AliasCache.TryGetIntegration(item.(string)); ok {
			delete(self.Spec, "integration")
			self.Spec["integrationId"] = value.Id.(interface{})
		}
	}
	self.Spec["resultMessage"] = self.Spec["message"]
	payload := &opslevel.CheckCustomEventCreateInput{}
	dataBytes, err := json.Marshal(self.Spec)
	cobra.CheckErr(err)
	json.Unmarshal(dataBytes, payload)
	return payload
}

func readCheckCreateInput() (*CheckCreateType, error) {
	readCreateConfigFile()
	evt := &CheckCreateType{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	evt.resolveAliases()
	return evt, nil
}
