package cmd

import (
	"fmt"
	"time"

	"github.com/opslevel/cli/client"

	"github.com/creasty/defaults"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var integrationId string

type Deployer struct {
	Email string `validate:"required" json:"email" default:"automation@opslevel.com"`
	Name  string `json:"name,omitempty"`
}

// Commit represents the commit being deployed
type Commit struct {
	SHA            string    `json:"sha,omitempty"`
	Message        string    `json:"message,omitempty"`
	Branch         string    `json:"branch,omitempty"`
	Date           time.Time `json:"date,omitempty"`
	CommitterName  string    `json:"committer_name,omitempty"`
	CommitterEmail string    `json:"committer_email,omitempty"`
	AuthorName     string    `json:"author_name,omitempty"`
	AuthorEmail    string    `json:"author_email,omitempty"`
	AuthoringDate  time.Time `json:"authoring_date,omitempty"`
}

// DeployRequest represents a structured request to the OpsLevel deploys webhook endpoint
type DeployEvent struct {
	Service      string    `validate:"required" json:"service"`
	Deployer     Deployer  `validate:"required" json:"deployer"`
	DeployedAt   time.Time `validate:"required" json:"deployed_at"`
	Description  string    `validate:"required" json:"description" default:"Event Created by OpsLevel CLI"`
	Environment  string    `json:"environment,omitempty"`
	DeployURL    string    `json:"deploy_url,omitempty"`
	DeployNumber string    `json:"deploy_number,omitempty"`
	Commit       Commit    `json:"commit,omitempty"`
	DedupID      string    `json:"dedup_id,omitempty"`
}

var deployCreateCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Create deployment events",
	Long:  "Create deployment events",
	Run: func(cmd *cobra.Command, args []string) {
		evt, err := readCreateConfigAsDeployEvent()
		cobra.CheckErr(err)
		log.Debug().Msgf("%#v", evt)
		var resp struct {
			Result string `json:"result"`
		}
		c := client.NewClient()
		c.Do("POST", fmt.Sprintf("/integrations/deploy/%s", integrationId), evt, &resp)
		// TODO: Check Response
	},
}

func init() {
	createCmd.AddCommand(deployCreateCmd)

	deployCreateCmd.Flags().StringVarP(&integrationId, "integration", "i", "", "The OpsLevel integration id")

	deployCreateCmd.Flags().StringP("service", "s", "", "The service alias for the event")
	deployCreateCmd.Flags().StringP("environment", "e", "", "The environment of the event")
	deployCreateCmd.Flags().StringP("number", "n", "", "The deploy number of the event")
	deployCreateCmd.Flags().String("url", "", "The deploy url of the event")
	deployCreateCmd.Flags().String("id", "", "The dedup id of the event")
	viper.BindPFlags(deployCreateCmd.Flags())
}

func readCreateConfigAsDeployEvent() (*DeployEvent, error) {
	readCreateConfigFile()
	evt := &DeployEvent{}
	viper.Unmarshal(&evt)
	if err := defaults.Set(evt); err != nil {
		return nil, err
	}
	evt.DeployedAt = time.Now()

	if service := viper.GetString("service"); service != "" {
		evt.Service = service
	}
	if environment := viper.GetString("environment"); environment != "" {
		evt.Environment = environment
	}
	if number := viper.GetString("number"); number != "" {
		evt.DeployNumber = number
	}
	if url := viper.GetString("url"); url != "" {
		evt.DeployURL = url
	}
	if id := viper.GetString("id"); id != "" {
		evt.DedupID = id
	}
	// TODO: fill in git commiter info automatically
	return evt, nil
}
