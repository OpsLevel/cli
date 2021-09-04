package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"
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
	path := fmt.Sprintf("%s", filename)

	var _, err = os.Stat(path)
	if os.IsExist(err) {
		removeErr := os.Remove(path)
		if removeErr != nil {
			panic(removeErr)
		}
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if makeExecutable {
		if err := os.Chmod(path, 0755); err != nil {
			panic(err)
		}
	}
	return file
}

// TODO: we need to remove blank lines to condense the terraform config - this is a hacky way we should find a better way
func templateConfig(tmpl string, a ...interface{}) string {
	return strings.ReplaceAll(fmt.Sprintf(tmpl, a...), "  \n", "")
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
	client := common.NewGraphClient()
	bash := newFile(fmt.Sprintf("%s/import.sh", directory), true)
	main := newFile(fmt.Sprintf("%s/main.tf", directory), false)
	teams := newFile(fmt.Sprintf("%s/opslevel_teams.tf", directory), false)
	repos := newFile(fmt.Sprintf("%s/opslevel_repos.tf", directory), false)
	rubric := newFile(fmt.Sprintf("%s/opslevel_rubric.tf", directory), false)
	filters := newFile(fmt.Sprintf("%s/opslevel_filters.tf", directory), false)

	defer bash.Close()
	defer main.Close()
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

	exportConstants(client, main)
	exportRepos(client, repos, bash)
	exportServices(client, bash, directory)
	exportTeams(client, teams, bash)
	exportFilters(client, filters, bash)
	exportRubric(client, rubric, bash)
	exportChecks(client, bash, directory)
	fmt.Println("Complete!")
}

// TODO: things that use this are susceptible to non-unique names
// maybe we can use some sort of global name registry and if duplicates happen add an increment
// we would likely need to tie the resource's ID to the generated string for future lookups for connected resources
func makeTerraformSlug(value string) string {
	return strings.ReplaceAll(slug.Make(value), "-", "_")
	//return strings.ReplaceAll(strings.ReplaceAll(slug.Make(value), "-", "_"), ":", "_")
}

func getIntegrationTerraformName(integration opslevel.Integration) string {
	return makeTerraformSlug(fmt.Sprintf("%s_%s", integration.Type, integration.Name))
}

func exportConstants(c *opslevel.Client, config *os.File) {
	lifecycleConfig := `data "opslevel_lifecycle" "%s" {
  id = "%s"
}
`
	lifecycles, err := c.ListLifecycles()
	cobra.CheckErr(err)
	for _, lifecycle := range lifecycles {
		config.WriteString(templateConfig(lifecycleConfig, lifecycle.Alias, lifecycle.Id))
	}

	tierConfig := `data "opslevel_tier" "%s" {
  id = "%s"
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
	integrations, err := c.ListIntegrations()
	cobra.CheckErr(err)
	for _, integration := range integrations {
		config.WriteString(templateConfig(integrationConfig, getIntegrationTerraformName(integration), integration.Id))
	}
}

func exportRepos(c *opslevel.Client, config *os.File, shell *os.File) {
	repoConfig := `data "opslevel_repository" "%s" {
  alias = "%s"
}
`
	repos, err := c.ListRepositories()
	cobra.CheckErr(err)
	for _, repo := range repos {
		config.WriteString(templateConfig(repoConfig, makeTerraformSlug(repo.DefaultAlias), repo.DefaultAlias))
	}
}

func flattenTags(tags []opslevel.Tag) string {
	tagStrings := []string{}
	for _, tag := range tags {
		tagStrings = append(tagStrings, fmt.Sprintf("%s:%s", tag.Key, tag.Value))
	}

	return strings.Join(tagStrings, "\", \"")
}

func flattenLifecycle(value opslevel.Lifecycle) string {
	if value.Id != nil {
		return fmt.Sprintf("lifecycle_alias = data.opslevel_lifecycle.%s.alias", value.Alias)
	}
	return ""
}

func flattenTier(value opslevel.Tier) string {
	if value.Id != nil {
		return fmt.Sprintf("tier_alias = data.opslevel_tier.%s.alias", value.Alias)
	}
	return ""
}

func flattenOwner(value opslevel.Team) string {
	if value.Id != nil {
		return fmt.Sprintf("owner_alias = opslevel_team.%s.alias", value.Alias)
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

  tags = ["%s"]
}
`
	serviceToolConfig := `resource "opslevel_service_tool" "%s" {
  service = data.opslevel_service.%s.id

  name = "%s"
  category = "%s"
  url = "%s"
  environment = "%s"
}
`

	serviceRepoConfig := `resource "opslevel_service_repository" "%s" {
  service = data.opslevel_service.%s.id
  repository = data.opslevel_repository.%s.id

  name = "%s"
  base_direcotry = "%s"
}
`
	services, err := c.ListServices()
	cobra.CheckErr(err)
	for _, service := range services {
		serviceMainAlias := makeTerraformSlug(service.Aliases[0])
		file := newFile(fmt.Sprintf("%s/opslevel_service_%s.tf", directory, serviceMainAlias), false)
		file.WriteString(templateConfig(serviceConfig, serviceMainAlias, service.Name, service.Description, service.Product, service.Framework, service.Language, flattenLifecycle(service.Lifecycle), flattenTier(service.Tier), flattenOwner(service.Owner), flattenTags(service.Tags.Nodes)))
		shell.WriteString(fmt.Sprintf("# Service: %s\n", serviceMainAlias))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_service.%s %s\n", serviceMainAlias, service.Id))
		for _, tool := range service.Tools.Nodes {
			toolTerraformName := makeTerraformSlug(fmt.Sprintf("%s_%s", serviceMainAlias, getToolTerraformName(tool)))
			file.WriteString(templateConfig(serviceToolConfig, toolTerraformName, serviceMainAlias, tool.DisplayName, tool.Category, tool.Url, tool.Environment))
			shell.WriteString(fmt.Sprintf("terraform import opslevel_service_tool.%s %s\n", toolTerraformName, tool.Id))
		}
		for _, edge := range service.Repositories.Edges {
			for _, serviceRepo := range edge.ServiceRepositories {
				repo := serviceRepo.Repository
				repoName := makeTerraformSlug(repo.DefaultAlias)
				serviceRepoTerraformName := fmt.Sprintf("%s_%s", serviceMainAlias, repoName)
				file.WriteString(templateConfig(serviceRepoConfig, serviceRepoTerraformName, serviceMainAlias, repoName, serviceRepo.DisplayName, serviceRepo.BaseDirectory))
				shell.WriteString(fmt.Sprintf("terraform import opslevel_service_repository.%s %s\n", serviceRepoTerraformName, serviceRepo.Id))
			}
		}
		file.Close()
		shell.WriteString("##########\n\n")
	}

}

