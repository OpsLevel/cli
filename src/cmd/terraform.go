package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/gosimple/slug"
	"github.com/spf13/cobra"
)

var importTerraformCmd = &cobra.Command{
	Use:        "terraform [Directory]",
	Short:      "Exports your account data to be controlled by terraform",
	Long:       `Writes a series of files to disk that enabled you to fully control your OpsLevel account via terraform`,
	Args:       cobra.MaximumNArgs(1),
	ArgAliases: []string{"Directory"},
	Run:        runExportTerraform,
}

func init() {
	exportCmd.AddCommand(importTerraformCmd)
}

func newFile(filename string, makeExecutable bool) *os.File {
	_, err := os.Stat(filename)
	if os.IsExist(err) {
		removeErr := os.Remove(filename)
		if removeErr != nil {
			panic(removeErr)
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	if makeExecutable {
		if err := os.Chmod(filename, 0o755); err != nil {
			panic(err)
		}
	}
	return file
}

func templateConfig(tmpl string, a ...interface{}) string {
	// TODO: it would be nice to remove blank lines to condense the terraform config - this is a hacky way that doesn't work
	// return strings.ReplaceAll(fmt.Sprintf(tmpl, a...), "  \n", "")
	return fmt.Sprintf(tmpl, a...)
}

func runExportTerraform(cmd *cobra.Command, args []string) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		path = "terraform"
	}
	directory, directoryErr := filepath.Abs(path)
	cobra.CheckErr(directoryErr)
	makeDirErr := os.MkdirAll(directory, os.ModePerm)
	cobra.CheckErr(makeDirErr)
	fmt.Printf("Writing files to: %s\n", directory)
	bash := newFile(fmt.Sprintf("%s/import.sh", directory), true)
	main := newFile(fmt.Sprintf("%s/main.tf", directory), false)
	constants := newFile(fmt.Sprintf("%s/opslevel_constants.tf", directory), false)
	teams := newFile(fmt.Sprintf("%s/opslevel_teams.tf", directory), false)
	repos := newFile(fmt.Sprintf("%s/opslevel_repos.tf", directory), false)
	rubric := newFile(fmt.Sprintf("%s/opslevel_rubric.tf", directory), false)
	filters := newFile(fmt.Sprintf("%s/opslevel_filters.tf", directory), false)

	defer bash.Close()
	defer main.Close()
	defer constants.Close()
	defer teams.Close()
	defer repos.Close()
	defer rubric.Close()
	defer filters.Close()

	main.WriteString(`terraform {
  required_providers {
    opslevel = {
      source  = "opslevel/opslevel"
    }
  }
}

provider "opslevel" {
}
`)
	bash.WriteString("#!/bin/sh\n\n")

	graphqlClient := getClientGQL()

	exportConstants(graphqlClient, constants)
	exportRepos(graphqlClient, repos, bash)
	exportServices(graphqlClient, bash, directory)
	exportTeams(graphqlClient, teams, bash)
	exportFilters(graphqlClient, filters, bash)
	exportRubric(graphqlClient, rubric, bash)
	exportChecks(graphqlClient, bash, directory)
	fmt.Println("Complete!")
}

// TODO: things that use this are susceptible to non-unique names
// maybe we can use some sort of global name registry and if duplicates happen add an increment
// we would likely need to tie the resource's ID to the generated string for future lookups for connected resources
func makeTerraformSlug(value string) string {
	return strings.ReplaceAll(slug.Make(value), "-", "_")
	// return strings.ReplaceAll(strings.ReplaceAll(slug.Make(value), "-", "_"), ":", "_")
}

func getIntegrationTerraformName(integration opslevel.IntegrationId) string {
	return makeTerraformSlug(fmt.Sprintf("%s_%s", integration.Type, integration.Name))
}

