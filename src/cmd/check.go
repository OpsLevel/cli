package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"
	"github.com/rs/zerolog/log"
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
		check, err := getClientGQL().GetCheck(args[0])
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
		list, err := getClientGQL().ListChecks()
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
		usePrompts := hasStdin() == false
		input, err := readCheckCreateInput()
		cobra.CheckErr(err)
		check, err := createCheck(*input, usePrompts)
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
		err := getClientGQL().DeleteCheck(key)
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
	Version string
	Kind    opslevel.CheckType
	Spec    map[string]interface{}
}

func (self *CheckCreateType) resolveCategoryAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := self.Spec["category"]; ok {
		delete(self.Spec, "category")
		if value, ok := opslevel.Cache.TryGetCategory(item.(string)); ok {
			self.Spec["categoryId"] = value.Id.(interface{})
			return nil
		} else {
			fmt.Printf("%s is not a valid category, please select a valid category\n", item.(string))
		}
	}
	if usePrompt {
		category, promptErr := common.PromptForCategories(client)
		if promptErr != nil {
			return promptErr
		}
		self.Spec["categoryId"] = category.Id.(interface{})
	} else {
		return fmt.Errorf("no valid value supplied for field 'category'")
	}
	return nil
}

func (self *CheckCreateType) resolveLevelAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := self.Spec["level"]; ok {
		delete(self.Spec, "level")
		if value, ok := opslevel.Cache.TryGetLevel(item.(string)); ok {
			self.Spec["levelId"] = value.Id.(interface{})
			return nil
		} else {
			fmt.Printf("%s is not a valid level, please select a valid level\n", item.(string))
		}
	}
	if usePrompt {
		level, promptErr := common.PromptForLevels(client)
		if promptErr != nil {
			return promptErr
		}
		self.Spec["levelId"] = level.Id.(interface{})
	} else {
		return fmt.Errorf("no valid value supplied for field 'level'")
	}
	return nil
}

func (self *CheckCreateType) resolveTeamAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := self.Spec["owner"]; ok {
		delete(self.Spec, "owner")
		if value, ok := opslevel.Cache.TryGetTeam(item.(string)); ok {
			self.Spec["ownerId"] = value.Id.(interface{})
			return nil
		} else {
			fmt.Printf("%s is not a valid team, please select a valid team\n", item.(string))
		}
	}
	if usePrompt {
		team, promptErr := common.PromptForTeam(client)
		if promptErr != nil {
			return promptErr
		}
		if team.Id != nil {
			self.Spec["ownerId"] = team.Id.(interface{})
		}
	} else {
		log.Warn().Msg("no value supplied for field 'owner'")
	}
	return nil
}

func (self *CheckCreateType) resolveFilterAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := self.Spec["filter"]; ok {
		delete(self.Spec, "filter")
		if value, ok := opslevel.Cache.TryGetFilter(item.(string)); ok {
			self.Spec["filterId"] = value.Id.(interface{})
			return nil
		} else {
			fmt.Printf("%s is not a valid filter, please select a valid filter\n", item.(string))
		}
	}
	if usePrompt {
		filter, promptErr := common.PromptForFilter(client)
		if promptErr != nil {
			return promptErr
		}
		if filter.Id != nil {
			self.Spec["filterId"] = filter.Id.(interface{})
		}
	} else {
		log.Warn().Msg("no value supplied for field 'filter'")
	}
	return nil
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

func (self *CheckCreateType) resolveIntegrationAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := self.Spec["integration"]; ok {
		delete(self.Spec, "integration")
		if value, ok := opslevel.Cache.TryGetIntegration(item.(string)); ok {
			self.Spec["integrationId"] = value.Id.(interface{})
			return nil
		} else {
			fmt.Printf("%s is not a valid integration, please select a valid integration\n", item.(string))
		}
	}
	if usePrompt {
		integration, promptErr := common.PromptForIntegration(client)
		if promptErr != nil {
			return promptErr
		}
		self.Spec["integrationId"] = integration.Id.(interface{})
	} else {
		return fmt.Errorf("no valid value supplied for field 'integration'")
	}
	return nil
}

