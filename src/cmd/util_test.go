package cmd_test

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/opslevel/cli/cmd"
	"gopkg.in/yaml.v2"
)

type Operation string

const (
	Assign   Operation = "assign"
	Create   Operation = "create"
	Delete   Operation = "delete"
	Example  Operation = "example"
	Get      Operation = "get"
	List     Operation = "list"
	Update   Operation = "update"
	Unassign Operation = "unassign"
)

// execute any OpsLevel CLI command
func execCmd(command Operation, resource string, inputs ...string) ([]byte, error) {
	cliArgs := []string{string(command), resource}
	cliArgs = append(cliArgs, inputs...)

	r, oldStdout, err := redirectStdout()
	defer r.Close()
	if err != nil {
		return nil, err
	}

	cmd.RootCmd.SetArgs(cliArgs)
	if err = cmd.RootCmd.Execute(); err != nil {
		return nil, err
	}

	return captureOutput(r, oldStdout)
}

// redirectStdout redirects os.Stdout to a pipe and returns the read and write ends of the pipe.
func redirectStdout() (*os.File, *os.File, error) {
	r, w, err := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w
	return r, oldStdout, err
}

// captureOutput reads from r until EOF and returns the result as a string.
func captureOutput(r *os.File, oldStdout *os.File) ([]byte, error) {
	w := os.Stdout
	os.Stdout = oldStdout
	w.Close()
	return io.ReadAll(r)
}

// convert a simple API response to a string
func asString(data []byte) string {
	return strings.TrimSpace(string(data))
}

// convert JSON response from API to OpsLevel resource
func jsonToResource[T any](jsonData []byte) (*T, error) {
	var resource T
	if err := json.Unmarshal(jsonData, &resource); err != nil {
		return nil, err
	}
	return &resource, nil
}

// write OpsLevel resource to YAML file for commands that read in a file
func writeToYaml(givenFileName string, opslevelResource any) error {
	yamlData, err := yaml.Marshal(&opslevelResource)
	if err != nil {
		return err
	}
	if err = os.WriteFile(givenFileName, yamlData, 0o644); err != nil {
		return err
	}
	return nil
}