// Given a field that could be a multiline string - this will return it with the correct formatting
func buildMultilineStringArg(fieldName string, fieldContents string) string {
	if len(fieldContents) > 0 {
		if len(strings.Split(fieldContents, "\n")) > 1 {
			if strings.HasSuffix(fieldContents, "\n") {
				fieldContents = strings.TrimSuffix(fieldContents, "\n")
				return fmt.Sprintf(`%s = <<-EOT
%s
EOT`, fieldName, fieldContents)
			} else {
				escaped, err := json.Marshal(fieldContents)
				if err != nil {
					fmt.Println(err)
				}
				return fmt.Sprintf("%s = %s", fieldName, string(escaped))
			}
		} else {
			return fmt.Sprintf("%s = %q", fieldName, fieldContents)
		}
	}
	return fieldContents
}

func exportConstants(c *opslevel.Client, config *os.File) {
	lifecycleConfig := `data "opslevel_lifecycle" "%s" {
  filter {
    field = "id"
    value = "%s"
  }
}
`
	lifecycles, err := c.ListLifecycles()
	cobra.CheckErr(err)
	for _, lifecycle := range lifecycles {
		config.WriteString(templateConfig(lifecycleConfig, lifecycle.Alias, lifecycle.Id))
	}

	tierConfig := `data "opslevel_tier" "%s" {
  filter {
    field = "id"
    value = "%s"
  }
}
`
	tiers, err := c.ListTiers()
	cobra.CheckErr(err)
	for _, tier := range tiers {
		config.WriteString(templateConfig(tierConfig, tier.Alias, tier.Id))
	}

	integrationConfig := `data "opslevel_integration" "%s" {
  filter {
    field = "id"
    value = "%s"
  }
}
`
	resp, err := c.ListIntegrations(nil)
	cobra.CheckErr(err)
	for _, integration := range resp.Nodes {
		config.WriteString(templateConfig(integrationConfig, getIntegrationTerraformName(integration.IntegrationId), integration.Id))
	}
}

func exportRepos(c *opslevel.Client, config *os.File, shell *os.File) {
	repoConfig := `data "opslevel_repository" "%s" {
  alias = "%s"
}
`
	resp, err := c.ListRepositories(nil)
	repos := resp.Nodes
	cobra.CheckErr(err)
	for _, repo := range repos {
		if repo.DefaultAlias == "" {
			continue
		}
		config.WriteString(templateConfig(repoConfig, makeTerraformSlug(repo.DefaultAlias), repo.DefaultAlias))
	}
}

func flattenAliases(aliases []string) string {
	return strings.Join(aliases, "\", \"")
}

func flattenTags(tags []opslevel.Tag) string {
	tagStrings := make([]string, len(tags))
	for _, tag := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s:%s", tag.Key, tag.Value))
	}

	return strings.Join(tagStrings, "\", \"")
}

func flattenLifecycle(value opslevel.Lifecycle) string {
	if value.Id != "" {
		return fmt.Sprintf("lifecycle_alias = data.opslevel_lifecycle.%s.alias", value.Alias)
	}
	return ""
}

func flattenTier(value opslevel.Tier) string {
	if value.Id != "" {
		return fmt.Sprintf("tier_alias = data.opslevel_tier.%s.alias", value.Alias)
	}
	return ""
}

func flattenOwner(value opslevel.TeamId) string {
	if value.Id != "" {
		return fmt.Sprintf("owner = opslevel_team.%s.alias", value.Alias)
	}
	return ""
}

func getToolTerraformName(value opslevel.Tool) string {
	return makeTerraformSlug(fmt.Sprintf("%s %s %s", value.Category, value.Environment, value.DisplayName))
}

