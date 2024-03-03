package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/opslevel/cli/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Commands

var getCheckCmd = &cobra.Command{
	Use:        "check ID",
	Short:      "Get details about a rubric check",
	Long:       `Get details about a rubric check`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		check, err := getClientGQL().GetCheck(opslevel.ID(args[0]))
		cobra.CheckErr(err)
		if isYamlOutput() {
			common.YamlPrint(marshalCheck(*check))
		} else {
			common.PrettyPrint(marshalCheck(*check))
		}
	},
}

var listCheckCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"checks"},
	Short:   "Lists the rubric checks",
	Long:    `Lists the rubric checks`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getClientGQL().ListChecks(nil)
		cobra.CheckErr(err)
		list := resp.Nodes
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
	Long:  `Create a rubric check`,
	Example: `
cat << EOF | opslevel create check -f -
Version: "1"
kind: "repo_grep"
spec:
  name: "new repo grep check"
  enabled: false
  category: "misc"
  level: "silver"
  notes: "some notes"
  directorySearch: false
  filePaths:
  - "Taskfile.yml"
  fileContentsPredicate:
    type: "exists"
    value: ""
EOF`,
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readCheckInput()
		cobra.CheckErr(err)
		usePrompts := true
		if cmd.Flag("noninteractive").Changed {
			usePrompts = false
		} else if dataFile == "-" {
			log.Warn().Msg("running noninteractively since using heredoc")
			usePrompts = false
		}
		check, err := createCheck(*input, usePrompts)
		cobra.CheckErr(err)
		fmt.Printf("Created Check '%s' with id '%s'\n", check.Name, check.Id)
	},
}

var checkUpdateCmd = &cobra.Command{
	Use:   "check",
	Short: "Update a rubric check",
	Long:  `Update a rubric check`,
	Example: `
cat << EOF | opslevel update check $CHECK_ID -f -
Version: "1"
kind: "repo_grep"
spec:
  name: "updated repo grep check"
  enabled: true
  fileContentsPredicate:
    type: "does_not_match_regex"
    value: "gofmt"
EOF`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"ID"},
	Run: func(cmd *cobra.Command, args []string) {
		input, err := readCheckInput()
		cobra.CheckErr(err)
		input.Spec["id"] = *opslevel.NewID(args[0])
		usePrompts := true
		if cmd.Flag("noninteractive").Changed {
			usePrompts = false
		} else if dataFile == "-" {
			log.Warn().Msg("running noninteractively since using heredoc")
			usePrompts = false
		}
		check, err := updateCheck(*input, usePrompts)
		cobra.CheckErr(err)
		common.JsonPrint(json.MarshalIndent(check, "", "    "))
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
		err := getClientGQL().DeleteCheck(opslevel.ID(key))
		cobra.CheckErr(err)
		fmt.Printf("deleted '%s' check\n", key)
	},
}

func init() {
	checkCreateCmd.Flags().Bool("noninteractive", false, "turns off automated prompts for fields missing data in the spec")
	createCmd.AddCommand(checkCreateCmd)
	checkUpdateCmd.Flags().Bool("noninteractive", false, "turns off automated prompts for fields missing data in the spec")
	updateCmd.AddCommand(checkUpdateCmd)
	getCmd.AddCommand(getCheckCmd)
	listCmd.AddCommand(listCheckCmd)
	deleteCmd.AddCommand(deleteCheckCmd)
}

// API requests

func createCheck(input CheckInputType, usePrompts bool) (*opslevel.Check, error) {
	var output *opslevel.Check
	var err error
	clientGQL := getClientGQL()
	opslevel.Cache.CacheCategories(clientGQL)
	opslevel.Cache.CacheLevels(clientGQL)
	opslevel.Cache.CacheTeams(clientGQL)
	opslevel.Cache.CacheFilters(clientGQL)
	opslevel.Cache.CacheIntegrations(clientGQL)
	err = input.resolveCategoryAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveLevelAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveTeamAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveFilterAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	if input.Kind == opslevel.CheckTypeGeneric {
		err = input.resolveIntegrationAliases(clientGQL, usePrompts)
		cobra.CheckErr(err)
	}

	checkData, err := opslevel.UnmarshalCheckCreateInput(input.Kind, toJson(input.Spec))
	cobra.CheckErr(err)
	output, err = clientGQL.CreateCheck(checkData)
	cobra.CheckErr(err)
	if output == nil {
		return nil, fmt.Errorf("unknown error - no check data returned")
	}
	return output, err
}

func updateCheck(input CheckInputType, usePrompts bool) (*opslevel.Check, error) {
	var output *opslevel.Check
	var err error
	clientGQL := getClientGQL()
	opslevel.Cache.CacheCategories(clientGQL)
	opslevel.Cache.CacheLevels(clientGQL)
	opslevel.Cache.CacheTeams(clientGQL)
	opslevel.Cache.CacheFilters(clientGQL)
	opslevel.Cache.CacheIntegrations(clientGQL)
	err = input.resolveCategoryAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveLevelAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveTeamAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	err = input.resolveFilterAliases(clientGQL, usePrompts)
	cobra.CheckErr(err)
	if input.Kind == opslevel.CheckTypeGeneric {
		err = input.resolveIntegrationAliases(clientGQL, usePrompts)
		cobra.CheckErr(err)
	}

	checkData, err := opslevel.UnmarshalCheckUpdateInput(input.Kind, toJson(input.Spec))
	cobra.CheckErr(err)
	output, err = clientGQL.UpdateCheck(checkData)
	cobra.CheckErr(err)
	if output == nil {
		return nil, fmt.Errorf("unknown error - no check data returned")
	}
	return output, err
}

