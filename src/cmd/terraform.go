package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
	"github.com/opslevel/cli/common"
	"github.com/opslevel/opslevel-go"
	"github.com/spf13/cobra"
)

// TODO: things that use slug.Make - are susceptible to non-unique names

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
	rubric := newFile(fmt.Sprintf("%s/opslevel_rubric.tf", directory), false)
	services := newFile(fmt.Sprintf("%s/opslevel_services.tf", directory), false)
	teams := newFile(fmt.Sprintf("%s/opslevel_teams.tf", directory), false)

	defer bash.Close()
	defer main.Close()
	defer rubric.Close()
	defer services.Close()
	defer teams.Close()

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

	exportRubric(client, rubric, bash)
	exportFilters(client, rubric, bash)
	exportServices(client, services, bash)
	exportTeams(client, teams, bash)
	fmt.Println("Complete!")
}

func exportRubric(c *opslevel.Client, config *os.File, shell *os.File) {
	categories, err := c.ListCategories()
	cobra.CheckErr(err)
	categoryConfig := `resource "opslevel_rubric_category" "%s" {
  name = "%s"
}
`
	for _, category := range categories {
		slug := slug.Make(category.Name)
		config.WriteString(fmt.Sprintf(categoryConfig, slug, category.Name))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_rubric_category.%s %s\n", slug, category.Id))
	}

	levelConfig := `resource "opslevel_rubric_level" "%s" {
  name = "%s"
  description = "%s"
}
`
	levels, err := c.ListLevels()
	cobra.CheckErr(err)
	for _, level := range levels {
		config.WriteString(fmt.Sprintf(levelConfig, level.Alias, level.Name, level.Description))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_rubric_level.%s %s\n", level.Alias, level.Id))
	}
}

func exportFilters(c *opslevel.Client, config *os.File, shell *os.File) {
	filterConfig := `resource "opslevel_filter" "%s" {
  name = "%s"
  connective = "%s"
}
`
	filters, err := c.ListFilters()
	cobra.CheckErr(err)
	for _, filter := range filters {
		slug := slug.Make(filter.Name)
		config.WriteString(fmt.Sprintf(filterConfig, slug, filter.Name, string(filter.Connective)))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_filter.%s %s\n", slug, filter.Id))
	}
}

func exportServices(c *opslevel.Client, config *os.File, shell *os.File) {
	serviceConfig := `resource "opslevel_service" "%s" {
  name = "%s"
  description = "%s"
  product = "%s"
  framework = "%s"
  language = "%s"
  lifecycle_alias = "%s"
  tier_alias = "%s"
  owner_alias = "%s"

  aliases = ["%s"]
}
`
	services, err := c.ListServices()
	cobra.CheckErr(err)
	for _, service := range services {
		config.WriteString(fmt.Sprintf(serviceConfig, service.Aliases[0], service.Name, service.Description, service.Product, service.Framework, service.Language, service.Lifecycle.Alias, service.Tier.Alias, service.Owner.Alias, strings.Join(service.Aliases, "\", \"")))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_service.%s %s\n", service.Aliases[0], service.Id))
	}
}

func exportTeams(c *opslevel.Client, config *os.File, shell *os.File) {
	teamConfig := `resource "opslevel_team" "%s" {
  name = "%s"
  manager_email = "%s"
  responsibilities = "%s"
}
`
	teams, err := c.ListTeams()
	cobra.CheckErr(err)
	for _, team := range teams {
		config.WriteString(fmt.Sprintf(teamConfig, team.Alias, team.Name, team.Manager.Email, team.Responsibilities))
		shell.WriteString(fmt.Sprintf("terraform import opslevel_team.%s %s\n", team.Alias, team.Id))
	}
}
