package cmd_test

import (
	"testing"

	"github.com/opslevel/cli/cmd"
	"github.com/opslevel/opslevel-go/v2024"
	"github.com/rocktavious/autopilot/v2023"
)

type MockResourceInput struct {
	Age  int    `json:"age" yaml:"age"`
	Age2 *int   `json:"age2" yaml:"age2"`
	Name string `json:"name" yaml:"name"`
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

func TestReadResourceInput(t *testing.T) {
	input := []byte(`
name: hello world
age: 50
age2: 60
`)
	act, err := cmd.ReadResourceInput[MockResourceInput]()
	autopilot.Ok(t, err)
	exp := MockResourceInput{Name: "hello world", Age: 50, Age2: opslevel.RefOf(60)}
	autopilot.Equals(t, exp, *act)
}
