package common_test

import (
	"io"
	"os"
	"strings"
	"testing"

	common "github.com/opslevel/cli/common"
	"github.com/rocktavious/autopilot"
)

func captureOutput() string {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	common.PrettyPrint("< > & alan was here & < >")

	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = stdOut

	return string(out)
}

func TestPrettyPrint(t *testing.T) {
	// Arrange
	trimmedTestString := strings.TrimRight(captureOutput(), "\n") // captureOuput adds an extra newline
	// Act
	// Assert
	autopilot.Equals(t, "\"< > & alan was here & < >\"", trimmedTestString)
}
