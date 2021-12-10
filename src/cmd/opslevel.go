package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/jsonschema"
	"github.com/opslevel/opslevel-go"
	"github.com/spf13/cobra"
)

type OpsLevelAlias struct{}

func (OpsLevelAlias) JSONSchemaType() *jsonschema.Type {
	return &jsonschema.Type{
		Type:      "string",
		MinLength: 1,
		MaxLength: 255,
	}
}

type OpsLevelTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OpsLevelTool struct {
	Name        string                `json:"name"`
	Category    opslevel.ToolCategory `json:"category"`
	Url         string                `json:"url"`
	Environment string                `json:"environment,omitempty"`
}

type OpsLevelRepository struct {
	Name     string `json:"name"`
	Path     string `json:"path,omitempty"`
	Provider string `json:"provider"` // TODO: this should be an enum
}

type OpsLevelServiceDependancy struct {
	Alias string `json:"alias"`
}

type OpsLevelService struct {
	Name         string                      `json:"name"`
	Description  string                      `json:"description,omitempty"`
	Owner        string                      `json:"owner,omitempty"`
	Lifecycle    string                      `json:"lifecycle,omitempty"`
	Tier         string                      `json:"tier,omitempty"`
	Product      string                      `json:"product,omitempty"`
	Language     string                      `json:"language,omitempty"`
	Framework    string                      `json:"framework,omitempty"`
	Aliases      []OpsLevelAlias             `json:"aliases,omitempty"`
	Tags         []OpsLevelTag               `json:"tags,omitempty"`
	Tools        []OpsLevelTool              `json:"tools,omitempty"`
	Repositories []OpsLevelRepository        `json:"repositories,omitempty"`
	Dependencies []OpsLevelServiceDependancy `json:"dependencies,omitempty"`
}

type OpsLevelConfig struct {
	Version    int                `json:"version"`
	Service    OpsLevelService    `json:"service" jsonschema:"oneof_required=service"`
	Repository OpsLevelRepository `json:"repository" jsonschema:"oneof_required=repository"`
}

var opslevelSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Print the jsonschema for opslevel.yml file",
	Long:  "Print the jsonschema for opslevel.yml file",
	Run: func(cmd *cobra.Command, args []string) {
		schema := jsonschema.Reflect(&OpsLevelConfig{})
		jsonBytes, jsonErr := json.MarshalIndent(schema, "", "  ")
		cobra.CheckErr(jsonErr)
		fmt.Println(string(jsonBytes))
	},
}

func init() {
	exportCmd.AddCommand(opslevelSchemaCmd)
}
