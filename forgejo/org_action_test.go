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
)

func TestCreateOrgActionSecret(t *testing.T) {
	log.Println("== TestCreateOrgActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_user", c)
	c.SetSudo(user.UserName)
	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "ActionOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	// create secret
	resp, err := c.CreateOrgActionSecret(newOrg.UserName, CreateSecretOption{Name: "test", Data: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// update secret
	resp, err = c.CreateOrgActionSecret(newOrg.UserName, CreateSecretOption{Name: "test", Data: "test2"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// list secrets
	secrets, _, err := c.ListOrgActionSecret(newOrg.UserName, ListOrgActionSecretOption{})
	require.NoError(t, err)
	assert.Len(t, secrets, 1)
}

func TestUpdateOrgActionSecret(t *testing.T) {
	log.Println("== TestUpdateOrgActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_update_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete existing org from previous test runs
	c.DeleteOrg("ActionUpdateOrg")

	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "ActionUpdateOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	// Pre-cleanup: delete any existing secret from previous test runs
	c.DeleteOrgActionSecret(newOrg.UserName, "test_update")

	// create secret first
	resp, err := c.CreateOrgActionSecret(newOrg.UserName, CreateSecretOption{Name: "test_update", Data: "initial_value"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// update secret using UpdateOrgActionSecret
	resp, err = c.UpdateOrgActionSecret(newOrg.UserName, "test_update", UpdateOrgActionSecretOption{
		Name:  "test_update",
		Value: "updated_value",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify update by listing secrets
	secrets, _, err := c.ListOrgActionSecret(newOrg.UserName, ListOrgActionSecretOption{})
	require.NoError(t, err)
	assert.Len(t, secrets, 1)

	// Cleanup
	c.DeleteOrgActionSecret(newOrg.UserName, "test_update")
}

func TestDeleteOrgActionSecret(t *testing.T) {
	log.Println("== TestDeleteOrgActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_delete_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete existing org from previous test runs
	c.DeleteOrg("ActionDeleteOrg")

	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "ActionDeleteOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	// create secret first
	resp, err := c.CreateOrgActionSecret(newOrg.UserName, CreateSecretOption{Name: "test_delete", Data: "delete_me"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// delete secret
	resp, err = c.DeleteOrgActionSecret(newOrg.UserName, "test_delete")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify deletion - list should be empty
	secrets, _, err := c.ListOrgActionSecret(newOrg.UserName, ListOrgActionSecretOption{})
	require.NoError(t, err)
	assert.Len(t, secrets, 0)
}
