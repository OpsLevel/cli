package cmd_test

import (
	"github.com/opslevel/cli/cmd"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rocktavious/autopilot/v2023"
	"testing"
)

type MockResourceInput struct {
	Name string `json:"name" yaml:"name"`
	Age  int    `json:"age" yaml:"age"`
	Age2 *int   `json:"age2" yaml:"age2"`
}

type MockResourceInputWithSchema struct {
	Name   string              `json:"name" yaml:"name"`
	Schema opslevel.JSONSchema `json:"schema" yaml:"schema"`
}

func TestSetResourceOnMap(t *testing.T) {
	input := []byte(`
name: hello world
age: 50
`)
	act, err := cmd.ReadResource[map[string]any](input)
	autopilot.Ok(t, err)
	exp := make(map[string]any)
	exp["name"] = "hello world"
	exp["age"] = 50
	autopilot.Equals(t, exp, *act)
}

func TestSetResourceOnStruct(t *testing.T) {
	input := []byte(`
name: hello world
age: 50
`)
	act, err := cmd.ReadResource[MockResourceInput](input)
	autopilot.Ok(t, err)
	exp := MockResourceInput{Name: "hello world", Age: 50}
	autopilot.Equals(t, exp, *act)
}

func TestSetResourceOnStructSetNullable(t *testing.T) {
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
name: hello world
schema: {"age": 50, "active": true}
`)
	act, err := cmd.ReadResource[MockResourceInputWithSchema](input)
	autopilot.Ok(t, err)
	jsonSchema := make(opslevel.JSONSchema)
	jsonSchema["age"] = 50
	jsonSchema["active"] = true
	exp := MockResourceInputWithSchema{Name: "hello world", Schema: jsonSchema}
	autopilot.Equals(t, exp, *act)
}

func TestSetResourceOnStructWithSchemaUsingYAML(t *testing.T) {
	type MockInput struct {
		Name   string              `json:"name" yaml:"name"`
		Schema opslevel.JSONSchema `json:"schema" yaml:"schema"`
	}
	input := []byte(`
name: hello world
schema:
  age: 50
  active: false
`)
	act, err := cmd.ReadResource[MockInput](input)
	autopilot.Ok(t, err)
	jsonSchema := make(opslevel.JSONSchema)
	jsonSchema["age"] = 50
	jsonSchema["active"] = false
	exp := MockInput{Name: "hello world", Schema: jsonSchema}
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
