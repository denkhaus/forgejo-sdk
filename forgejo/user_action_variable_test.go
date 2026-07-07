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

func TestCreateUserActionVariable(t *testing.T) {
	log.Println("== TestCreateUserActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_variable_user", c)
	c.SetSudo(user.UserName)

	testValue := "test_value"

	// create variable
	resp, err := c.CreateUserActionVariable("test", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// list variables
	variables, _, err := c.ListUserActionVariables(ListUserActionVariablesOption{})
	require.NoError(t, err)
	assert.Len(t, variables, 1)
	// Variable names are uppercased by Forgejo
	assert.Equal(t, "TEST", variables[0].Name)
	// User variables DO return data in list operations (unlike repo/org variables)
	assert.Equal(t, testValue, variables[0].Data)

	// get variable
	variable, _, err := c.GetUserActionVariable("TEST")
	require.NoError(t, err)
	assert.NotNil(t, variable)
	assert.Equal(t, "TEST", variable.Name)
	assert.Equal(t, testValue, variable.Data)

	// update variable
	updatedValue := "updated_value"
	resp, err = c.UpdateUserActionVariable("TEST", models.UpdateVariableOption{
		Value: &updatedValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify update
	variable, _, err = c.GetUserActionVariable("TEST")
	require.NoError(t, err)
	assert.Equal(t, updatedValue, variable.Data)

	// delete variable
	resp, err = c.DeleteUserActionVariable("TEST")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify deletion
	variables, _, err = c.ListUserActionVariables(ListUserActionVariablesOption{})
	require.NoError(t, err)
	assert.Empty(t, variables)
}

func TestListUserActionVariables(t *testing.T) {
	log.Println("== TestListUserActionVariables ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_variable_list_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete any existing variables from previous test runs
	_, _ = c.DeleteUserActionVariable("VAR1")
	_, _ = c.DeleteUserActionVariable("VAR2")
	_, _ = c.DeleteUserActionVariable("VAR3")

	// Cleanup: delete variables at the end of the test
	t.Cleanup(func() {
		_, _ = c.DeleteUserActionVariable("VAR1")
		_, _ = c.DeleteUserActionVariable("VAR2")
		_, _ = c.DeleteUserActionVariable("VAR3")
	})

	// create multiple variables
	value1 := "value1"
	value2 := "value2"
	value3 := "value3"

	_, err := c.CreateUserActionVariable("var1", models.CreateVariableOption{
		Value: &value1,
	})
	require.NoError(t, err)

	_, err = c.CreateUserActionVariable("var2", models.CreateVariableOption{
		Value: &value2,
	})
	require.NoError(t, err)

	_, err = c.CreateUserActionVariable("var3", models.CreateVariableOption{
		Value: &value3,
	})
	require.NoError(t, err)

	// list variables
	variables, _, err := c.ListUserActionVariables(ListUserActionVariablesOption{})
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

func TestGetUserActionVariable(t *testing.T) {
	log.Println("== TestGetUserActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_variable_get_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete any existing variable from previous test runs
	_, _ = c.DeleteUserActionVariable("TEST_GET")

	// Cleanup: delete variable at the end of the test
	t.Cleanup(func() {
		_, _ = c.DeleteUserActionVariable("TEST_GET")
	})

	testValue := "get_test_value"

	// create variable
	_, err := c.CreateUserActionVariable("test_get", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)

	// get variable
	variable, _, err := c.GetUserActionVariable("TEST_GET")
	require.NoError(t, err)
	assert.NotNil(t, variable)
	assert.Equal(t, "TEST_GET", variable.Name)
	assert.Equal(t, testValue, variable.Data)
}

func TestUpdateUserActionVariable(t *testing.T) {
	log.Println("== TestUpdateUserActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_variable_update_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete any existing variable from previous test runs
	_, _ = c.DeleteUserActionVariable("TEST_UPDATE")

	// Cleanup: delete variable at the end of the test
	t.Cleanup(func() {
		_, _ = c.DeleteUserActionVariable("TEST_UPDATE")
	})

	originalValue := "original_value"
	newValue := "new_value"

	// create variable
	_, err := c.CreateUserActionVariable("test_update", models.CreateVariableOption{
		Value: &originalValue,
	})
	require.NoError(t, err)

	// update variable with new value
	resp, err := c.UpdateUserActionVariable("TEST_UPDATE", models.UpdateVariableOption{
		Value: &newValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify update
	variable, _, err := c.GetUserActionVariable("TEST_UPDATE")
	require.NoError(t, err)
	assert.Equal(t, newValue, variable.Data)
}

func TestDeleteUserActionVariable(t *testing.T) {
	log.Println("== TestDeleteUserActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_variable_delete_user", c)
	c.SetSudo(user.UserName)

	testValue := "delete_test_value"

	// create variable
	_, err := c.CreateUserActionVariable("test_delete", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)

	// delete variable
	resp, err := c.DeleteUserActionVariable("TEST_DELETE")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify deletion - get should fail
	_, resp, err = c.GetUserActionVariable("TEST_DELETE")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUserActionVariableErrorHandling(t *testing.T) {
	log.Println("== TestUserActionVariableErrorHandling ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_variable_error_user", c)
	c.SetSudo(user.UserName)

	// test getting non-existent variable
	_, resp, err := c.GetUserActionVariable("non_existent")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// test updating non-existent variable
	testValue := "test"
	resp, err = c.UpdateUserActionVariable("non_existent", models.UpdateVariableOption{
		Value: &testValue,
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// test deleting non-existent variable
	resp, err = c.DeleteUserActionVariable("non_existent")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
