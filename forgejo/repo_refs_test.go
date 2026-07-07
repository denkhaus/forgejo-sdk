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

func Test_RepoRefs(t *testing.T) {
	log.Println("== Test_RepoRefs ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestRepoRefs", c)

	// GetRepoRefs — safe GET. An auto-init repo has at least one branch
	// ("main") under refs/heads.
	refs, _, err := c.GetRepoRefs(repo.Owner.UserName, repo.Name, "heads")
	require.NoError(t, err)
	assert.NotEmpty(t, refs)

	// GetRepoRef — fetch the "main" branch ref directly (full ref path).
	ref, _, err := c.GetRepoRef(repo.Owner.UserName, repo.Name, "refs/heads/main")
	if err == nil {
		assert.NotNil(t, ref)
		assert.NotEmpty(t, ref.Ref)
	}
}
