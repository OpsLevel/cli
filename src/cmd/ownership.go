package cmd

import (
	"fmt"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/cobra"
	git "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
	"time"
)

var ownershipCreateCmd = &cobra.Command{
	Use:   "ownership",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Msgf("Lets build Ownership!")

		emails, err := getEmailsFromRepo("/Users/rocktavious/Projects/opslevel-website", 5)
		if err != nil {
			return err
		}

		teamLookup := map[string]string{}
		teams, err := getClientGQL().ListTeams()
		for _, team := range teams {
			for _, member := range team.Members.Nodes {
				teamLookup[member.Email] = team.Alias
			}
		}

		probability := map[string]int{}
		for _, email := range emails {
			if team, ok := teamLookup[email]; ok {
				probability[team] += 1
			}
		}

		for team, value := range probability {
			fmt.Printf("%s = %d\n", team, value)
		}

		return nil
	},
}

func init() {
	createCmd.AddCommand(ownershipCreateCmd)
}

func getEmailsFromRepo(repositoryPath string, monthsAgo int) ([]string, error) {
	var output []string

	r, err := git.PlainOpen(repositoryPath)
	if err != nil {
		return output, err
	}

	ref, err := r.Head()
	if err != nil {
		return output, err
	}

	until := time.Now()
	since := until.AddDate(0, -monthsAgo , 0)
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash(), Since: &since, Until: &until})

	emails := map[string]bool{}
	err = cIter.ForEach(func(c *object.Commit) error {
		if _, ok := emails[c.Author.Email]; !ok {
			emails[c.Author.Email] = true
		}
		if _, ok := emails[c.Committer.Email]; !ok {
			emails[c.Committer.Email] = true
		}
		return nil
	})

	for email, _ := range emails {
		output = append(output, email)
	}
	return output, nil
}