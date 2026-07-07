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

// TestRepoIssueConfigV15 exercises fetching and validating repo issue config.
func TestRepoIssueConfigV15(t *testing.T) {
	log.Println("== TestRepoIssueConfigV15 ==")
	c := newTestClient()

	repo, err := createTestRepo(t, "v15-issueconfig-repo", c)
	require.NoError(t, err)

	cfg, _, err := c.GetRepoIssueConfig(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	v, _, err := c.ValidateRepoIssueConfig(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	assert.NotNil(t, v)
}
