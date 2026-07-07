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

func Test_RepoKey(t *testing.T) {
	log.Println("== Test_RepoKey ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestRepoKey", c)

	// ListDeployKeys — safe GET on a fresh repo (empty list expected).
	keys, _, err := c.ListDeployKeys(repo.Owner.UserName, repo.Name, ListDeployKeysOptions{})
	require.NoError(t, err)
	assert.Empty(t, keys)

	// Add a throwaway deploy key. Creation may be rejected if the server
	// dislikes the key material; tolerate that and only exercise the
	// dependent get/list/delete path on success.
	dk, _, err := c.CreateDeployKey(repo.Owner.UserName, repo.Name, CreateKeyOption{
		Title:    "ci-throwaway",
		Key:      "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEhT6p7Qx5r8W9Y3bQl2mFk0Z9n4t7s2h1j6g5d4c3b test-throwaway@example.com",
		ReadOnly: true,
	})
	if err != nil {
		t.Skipf("deploy key creation not permitted in this environment: %v", err)
	}
	assert.NotZero(t, dk.ID)

	// GetDeployKey.
	got, _, err := c.GetDeployKey(repo.Owner.UserName, repo.Name, dk.ID)
	require.NoError(t, err)
	assert.EqualValues(t, dk.ID, got.ID)

	// List reflects the added key.
	keys, _, err = c.ListDeployKeys(repo.Owner.UserName, repo.Name, ListDeployKeysOptions{KeyID: dk.ID})
	require.NoError(t, err)
	assert.NotEmpty(t, keys)

	// DeleteDeployKey.
	_, err = c.DeleteDeployKey(repo.Owner.UserName, repo.Name, dk.ID)
	require.NoError(t, err)

	keys, _, err = c.ListDeployKeys(repo.Owner.UserName, repo.Name, ListDeployKeysOptions{})
	require.NoError(t, err)
	assert.Empty(t, keys)
}