func exportServices(c *opslevel.Client, shell *os.File, directory string) {
	serviceConfig := `resource "opslevel_service" "%s" {
  name = "%s"
  description = "%s"
  product = "%s"
  framework = "%s"
  language = "%s"
  %s
  %s
  %s

  %s
  %s
}
`
	serviceToolConfig := `resource "opslevel_service_tool" "%s" {
  service = opslevel_service.%s.id

  name = "%s"
  category = "%s"
  url = "%s"
  environment = "%s"
}
`

	serviceRepoConfig := `resource "opslevel_service_repository" "%s" {
  service = opslevel_service.%s.id
  repository = data.opslevel_repository.%s.id

  name = "%s"
  base_directory = "%s"
}
`
	resp, err := c.ListServices(nil)
	services := resp.Nodes
	cobra.CheckErr(err)
	for _, service := range services {
		serviceMainAlias := makeTerraformSlug(service.Aliases[0])
		file := newFile(fmt.Sprintf("%s/opslevel_service_%s.tf", directory, serviceMainAlias), false)
		aliases := flattenAliases(service.Aliases)
		if len(aliases) > 0 {
			aliases = fmt.Sprintf("aliases = [\"%s\"]", aliases)
		}
		serviceTags, err := service.GetTags(c, nil)
		cobra.CheckErr(err)
		tags := flattenTags(serviceTags.Nodes)
		if len(tags) > 0 {
			tags = fmt.Sprintf("tags = [\"%s\"]", tags)
		}
		file.WriteString(templateConfig(serviceConfig, serviceMainAlias, service.Name, service.Description, service.Product, service.Framework, service.Language, flattenLifecycle(service.Lifecycle), flattenTier(service.Tier), flattenOwner(service.Owner), aliases, tags))
		shell.WriteString(fmt.Sprintf("# Service: %s\n", serviceMainAlias))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_service.%s %s\n", serviceMainAlias, service.Id))
		for _, tool := range service.Tools.Nodes {
			toolTerraformName := makeTerraformSlug(fmt.Sprintf("%s_%s", serviceMainAlias, getToolTerraformName(tool)))
			file.WriteString(templateConfig(serviceToolConfig, toolTerraformName, serviceMainAlias, tool.DisplayName, tool.Category, tool.Url, tool.Environment))
			shell.WriteString(fmt.Sprintf("terraform import opslevel_service_tool.%s %s:%s\n", toolTerraformName, service.Id, tool.Id))
		}
		for _, edge := range service.Repositories.Edges {
			for _, serviceRepo := range edge.ServiceRepositories {
				repo := serviceRepo.Repository
				repoName := makeTerraformSlug(repo.DefaultAlias)
				serviceRepoTerraformName := fmt.Sprintf("%s_%s", serviceMainAlias, repoName)
				file.WriteString(templateConfig(serviceRepoConfig, serviceRepoTerraformName, serviceMainAlias, repoName, serviceRepo.DisplayName, serviceRepo.BaseDirectory))
				shell.WriteString(fmt.Sprintf("terraform import opslevel_service_repository.%s %s:%s\n", serviceRepoTerraformName, service.Id, serviceRepo.Id))
			}
		}
		file.Close()
		shell.WriteString("##########\n\n")
	}
}

func getMembershipsAsTerraformConfig(members []opslevel.TeamMembership) string {
	memberConfig := `
  member {
    email = "%s"
    role = "%s"
  }`

	membersBody := strings.Builder{}
	for _, member := range members {
		membersBody.WriteString(fmt.Sprintf(memberConfig, member.User.Email, member.Role))
	}
	return membersBody.String()
}

func exportTeams(c *opslevel.Client, config *os.File, shell *os.File) {
	shell.WriteString("# Teams\n")

	teamConfig := `resource "opslevel_team" "%s" {%s
}

`
	resp, err := c.ListTeams(nil)
	teams := resp.Nodes
	cobra.CheckErr(err)
	teamBody := strings.Builder{}
	for _, team := range teams {
		aliases := flattenAliases(team.Aliases)
		teamBody.WriteString(fmt.Sprintf("\n  aliases = [\"%s\"]", aliases))
		teamBody.WriteString(fmt.Sprintf("\n  name = \"%s\"", team.Name))

		membersOutput := getMembershipsAsTerraformConfig(team.Memberships.Nodes)
		teamBody.WriteString(membersOutput)
		if len(team.ParentTeam.Alias) > 0 {
			teamBody.WriteString(fmt.Sprintf("\n  parent = [\"%s\"]", team.ParentTeam.Alias))
		}
		if len(team.Responsibilities) > 0 {
			responsibilities := buildMultilineStringArg("responsibilities", team.Responsibilities)
			teamBody.WriteString(fmt.Sprintf("\n  %s", responsibilities))
		}

		config.WriteString(templateConfig(
			teamConfig,
			team.Alias,
			teamBody.String(),
		))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_team.%s %s\n", team.Alias, team.Id))
		teamBody.Reset()
	}
	shell.WriteString("##########\n\n")
}

