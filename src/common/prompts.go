package common

import (
	"fmt"
	"sort"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/manifoldco/promptui"
)

func PromptForCategories(client *opslevel.Client) (*opslevel.Category, error) {
	resp, err := client.ListCategories(nil)
	if err != nil {
		return nil, err
	}
	list := resp.Nodes
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf("%s {{ .Name | cyan }}", promptui.IconSelect),
		Inactive: "    {{ .Name | cyan }}",
		Selected: fmt.Sprintf("%s {{ .Name | faint }}", promptui.IconGood),
	}

	prompt := promptui.Select{
		Label:     "Select Category",
		Items:     list,
		Templates: templates,
		Size:      MinInt(6, len(list)),
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &list[index], nil
}

func PromptForLevels(client *opslevel.Client) (*opslevel.Level, error) {
	list, err := client.ListLevels()
	if err != nil {
		return nil, err
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf("%s {{ .Name | cyan }} ({{ .Index | red }})", promptui.IconSelect),
		Inactive: "    {{ .Name | cyan }} ({{ .Index | red }})",
		Selected: fmt.Sprintf("%s {{ .Name | faint }} ({{ .Index | red }})", promptui.IconGood),
		Details: `
		{{ "Alias:" | faint }}	{{ .Alias }}
		{{ "Description:" | faint }}	{{ .Description }}`,
	}

	var filteredList []opslevel.Level
	for _, val := range list {
		if val.Index != 0 {
			filteredList = append(filteredList, val)
		}
	}
	sort.Slice(filteredList, func(i, j int) bool { return filteredList[i].Index < filteredList[j].Index })

	prompt := promptui.Select{
		Label:     "Select Level",
		Items:     filteredList,
		Templates: templates,
		Size:      MinInt(6, len(filteredList)),
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &filteredList[index], nil
}

func PromptForFilter(client *opslevel.Client) (*opslevel.Filter, error) {
	resp, err := client.ListFilters(nil)
	if err != nil {
		return nil, err
	}
	list := resp.Nodes
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf("%s {{ .Name | cyan }}", promptui.IconSelect),
		Inactive: "    {{ .Name | cyan }}",
		Selected: fmt.Sprintf("%s {{ .Name | faint }}", promptui.IconGood),
	}

	noneValue := opslevel.Filter{
		FilterId: opslevel.FilterId{
			Name: "None",
			Id:   "",
		},
	}
	list = append([]opslevel.Filter{noneValue}, list...)

	prompt := promptui.Select{
		Label:     "Select Filter",
		Items:     list,
		Templates: templates,
		Size:      MinInt(6, len(list)),
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &list[index], nil
}

func PromptForTeam(client *opslevel.Client) (*opslevel.Team, error) {
	resp, err := client.ListTeams(nil)
	if err != nil {
		return nil, err
	}
	list := resp.Nodes

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf("%s {{ .Name | cyan }}", promptui.IconSelect),
		Inactive: "    {{ .Name | cyan }}",
		Selected: fmt.Sprintf("%s {{ .Name | faint }}", promptui.IconGood),
		Details: `
		{{ "Aliases:" | faint }}	{{ .Aliases }}
		{{ "Manager:" | faint }}	{{ .Manager.Email }}`,
	}

	noneValue := opslevel.Team{
		Name: "None",
	}
	list = append([]opslevel.Team{noneValue}, list...)

	prompt := promptui.Select{
		Label:     "Select Team",
		Items:     list,
		Templates: templates,
		Size:      MinInt(6, len(list)),
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &list[index], nil
}

func PromptForIntegration(client *opslevel.Client) (*opslevel.Integration, error) {
	resp, err := client.ListIntegrations(nil)
	if err != nil {
		return nil, err
	}
	list := resp.Nodes
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   fmt.Sprintf("%s {{ .Name | cyan }}", promptui.IconSelect),
		Inactive: "    {{ .Name | cyan }}",
		Selected: fmt.Sprintf("%s {{ .Name | faint }}", promptui.IconGood),
		Details: `
		{{ "Type:" | faint }}	{{ .Type }}`,
	}

	filteredList := make([]opslevel.Integration, 0)
	for _, val := range list {
		if val.Type == "generic" {
			filteredList = append(filteredList, val)
		}
	}

	prompt := promptui.Select{
		Label:     "Select Integration",
		Items:     filteredList,
		Templates: templates,
		Size:      MinInt(6, len(filteredList)),
	}

	index, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &filteredList[index], nil
}
