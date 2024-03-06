package cmd_test

import (
	"fmt"
	"testing"

	"github.com/opslevel/cli/cmd"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rocktavious/autopilot/v2023"
)

var expectedSchema = opslevel.JSONSchema{
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

func TestReadAssignPropertyInputValueUsingYAMLValue(t *testing.T) {
	type TestCase struct {
		Input         any
		ExpectedInput string
	}
	testCases := []TestCase{
		{"hello world", "\"hello world\""},
		{true, "true"},
		{false, "false"},
		{50, "50"},
		{0, "0"},
	}

	for _, tc := range testCases {
		input := []byte(fmt.Sprintf(`
definition:
  alias: propertyDef
owner:
  alias: propertyOwner
value: %v
`, tc.Input))
		act, err := cmd.ReadPropertyAssignInput(input)
		autopilot.Ok(t, err)
		exp := opslevel.PropertyInput{
			Definition: *opslevel.NewIdentifier("propertyDef"),
			Owner:      *opslevel.NewIdentifier("propertyOwner"),
			Value:      opslevel.JsonString(tc.ExpectedInput),
		}
		autopilot.Equals(t, exp, *act)
	}
}

func TestReadAssignPropertyInputValueUsingYAMLObject(t *testing.T) {
	input := []byte(`
definition:
  alias: propertyDef
owner:
  alias: propertyOwner
value:
  key1: val1
  key2:
    key3: val2
    array:
      - val3
      - val4
`)
	act, err := cmd.ReadPropertyAssignInput(input)
	autopilot.Ok(t, err)
	exp := opslevel.PropertyInput{
		Definition: *opslevel.NewIdentifier("propertyDef"),
		Owner:      *opslevel.NewIdentifier("propertyOwner"),
		Value:      "{\"key1\":\"val1\",\"key2\":{\"array\":[\"val3\",\"val4\"],\"key3\":\"val2\"}}",
	}
	autopilot.Equals(t, exp, *act)
}

func TestReadAssignPropertyInputValueUsingYAMLList(t *testing.T) {
	input := []byte(`
definition:
  alias: propertyDef
owner:
  alias: propertyOwner
value:
  - val1
  - val2
`)
	act, err := cmd.ReadPropertyAssignInput(input)
	autopilot.Ok(t, err)
	exp := opslevel.PropertyInput{
		Definition: *opslevel.NewIdentifier("propertyDef"),
		Owner:      *opslevel.NewIdentifier("propertyOwner"),
		Value:      "[\"val1\",\"val2\"]",
	}
	autopilot.Equals(t, exp, *act)
}

func TestReadPropertyDefinitionInputSchemaUsingYAML(t *testing.T) {
	input := []byte(`
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
		Schema: &expectedSchema,
	}
	autopilot.Equals(t, exp, *act)
}

func TestReadPropertyDefinitionInputSchemaUsingJSON(t *testing.T) {
	input := []byte(`
schema:
  {
      "type": "object",
      "required": [
          "name"
      ],
      "properties": {
          "age": {
              "type": "number"
          },
          "name": {
              "type": "string"
          },
      }
  }
`)
	act, err := cmd.ReadPropertyDefinitionInput(input)
	autopilot.Ok(t, err)
	exp := opslevel.PropertyDefinitionInput{
		Schema: &expectedSchema,
	}
	autopilot.Equals(t, exp, *act)
}
