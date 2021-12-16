package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/alecthomas/jsonschema"
	"github.com/spf13/cobra"
)

type OpsLevelTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type OpsLevelTool struct {
	Name        string `json:"name"`
	Category    string `json:"category"` // TODO: this should be an enum
	Url         string `json:"url"`
	Environment string `json:"environment,omitempty"`
}

type OpsLevelRepository struct {
	Name     string `json:"name"`
	Path     string `json:"path,omitempty"`
	Provider string `json:"provider"` // TODO: this should be an enum
}

type OpsLevelService struct {
	Name         string               `json:"name"`
	Description  string               `json:"description"`
	Owner        string               `json:"owner"`
	Lifecycle    string               `json:"lifecycle"`
	Tier         string               `json:"tier"`
	Product      string               `json:"product"`
	Language     string               `json:"language"`
	Framework    string               `json:"framework"`
	Aliases      []string             `json:"aliases"` // JQ expressions that return a single string or a []string
	Tags         []OpsLevelTag        `json:"tags"`
	Tools        []OpsLevelTool       `json:"tools"`        // JQ expressions that return a single map[string]string or a []map[string]string
	Repositories []OpsLevelRepository `json:"repositories"` // JQ expressions that return a single string or []string or map[string]string or a []map[string]string
	// TODO: Dependencies
}

type OpsLevelConfig struct {
	Version int             `json:"version"`
	Service OpsLevelService `json:"service"`
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
