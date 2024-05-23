package cmd

import (
	"encoding/json"
	"os"
	"time"

	"github.com/opslevel/opslevel-go/v2024"

	"github.com/creasty/defaults"
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var integrationUrl string

type Deployer struct {
	Email string `validate:"required" json:"email" default:"automation@opslevel.com"`
	Name  string `json:"name,omitempty"`
}

// Commit represents the commit being deployed
type Commit struct {
	SHA            string     `json:"sha,omitempty"`
	Message        string     `json:"message,omitempty"`
	Branch         string     `json:"branch,omitempty"`
	Date           *time.Time `json:"date,omitempty"`
	CommitterName  string     `json:"committer_name,omitempty" yaml:"committer-name"`
	CommitterEmail string     `json:"committer_email,omitempty" yaml:"committer-email"`
	AuthorName     string     `json:"author_name,omitempty" yaml:"author-name"`
	AuthorEmail    string     `json:"author_email,omitempty" yaml:"author-email"`
	AuthoringDate  *time.Time `json:"authoring_date,omitempty" yaml:"authoring-date"`
}

// DeployEvent represents a structured request to the OpsLevel deploys webhook endpoint
type DeployEvent struct {
	Service      string    `validate:"required" json:"service"`
	Deployer     Deployer  `validate:"required" json:"deployer"`
	DeployedAt   time.Time `validate:"required" json:"deployed_at" yaml:"deployed-at"`
	Description  string    `validate:"required" json:"description" default:"Event Created by OpsLevel CLI"`
	Environment  string    `json:"environment,omitempty"`
	DeployURL    string    `json:"deploy_url,omitempty" yaml:"deploy-url"`
	DeployNumber string    `json:"deploy_number,omitempty" yaml:"deploy-number"`
	Commit       Commit    `json:"commit,omitempty"`
	DedupID      string    `json:"dedup_id,omitempty" yaml:"dedup-id"`
}

var deployCreateCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create deployment events (report deploy)",
	Long:  "Create deployment events (report deploy)",
	Run: func(cmd *cobra.Command, args []string) {
		if integrationUrl == "" {
			log.Error().Msg("Please provide '--integration-url' to send the deployment information to")
			os.Exit(1)
		}

		evt, err := readCreateConfigAsDeployEvent()
		cobra.CheckErr(err)
		if dryrun := viper.GetBool("dry-run"); dryrun {
			b, _ := json.Marshal(evt)
			log.Info().Msgf("%s", string(b))
			return
		}
		body, err := json.Marshal(evt)
		cobra.CheckErr(err)
		response := &opslevel.RestResponse{}
		resp, err := getClientRest().R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			SetResult(response).
			Post(integrationUrl)
		cobra.CheckErr(err)
		if !resp.IsSuccess() {
			log.Error().Msg(resp.String())
			os.Exit(1)
		}
		log.Info().Msgf("Successfully registered deploy event for '%s'", evt.Service)
	},
}