func exportRubric(c *opslevel.Client, config *os.File, shell *os.File) {
	shell.WriteString("# Rubric\n")

	resp, err := c.ListCategories(nil)
	cobra.CheckErr(err)
	categories := resp.Nodes
	categoryConfig := `resource "opslevel_rubric_category" "%s" {
  name = "%s"
}
`
	for _, category := range categories {
		categoryTerraformName := makeTerraformSlug(category.Name)
		config.WriteString(templateConfig(categoryConfig, categoryTerraformName, category.Name))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_rubric_category.%s %s\n", categoryTerraformName, category.Id))
	}

	levelConfig := `resource "opslevel_rubric_level" "%s" {
  name = "%s"
  description = "%s"
  index = %d
}
`
	levels, err := c.ListLevels()
	cobra.CheckErr(err)
	for _, level := range levels {
		config.WriteString(templateConfig(levelConfig, level.Alias, level.Name, level.Description, level.Index))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_rubric_level.%s %s\n", level.Alias, level.Id))
	}

	shell.WriteString("##########\n\n")
}

func exportFilters(c *opslevel.Client, config *os.File, shell *os.File) {
	shell.WriteString("# Filters\n")
	filterConfig := `resource "opslevel_filter" "%s" {
  name = "%s"
  connective = "%s"
  %s
}
`
	resp, err := c.ListFilters(nil)
	cobra.CheckErr(err)
	for _, filter := range resp.Nodes {
		predicates := ""
		filterTerraformName := makeTerraformSlug(filter.Name)
		for _, predicate := range filter.Predicates {
			predicates += flattenFilterPredicate(&predicate)
		}
		config.WriteString(templateConfig(filterConfig, filterTerraformName, filter.Name, string(filter.Connective), predicates))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_filter.%s %s\n", filterTerraformName, filter.Id))
	}

	shell.WriteString("##########\n\n")
}

func flattenCheckOwner(value opslevel.CheckOwner) string {
	if value.Team.Id != "" {
		return fmt.Sprintf("owner = opslevel_team.%s.id", value.Team.Alias)
	}
	return ""
}

func flattenCheckFilter(value opslevel.Filter) string {
	if value.Id != "" {
		return fmt.Sprintf("filter = opslevel_filter.%s.id", makeTerraformSlug(value.Name))
	}
	return ""
}

func flattenPredicate(key string, value *opslevel.Predicate) string {
	config := `
  %s {
    type = "%s"
    %s
  }
`
	if value != nil {
		return templateConfig(config, key, value.Type, buildMultilineStringArg("value", strings.ReplaceAll(value.Value, "\"", "\\\"")))
	}
	return ""
}

func flattenFilterPredicate(value *opslevel.FilterPredicate) string {
	config := `
  predicate {
    key = "%s"
    key_data = "%s"
    type = "%s"
    %s
  }
`
	if value != nil {
		return templateConfig(config, value.Key, value.KeyData, value.Type, buildMultilineStringArg("value", value.Value))
	}
	return ""
}

func flattenUpdateFrequency(value *opslevel.ManualCheckFrequency) string {
	config := `
  update_frequency {
    starting_data = "%s"
    time_scale = "%s"
    value = %d
  }
`
	if value != nil {
		return templateConfig(config, value.StartingDate.Format(time.RFC3339), value.FrequencyTimeScale, value.FrequencyValue)
	}
	return ""
}