// Resolving foreign keys

func (checkInputType *CheckInputType) resolveCategoryAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := checkInputType.Spec["category"]; ok {
		delete(checkInputType.Spec, "category")
		if value, ok := opslevel.Cache.TryGetCategory(item.(string)); ok {
			checkInputType.Spec["categoryId"] = value.Id
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
		checkInputType.Spec["categoryId"] = category.Id
	} else {
		if checkInputType.IsUpdateInput() {
			return nil
		}

		return fmt.Errorf("no valid value supplied for field 'category'")
	}
	return nil
}

func (checkInputType *CheckInputType) resolveLevelAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := checkInputType.Spec["level"]; ok {
		delete(checkInputType.Spec, "level")
		if value, ok := opslevel.Cache.TryGetLevel(item.(string)); ok {
			checkInputType.Spec["levelId"] = value.Id
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
		checkInputType.Spec["levelId"] = level.Id
	} else {
		if checkInputType.IsUpdateInput() {
			return nil
		}

		return fmt.Errorf("no valid value supplied for field 'level'")
	}
	return nil
}

func (checkInputType *CheckInputType) resolveTeamAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := checkInputType.Spec["owner"]; ok {
		delete(checkInputType.Spec, "owner")
		if value, ok := opslevel.Cache.TryGetTeam(item.(string)); ok {
			checkInputType.Spec["ownerId"] = value.Id
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
		if team.Id != "" {
			checkInputType.Spec["ownerId"] = team.Id
		}
	} else {
		log.Warn().Msg("no value supplied for field 'owner'")
	}
	return nil
}

func (checkInputType *CheckInputType) resolveFilterAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := checkInputType.Spec["filter"]; ok {
		delete(checkInputType.Spec, "filter")
		if value, ok := opslevel.Cache.TryGetFilter(item.(string)); ok {
			checkInputType.Spec["filterId"] = value.Id
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
		if filter.Id != "" {
			checkInputType.Spec["filterId"] = filter.Id
		}
	} else {
		log.Warn().Msg("no value supplied for field 'filter'")
	}
	return nil
}

func (checkInputType *CheckInputType) resolveIntegrationAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := checkInputType.Spec["integration"]; ok {
		delete(checkInputType.Spec, "integration")
		if value, ok := opslevel.Cache.TryGetIntegration(item.(string)); ok {
			checkInputType.Spec["integrationId"] = value.Id
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
		checkInputType.Spec["integrationId"] = integration.Id
	} else {
		if checkInputType.IsUpdateInput() {
			return nil
		}

		return fmt.Errorf("no valid value supplied for field 'integration'")
	}
	return nil
}

// Serialization

func toJson(data map[string]interface{}) []byte {
	output, err := json.Marshal(data)
	cobra.CheckErr(err)
	return output
}

func marshalCheck(check opslevel.Check) *CheckInputType {
	output := &CheckInputType{
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
	if check.Filter.Id != "" {
		output.Spec["filter"] = check.Filter.Alias()
	}
	if check.Owner.Team.Id != "" {
		output.Spec["owner"] = check.Owner.Team.Alias
	}
	switch check.Type {
	case opslevel.CheckTypeHasOwner:

	case opslevel.CheckTypeHasRecentDeploy:
		output.Spec["days"] = check.HasRecentDeployCheckFragment.Days

	case opslevel.CheckTypeServiceProperty:
		output.Spec["serviceProperty"] = check.ServicePropertyCheckFragment.Property
		output.Spec["propertyValuePredicate"] = check.ServicePropertyCheckFragment.Predicate

	case opslevel.CheckTypeHasServiceConfig:

	case opslevel.CheckTypeHasDocumentation:
		output.Spec["documentType"] = check.HasDocumentationCheckFragment.DocumentType
		output.Spec["documentSubtype"] = check.HasDocumentationCheckFragment.DocumentSubtype

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

	case opslevel.CheckTypeRepoGrep:
		output.Spec["directorySearch"] = check.RepositoryGrepCheckFragment.DirectorySearch
		output.Spec["filePaths"] = check.RepositoryGrepCheckFragment.Filepaths
		output.Spec["fileContentsPredicate"] = check.RepositoryGrepCheckFragment.FileContentsPredicate

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

	case opslevel.CheckTypeAlertSourceUsage:
		output.Spec["alertSourceNamePredicate"] = check.AlertSourceUsageCheckFragment.AlertSourceNamePredicate
		output.Spec["alertSourceType"] = check.AlertSourceUsageCheckFragment.AlertSourceType

	case opslevel.CheckTypeGitBranchProtection:

	case opslevel.CheckTypeServiceDependency:

	}

	return output
}

// Types

var CheckConfigCurrentVersion = "1"

type CheckInputType struct {
	Version string `yaml:"Version"`
	Kind    opslevel.CheckType
	Spec    map[string]interface{}
}

func (checkInputType *CheckInputType) IsUpdateInput() bool {
	_, ok := checkInputType.Spec["id"]
	return ok
}

func readCheckInput() (*CheckInputType, error) {
	input, err := readResourceInput[CheckInputType]()
	if err != nil {
		return nil, err
	}
	if input == nil {
		return nil, fmt.Errorf("unexpected nil input")
	}
	if input.Version != CheckConfigCurrentVersion {
		return nil, fmt.Errorf("supported config version is '%s' but found '%s' | please update config file",
			CheckConfigCurrentVersion, input.Version)
	}
	return input, nil
}
