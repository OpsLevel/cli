package cmd_test

import (
	"fmt"
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

func TestReadAssignPropertyInputValue(t *testing.T) {
	type TestCase struct {
		Input         any
		ExpectedInput string
	}
	testCases := []TestCase{
		{"hello world", `"hello world"`},
		{true, "true"},
		{false, "false"},
		{50, "50"},
		{0, "0"},
		{49.99, "49.99"},
		{0.0, "0"},
	}

	for _, tc := range testCases {
		input := []byte(fmt.Sprintf(`
definition:
  alias: propertyDef
owner:
  alias: propertyOwner
value: %v
`, tc.Input))
		act, err := cmd.ReadPropertyInput(input)
		autopilot.Ok(t, err)
		exp := opslevel.PropertyInput{
			Definition: *opslevel.NewIdentifier("propertyDef"),
			Owner:      *opslevel.NewIdentifier("propertyOwner"),
			Value:      opslevel.JsonString(tc.ExpectedInput),
		}
		autopilot.Equals(t, exp, *act)
	}
}

func TestReadAssignPropertyInputValueIsJSON(t *testing.T) {
	input := []byte(`
definition:
  alias: propertyDef
owner:
  alias: propertyOwner
value: |
  {
      "key1": "val1",
      "key2": {
          "key3": "val2",
          "array": [
              "val3",
              "val4"
          ]
      }
  }
`)
	act, err := cmd.ReadPropertyInput(input)
	autopilot.Ok(t, err)
	exp := opslevel.PropertyInput{
		Definition: *opslevel.NewIdentifier("propertyDef"),
		Owner:      *opslevel.NewIdentifier("propertyOwner"),
		//Value:      `{"key1":"val1","key2":{"array":["val3","val4"],"key3":"val2"}}`,
		Value: "\"{\\n    \\\"key1\\\": \\\"val1\\\",\\n    \\\"key2\\\": {\\n        \\\"key3\\\": \\\"val2\\\",\\n        \\\"array\\\": [\\n            \\\"val3\\\",\\n            \\\"val4\\\"\\n        ]\\n    }\\n}\\n\"",
	}
	autopilot.Equals(t, exp, *act)
}

func TestReadAssignPropertyInputValueIsYAML(t *testing.T) {
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
	act, err := cmd.ReadPropertyInput(input)
	autopilot.Ok(t, err)
	exp := opslevel.PropertyInput{
		Definition: *opslevel.NewIdentifier("propertyDef"),
		Owner:      *opslevel.NewIdentifier("propertyOwner"),
		Value:      `{"key1":"val1","key2":{"array":["val3","val4"],"key3":"val2"}}`,
	}
	autopilot.Equals(t, exp, *act)
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
