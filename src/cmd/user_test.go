package cmd_test

import (
	"strings"
	"testing"

	"github.com/opslevel/cli/cmd"
	"github.com/opslevel/opslevel-go/v2024"
)

const (
	defaultUserRole = opslevel.UserRoleUser
	userFileName    = "test_user.yaml"
	userName        = "CLI Test User"
)

func Test_UserCRUD(t *testing.T) {
	expectedUser := opslevel.User{
		UserId: opslevel.UserId{Email: "testcli+pat@opslevel.com"},
		Name:   userName,
	}
	cliArgs := []string{expectedUser.Email, expectedUser.Name}
	cmd.RootCmd.SetArgs(cliArgs)

	// Create User
	createOutput, err := execCmd(Create, "user", cliArgs...)
	if err != nil {
		t.Errorf("Create 'user' failed, got error: %v", err)
	}
	userId := asString(createOutput)

	// Get User
	getOutput, err := execCmd(Get, "user", userId)
	if err != nil {
		t.Errorf("Get 'user' failed, got error: %v", err)
	}

	createdUser := jsonToResource[opslevel.User](getOutput)
	if createdUser.Name != expectedUser.Name ||
		createdUser.Email != expectedUser.Email ||
		string(createdUser.Role) != string(defaultUserRole) ||
		!strings.HasPrefix(createdUser.HTMLUrl, "https://app.opslevel.com/users/") {
		t.Errorf("Create 'user' failed, expected user '%+v' but got '%+v'", expectedUser, createdUser)
	}

	expectedUpdatedUser := opslevel.User{
		UserId: createdUser.UserId,
		Name:   createdUser.Name,
		Role:   opslevel.UserRoleTeamMember,
	}
	if err := writeToYaml(userFileName, expectedUpdatedUser); err != nil {
		t.Errorf("Error while writing '%v' to file '%s': %v", expectedUpdatedUser, userFileName, err)
	}

	// Store Update User stuff to "file"
	cliArgs = []string{expectedUser.Email, "-f", userFileName}
	updateOutput, err := execCmd(Update, "user", cliArgs...)
	if err != nil {
		t.Errorf("Update 'user' failed, got error: %v", err)
	}
	updatedUserId := asString(updateOutput)
	if string(createdUser.Id) != updatedUserId {
		t.Errorf("Update 'user' failed, expected returned ID '%s' but got '%s'", string(createdUser.Id), updatedUserId)
	}

	// Delete User
	if _, err = execCmd(Delete, "user", string(createdUser.Id)); err != nil {
		t.Error(err)
	}
}
