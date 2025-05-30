package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/types"
	"github.com/opslevel/opslevel-go/v2025"
	"github.com/relvacode/iso8601"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type regoInput struct {
	Files []string               `json:"files"`
	Data  map[string]interface{} `json:"data"`
}

type gitlabResponse struct {
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Language    string             `json:"language,omitempty"`
	Languages   map[string]float64 `json:"languages,omitempty"`
}

type githubResponse struct {
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Language    string             `json:"language,omitempty"`
	Languages   map[string]float64 `json:"languages,omitempty"`
}

type commonRepoMetadata struct {
	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	Language    string             `json:"language,omitempty"`
	Languages   map[string]float64 `json:"languages,omitempty"`
}

type astMarshallable interface {
	commonRepoMetadata | opslevel.Level
}

func (r *githubResponse) toRepoMetadata() commonRepoMetadata {
	return commonRepoMetadata{
		Name:        r.Name,
		Description: r.Description,
		Language:    r.Language,
		Languages:   r.Languages,
	}
}

func (r *gitlabResponse) toRepoMetadata() commonRepoMetadata {
	return commonRepoMetadata{
		Name:        r.Name,
		Description: r.Description,
		Language:    r.Language,
		Languages:   r.Languages,
	}
}

func toASTValue[T astMarshallable](input T) (ast.Value, error) {
	marshalData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	reader := strings.NewReader(string(marshalData))
	v, err := ast.ValueFromReader(reader)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func getKeyFromMapByMaxValue(m map[string]float64) string {
	maxValue := float64(0)
	var output string
	for k, v := range m {
		if v > maxValue {
			maxValue = v
			output = k
		}
	}
	return output
}

func convertToPercentage(m map[string]float64) map[string]float64 {
	total := float64(0)
	output := make(map[string]float64)
	for _, v := range m {
		total += v
	}

	for k, v := range m {
		output[k] = math.Floor(v/total*100*100) / 100
	}
	return output
}

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Invoked along with a Rego policy to create and output JSON",
	Long: `We are extending OPA Rego to work with OpsLevel constructs:

Examples:
    opslevel run policy -f policy.rego | jq
    opslevel run policy -f policy.rego -i /tmp/input.json -o ./output.json
	`,
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		filePath, err := flags.GetString("file")
		cobra.CheckErr(err)
		inputFilePath, err := flags.GetString("input")
		cobra.CheckErr(err)
		inputJSON := &map[string]interface{}{}
		if inputFilePath != "" {
			inputRead, err := os.ReadFile(inputFilePath)
			cobra.CheckErr(err)
			cobra.CheckErr(json.Unmarshal(inputRead, inputJSON))
		}
		outputFilePath, err := flags.GetString("output")
		cobra.CheckErr(err)
		policy, err := os.ReadFile(filePath)
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
		input.Data = *inputJSON
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
			rego.Function1(
				&rego.Function{
					Name:    "opslevel.repo.gitlab",
					Decl:    types.NewFunction(types.Args(types.S), types.A),
					Memoize: true,
				},
				RegoFuncGetGitlabRepo),
			rego.Function1(
				&rego.Function{
					Name: "opslevel.service_maturity_level",
					Decl: types.NewFunction(types.Args(types.S), types.A),
				},
				RegoFuncGetMaturity),
			rego.Function2(
				&rego.Function{
					Name: "opslevel.time.diff",
					Decl: types.NewFunction(types.Args(types.S, types.S), types.A),
				},
				RegoFuncTimeDiff),
			rego.Input(input),
		)
		rs, err := rego.Eval(context.Background())
		cobra.CheckErr(err)
		b, err := json.Marshal(rs[0].Expressions[0].Value) // TODO: need more advanced handling of multiple things in json and reading from stdin
		cobra.CheckErr(err)

		if outputFilePath == "-" {
			fmt.Println(string(b))
		} else {
			main := newFile(outputFilePath, false)
			defer main.Close()
			main.WriteString(string(b))
		}
	},
}