func exportTeams(c *opslevel.Client, config *os.File, shell *os.File) {
	shell.WriteString("# Teams\n")

	teamConfig := `resource "opslevel_team" "%s" {
  name = "%s"
  manager_email = "%s"
  responsibilities = <<-EOT
%s
EOT
}
`
	teams, err := c.ListTeams()
	cobra.CheckErr(err)
	for _, team := range teams {
		config.WriteString(templateConfig(teamConfig, team.Alias, team.Name, team.Manager.Email, team.Responsibilities))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_team.%s %s\n", team.Alias, team.Id))
	}
	shell.WriteString("##########\n\n")
}

func exportRubric(c *opslevel.Client, config *os.File, shell *os.File) {
	shell.WriteString("# Rubric\n")

	categories, err := c.ListCategories()
	cobra.CheckErr(err)
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
}
`
	filters, err := c.ListFilters()
	cobra.CheckErr(err)
	for _, filter := range filters {
		filterTerraformName := makeTerraformSlug(filter.Name)
		config.WriteString(templateConfig(filterConfig, filterTerraformName, filter.Name, string(filter.Connective)))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_filter.%s %s\n", filterTerraformName, filter.Id))
	}

	shell.WriteString("##########\n\n")
}

func flattenCheckOwner(value opslevel.CheckOwner) string {
	if value.Team.Id != nil {
		return fmt.Sprintf("owner = opslevel_team.%s.id", value.Team.Alias)
	}
	return ""
}

func flattenCheckFilter(value opslevel.Filter) string {
	if value.Id != nil {
		return fmt.Sprintf("filter = opslevel_filter.%s.id", makeTerraformSlug(value.Name))
	}
	return ""
}

func flattenCheckPredicate(key string, value *opslevel.Predicate) string {
	config := `
  %s {
    type = "%s"
    value = "%s"
  }
`
	if value != nil {
		return templateConfig(config, key, value.Type, strings.ReplaceAll(value.Value, "\"", "\\\""))
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

	baseCheckConfig := `resource "opslevel_check_%s" "%s" {
  name = "%s"
  enabled = %v
  category = opslevel_rubric_category.%s.id
  level = opslevel_rubric_level.%s.id
  %s
  %s
  %s
  notes = <<-EOT
%s
EOT
}
`
	customEventCheckConfig := `integration = data.opslevel_integration.%s.id
  service_selector = "%s"
  success_condition = "%s"
  message = <<-EOT
