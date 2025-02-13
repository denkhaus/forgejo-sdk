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
