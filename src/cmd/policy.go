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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	}
	if err := ast.As(b.Value, &repo); err != nil {
		return nil, err
	}

	if org == "" {
		log.Error().Msgf("opslevel.repo.github(\"%s\", \"%s\") failed: Please provide a valid org", org, repo)
		return nil, nil
	}

	if repo == "" {
		log.Error().Msgf("opslevel.repo.github(\"%s\", \"%s\") failed: Please provide a valid repo", org, repo)
		return nil, nil
	}

	githubToken := viper.GetString("github-token")
	authorizationHeader := fmt.Sprintf("token %s", githubToken)
	githubAPIUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v", org, repo)

	response, err := getClientRest().R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("Authorization", authorizationHeader).
		Get(githubAPIUrl)
	cobra.CheckErr(err)

	if response.IsError() == true {
		log.Error().Msgf("error requesting Github repo metadata. CODE: %d: REASON: %s", response.StatusCode(), response)
		return nil, nil
	}

	reader := strings.NewReader(response.String())
	v, err := ast.ValueFromReader(reader)
	return ast.NewTerm(v), nil
}

func init() {
	runCmd.AddCommand(policyCmd)

	policyCmd.Flags().StringP("file", "f", "-", "File to read Rego policy from. Defaults to reading from stdin.")
	policyCmd.PersistentFlags().String("github-token", "", "The Github API token to use when calling opslevel.repo.github function within a Rego policy. Overrides environment variable 'GITHUB_API_TOKEN'")

	viper.BindPFlags(policyCmd.PersistentFlags())
	viper.BindEnv("github-token", "GITHUB_API_TOKEN")
}
