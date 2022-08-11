package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/rs/zerolog/log"
	cobra "github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type regoInput struct {
	Files []string `json:"files"`
}

type GithubRepo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Language    string `json:"language"`
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
			rego.Function1(
				&rego.Function{
					Name: "opslevel.read_file",
					Decl: types.NewFunction(types.Args(types.S), types.S),
				},
				RegoFuncReadFile),
			rego.Function2(
				&rego.Function{
					Name:    "opslevel.repo.github",
					Decl:    types.NewFunction(types.Args(types.S, types.S), types.A),
					Memoize: true,
				},
				RegoFuncGetGithubRepo),
			rego.Input(input),
		)
		rs, err := rego.Eval(context.Background())
		cobra.CheckErr(err)
		b, err := json.Marshal(rs[0].Expressions[0].Value) //TODO: need more advanced handling of multiple things in json and reading from stdin
		cobra.CheckErr(err)
		fmt.Println(string(b))

	},
}

func RegoFuncReadFile(ctx rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
	if str, ok := a.Value.(ast.String); ok {
		if _, err := os.Stat(string(str)); err != nil {
			log.Warn().Msgf("%s", err)
		} else {
			file, err := os.Open(string(str))
			defer file.Close()
			cobra.CheckErr(err)

			var lines []*ast.Term
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				lines = append(lines, ast.StringTerm(scanner.Text()))
			}
			return ast.ArrayTerm(lines...), nil
		}
	}
	return nil, nil
}

func RegoFuncGetGithubRepo(ctx rego.BuiltinContext, a, b *ast.Term) (*ast.Term, error) {

	var org, repo string

	if err := ast.As(a.Value, &org); err != nil {
		return nil, err
	} else if ast.As(b.Value, &repo); err != nil {
		return nil, err
	}

	flags := runCmd.Flags()
	githubToken, err := flags.GetString("github-token")
	cobra.CheckErr(err)

	response, err := getClientRest().R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("Authorization", fmt.Sprintf("token %v", githubToken)).
		Get(fmt.Sprintf("https://api.github.com/repos/%v/%v", org, repo))
	cobra.CheckErr(err)

	if response.IsError() == true {
		log.Error().Msgf("%d: %s", response.StatusCode(), response)
		return nil, err
	}

	reader := strings.NewReader(response.String())
	v, err := ast.ValueFromReader(reader)
	return ast.NewTerm(v), nil
}

func init() {
	runCmd.AddCommand(policyCmd)

	policyCmd.Flags().StringP("file", "f", "-", "File to read Rego policy from. Defaults to reading from stdin.")
}
