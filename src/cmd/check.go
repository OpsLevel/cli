package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/opslevel/opslevel-go/v2023"

	"github.com/creasty/defaults"
	"github.com/opslevel/cli/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var checkCreateCmd = &cobra.Command{
	Use:     "check",
	Short:   "Create a rubric check",
	Long:    `Create a rubric check`,
	Example: `opslevel create check -f my_cec.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		usePrompts := !hasStdin()
		input, err := readCheckCreateInput()
		cobra.CheckErr(err)
		clientGQL := getClientGQL()
		opslevel.Cache.CacheCategories(clientGQL)
		opslevel.Cache.CacheLevels(clientGQL)
		opslevel.Cache.CacheTeams(clientGQL)
		opslevel.Cache.CacheFilters(clientGQL)
		check, err := createCheck(*input, usePrompts, _dryRun)
		cobra.CheckErr(err)
		fmt.Printf("Created Check '%s' with id '%s'\n", check.Name, check.Id)
	},
}

var importCheckCmd = &cobra.Command{
	Use:        "check CSV_FILEPATH",
	Short:      "Import a CSV of check definitions",
	Long:       `Import a CSV of check definitions`,
	Example:    `opslevel import check data.csv`,
	Aliases:    []string{"checks"},
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"CSV_FILEPATH"},
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		reader, err := common.ReadCSVFile(filepath)
		cobra.CheckErr(err)
		clientGQL := getClientGQL()
		opslevel.Cache.CacheCategories(clientGQL)
		opslevel.Cache.CacheLevels(clientGQL)
		opslevel.Cache.CacheTeams(clientGQL)
		opslevel.Cache.CacheFilters(clientGQL)
		for reader.Rows() {
			spec := marshalFromCSV(reader)
			check, err := createCheck(spec, false, _dryRun)
			cobra.CheckErr(err)
			fmt.Printf("Created Check '%s' with id '%s'\n", check.Name, check.Id)
		}
	},
}

var exportCheckCmd = &cobra.Command{
	Use:        "check CSV_FILEPATH",
	Short:      "Export a CSV of check definitions",
	Long:       `Export a CSV of check definitions`,
	Example:    `opslevel export check data.csv`,
	Aliases:    []string{"checks"},
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"CSV_FILEPATH"},
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		output := newFile(filepath, false)
		defer output.Close()
		output.WriteString(strings.Join(getCheckCSVHeaders(), "\t") + "\n")
		client := getClientGQL()
		resp, err := client.ListChecks(nil)
		cobra.CheckErr(err)
		for _, check := range resp.Nodes {
			marshalToCSV(output, check)
		}
	},
}

var getCheckCmd = &cobra.Command{
	Use:        "check ID",
	Short:      "Get details about a rubic check",
	Long:       `Get details about a rubic check`,
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
	createCmd.AddCommand(checkCreateCmd)
	importCmd.AddCommand(importCheckCmd)
	exportCmd.AddCommand(exportCheckCmd)
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
			self.Spec["categoryId"] = value.Id
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
		self.Spec["categoryId"] = category.Id
	} else {
		return fmt.Errorf("no valid value supplied for field 'category'")
	}
	return nil
}

func (self *CheckCreateType) resolveLevelAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := self.Spec["level"]; ok {
		delete(self.Spec, "level")
		if value, ok := opslevel.Cache.TryGetLevel(item.(string)); ok {
			self.Spec["levelId"] = value.Id
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
		self.Spec["levelId"] = level.Id
	} else {
		return fmt.Errorf("no valid value supplied for field 'level'")
	}
	return nil
}

func (self *CheckCreateType) resolveTeamAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := self.Spec["owner"]; ok {
		delete(self.Spec, "owner")
		if value, ok := opslevel.Cache.TryGetTeam(item.(string)); ok {
			self.Spec["ownerId"] = value.Id
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
			self.Spec["ownerId"] = team.Id
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
			self.Spec["filterId"] = value.Id
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
			self.Spec["filterId"] = filter.Id
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

func (self *CheckCreateType) AsHasRecentDeployCreateInput() *opslevel.CheckHasRecentDeployCreateInput {
	payload := &opslevel.CheckHasRecentDeployCreateInput{}
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

func (self *CheckCreateType) AsHasDocumentationCreateInput() *opslevel.CheckHasDocumentationCreateInput {
	payload := &opslevel.CheckHasDocumentationCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsRepositoryIntegratedCreateInput() *opslevel.CheckRepositoryIntegratedCreateInput {
	payload := &opslevel.CheckRepositoryIntegratedCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsToolUsageCreateInput() *opslevel.CheckToolUsageCreateInput {
	payload := &opslevel.CheckToolUsageCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsTagDefinedCreateInput() *opslevel.CheckTagDefinedCreateInput {
	payload := &opslevel.CheckTagDefinedCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsRepositoryFileCreateInput() *opslevel.CheckRepositoryFileCreateInput {
	payload := &opslevel.CheckRepositoryFileCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsRepositoryGrepCreateInput() *opslevel.CheckRepositoryGrepCreateInput {
	payload := &opslevel.CheckRepositoryGrepCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsRepositorySearchCreateInput() *opslevel.CheckRepositorySearchCreateInput {
	payload := &opslevel.CheckRepositorySearchCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsManualCreateInput() *opslevel.CheckManualCreateInput {
	payload := &opslevel.CheckManualCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsAlertSourceUsageCreateInput() *opslevel.CheckAlertSourceUsageCreateInput {
	payload := &opslevel.CheckAlertSourceUsageCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsGitBranchProtectionCreateInput() *opslevel.CheckGitBranchProtectionCreateInput {
	payload := &opslevel.CheckGitBranchProtectionCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) AsServiceDependencyCreateInput() *opslevel.CheckServiceDependencyCreateInput {
	payload := &opslevel.CheckServiceDependencyCreateInput{}
	json.Unmarshal(toJson(self.Spec), payload)
	return payload
}

func (self *CheckCreateType) resolveIntegrationAliases(client *opslevel.Client, usePrompt bool) error {
	if item, ok := self.Spec["integration"]; ok {
		delete(self.Spec, "integration")
		if value, ok := opslevel.Cache.TryGetIntegration(item.(string)); ok {
			self.Spec["integrationId"] = value.Id
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
		self.Spec["integrationId"] = integration.Id
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

func createCheck(input CheckCreateType, usePrompts bool, dryRun bool) (*opslevel.Check, error) {
	var output *opslevel.Check
	var err error
	clientGQL := getClientGQL()
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
		if dryRun {
			output, err = clientGQL.CreateCheckServiceOwnership(*input.AsServiceOwnershipCreateInput())
		} else {
			log.Info().Type("check", input.AsServiceOwnershipCreateInput()).Msg("")
		}

	case opslevel.CheckTypeHasRecentDeploy:
		if dryRun {
			output, err = clientGQL.CreateCheckHasRecentDeploy(*input.AsHasRecentDeployCreateInput())
		} else {
			log.Info().Type("check", input.AsHasRecentDeployCreateInput()).Msg("")
		}

	case opslevel.CheckTypeServiceProperty:
		if dryRun {
			output, err = clientGQL.CreateCheckServiceProperty(*input.AsServicePropertyCreateInput())
		} else {
			log.Info().Type("check", input.AsServicePropertyCreateInput()).Msg("")
		}

	case opslevel.CheckTypeHasServiceConfig:
		if dryRun {
			output, err = clientGQL.CreateCheckServiceConfiguration(*input.AsServiceConfigurationCreateInput())
		} else {
			log.Info().Type("check", input.AsServiceConfigurationCreateInput()).Msg("")
		}

	case opslevel.CheckTypeHasDocumentation:
		if dryRun {
			output, err = clientGQL.CreateCheckHasDocumentation(*input.AsHasDocumentationCreateInput())
		} else {
			log.Info().Type("check", input.AsHasDocumentationCreateInput()).Msg("")
		}

	case opslevel.CheckTypeHasRepository:
		if dryRun {
			output, err = clientGQL.CreateCheckRepositoryIntegrated(*input.AsRepositoryIntegratedCreateInput())
		} else {
			log.Info().Type("check", input.AsRepositoryIntegratedCreateInput()).Msg("")
		}

	case opslevel.CheckTypeToolUsage:
		if dryRun {
			output, err = clientGQL.CreateCheckToolUsage(*input.AsToolUsageCreateInput())
		} else {
			log.Info().Type("check", input.AsToolUsageCreateInput()).Msg("")
		}

	case opslevel.CheckTypeTagDefined:
		if dryRun {
			output, err = clientGQL.CreateCheckTagDefined(*input.AsTagDefinedCreateInput())
		} else {
			log.Info().Type("check", input.AsTagDefinedCreateInput()).Msg("")
		}

	case opslevel.CheckTypeRepoFile:
		if dryRun {
			output, err = clientGQL.CreateCheckRepositoryFile(*input.AsRepositoryFileCreateInput())
		} else {
			log.Info().Type("check", input.AsRepositoryFileCreateInput()).Msg("")
		}

	case opslevel.CheckTypeRepoGrep:
		if dryRun {
			output, err = clientGQL.CreateCheckRepositoryGrep(*input.AsRepositoryGrepCreateInput())
		} else {
			log.Info().Type("check", input.AsRepositoryGrepCreateInput()).Msg("")
		}

	case opslevel.CheckTypeRepoSearch:
		if dryRun {
			output, err = clientGQL.CreateCheckRepositorySearch(*input.AsRepositorySearchCreateInput())
		} else {
			log.Info().Type("check", input.AsRepositorySearchCreateInput()).Msg("")
		}

	case opslevel.CheckTypeManual:
		if dryRun {
			output, err = clientGQL.CreateCheckManual(*input.AsManualCreateInput())
		} else {
			log.Info().Type("check", input.AsManualCreateInput()).Msg("")
		}

	case opslevel.CheckTypeGeneric:
		opslevel.Cache.CacheIntegrations(clientGQL)
		err = input.resolveIntegrationAliases(clientGQL, usePrompts)
		cobra.CheckErr(err)
		if dryRun {
			output, err = clientGQL.CreateCheckCustomEvent(*input.AsCustomEventCreateInput())
		} else {
			log.Info().Type("check", input.AsCustomEventCreateInput()).Msg("")
		}

	case opslevel.CheckTypeAlertSourceUsage:
		if dryRun {
			output, err = clientGQL.CreateCheckAlertSourceUsage(*input.AsAlertSourceUsageCreateInput())
		} else {
			log.Info().Type("check", input.AsAlertSourceUsageCreateInput()).Msg("")
		}

	case opslevel.CheckTypeGitBranchProtection:
		if dryRun {
			output, err = clientGQL.CreateCheckGitBranchProtection(*input.AsGitBranchProtectionCreateInput())
		} else {
			log.Info().Type("check", input.AsGitBranchProtectionCreateInput()).Msg("")
		}

	case opslevel.CheckTypeServiceDependency:
		if dryRun {
			output, err = clientGQL.CreateCheckServiceDependency(*input.AsServiceDependencyCreateInput())
		} else {
			log.Info().Type("check", input.AsServiceDependencyCreateInput()).Msg("")
		}
	}
	cobra.CheckErr(err)
	if output == nil {
		return nil, fmt.Errorf("unknown error - no check data returned")
	}
	return output, err
}

func marshalPredicateInput(predicate string) opslevel.PredicateInput {
	output := opslevel.PredicateInput{}
	cobra.CheckErr(json.Unmarshal([]byte(predicate), &output))
	return output
}

func marshalManualCheckFrequencyInput(predicate string) opslevel.ManualCheckFrequencyInput {
	output := opslevel.ManualCheckFrequencyInput{}
	cobra.CheckErr(json.Unmarshal([]byte(predicate), &output))
	return output
}

func marshalFromCSV(reader *common.CSVReader) CheckCreateType {
	checkCreate := CheckCreateType{
		Spec: map[string]interface{}{
			"name":     reader.Text("Name"),
			"enabled":  true,
			"category": reader.Text("Category"),
			"level":    reader.Text("Level"),
			"filter":   reader.Text("Filter"),
			"owner":    reader.Text("Owner"),
			"notes":    reader.Text("Notes"),
		},
	}
	switch reader.Text("Type") {
	case "Alert Source Usage":
		checkCreate.Kind = opslevel.CheckTypeAlertSourceUsage
		checkCreate.Spec["alertSourceNamePredicate"] = marshalPredicateInput(reader.Text("Alert Source Name Predicate"))
		checkCreate.Spec["alertSourceType"] = reader.Text("Alert Source Type")
	case "Custom Event":
		checkCreate.Kind = opslevel.CheckTypeGeneric
		checkCreate.Spec["integration"] = reader.Text("Integration")
		checkCreate.Spec["serviceSelector"] = reader.Text("Service Selector")
		checkCreate.Spec["successCondition"] = reader.Text("Success Condition")
		checkCreate.Spec["message"] = reader.Text("Message")
	case "Git Branch Protection":
		checkCreate.Kind = opslevel.CheckTypeGitBranchProtection
	case "Has Documentation":
		checkCreate.Kind = opslevel.CheckTypeHasDocumentation
		checkCreate.Spec["documentType"] = reader.Text("Document Type")
		checkCreate.Spec["documentSubtype"] = reader.Text("Document Subtype")
	case "Has Owner":
		checkCreate.Kind = opslevel.CheckTypeHasOwner
		checkCreate.Spec["requireContactMethod"] = reader.Bool("Require Contact Method")
		checkCreate.Spec["contactMethod"] = reader.Text("Contact Method")
		checkCreate.Spec["tagKey"] = reader.Text("Tag Key")
		checkCreate.Spec["tagPredicate"] = marshalPredicateInput(reader.Text("Tag Predicate"))
	case "Has Recent Deploy":
		checkCreate.Kind = opslevel.CheckTypeHasRecentDeploy
		checkCreate.Spec["days"] = reader.Int("Days")
	case "Has Repository":
		checkCreate.Kind = opslevel.CheckTypeHasRepository
	case "Has Service Config":
		checkCreate.Kind = opslevel.CheckTypeHasServiceConfig
	case "Manual":
		checkCreate.Kind = opslevel.CheckTypeManual
		checkCreate.Spec["updateFrequency"] = marshalManualCheckFrequencyInput(reader.Text("Update Frequency"))
		checkCreate.Spec["updateRequiresComment"] = reader.Bool("Require Update Comment")
	case "Repository File":
		checkCreate.Kind = opslevel.CheckTypeRepoFile
		checkCreate.Spec["directorySearch"] = reader.Bool("Directory Search")
		checkCreate.Spec["filePaths"] = strings.Split(reader.Text("File Paths"), ";")
		checkCreate.Spec["fileContentsPredicate"] = marshalPredicateInput(reader.Text("File Contents Predicate"))
		checkCreate.Spec["useAbsoluteRootPath"] = reader.Bool("Use Absolute Root Path")
	case "Repository Grep":
		checkCreate.Kind = opslevel.CheckTypeRepoGrep
		checkCreate.Spec["directorySearch"] = reader.Bool("Directory Search")
		checkCreate.Spec["filePaths"] = strings.Split(reader.Text("File Paths"), ";")
		checkCreate.Spec["fileContentsPredicate"] = marshalPredicateInput(reader.Text("File Contents Predicate"))
	case "Repository Search":
		checkCreate.Kind = opslevel.CheckTypeRepoSearch
		checkCreate.Spec["fileExtensions"] = strings.Split(reader.Text("File Extensions"), ";")
		checkCreate.Spec["fileContentsPredicate"] = marshalPredicateInput(reader.Text("File Contents Predicate"))
	case "Service Dependency":
		checkCreate.Kind = opslevel.CheckTypeServiceDependency
	case "Service Property":
		checkCreate.Kind = opslevel.CheckTypeServiceProperty
		checkCreate.Spec["serviceProperty"] = reader.Text("Service Property")
		checkCreate.Spec["propertyValuePredicate"] = marshalPredicateInput(reader.Text("Service Property Predicate"))
	case "Tag Defined":
		checkCreate.Kind = opslevel.CheckTypeTagDefined
		checkCreate.Spec["tagKey"] = reader.Text("Tag Key")
		checkCreate.Spec["tagPredicate"] = marshalPredicateInput(reader.Text("Tag Predicate"))
	case "Tool Usage":
		checkCreate.Kind = opslevel.CheckTypeToolUsage
		checkCreate.Spec["toolCategory"] = reader.Text("Tool Category")
		checkCreate.Spec["toolNamePredicate"] = marshalPredicateInput(reader.Text("Tool Name Predicate"))
		checkCreate.Spec["toolUrlPredicate"] = marshalPredicateInput(reader.Text("Tool URL Predicate"))
		checkCreate.Spec["environmentPredicate"] = marshalPredicateInput(reader.Text("Tool Environment Predicate"))
	}

	return checkCreate
}

func marshalString(input string) string {
	data, err := json.Marshal(input)
	cobra.CheckErr(err)
	return string(data)
}

func marshalPredicate(predicate *opslevel.Predicate) string {
	if predicate == nil {
		return ""
	}
	data, err := json.Marshal(predicate)
	cobra.CheckErr(err)
	return string(data)
}

func getCheckCSVHeaders() []string {
	return []string{
		"Name",
		"Enabled",
		"Category",
		"Level",
		"Filter",
		"Owner",
		"Notes",
		"Alert Source Name Predicate",
		"Alert Source Type",
		"Integration",
		"Service Selector",
		"Success Condition",
		"Message",
		"Document Type",
		"Document Subtype",
		"Require Contact Method",
		"Contact Method",
		"Tag Key",
		"Tag Predicate",
		"Days",
		"Update Frequency",
		"Require Update Comment",
		"Directory Search",
		"File Extensions",
		"File Paths",
		"File Contents Predicate",
		"Use Absolute Root Path",
		"Service Property",
		"Service Property Predicate",
		"Tool Category",
		"Tool Name Predicate",
		"Tool URL Predicate",
		"Tool Environment Predicate",
	}
}

func marshalToCSV(output *os.File, check opslevel.Check) {
	headers := getCheckCSVHeaders()
	row := map[string]string{}
	for _, key := range headers {
		row[key] = ""
	}
	row["Name"] = marshalString(check.Name)
	row["Enabled"] = strconv.FormatBool(check.Enabled)
	row["Category"] = check.Category.Alias()
	row["Level"] = check.Level.Alias
	row["Filter"] = check.Filter.Alias()
	row["Owner"] = check.Owner.Team.Alias
	row["Notes"] = marshalString(check.Notes)

	switch check.Type {
	case opslevel.CheckTypeAlertSourceUsage:
		row["Alert Source Name Predicate"] = marshalPredicate(&check.AlertSourceNamePredicate)
		row["Alert Source Type"] = string(check.AlertSourceType)
	case opslevel.CheckTypeGeneric:
		row["Integration"] = check.Integration.Alias()
		row["Service Selector"] = marshalString(check.ServiceSelector)
		row["Success Condition"] = marshalString(check.SuccessCondition)
		row["Message"] = marshalString(check.ResultMessage)
	case opslevel.CheckTypeHasDocumentation:
		row["Document Type"] = string(check.DocumentType)
		row["Document Subtype"] = string(check.DocumentSubtype)
	case opslevel.CheckTypeHasOwner:
		row["Require Contact Method"] = strconv.FormatBool(*check.RequireContactMethod)
		row["Contact Method"] = string(*check.ContactMethod)
		row["Tag Key"] = check.TagKey
		row["Tag Predicate"] = marshalPredicate(check.TagPredicate)
	case opslevel.CheckTypeHasRecentDeploy:
		row["Days"] = strconv.FormatInt(int64(check.Days), 10)
	case opslevel.CheckTypeManual:
		data, err := json.Marshal(check.UpdateFrequency)
		cobra.CheckErr(err)
		row["Update Frequency"] = string(data)
		row["Require Update Comment"] = strconv.FormatBool(check.UpdateRequiresComment)
	case opslevel.CheckTypeRepoFile:
		row["Directory Search"] = strconv.FormatBool(check.RepositoryFileCheckFragment.DirectorySearch)
		row["File Paths"] = strings.Join(check.RepositoryFileCheckFragment.Filepaths, ";")
		row["File Contents Predicate"] = marshalPredicate(check.RepositoryFileCheckFragment.FileContentsPredicate)
		row["Use Absolute Root Path"] = strconv.FormatBool(check.UseAbsoluteRoot)
	case opslevel.CheckTypeRepoGrep:
		row["Directory Search"] = strconv.FormatBool(check.RepositoryGrepCheckFragment.DirectorySearch)
		row["File Paths"] = strings.Join(check.RepositoryGrepCheckFragment.Filepaths, ";")
		row["File Contents Predicate"] = marshalPredicate(check.RepositoryGrepCheckFragment.FileContentsPredicate)
	case opslevel.CheckTypeRepoSearch:
		row["File Extensions"] = strings.Join(check.RepositorySearchCheckFragment.FileExtensions, ";")
		row["File Contents Predicate"] = marshalPredicate(&check.RepositorySearchCheckFragment.FileContentsPredicate)
	case opslevel.CheckTypeServiceProperty:
		row["Service Property"] = string(check.ServicePropertyCheckFragment.Property)
		row["Service Property Predicate"] = marshalPredicate(check.ServicePropertyCheckFragment.Predicate)
	case opslevel.CheckTypeTagDefined:
		row["Tag Key"] = check.TagKey
		row["Tag Predicate"] = marshalPredicate(check.TagPredicate)
	case opslevel.CheckTypeToolUsage:
		row["Tool Category"] = string(check.ToolCategory)
		row["Tool Name Predicate"] = marshalPredicate(check.ToolNamePredicate)
		row["Tool URL Predicate"] = marshalPredicate(check.ToolUrlPredicate)
		row["Tool Environment Predicate"] = marshalPredicate(check.ToolUsageCheckFragment.EnvironmentPredicate)
	}

	var values []string
	for _, key := range headers {
		values = append(values, row[key])
	}
	output.WriteString(strings.Join(values, "\t") + "\n")
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
		output.Spec["useAbsoluteRoot"] = check.RepositoryFileCheckFragment.UseAbsoluteRoot

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

var CheckConfigCurrentVersion = "1"

type ConfigVersion struct {
	Version string
}

func readCheckCreateInput() (*CheckCreateType, error) {
	readInputConfig()
	// Validate Version
	v := &ConfigVersion{}
	viper.Unmarshal(&v)
	if v.Version != CheckConfigCurrentVersion {
		return nil, fmt.Errorf("Supported config version is '%s' but found '%s' | Please update config file", CheckConfigCurrentVersion, v.Version)
	}
	// Unmarshall
	evt := &CheckCreateType{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	return evt, nil
}
