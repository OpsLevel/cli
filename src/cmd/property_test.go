package cmd_test

import (
	"testing"

	"github.com/opslevel/cli/cmd"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rocktavious/autopilot/v2023"
)

var expSchemaMap = &opslevel.JSONSchema{
	"type":     "object",
	"required": []any{"name"},
	"properties": map[string]any{
		"name": map[string]any{
			"type": "string",
		},
		"age": map[string]any{
			"type": "number",
		},
	},
}

func TestReadPropertyDefinitionInputSchemaIsJSON(t *testing.T) {
	input := []byte(`
name: hello world
schema: |
  {
      "type": "object",
      "required": [
          "name"
      ],
      "properties": {
          "name": {
              "type": "string"
          },
          "age": {
              "type": "number"
          }
      }
  }
`)
	act, err := cmd.ReadPropertyDefinitionInput(input)
	autopilot.Ok(t, err)
	exp := opslevel.PropertyDefinitionInput{
		Name:   opslevel.RefOf("hello world"),
		Schema: expSchemaMap,
	}
	autopilot.Equals(t, exp, *act)
}

func TestReadPropertyDefinitionInputSchemaIsYAML(t *testing.T) {
	input := []byte(`
name: hello world
schema:
  type: object
  required:
    - name
  properties:
    name:
      type: string
    age:
      type: number
`)
	act, err := cmd.ReadPropertyDefinitionInput(input)
	autopilot.Ok(t, err)
	exp := opslevel.PropertyDefinitionInput{
		Name:   opslevel.RefOf("hello world"),
		Schema: expSchemaMap,
	}
	autopilot.Equals(t, exp, *act)
}
