package common

import (
	"slices"
	"strings"
)

var hardcodedAliases = map[string][]string{
	"Category":           {"cat"},
	"Document":           {"doc"},
	"Integration":        {"int"},
	"Property":           {"prop"},
	"PropertyDefinition": {"pd"},
	"Repository":         {"repo"},
	"Scorecard":          {"sc"},
	"Service":            {"svc"},
	"System":             {"sys"},
	"TriggerDefinition":  {"td"},
}

func addAlias(key string, value string) {
	if !slices.Contains(hardcodedAliases[key], value) {
		hardcodedAliases[key] = append(hardcodedAliases[key], value)
	}
}

// GetAliases will return a slice of variations of a PascalCase string like the
// lowercase, snake_case, kebab-case, singular and plural versions
func GetAliases(PascalCase string) []string {
	if _, ok := hardcodedAliases[PascalCase]; !ok {
		hardcodedAliases[PascalCase] = []string{}
	}
	addAlias(PascalCase, strings.ToLower(PascalCase))
	addAlias(PascalCase, SnakeCase(PascalCase))
	addAlias(PascalCase, KebabCase(PascalCase))
	for _, alias := range hardcodedAliases[PascalCase] {
		addAlias(PascalCase, Pluralize(alias))
	}
	slices.Sort(hardcodedAliases[PascalCase])
	return hardcodedAliases[PascalCase]
}