func exportChecks(c *opslevel.Client, shell *os.File, directory string) {
	shell.WriteString("# Checks\n")
	// TODO: If we use golang templating here we can easily remove all the extra newlines
	baseCheckConfig := `resource "opslevel_check_%s" "%s" {
  name = %q
  enabled = %v
  category = opslevel_rubric_category.%s.id
  level = opslevel_rubric_level.%s.id
  %s
  %s
  %s
  %s
}
`
	customEventCheckConfig := `integration = data.opslevel_integration.%s.id
  service_selector = %q
  success_condition = %q
  %s`
	manualCheckConfig := `%s
  update_requires_comment = %v`
	repoFileCheckConfig := `directory_search = %v
  filepaths = ["%s"]
  %s`
	repoGrepCheckConfig := `directory_search = %v
  filepaths = ["%s"]
  %s`
	repoSearchCheckConfig := `%s
  %s`
	servicePropertyCheckConfig := `property = "%s"
  %s`
	tagDefinedCheckConfig := `tag_key = "%s"
  %s`
	toolUsageCheckConfig := `tool_category = "%s"
  %s
  %s`
	hasDocumentationCheckConfig := `document_type = "%s"
  document_subtype = "%s"`

	hasRecentDeployCheckConfig := `days = "%d"`

	alertSourceUsageCheckConfig := `%s
  alert_source_type = "%s"`

	alertSourceUsageCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_alert_source_usage.tf", directory), false)
	customEventCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_custom_event.tf", directory), false)
	hasRecentDeployCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_has_recent_deploy.tf", directory), false)
	manualCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_manual.tf", directory), false)
	repoFileCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_repository_file.tf", directory), false)
	repoGrepCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_repository_grep.tf", directory), false)
	repoIntegratedCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_repository_integrated.tf", directory), false)
	repoSearchCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_repository_search.tf", directory), false)
	serviceConfigCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_service_configuration.tf", directory), false)
	serviceOwnershipCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_service_ownership.tf", directory), false)
	servicePropertyCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_service_property.tf", directory), false)
	tagDefinedCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_tag_defined.tf", directory), false)
	toolUsageCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_tool_usage.tf", directory), false)
	hasDocumentationCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_has_documentation.tf", directory), false)
	gitBranchProtectionCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_git_branch_protection.tf", directory), false)
	serviceDependencyCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_service_dependency.tf", directory), false)

	defer alertSourceUsageCheckFile.Close()
	defer customEventCheckFile.Close()
	defer hasRecentDeployCheckFile.Close()
	defer manualCheckFile.Close()
	defer repoFileCheckFile.Close()
	defer repoGrepCheckFile.Close()
	defer repoIntegratedCheckFile.Close()
	defer repoSearchCheckFile.Close()
	defer serviceConfigCheckFile.Close()
	defer serviceOwnershipCheckFile.Close()
	defer servicePropertyCheckFile.Close()
	defer tagDefinedCheckFile.Close()
	defer toolUsageCheckFile.Close()
	defer hasDocumentationCheckFile.Close()
	defer gitBranchProtectionCheckFile.Close()
	defer serviceDependencyCheckFile.Close()

	var activeFile *os.File
	resp, err := c.ListChecks(nil)
	cobra.CheckErr(err)
	for _, check := range resp.Nodes {
		checkTerraformName := makeTerraformSlug(check.Name)
		checkTypeTerraformName := ""
		checkExtras := ""
		switch check.Type {
		case opslevel.CheckTypeAlertSourceUsage:
			casted := check.AlertSourceUsageCheckFragment
			activeFile = alertSourceUsageCheckFile
			checkTypeTerraformName = "alert_source_usage"
			checkExtras = templateConfig(alertSourceUsageCheckConfig, flattenPredicate("alert_source_name_predicate", casted.AlertSourceNamePredicate), casted.AlertSourceType)
		case opslevel.CheckTypeGeneric:
			casted := check.CustomEventCheckFragment
			activeFile = customEventCheckFile
			checkTypeTerraformName = "custom_event"
			checkExtras = templateConfig(customEventCheckConfig, getIntegrationTerraformName(casted.Integration), casted.ServiceSelector, casted.SuccessCondition, buildMultilineStringArg("message", casted.ResultMessage))
		case opslevel.CheckTypeHasRecentDeploy:
			casted := check.HasRecentDeployCheckFragment
			activeFile = hasRecentDeployCheckFile
			checkTypeTerraformName = "has_recent_deploy"
			checkExtras = templateConfig(hasRecentDeployCheckConfig, casted.Days)
		case opslevel.CheckTypeManual:
			casted := check.ManualCheckFragment
			activeFile = manualCheckFile
			checkTypeTerraformName = "manual"
			checkExtras = templateConfig(manualCheckConfig, flattenUpdateFrequency(casted.UpdateFrequency), casted.UpdateRequiresComment)
		case opslevel.CheckTypeRepoFile:
			casted := check.RepositoryFileCheckFragment
			activeFile = repoFileCheckFile
			checkTypeTerraformName = "repository_file"
			checkExtras = templateConfig(repoFileCheckConfig, casted.DirectorySearch, strings.Join(casted.Filepaths, "\", \""), flattenPredicate("file_contents_predicate", casted.FileContentsPredicate))
		case opslevel.CheckTypeRepoGrep:
			casted := check.RepositoryGrepCheckFragment
			activeFile = repoGrepCheckFile
			checkTypeTerraformName = "repository_grep"
			checkExtras = templateConfig(repoGrepCheckConfig, casted.DirectorySearch, strings.Join(casted.Filepaths, "\", \""), flattenPredicate("file_contents_predicate", &casted.FileContentsPredicate))
		case opslevel.CheckTypeHasRepository:
			activeFile = repoIntegratedCheckFile
			checkTypeTerraformName = "repository_integrated"
		case opslevel.CheckTypeRepoSearch:
			casted := check.RepositorySearchCheckFragment
			activeFile = repoSearchCheckFile
			checkTypeTerraformName = "repository_search"
			fileExtensions := ""
			if len(casted.FileExtensions) > 0 {
				fileExtensions = fmt.Sprintf(`file_extensions = ["%s"]`, strings.Join(casted.FileExtensions, "\", \""))
			}
			checkExtras = templateConfig(repoSearchCheckConfig, fileExtensions, flattenPredicate("file_contents_predicate", &casted.FileContentsPredicate))
		case opslevel.CheckTypeHasServiceConfig:
			activeFile = serviceConfigCheckFile
			checkTypeTerraformName = "service_configuration"
		case opslevel.CheckTypeHasOwner:
			activeFile = serviceOwnershipCheckFile
			checkTypeTerraformName = "service_ownership"
		case opslevel.CheckTypeServiceProperty:
			casted := check.ServicePropertyCheckFragment
			activeFile = servicePropertyCheckFile
			checkTypeTerraformName = "service_property"
			checkExtras = templateConfig(servicePropertyCheckConfig, casted.Property, flattenPredicate("predicate", casted.Predicate))
		case opslevel.CheckTypeTagDefined:
			casted := check.TagDefinedCheckFragment
			activeFile = tagDefinedCheckFile
			checkTypeTerraformName = "tag_defined"
			checkExtras = templateConfig(tagDefinedCheckConfig, casted.TagKey, flattenPredicate("tag_predicate", casted.TagPredicate))
		case opslevel.CheckTypeToolUsage:
			casted := check.ToolUsageCheckFragment
			activeFile = toolUsageCheckFile
			checkTypeTerraformName = "tool_usage"
			checkExtras = templateConfig(toolUsageCheckConfig, casted.ToolCategory, flattenPredicate("tool_name_predicate", casted.ToolNamePredicate), flattenPredicate("environment_predicate", casted.EnvironmentPredicate))
		case opslevel.CheckTypeHasDocumentation:
			casted := check.HasDocumentationCheckFragment
			activeFile = hasDocumentationCheckFile
			checkTypeTerraformName = "has_documentation"
			checkExtras = templateConfig(hasDocumentationCheckConfig, casted.DocumentType, casted.DocumentSubtype)
		case opslevel.CheckTypeGitBranchProtection:
			activeFile = gitBranchProtectionCheckFile
			checkTypeTerraformName = "git_branch_protection"
		case opslevel.CheckTypeServiceDependency:
			activeFile = serviceDependencyCheckFile
			checkTypeTerraformName = "service_dependency"
		default:
			continue
		}

		if activeFile == nil {
			continue
		}
		checkConfig := templateConfig(baseCheckConfig, checkTypeTerraformName, checkTerraformName, check.Name, check.Enabled, makeTerraformSlug(check.Category.Name), check.Level.Alias, flattenCheckOwner(check.Owner), flattenCheckFilter(check.Filter), checkExtras, buildMultilineStringArg("notes", check.Notes))
		activeFile.WriteString(checkConfig)
		shell.WriteString(fmt.Sprintf("terraform import opslevel_check_%s.%s %s\n", checkTypeTerraformName, checkTerraformName, check.Id))
	}

	shell.WriteString("##########\n")
}
