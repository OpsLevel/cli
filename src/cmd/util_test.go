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

	r, oldStdout := redirectStdout()
	defer r.Close()

	cmd.RootCmd.SetArgs(cliArgs)
	err := cmd.RootCmd.Execute()

	output := captureOutput(r, oldStdout)
	return output, err
}

// redirectStdout redirects os.Stdout to a pipe and returns the read and write ends of the pipe.
func redirectStdout() (*os.File, *os.File) {
	r, w, _ := os.Pipe()
	oldStdout := os.Stdout
	os.Stdout = w
	return r, oldStdout
}

// captureOutput reads from r until EOF and returns the result as a string.
func captureOutput(r *os.File, oldStdout *os.File) []byte {
	w := os.Stdout
	os.Stdout = oldStdout
	w.Close()
	gotOutput, _ := io.ReadAll(r)
	return gotOutput
	// return strings.TrimSpace(string(gotOutput))
}

// convert a simple API response to a string
func asString(data []byte) string {
	return strings.TrimSpace(string(data))
}

// convert JSON response from API to OpsLevel resource
func jsonToResource[T any](jsonData []byte) *T {
	var resource T
	if err := json.Unmarshal(jsonData, &resource); err != nil {
		return nil
	}
	return &resource
}

// write OpsLevel resource to YAML file for commands that read in a file
func writeToYaml(userFileName string, opslevelResource any) error {
	yamlData, err := yaml.Marshal(&opslevelResource)
	if err != nil {
		return err
	}
	if err = os.WriteFile(userFileName, yamlData, 0o644); err != nil {
		return err
	}
	return nil
}
