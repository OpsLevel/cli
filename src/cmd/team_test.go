package cmd

import (
	"testing"

	"github.com/opslevel/opslevel-go/v2024"
)

const (
	teamFileName = "test_team.yaml"
	teamName     = "CLI Test Team"
)

func Test_TeamCRUD(t *testing.T) {
	teamToCreate := opslevel.TeamCreateInput{
		Name:             teamName,
		Responsibilities: opslevel.RefOf("all the things"),
	}
	if err := writeToYaml(teamFileName, teamToCreate); err != nil {
		t.Errorf("Error while writing '%v' to file '%s': %v", teamToCreate, teamFileName, err)
	}

	cliArgs := []string{teamToCreate.Name, "-f", teamFileName}
	rootCmd.SetArgs(cliArgs)

	// Create Team
	createOutput, err := execCmd(Create, "team", cliArgs...)
	if err != nil {
		t.Errorf("Create 'team' failed, got error: %v", err)
	}
	teamId := asString(createOutput)

	// Get Team
	getOutput, err := execCmd(Get, "team", teamId)
	if err != nil {
		t.Errorf("Get 'team' failed, got error: %v", err)
	}

	createdTeam := jsonToResource[opslevel.Team](getOutput)
	if createdTeam.Name != teamToCreate.Name ||
		createdTeam.Responsibilities != *teamToCreate.Responsibilities {
		t.Errorf("Create 'team' failed, expected team '%+v' but got '%+v'", teamToCreate, createdTeam)
	}

	teamToUpdate := opslevel.TeamCreateInput{
		Name:             createdTeam.Name,
		Responsibilities: opslevel.RefOf("new things"),
	}
	if err := writeToYaml(teamFileName, teamToUpdate); err != nil {
		t.Errorf("Error while writing '%v' to file '%s': %v", teamToUpdate, teamFileName, err)
	}

	// Store Update Team stuff to "file"
	cliArgs = []string{teamId, "-f", teamFileName}
	updateOutput, err := execCmd(Update, "team", cliArgs...)
	if err != nil {
		t.Errorf("Update 'team' failed, got error: %v", err)
	}
	updatedTeamId := asString(updateOutput)
	if string(createdTeam.Id) != updatedTeamId {
		t.Errorf("Update 'team' failed, expected returned ID '%s' but got '%s'", string(createdTeam.Id), updatedTeamId)
	}

	// Delete Team
	if _, err = execCmd(Delete, "team", string(createdTeam.Id)); err != nil {
		t.Error(err)
	}
}