func RegoFuncReadFile(ctx rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
	if str, ok := a.Value.(ast.String); ok {
		if _, err := os.Stat(string(str)); err != nil {
			log.Warn().Msgf("%s", err)
		} else {
			file, err := os.Open(string(str))
			defer file.Close()
			if err != nil {
				log.Error().Err(err).Msg("")
				return nil, err
			}

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
		log.Error().Err(err).Msg("")
		return nil, err
	}
	if err := ast.As(b.Value, &repo); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if org == "" {
		err := fmt.Errorf("opslevel.repo.github(\"%s\", \"%s\") failed: Please provide a valid org", org, repo)
		log.Error().Err(err).Msgf("")
		return nil, err
	}

	if repo == "" {
		err := fmt.Errorf("opslevel.repo.github(\"%s\", \"%s\") failed: Please provide a valid repo", org, repo)
		log.Error().Err(err).Msg("")
		return nil, err
	}

	githubToken := viper.GetString("github-token")
	authorizationHeader := fmt.Sprintf("token %s", githubToken)
	githubAPIUrl := fmt.Sprintf("https://api.github.com/repos/%v/%v", org, repo)

	var result githubResponse

	response, err := getClientRest().R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("Authorization", authorizationHeader).
		SetResult(&result).
		Get(githubAPIUrl)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if response.IsError() {
		err := fmt.Errorf("error requesting Github repo metadata. CODE: %d: REASON: %s", response.StatusCode(), response)
		log.Error().Err(err).Msgf("")
		return nil, err
	}

	languagesResponse, err := getClientRest().R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("Authorization", authorizationHeader).
		SetResult(&result.Languages).
		Get(githubAPIUrl + "/languages")
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if languagesResponse.IsError() {
		err := fmt.Errorf("error requesting Github repo languages. CODE: %d: REASON: %s", response.StatusCode(), response)
		log.Error().Err(err).Msgf("")
		return nil, err
	}

	result.Languages = convertToPercentage(result.Languages)
	repoMetadata := result.toRepoMetadata()

	v, err := toASTValue(repoMetadata)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return ast.NewTerm(v), nil
}

func RegoFuncGetGitlabRepo(ctx rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
	var path string
	if err := ast.As(a.Value, &path); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if path == "" {
		err := fmt.Errorf("opslevel.repo.gitlab(\"%s\") failed: Please provide a valid Gitlab project path", path)
		log.Error().Err(err).Msgf("")
		return nil, err
	}

	escapedPath := url.QueryEscape(path)
	gitlabToken := viper.GetString("gitlab-token")
	gitlabAPIUrl := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s", escapedPath)

	var result gitlabResponse

	response, err := getClientRest().R().
		SetHeader("PRIVATE-TOKEN", gitlabToken).
		SetResult(&result).
		Get(gitlabAPIUrl)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if response.IsError() {
		err := fmt.Errorf("error requesting Gitlab repo metadata. CODE: %d: REASON: %s", response.StatusCode(), response)
		log.Error().Err(err).Msgf("")
		return nil, err
	}

	languagesResponse, err := getClientRest().R().
		SetHeader("PRIVATE-TOKEN", gitlabToken).
		SetResult(&result.Languages).
		Get(gitlabAPIUrl + "/languages")
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if languagesResponse.IsError() {
		err := fmt.Errorf("error requesting Gitlab repo languages. CODE: %d: REASON: %s", response.StatusCode(), response)
		log.Error().Err(err).Msgf("")
		return nil, err
	}

	result.Language = getKeyFromMapByMaxValue(result.Languages)
	repoMetadata := result.toRepoMetadata()
	v, err := toASTValue(repoMetadata)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return ast.NewTerm(v), nil
}

func RegoFuncGetMaturity(ctx rego.BuiltinContext, a *ast.Term) (*ast.Term, error) {
	var alias string
	if err := ast.As(a.Value, &alias); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	if alias == "" {
		err := fmt.Errorf("opslevel.service_maturity_level(\"%s\") failed: Please provide a valid alias", alias)
		log.Error().Err(err).Msgf("")
		return nil, nil
	}

	client := getClientGQL()
	service, err := client.GetServiceMaturityWithAlias(alias)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	v, err := toASTValue(service.MaturityReport.OverallLevel)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	return ast.NewTerm(v), nil
}

func RegoFuncTimeDiff(ctx rego.BuiltinContext, s *ast.Term, e *ast.Term) (*ast.Term, error) {
	var start string
	if err := ast.As(s.Value, &start); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	var end string
	if err := ast.As(e.Value, &end); err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}

	startTime, err := iso8601.ParseString(start)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	endTime, err := iso8601.ParseString(end)
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, err
	}
	duration := endTime.Sub(startTime)
	return ast.FloatNumberTerm(duration.Seconds()), nil
}

func init() {
	runCmd.AddCommand(policyCmd)

	policyCmd.Flags().StringP("file", "f", "-", "File to read Rego policy from. Defaults to reading from stdin.")
	policyCmd.Flags().StringP("input", "i", "", "File to read extra JSON data input to be used in Rego policy. Defaults to not reading anything.")
	policyCmd.Flags().StringP("output", "o", "-", "File to write Rego policy output to. Defaults to writing to stdout.")
	policyCmd.PersistentFlags().String("github-token", "", "The Github API token to use when calling opslevel.repo.github function within a Rego policy. Overrides environment variable 'GITHUB_API_TOKEN'")
	policyCmd.PersistentFlags().String("gitlab-token", "", "The Gitlab API token to use when calling opslevel.repo.gitlab function within a Rego policy. Overrides environment variable 'GITLAB_API_TOKEN'")

	viper.BindPFlags(policyCmd.PersistentFlags())
	viper.BindEnv("github-token", "GITHUB_API_TOKEN")
	viper.BindEnv("gitlab-token", "GITLAB_API_TOKEN")
}
