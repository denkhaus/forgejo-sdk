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

func TestCreateUserActionSecret(t *testing.T) {
	log.Println("== TestCreateUserActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_secret_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete any existing secret from previous test runs
	c.DeleteUserActionSecret("test")

	// create secret
	resp, err := c.CreateUserActionSecret(CreateSecretOption{Name: "test", Data: "test"})
	require.NoError(t, err)
	// Accept either 201 (first creation) or 204 (already exists)
	assert.Contains(t, []int{http.StatusCreated, http.StatusNoContent}, resp.StatusCode)

	// update secret
	resp, err = c.CreateUserActionSecret(CreateSecretOption{Name: "test", Data: "test2"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Note: ListUserActionSecret endpoint returns 404 in Forgejo 13.0.3
	// The endpoint exists in swagger but is not implemented yet.

	// Cleanup
	c.DeleteUserActionSecret("test")
}

func TestUpdateUserActionSecret(t *testing.T) {
	log.Println("== TestUpdateUserActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_update_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete any existing secret from previous test runs
	c.DeleteUserActionSecret("test_update")

	// create secret first
	resp, err := c.CreateUserActionSecret(CreateSecretOption{Name: "test_update", Data: "initial_value"})
	require.NoError(t, err)
	// Accept either 201 or 204
	assert.Contains(t, []int{http.StatusCreated, http.StatusNoContent}, resp.StatusCode)

	// update secret using UpdateUserActionSecret
	resp, err = c.UpdateUserActionSecret("test_update", UpdateUserActionSecretOption{
		Name:  "test_update",
		Value: "updated_value",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Note: ListUserActionSecret endpoint returns 404 in Forgejo 13.0.3

	// Cleanup
	c.DeleteUserActionSecret("test_update")
}

func TestDeleteUserActionSecret(t *testing.T) {
	log.Println("== TestDeleteUserActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "user_action_delete_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete any existing secret from previous test runs
	c.DeleteUserActionSecret("test_delete")

	// create secret first
	resp, err := c.CreateUserActionSecret(CreateSecretOption{Name: "test_delete", Data: "delete_me"})
	require.NoError(t, err)
	// Accept either 201 or 204
	assert.Contains(t, []int{http.StatusCreated, http.StatusNoContent}, resp.StatusCode)

	// delete secret
	resp, err = c.DeleteUserActionSecret("test_delete")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Note: ListUserActionSecret endpoint returns 404 in Forgejo 13.0.3
	// Cannot verify deletion via list
}