func init() {
	createCmd.AddCommand(deployCreateCmd)

	deployCreateCmd.Flags().Bool("dry-run", false, "if true data will be logged and not sent to the integration-url (OPSLEVEL_DRY_RUN)")
	viper.BindEnv("dry-run", "OPSLEVEL_DRY_RUN", "OL_DRY_RUN")

	deployCreateCmd.Flags().String("commit-message", "", "git commit message associated with the event (OPSLEVEL_DEPLOYER_EMAIL)")
	viper.BindEnv("commit-message", "OPSLEVEL_COMMIT_MESSAGE", "OL_COMMIT_MESSAGE")

	deployCreateCmd.Flags().String("commit-sha", "", "git commit sha associated with the event (OPSLEVEL_DEPLOYER_NAME)")
	viper.BindEnv("commit-sha", "OPSLEVEL_COMMIT_SHA", "OL_COMMIT_SHA")

	deployCreateCmd.Flags().String("dedup-id", "", "dedup id of the event (OPSLEVEL_DEDUP_ID)")
	viper.BindEnv("dedup-id", "OPSLEVEL_DEDUP_ID", "OL_DEDUP_ID")

	deployCreateCmd.Flags().String("deploy-number", "", "deploy number of the event (OPSLEVEL_DEPLOY_NUMBER)")
	viper.BindEnv("deploy-number", "OPSLEVEL_DEPLOY_NUMBER", "OL_DEPLOY_NUMBER")

	deployCreateCmd.Flags().String("deploy-url", "", "url the event will link back to (OPSLEVEL_DEPLOY_URL)")
	viper.BindEnv("deploy-url", "OPSLEVEL_DEPLOY_URL", "OL_DEPLOY_URL")

	deployCreateCmd.Flags().String("deployer-email", "", "deployer email who created the event (OPSLEVEL_DEPLOYER_EMAIL)")
	viper.BindEnv("deployer-email", "OPSLEVEL_DEPLOYER_EMAIL", "OL_DEPLOYER_EMAIL")

	deployCreateCmd.Flags().String("deployer-name", "", "deployer name who created the event (OPSLEVEL_DEPLOYER_NAME)")
	viper.BindEnv("deployer-name", "OPSLEVEL_DEPLOYER_NAME", "OL_DEPLOYER_NAME")

	deployCreateCmd.Flags().String("environment", "", "environment name of the event (OPSLEVEL_ENVIRONMENT)")
	viper.BindEnv("environment", "OPSLEVEL_ENVIRONMENT", "OL_ENVIRONMENT")

	deployCreateCmd.Flags().String("git-path", "./", "relative path to grab the git commit info from (if git repo is found overrides all commit details)")
	viper.BindEnv("git-path", "OPSLEVEL_GIT_PATH", "OL_GIT_PATH")

	deployCreateCmd.Flags().StringP("description", "d", "", "description of the event (OPSLEVEL_DESCRIPTION)")
	viper.BindEnv("description", "OPSLEVEL_DESCRIPTION", "OL_DESCRIPTION")

	deployCreateCmd.Flags().StringP("service", "s", "", "service alias for the event (OPSLEVEL_SERVICE)")
	viper.BindEnv("service", "OPSLEVEL_SERVICE", "OL_SERVICE")

	deployCreateCmd.Flags().StringVarP(&integrationUrl, "integration-url", "i", "", "OpsLevel integration url (OPSLEVEL_INTEGRATION_URL)")
	viper.BindEnv("integration-url", "OPSLEVEL_INTEGRATION_URL", "OL_INTEGRATION_URL")

	viper.BindPFlags(deployCreateCmd.Flags())
}

func readCreateConfigAsDeployEvent() (*DeployEvent, error) {
	readInputConfig()
	evt := &DeployEvent{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	evt.DeployedAt = time.Now().UTC()

	fillWithOverrides(evt)
	fillGitInfo(evt)
	return evt, nil
}

func fillWithOverrides(evt *DeployEvent) {
	if service := viper.GetString("service"); service != "" {
		evt.Service = service
	}
	if description := viper.GetString("description"); description != "" {
		evt.Description = description
	}
	if environment := viper.GetString("environment"); environment != "" {
		evt.Environment = environment
	}
	if number := viper.GetString("deploy-number"); number != "" {
		evt.DeployNumber = number
	}
	if url := viper.GetString("deploy-url"); url != "" {
		evt.DeployURL = url
	}
	if id := viper.GetString("dedup-id"); id != "" {
		evt.DedupID = id
	}
	if name := viper.GetString("deployer-name"); name != "" {
		evt.Deployer.Name = name
	}
	if email := viper.GetString("deployer-email"); email != "" {
		evt.Deployer.Email = email
	}
	if sha := viper.GetString("commit-sha"); sha != "" {
		evt.Commit.SHA = sha
	}
	if message := viper.GetString("commit-message"); message != "" {
		evt.Commit.Message = message
	}
}

func fillGitInfo(evt *DeployEvent) {
	var err error
	r, err := git.PlainOpen(viper.GetString("git-path"))
	if err != nil {
		log.Debug().Msgf("Failed to open git repo: '%s'", viper.GetString("git-path"))
		return
	}
	ref, err := r.Head()
	if err != nil {
		log.Debug().Msg("Failed to get HEAD of git repo")
		return
	}
	hash := ref.Hash()
	commit, err := r.CommitObject(hash)
	if err != nil {
		log.Debug().Msg("Failed to read 'CommitObject' from hash of HEAD of git repo")
		return
	}
	evt.Commit = Commit{
		SHA:            hash.String(),
		Message:        commit.Message,
		Date:           &commit.Committer.When,
		CommitterName:  commit.Committer.Name,
		CommitterEmail: commit.Committer.Email,
		AuthorName:     commit.Author.Name,
		AuthorEmail:    commit.Author.Email,
		AuthoringDate:  &commit.Author.When,
	}
}
