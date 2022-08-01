package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
)

type regoInput struct {
	Files []string `json:"files"`
}

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Invoked along with a Rego policy to create and output JSON",
	Long: `We are extending OPA Rego to work with OpsLevel constructs:

opslevel run policy -f policy.rego | jq
	`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		filePath, err := flags.GetString("file")
		cobra.CheckErr(err)
		policy, err := ioutil.ReadFile(filePath)
		cobra.CheckErr(err)
		input := regoInput{}
		err = filepath.Walk(".",
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				input.Files = append(input.Files, path)
				return nil
			})
		cobra.CheckErr(err)
		rego := rego.New(
			rego.Query("data.opslevel"),
			rego.Module("test.rego",
				string(policy),
			),
			rego.Input(input),
		)
		rs, err := rego.Eval(context.Background())
		cobra.CheckErr(err)
		b, err := json.Marshal(rs[0].Expressions[0].Value) //TODO: need more advanced handling of multiple things in json and reading from stdin
		cobra.CheckErr(err)
		fmt.Println(string(b))

	},
}

func init() {
	runCmd.AddCommand(policyCmd)

	policyCmd.Flags().StringP("file", "f", "-", "File to read Rego policy from. Defaults to reading from stdin.")
}
