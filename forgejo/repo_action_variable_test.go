// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2024 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

func TestCreateRepoActionVariable(t *testing.T) {
	log.Println("== TestCreateRepoActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_variable_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-repo-variable",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	testValue := "test_value"

	// create variable
	resp, err := c.CreateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "test", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// list variables
	variables, _, err := c.ListRepoActionVariables(newRepo.Owner.UserName, newRepo.Name, ListRepoActionVariablesOption{})
	require.NoError(t, err)
	assert.Len(t, variables, 1)
	// Variable names are uppercased by Forgejo
	assert.Equal(t, "TEST", variables[0].Name)
	// Data field is NOT returned in list operations (empty for security)
	assert.Empty(t, variables[0].Data)

	// get variable
	variable, _, err := c.GetRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST")
	require.NoError(t, err)
	assert.NotNil(t, variable)
	assert.Equal(t, "TEST", variable.Name)
	assert.Equal(t, testValue, variable.Data)

	// update variable
	updatedValue := "updated_value"
	resp, err = c.UpdateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST", models.UpdateVariableOption{
		Value: &updatedValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify update
	variable, _, err = c.GetRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST")
	require.NoError(t, err)
	assert.Equal(t, updatedValue, variable.Data)

	// delete variable
	resp, err = c.DeleteRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify deletion
	variables, _, err = c.ListRepoActionVariables(newRepo.Owner.UserName, newRepo.Name, ListRepoActionVariablesOption{})
	require.NoError(t, err)
	assert.Len(t, variables, 0)
}

func TestListRepoActionVariables(t *testing.T) {
	log.Println("== TestListRepoActionVariables ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_variable_list_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-repo-list",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// create multiple variables
	value1 := "value1"
	value2 := "value2"
	value3 := "value3"

	_, err = c.CreateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "var1", models.CreateVariableOption{
		Value: &value1,
	})
	require.NoError(t, err)

	_, err = c.CreateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "var2", models.CreateVariableOption{
		Value: &value2,
	})
	require.NoError(t, err)

	_, err = c.CreateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "var3", models.CreateVariableOption{
		Value: &value3,
	})
	require.NoError(t, err)

	// list variables
	variables, _, err := c.ListRepoActionVariables(newRepo.Owner.UserName, newRepo.Name, ListRepoActionVariablesOption{})
	require.NoError(t, err)
	assert.Len(t, variables, 3)

	// verify variable names (names are uppercased by Forgejo)
	variableNames := make(map[string]bool)
	for _, v := range variables {
		variableNames[v.Name] = true
	}
	assert.True(t, variableNames["VAR1"])
	assert.True(t, variableNames["VAR2"])
	assert.True(t, variableNames["VAR3"])
}

func TestGetRepoActionVariable(t *testing.T) {
	log.Println("== TestGetRepoActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_variable_get_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-repo-get",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	testValue := "get_test_value"

	// create variable
	_, err = c.CreateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "test_get", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)

	// get variable
	variable, _, err := c.GetRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST_GET")
	require.NoError(t, err)
	assert.NotNil(t, variable)
	assert.Equal(t, "TEST_GET", variable.Name)
	assert.Equal(t, testValue, variable.Data)
}

func TestUpdateRepoActionVariable(t *testing.T) {
	log.Println("== TestUpdateRepoActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_variable_update_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-repo-update",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	originalValue := "original_value"
	newValue := "new_value"

	// create variable
	_, err = c.CreateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "test_update", models.CreateVariableOption{
		Value: &originalValue,
	})
	require.NoError(t, err)

	// update variable with new value
	resp, err := c.UpdateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST_UPDATE", models.UpdateVariableOption{
		Value: &newValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify update
	variable, _, err := c.GetRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST_UPDATE")
	require.NoError(t, err)
	assert.Equal(t, newValue, variable.Data)
}

func TestDeleteRepoActionVariable(t *testing.T) {
	log.Println("== TestDeleteRepoActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_variable_delete_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-repo-delete",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	testValue := "delete_test_value"

	// create variable
	_, err = c.CreateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "test_delete", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)

	// delete variable
	resp, err := c.DeleteRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST_DELETE")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify deletion - get should fail
	_, resp, err = c.GetRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "TEST_DELETE")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestRepoActionVariableErrorHandling(t *testing.T) {
	log.Println("== TestRepoActionVariableErrorHandling ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_variable_error_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-repo-error",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// test getting non-existent variable
	_, resp, err := c.GetRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "non_existent")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// test updating non-existent variable
	testValue := "test"
	resp, err = c.UpdateRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "non_existent", models.UpdateVariableOption{
		Value: &testValue,
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// test deleting non-existent variable
	resp, err = c.DeleteRepoActionVariable(newRepo.Owner.UserName, newRepo.Name, "non_existent")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
