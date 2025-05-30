package cmd_test

import (
	"fmt"
	"testing"

	"github.com/opslevel/cli/cmd"
	"github.com/opslevel/opslevel-go/v2025"
)

const (
	teamFileName = "test_team.yaml"
	teamName     = "CLI Test Team"
)

func Test_TeamCRUD(t *testing.T) {
	// Create Team
	teamToCreate := opslevel.TeamCreateInput{
		Name:             teamName,
		Responsibilities: opslevel.RefOf("all the things"),
	}
	teamId, err := createTeam(teamToCreate)
	if err != nil {
		t.Fatal(err)
	}

	// Get Team
	createdTeam, err := getTeam(teamId)
	if err != nil {
		t.Fatal(err)
	}
	if createdTeam.Name != teamToCreate.Name ||
		createdTeam.Responsibilities != *teamToCreate.Responsibilities {
		t.Errorf("Create 'team' failed, expected team '%+v' but got '%+v'", teamToCreate, createdTeam)
	}

	// Update Team
	teamToUpdate := opslevel.TeamUpdateInput{
		Name:             opslevel.RefOf(createdTeam.Name),
		Responsibilities: opslevel.RefOf("new things"),
	}
	updatedTeamId, err := updateTeam(teamId, teamToUpdate)
	if err != nil {
		_ = deleteTeam(string(createdTeam.Id))
		t.Fatal(err)
	}
	if string(createdTeam.Id) != updatedTeamId {
		t.Errorf("Update 'team' failed, expected returned ID '%s' but got '%s'", string(createdTeam.Id), updatedTeamId)
	}

	// Delete Team
	if err = deleteTeam(string(createdTeam.Id)); err != nil {
		t.Errorf("Delete 'team' failed, got error '%s'", err)
	}
}

func createTeam(teamToCreate opslevel.TeamCreateInput) (string, error) {
	if err := writeToYaml(teamFileName, teamToCreate); err != nil {
		return "", fmt.Errorf("Error while writing '%v' to file '%s': %v", teamToCreate, teamFileName, err)
	}

	cliArgs := []string{teamToCreate.Name, "-f", teamFileName}
	cmd.RootCmd.SetArgs(cliArgs)

	// Create Team
	createOutput, err := execCmd(Create, "team", cliArgs...)
	if err != nil {
		return "", fmt.Errorf("Create 'team' failed, got error: %v", err)
	}
	return asString(createOutput), nil
}

func getTeam(teamId string) (*opslevel.Team, error) {
	getOutput, err := execCmd(Get, "team", teamId)
	if err != nil {
		return nil, fmt.Errorf("Get 'team' failed, got error: %v", err)
	}

	createdTeam, err := jsonToResource[opslevel.Team](getOutput)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert JSON from API to 'opslevel.Team' struct")
	}
	return createdTeam, err
}

func updateTeam(teamId string, teamToUpdate opslevel.TeamUpdateInput) (string, error) {
	if err := writeToYaml(teamFileName, teamToUpdate); err != nil {
		return "", fmt.Errorf("Error while writing '%v' to file '%s': %v", teamToUpdate, teamFileName, err)
	}

	// Store Update Team stuff to "file"
	cliArgs := []string{teamId, "-f", teamFileName}
	updateOutput, err := execCmd(Update, "team", cliArgs...)
	if err != nil {
		return "", fmt.Errorf("Update 'team' failed, got error: %v", err)
	}
	return asString(updateOutput), nil
}

func deleteTeam(teamId string) error {
	_, err := execCmd(Delete, "team", teamId)
	return err
}
