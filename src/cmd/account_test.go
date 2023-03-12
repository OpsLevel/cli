package cmd

import (
	"bytes"
	"flag"
	"fmt"
	ol "github.com/opslevel/opslevel-go/v2023"
	"github.com/rocktavious/autopilot/v2022"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	log.Logger = log.Output(output)
	flag.Parse()
	teardown := autopilot.Setup()
	defer teardown()
	os.Exit(m.Run())
}

func Templated(input string) string {
	response, err := autopilot.Templater.Use(input)
	if err != nil {
		panic(err)
	}
	return response
}

func TemplatedResponse(response string) autopilot.ResponseWriter {
	return func(w http.ResponseWriter) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, Templated(response))
	}
}

func ABetterTestClient(endpoint string, response string) *ol.Client {
	return ol.NewGQLClient(ol.SetAPIToken("x"), ol.SetMaxRetries(0), ol.SetURL(autopilot.RegisterEndpoint(fmt.Sprintf("/LOCAL_TESTING/%s", endpoint),
		TemplatedResponse(response),
		autopilot.SkipRequestValidation())))
}

func ExecuteCommand(cmd *cobra.Command, args ...string) (string, error) {
	buffer := new(bytes.Buffer)
	cmd.SetOut(buffer)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return strings.TrimSpace(buffer.String()), err
}

func Test_LifecycleCmd_Json(t *testing.T) {
	// Assign
	response := `{"data": {
	"account": {
		"lifecycles": [
			{{ template "lifecycle-pre-alpha" }},
			{{ template "lifecycle-alpha" }},
			{{ template "lifecycle-beta" }},
			{{ template "lifecycle-ga" }},
			{{ template "lifecycle-eol" }}
		]
	}
}}`
	exp := Templated(`[
			{{ template "lifecycle-pre-alpha" }},
			{{ template "lifecycle-alpha" }},
			{{ template "lifecycle-beta" }},
			{{ template "lifecycle-ga" }},
			{{ template "lifecycle-eol" }}
	]`)
	client := ABetterTestClient("account/lifecycle_json", response)
	// Act
	out, err := ExecuteCommand(NewLifecycleCmd(client), "-o", "json")
	// Assert
	autopilot.Ok(t, err)
	autopilot.Equals(t, exp, out)
}
