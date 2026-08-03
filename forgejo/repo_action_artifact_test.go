// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoActionArtifacts(t *testing.T) {
	log.Println("== TestRepoActionArtifacts ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_artifact_user", c)
	c.SetSudo(user.UserName)

	// pre-cleanup then create a fresh repo
	_, _ = c.DeleteRepo(user.UserName, "artifacts-test")
	repo, _, err := c.CreateRepo(CreateRepoOption{Name: "artifacts-test"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteRepo(user.UserName, repo.Name) }()

	owner, name := user.UserName, repo.Name

	// listing an empty repo's artifacts must succeed and be empty
	arts, _, err := c.ListActionArtifacts(owner, name, ListActionArtifactsOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.Empty(t, arts)

	// the name filter is accepted by the list endpoint
	_, _, err = c.ListActionArtifacts(owner, name, ListActionArtifactsOption{Name: "no-such-artifact", ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)

	// get / download / delete a non-existent artifact -> error
	_, _, err = c.GetActionArtifact(owner, name, 999999)
	require.Error(t, err)

	_, err = c.DeleteActionArtifact(owner, name, 999999)
	require.Error(t, err)

	r, _, err := c.DownloadActionArtifact(owner, name, 999999)
	if r != nil {
		r.Close()
	}
	assert.Error(t, err)
}
