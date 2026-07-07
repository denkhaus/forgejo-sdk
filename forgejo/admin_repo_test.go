// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_AdminRepo exercises AdminCreateRepo, creating a repository owned by a
// test user on their behalf, then cleaning it up.
func Test_AdminRepo(t *testing.T) {
	log.Println("== Test_AdminRepo ==")
	c := newTestClient()

	user := createTestUser(t, "admin-repo-user", c)
	repoName := "AdminCreatedRepo"
	// Ensure a clean slate if a previous run left the repo behind.
	_, _ = c.DeleteRepo(user.UserName, repoName)

	repo, _, err := c.AdminCreateRepo(user.UserName, CreateRepoOption{
		Name:        repoName,
		Description: "repo created via admin API",
		AutoInit:    true,
		Private:     false,
	})
	require.NoError(t, err)
	if assert.NotNil(t, repo) {
		assert.EqualValues(t, repoName, repo.Name)
		assert.EqualValues(t, user.UserName, repo.Owner.UserName)
	}

	// cleanup: remove the repo we created (ignore errors if it is already gone).
	if repo != nil {
		_, _ = c.DeleteRepo(user.UserName, repoName)
	}
}
