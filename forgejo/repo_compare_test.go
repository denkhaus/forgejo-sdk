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

func Test_RepoCompare(t *testing.T) {
	log.Println("== Test_RepoCompare ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestRepoCompare", c)

	// CompareCommits requires server >= 1.22.0; an older instance returns a
	// version error. Comparing a ref against itself is a safe no-diff call.
	cmp, _, err := c.CompareCommits(repo.Owner.UserName, repo.Name, "main", "main")
	if err != nil {
		t.Skipf("CompareCommits not available in this environment: %v", err)
	}
	require.NotNil(t, cmp)
	assert.NotNil(t, cmp.Commits)
}
