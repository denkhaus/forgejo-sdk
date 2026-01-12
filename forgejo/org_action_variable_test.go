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

func TestCreateOrgActionVariable(t *testing.T) {
	log.Println("== TestCreateOrgActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_variable_user", c)
	c.SetSudo(user.UserName)
	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "TestVariableOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	testValue := "test_value"

	// create variable
	resp, err := c.CreateOrgActionVariable(newOrg.UserName, "test", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// list variables
	variables, _, err := c.ListOrgActionVariables(newOrg.UserName, ListOrgActionVariablesOption{})
	require.NoError(t, err)
	assert.Len(t, variables, 1)
	// Variable names are uppercased by Forgejo
	assert.Equal(t, "TEST", variables[0].Name)
	// Data field IS returned in org list operations
	assert.Equal(t, testValue, variables[0].Data)

	// get variable
	variable, _, err := c.GetOrgActionVariable(newOrg.UserName, "TEST")
	require.NoError(t, err)
	assert.NotNil(t, variable)
	assert.Equal(t, "TEST", variable.Name)
	assert.Equal(t, testValue, variable.Data)

	// update variable
	updatedValue := "updated_value"
	resp, err = c.UpdateOrgActionVariable(newOrg.UserName, "TEST", models.UpdateVariableOption{
		Value: &updatedValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify update
	variable, _, err = c.GetOrgActionVariable(newOrg.UserName, "TEST")
	require.NoError(t, err)
	assert.Equal(t, updatedValue, variable.Data)

	// delete variable
	resp, err = c.DeleteOrgActionVariable(newOrg.UserName, "TEST")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify deletion
	variables, _, err = c.ListOrgActionVariables(newOrg.UserName, ListOrgActionVariablesOption{})
	require.NoError(t, err)
	assert.Len(t, variables, 0)
}

func TestListOrgActionVariables(t *testing.T) {
	log.Println("== TestListOrgActionVariables ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_variable_list_user", c)
	c.SetSudo(user.UserName)
	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "TestVariableListOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	// create multiple variables
	value1 := "value1"
	value2 := "value2"
	value3 := "value3"

	_, err = c.CreateOrgActionVariable(newOrg.UserName, "var1", models.CreateVariableOption{
		Value: &value1,
	})
	require.NoError(t, err)

	_, err = c.CreateOrgActionVariable(newOrg.UserName, "var2", models.CreateVariableOption{
		Value: &value2,
	})
	require.NoError(t, err)

	_, err = c.CreateOrgActionVariable(newOrg.UserName, "var3", models.CreateVariableOption{
		Value: &value3,
	})
	require.NoError(t, err)

	// list variables
	variables, _, err := c.ListOrgActionVariables(newOrg.UserName, ListOrgActionVariablesOption{})
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

func TestGetOrgActionVariable(t *testing.T) {
	log.Println("== TestGetOrgActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_variable_get_user", c)
	c.SetSudo(user.UserName)
	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "TestVariableGetOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	testValue := "get_test_value"

	// create variable
	_, err = c.CreateOrgActionVariable(newOrg.UserName, "test_get", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)

	// get variable
	variable, _, err := c.GetOrgActionVariable(newOrg.UserName, "TEST_GET")
	require.NoError(t, err)
	assert.NotNil(t, variable)
	assert.Equal(t, "TEST_GET", variable.Name)
	assert.Equal(t, testValue, variable.Data)
}

func TestUpdateOrgActionVariable(t *testing.T) {
	log.Println("== TestUpdateOrgActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_variable_update_user", c)
	c.SetSudo(user.UserName)
	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "TestVariableUpdateOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	originalValue := "original_value"
	newValue := "new_value"

	// create variable
	_, err = c.CreateOrgActionVariable(newOrg.UserName, "test_update", models.CreateVariableOption{
		Value: &originalValue,
	})
	require.NoError(t, err)

	// update variable with new value
	resp, err := c.UpdateOrgActionVariable(newOrg.UserName, "TEST_UPDATE", models.UpdateVariableOption{
		Value: &newValue,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify update
	variable, _, err := c.GetOrgActionVariable(newOrg.UserName, "TEST_UPDATE")
	require.NoError(t, err)
	assert.Equal(t, newValue, variable.Data)
}

func TestDeleteOrgActionVariable(t *testing.T) {
	log.Println("== TestDeleteOrgActionVariable ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_variable_delete_user", c)
	c.SetSudo(user.UserName)
	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "TestVariableDeleteOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	testValue := "delete_test_value"

	// create variable
	_, err = c.CreateOrgActionVariable(newOrg.UserName, "test_delete", models.CreateVariableOption{
		Value: &testValue,
	})
	require.NoError(t, err)

	// delete variable
	resp, err := c.DeleteOrgActionVariable(newOrg.UserName, "TEST_DELETE")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify deletion - get should fail
	_, resp, err = c.GetOrgActionVariable(newOrg.UserName, "TEST_DELETE")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestOrgActionVariableErrorHandling(t *testing.T) {
	log.Println("== TestOrgActionVariableErrorHandling ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_variable_error_user", c)
	c.SetSudo(user.UserName)
	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "TestVariableErrorOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	// test getting non-existent variable
	_, resp, err := c.GetOrgActionVariable(newOrg.UserName, "non_existent")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// test updating non-existent variable
	testValue := "test"
	resp, err = c.UpdateOrgActionVariable(newOrg.UserName, "non_existent", models.UpdateVariableOption{
		Value: &testValue,
	})
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// test deleting non-existent variable
	resp, err = c.DeleteOrgActionVariable(newOrg.UserName, "non_existent")
	require.Error(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