%s
EOT`
	manualCheckConfig := `%s
  update_requires_comment = %v`
	repoFileCheckConfig := `directory_search = %v
  filepaths = ["%s"]
  %s`
	repoSearchCheckConfig := `file_extensions = ["%s"]
  %s`
	servicePropertyCheckConfig := `property = "%s"
  %s`
	tagDefinedCheckConfig := `tag_key = "%s"
  %s`
	toolUsageCheckConfig := `tool_category = "%s"
  %s
  %s`

	customEventCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_custom_event.tf", directory), false)
	manualCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_manual.tf", directory), false)
	repoFileCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_repository_file.tf", directory), false)
	repoIntegratedCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_repository_integrated.tf", directory), false)
	repoSearchCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_repository_search.tf", directory), false)
	serviceConfigCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_service_configuration.tf", directory), false)
	serviceOwnershipCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_service_ownership.tf", directory), false)
	servicePropertyCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_service_property.tf", directory), false)
	tagDefinedCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_tag_defined.tf", directory), false)
	toolUsageCheckFile := newFile(fmt.Sprintf("%s/opslevel_checks_tool_usage.tf", directory), false)

	defer customEventCheckFile.Close()
	defer manualCheckFile.Close()
	defer repoFileCheckFile.Close()
	defer repoIntegratedCheckFile.Close()
	defer repoSearchCheckFile.Close()
	defer serviceConfigCheckFile.Close()
	defer serviceOwnershipCheckFile.Close()
	defer servicePropertyCheckFile.Close()
	defer tagDefinedCheckFile.Close()
	defer toolUsageCheckFile.Close()

	var activeFile *os.File
	checks, err := c.ListChecks()
	cobra.CheckErr(err)
	for _, check := range checks {
		checkTerraformName := makeTerraformSlug(check.Name)
		checkTypeTerraformName := ""
		checkExtras := ""
		switch check.Type {
		case opslevel.CheckTypeGeneric:
			casted := check.CustomEventCheckFragment
			activeFile = customEventCheckFile
			checkTypeTerraformName = "custom_event"
			checkExtras = templateConfig(customEventCheckConfig, makeTerraformSlug(casted.Integration.Name), casted.ServiceSelector, casted.SuccessCondition, casted.ResultMessage)
		case opslevel.CheckTypeManual:
			casted := check.ManualCheckFragment
			activeFile = manualCheckFile
			checkTypeTerraformName = "manual"
			checkExtras = templateConfig(manualCheckConfig, flattenUpdateFrequency(casted.UpdateFrequency), casted.UpdateRequiresComment)
		case opslevel.CheckTypeRepoFile:
			casted := check.RepositoryFileCheckFragment
			activeFile = repoFileCheckFile
			checkTypeTerraformName = "repository_file"
			checkExtras = templateConfig(repoFileCheckConfig, casted.DirectorySearch, strings.Join(casted.Filepaths, "\", \""), flattenCheckPredicate("file_contents_predicate", casted.FileContentsPredicate))
		case opslevel.CheckTypeHasRepository:
			activeFile = repoIntegratedCheckFile
			checkTypeTerraformName = "repository_integrated"
		case opslevel.CheckTypeRepoSearch:
			casted := check.RepositorySearchCheckFragment
			activeFile = repoSearchCheckFile
			checkTypeTerraformName = "repository_search"
			checkExtras = templateConfig(repoSearchCheckConfig, strings.Join(casted.FileExtensions, "\", \""), flattenCheckPredicate("file_contents_predicate", &casted.FileContentsPredicate))
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
			checkExtras = templateConfig(servicePropertyCheckConfig, casted.Property, flattenCheckPredicate("predicate", casted.Predicate))
		case opslevel.CheckTypeTagDefined:
			casted := check.TagDefinedCheckFragment
			activeFile = tagDefinedCheckFile
			checkTypeTerraformName = "tag_defined"
			checkExtras = templateConfig(tagDefinedCheckConfig, casted.TagKey, flattenCheckPredicate("tag_predicate", casted.TagPredicate))
		case opslevel.CheckTypeToolUsage:
			casted := check.ToolUsageCheckFragment
			activeFile = toolUsageCheckFile
			checkTypeTerraformName = "tool_usage"
			checkExtras = templateConfig(toolUsageCheckConfig, casted.ToolCategory, flattenCheckPredicate("tool_name_predicate", casted.ToolNamePredicate), flattenCheckPredicate("environment_predicate", casted.EnvironmentPredicate))
		default:
			continue
		}

		if activeFile == nil {
			continue
		}
		checkConfig := templateConfig(baseCheckConfig, checkTypeTerraformName, checkTerraformName, check.Name, check.Enabled, makeTerraformSlug(check.Category.Name), check.Level.Alias, flattenCheckOwner(check.Owner), flattenCheckFilter(check.Filter), checkExtras, check.Notes)
		activeFile.WriteString(checkConfig)
		shell.WriteString(fmt.Sprintf("terraform import opslevel_check_%s %s\n", checkTerraformName, check.Id))
	}

	shell.WriteString("##########\n")
}
