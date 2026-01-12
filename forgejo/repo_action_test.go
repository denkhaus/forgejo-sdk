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

func TestCreateRepoActionSecret(t *testing.T) {
	log.Println("== TestCreateRepoActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// create secret
	resp, err := c.CreateRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, CreateSecretOption{Name: "test", Data: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// update secret
	resp, err = c.CreateRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, CreateSecretOption{Name: "test", Data: "test2"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// list secrets
	secrets, _, err := c.ListRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, ListRepoActionSecretOption{})
	require.NoError(t, err)
	assert.Len(t, secrets, 1)
}

func TestUpdateRepoActionSecret(t *testing.T) {
	log.Println("== TestUpdateRepoActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_update_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete existing repo from previous test runs
	c.DeleteRepo(user.UserName, "test-update")

	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-update",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// Pre-cleanup: delete any existing secret from previous test runs
	c.DeleteRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, "test_update")

	// create secret first
	resp, err := c.CreateRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, CreateSecretOption{Name: "test_update", Data: "initial_value"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// update secret using UpdateRepoActionSecret
	resp, err = c.UpdateRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, "test_update", UpdateRepoActionSecretOption{
		Name:  "test_update",
		Value: "updated_value",
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify update by listing secrets
	secrets, _, err := c.ListRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, ListRepoActionSecretOption{})
	require.NoError(t, err)
	assert.Len(t, secrets, 1)

	// Cleanup
	c.DeleteRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, "test_update")
}

func TestDeleteRepoActionSecret(t *testing.T) {
	log.Println("== TestDeleteRepoActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_delete_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete existing repo from previous test runs
	c.DeleteRepo(user.UserName, "test-delete")

	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-delete",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// create secret first
	resp, err := c.CreateRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, CreateSecretOption{Name: "test_delete", Data: "delete_me"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// delete secret
	resp, err = c.DeleteRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, "test_delete")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// verify deletion - list should be empty
	secrets, _, err := c.ListRepoActionSecret(newRepo.Owner.UserName, newRepo.Name, ListRepoActionSecretOption{})
	require.NoError(t, err)
	assert.Len(t, secrets, 0)
}
