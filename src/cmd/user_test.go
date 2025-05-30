package cmd_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/opslevel/cli/cmd"
	"github.com/opslevel/opslevel-go/v2025"
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
	// Create User
	userId, err := createUser(expectedUser)
	if err != nil {
		t.Fatal(err)
	}

	// Get User
	createdUser, err := getUser(userId)
	if err != nil {
		t.Fatal(err)
	}
	if createdUser.Name != expectedUser.Name ||
		createdUser.Email != expectedUser.Email ||
		string(createdUser.Role) != string(defaultUserRole) ||
		!strings.HasPrefix(createdUser.HTMLUrl, "https://app.opslevel.com/users/") {
		t.Errorf("Create 'user' failed, expected user '%+v' but got '%+v'", expectedUser, createdUser)
	}

	// Update User
	expectedUpdatedUser := opslevel.User{
		UserId: createdUser.UserId,
		Name:   createdUser.Name,
		Role:   opslevel.UserRoleTeamMember,
	}
	updatedUserId, err := updateUser(string(createdUser.Id), expectedUpdatedUser)
	if err != nil {
		t.Fatal(err)
	}
	if string(createdUser.Id) != updatedUserId {
		t.Errorf("Update 'user' failed, expected returned ID '%s' but got '%s'", string(createdUser.Id), updatedUserId)
	}

	// Delete User
	if err = deleteUser(string(createdUser.Id)); err != nil {
		t.Error(err)
	}
}

func createUser(expectedUser opslevel.User) (string, error) {
	cliArgs := []string{expectedUser.Email, expectedUser.Name}
	cmd.RootCmd.SetArgs(cliArgs)

	// Create User
	createOutput, err := execCmd(Create, "user", cliArgs...)
	if err != nil {
		return "", fmt.Errorf("Create 'user' failed, got error: %v", err)
	}
	return asString(createOutput), nil
}

func getUser(userId string) (*opslevel.User, error) {
	getOutput, err := execCmd(Get, "user", userId)
	if err != nil {
		return nil, fmt.Errorf("Get 'user' failed, got error: %v", err)
	}

	createdUser, err := jsonToResource[opslevel.User](getOutput)
	if err != nil {
		return nil, fmt.Errorf("Failed to convert JSON from API to 'opslevel.User' struct")
	}
	return createdUser, err
}

func updateUser(userId string, userToUpdate opslevel.User) (string, error) {
	if err := writeToYaml(userFileName, userToUpdate); err != nil {
		return "", fmt.Errorf("Error while writing '%v' to file '%s': %v", userToUpdate, userFileName, err)
	}

	// Store Update User stuff to "file"
	cliArgs := []string{userId, "-f", userFileName}
	updateOutput, err := execCmd(Update, "user", cliArgs...)
	if err != nil {
		return "", fmt.Errorf("Update 'user' failed, got error: %v", err)
	}
	return asString(updateOutput), nil
}

func deleteUser(userId string) error {
	_, err := execCmd(Delete, "user", userId)
	return err
}
