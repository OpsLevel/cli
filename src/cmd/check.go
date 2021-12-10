package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/gosimple/slug"
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
		if isYamlOutput() {
			common.YamlPrint(marshalCheck(*check))
		} else {
			common.PrettyPrint(check)
		}
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
			w := common.NewTabWriter("NAME", "TYPE", "ID")
			for _, item := range list {
				fmt.Fprintf(w, "%s\t%s\t%s\t\n", item.Name, item.Type, item.Id)
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
		check, err := createCheck(*input)
		cobra.CheckErr(err)
		fmt.Printf("Created Check '%s' with id '%s'\n", check.Name, check.Id)
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

func toJson(data map[string]interface{}) []byte {
	output, err := json.Marshal(data)
	cobra.CheckErr(err)
	return output
}

func (self *CheckCreateType) AsServiceOwnershipCreateInput() *opslevel.CheckServiceOwnershipCreateInput {
	payload := &opslevel.CheckServiceOwnershipCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsServicePropertyCreateInput() *opslevel.CheckServicePropertyCreateInput {
	payload := &opslevel.CheckServicePropertyCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsServiceConfigurationCreateInput() *opslevel.CheckServiceConfigurationCreateInput {
	payload := &opslevel.CheckServiceConfigurationCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsRepositoryFileCreateInput() *opslevel.CheckRepositoryFileCreateInput {
	payload := &opslevel.CheckRepositoryFileCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsRepositoryIntegratedCreateInput() *opslevel.CheckRepositoryIntegratedCreateInput {
	payload := &opslevel.CheckRepositoryIntegratedCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsRepositorySearchCreateInput() *opslevel.CheckRepositorySearchCreateInput {
	payload := &opslevel.CheckRepositorySearchCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsTagDefinedCreateInput() *opslevel.CheckTagDefinedCreateInput {
	payload := &opslevel.CheckTagDefinedCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsToolUsageCreateInput() *opslevel.CheckToolUsageCreateInput {
	payload := &opslevel.CheckToolUsageCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsManualCreateInput() *opslevel.CheckManualCreateInput {
	payload := &opslevel.CheckManualCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsCustomEventCreateInput() *opslevel.CheckCustomEventCreateInput {
	if item, ok := self.Spec["integration"]; ok {
		if value, ok := common.AliasCache.TryGetIntegration(item.(string)); ok {
			delete(self.Spec, "integration")
			self.Spec["integrationId"] = value.Id.(interface{})
		}
	}
	self.Spec["resultMessage"] = self.Spec["message"]
	payload := &opslevel.CheckCustomEventCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func createCheck(input CheckCreateType) (*opslevel.Check, error) {
	var output *opslevel.Check
	var err error
	switch input.Kind {
	case opslevel.CheckTypeHasOwner:
		output, err = graphqlClient.CreateCheckServiceOwnership(*input.AsServiceOwnershipCreateInput())

	case opslevel.CheckTypeServiceProperty:
		output, err = graphqlClient.CreateCheckServiceProperty(*input.AsServicePropertyCreateInput())

	case opslevel.CheckTypeHasServiceConfig:
		output, err = graphqlClient.CreateCheckServiceConfiguration(*input.AsServiceConfigurationCreateInput())

	case opslevel.CheckTypeHasRepository:
		output, err = graphqlClient.CreateCheckRepositoryIntegrated(*input.AsRepositoryIntegratedCreateInput())

	case opslevel.CheckTypeToolUsage:
		output, err = graphqlClient.CreateCheckToolUsage(*input.AsToolUsageCreateInput())

	case opslevel.CheckTypeTagDefined:
		output, err = graphqlClient.CreateCheckTagDefined(*input.AsTagDefinedCreateInput())

	case opslevel.CheckTypeRepoFile:
		output, err = graphqlClient.CreateCheckRepositoryFile(*input.AsRepositoryFileCreateInput())

	case opslevel.CheckTypeRepoSearch:
		output, err = graphqlClient.CreateCheckRepositorySearch(*input.AsRepositorySearchCreateInput())

	case opslevel.CheckTypeManual:
		output, err = graphqlClient.CreateCheckManual(*input.AsManualCreateInput())

	case opslevel.CheckTypeGeneric:
		output, err = graphqlClient.CreateCheckCustomEvent(*input.AsCustomEventCreateInput())
	}
	if output == nil {
		return nil, fmt.Errorf("")
	}
	return output, err
}

func marshalCheck(check opslevel.Check) *CheckCreateType {
	output := &CheckCreateType{
		Kind: check.Type,
		Spec: map[string]interface{}{
			"name":     check.Name,
			"enabled":  check.Enabled,
			"category": slug.Make(check.Category.Name),
			"level":    check.Level.Alias,
			"notes":    check.Notes,
		},
	}
	if check.Filter.Id != nil {
		output.Spec["filter"] = slug.Make(check.Filter.Name)
	}
	if check.Owner.Team.Id != nil {
		output.Spec["owner"] = check.Owner.Team.Alias
	}
	switch check.Type {
	case opslevel.CheckTypeHasOwner:

	case opslevel.CheckTypeServiceProperty:
		output.Spec["serviceProperty"] = check.ServicePropertyCheckFragment.Property
		output.Spec["propertyValuePredicate"] = check.ServicePropertyCheckFragment.Predicate

	case opslevel.CheckTypeHasServiceConfig:

	case opslevel.CheckTypeHasRepository:

	case opslevel.CheckTypeToolUsage:
		output.Spec["toolCategory"] = check.ToolUsageCheckFragment.ToolCategory
		output.Spec["toolNamePredicate"] = check.ToolUsageCheckFragment.ToolNamePredicate
		output.Spec["environmentPredicate"] = check.ToolUsageCheckFragment.EnvironmentPredicate

	case opslevel.CheckTypeTagDefined:
		output.Spec["tagKey"] = check.TagDefinedCheckFragment.TagKey
		output.Spec["tagPredicate"] = check.TagDefinedCheckFragment.TagPredicate

	case opslevel.CheckTypeRepoFile:
		output.Spec["directorySearch"] = check.RepositoryFileCheckFragment.DirectorySearch
		output.Spec["filePaths"] = check.RepositoryFileCheckFragment.Filepaths
		output.Spec["fileContentsPredicate"] = check.RepositoryFileCheckFragment.FileContentsPredicate

	case opslevel.CheckTypeRepoSearch:
		output.Spec["fileExtensions"] = check.RepositorySearchCheckFragment.FileExtensions
		output.Spec["fileContentsPredicate"] = check.RepositorySearchCheckFragment.FileContentsPredicate

	case opslevel.CheckTypeCustom:

	case opslevel.CheckTypePayload:

	case opslevel.CheckTypeManual:
		output.Spec["updateFrequency"] = check.ManualCheckFragment.UpdateFrequency
		output.Spec["updateRequiresComment"] = check.ManualCheckFragment.UpdateRequiresComment

	case opslevel.CheckTypeGeneric:
		output.Spec["integration"] = fmt.Sprintf("%s-%s", slug.Make(check.CustomEventCheckFragment.Integration.Type), slug.Make(check.CustomEventCheckFragment.Integration.Name))
		output.Spec["serviceSelector"] = check.CustomEventCheckFragment.ServiceSelector
		output.Spec["successCondition"] = check.CustomEventCheckFragment.SuccessCondition
		output.Spec["message"] = check.CustomEventCheckFragment.ResultMessage
	}

	return output
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
