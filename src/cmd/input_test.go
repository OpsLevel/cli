package cmd_test

import (
	"fmt"
	"testing"

	"github.com/opslevel/cli/cmd"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rocktavious/autopilot/v2023"
)

type MockResourceInput struct {
	Age    int                 `json:"age" yaml:"age"`
	Age2   *int                `json:"age2" yaml:"age2"`
	Name   string              `json:"name" yaml:"name"`
	Schema opslevel.JSONSchema `json:"schema" yaml:"schema"`
	Value  opslevel.JsonString `json:"value" yaml:"value"`
}

func TestSetResourceOnMap(t *testing.T) {
	input := []byte(`
name: hello world
age: 50
pets:
  dog:
    age: 2
    name: bella
  cat:
    age: 1
    name: daisy
`)
	act, err := cmd.ReadResource[map[string]any](input)
	autopilot.Ok(t, err)
	exp := map[string]any{
		"name": "hello world",
		"age":  50,
		"pets": map[string]any{
			"dog": map[string]any{
				"age":  2,
				"name": "bella",
			},
			"cat": map[string]any{
				"age":  1,
				"name": "daisy",
			},
		},
	}
	autopilot.Equals(t, exp, *act)
}

func TestSetResourceOnStruct(t *testing.T) {
	input := []byte(`
name: hello world
age: 50
age2: 60
`)
	act, err := cmd.ReadResource[MockResourceInput](input)
	autopilot.Ok(t, err)
	exp := MockResourceInput{Name: "hello world", Age: 50, Age2: opslevel.RefOf(60)}
	autopilot.Equals(t, exp, *act)
}

func TestSetResourceOnStructWithSchemaUsingJSON(t *testing.T) {
	input := []byte(`
schema: |
  {
      "active": true,
      "age": 50
  }
`)
	act, err := cmd.ReadResourceHandleJSONFields[MockResourceInput](input)
	autopilot.Ok(t, err)
	exp := MockResourceInput{Schema: opslevel.JSONSchema{"active": true, "age": float64(50)}}
	autopilot.Equals(t, exp, *act)
}

func TestSetResourceOnStructWithSchemaUsingYAML(t *testing.T) {
	input := []byte(`
schema:
  active: true
  age: 50
`)
	act, err := cmd.ReadResourceHandleJSONFields[MockResourceInput](input)
	autopilot.Ok(t, err)
	exp := MockResourceInput{Schema: map[string]any{"active": true, "age": 50}}
	autopilot.Equals(t, exp, *act)
}

func TestSetResourceOnStructWithValueUsingYAMLValue(t *testing.T) {
	type TestCase struct {
		Value         any
		ExpectedValue string
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
value: %v
`, tc.Value))
		act, err := cmd.ReadResourceHandleJSONFields[MockResourceInput](input)
		autopilot.Ok(t, err)
		exp := MockResourceInput{Value: opslevel.JsonString(tc.ExpectedValue)}
		autopilot.Equals(t, exp, *act)
	}
}

func TestSetResourceOnStructWithValueUsingYAMLObject(t *testing.T) {
	input := []byte(`
value:
 key1: val1
 key2:
   key3: val2
   array:
     - val3
     - val4
`)
	act, err := cmd.ReadResourceHandleJSONFields[MockResourceInput](input)
	autopilot.Ok(t, err)
	exp := MockResourceInput{Value: "{\"key1\":\"val1\",\"key2\":{\"array\":[\"val3\",\"val4\"],\"key3\":\"val2\"}}"}
	autopilot.Equals(t, exp, *act)
}

func TestSetResourceOnStructWithValueUsingYAMLList(t *testing.T) {
	input := []byte(`
value:
  - val1
  - val2
`)
	act, err := cmd.ReadResourceHandleJSONFields[MockResourceInput](input)
	autopilot.Ok(t, err)
	exp := MockResourceInput{Value: "[\"val1\",\"val2\"]"}
	autopilot.Equals(t, exp, *act)
}

func TestReadResourceInput(t *testing.T) {
	input := []byte(`
name: hello world
age: 50
age2: 60
`)
	act, err := cmd.ReadResourceInput[MockResourceInput](input)
	autopilot.Ok(t, err)
	exp := MockResourceInput{Name: "hello world", Age: 50, Age2: opslevel.RefOf(60)}
	autopilot.Equals(t, exp, *act)
}