func (self *CheckCreateType) AsCustomEventCreateInput() *opslevel.CheckCustomEventCreateInput {
	self.Spec["resultMessage"] = self.Spec["message"]
	payload := &opslevel.CheckCustomEventCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func createCheck(input CheckCreateType, usePrompts bool) (*opslevel.Check, error) {
	var output *opslevel.Check
	var err error
	clientGQL := getClientGQL()
	opslevel.Cache.CacheCategories(clientGQL)
	opslevel.Cache.CacheLevels(clientGQL)
	opslevel.Cache.CacheTeams(clientGQL)
	opslevel.Cache.CacheFilters(clientGQL)
	err = input.resolveCategoryAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveLevelAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveTeamAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveFilterAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	switch input.Kind {
	case opslevel.CheckTypeHasOwner:
		output, err = clientGQL.CreateCheckServiceOwnership(*input.AsServiceOwnershipCreateInput())

	case opslevel.CheckTypeServiceProperty:
		output, err = clientGQL.CreateCheckServiceProperty(*input.AsServicePropertyCreateInput())

	case opslevel.CheckTypeHasServiceConfig:
		output, err = clientGQL.CreateCheckServiceConfiguration(*input.AsServiceConfigurationCreateInput())

	case opslevel.CheckTypeHasRepository:
		output, err = clientGQL.CreateCheckRepositoryIntegrated(*input.AsRepositoryIntegratedCreateInput())

	case opslevel.CheckTypeToolUsage:
		output, err = clientGQL.CreateCheckToolUsage(*input.AsToolUsageCreateInput())

	case opslevel.CheckTypeTagDefined:
		output, err = clientGQL.CreateCheckTagDefined(*input.AsTagDefinedCreateInput())

	case opslevel.CheckTypeRepoFile:
		output, err = clientGQL.CreateCheckRepositoryFile(*input.AsRepositoryFileCreateInput())

	case opslevel.CheckTypeRepoSearch:
		output, err = clientGQL.CreateCheckRepositorySearch(*input.AsRepositorySearchCreateInput())

	case opslevel.CheckTypeManual:
		output, err = clientGQL.CreateCheckManual(*input.AsManualCreateInput())

	case opslevel.CheckTypeGeneric:
		opslevel.Cache.CacheIntegrations(clientGQL)
		err = input.resolveIntegrationAliases(clientGQL, usePrompts)
		cobra.CheckErr(err)
		output, err = clientGQL.CreateCheckCustomEvent(*input.AsCustomEventCreateInput())
	}
	cobra.CheckErr(err)
	if output == nil {
		return nil, fmt.Errorf("unknown error - no check data returned")
	}
	return output, err
}

func marshalCheck(check opslevel.Check) *CheckCreateType {
	output := &CheckCreateType{
		Version: "1",
		Kind:    check.Type,
		Spec: map[string]interface{}{
			"name":     check.Name,
			"enabled":  check.Enabled,
			"category": check.Category.Alias(),
			"level":    check.Level.Alias,
			"notes":    check.Notes,
		},
	}
	if check.Filter.Id != nil {
		output.Spec["filter"] = check.Filter.Alias()
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
		output.Spec["integration"] = check.CustomEventCheckFragment.Integration.Alias()
		output.Spec["serviceSelector"] = check.CustomEventCheckFragment.ServiceSelector
		output.Spec["successCondition"] = check.CustomEventCheckFragment.SuccessCondition
		output.Spec["message"] = check.CustomEventCheckFragment.ResultMessage
	}

	return output
}

var CheckConfigCurrentVersion = "1"

type ConfigVersion struct {
	Version string
}

func readCheckCreateInput() (*CheckCreateType, error) {
	readCreateConfigFile()
	// Validate Version
	v := &ConfigVersion{}
	viper.Unmarshal(&v)
	if v.Version != CheckConfigCurrentVersion {
		return nil, errors.New(fmt.Sprintf("Supported config version is '%s' but found '%s' | Please update config file", CheckConfigCurrentVersion, v.Version))
	}
	// Unmarshall
	evt := &CheckCreateType{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
