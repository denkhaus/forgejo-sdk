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

// TestIssuePinningV15 exercises issue pinning (Forgejo v15+).
func TestIssuePinningV15(t *testing.T) {
	log.Println("== TestIssuePinningV15 ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "v15-pin-repo", c)
	require.NoError(t, err)

	allowed, _, err := c.NewIssuePinsAllowed(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	require.NotNil(t, allowed)

	issue := createTestIssue(t, c, repo.Name, "pin-target", "", nil, nil, 0, nil, false, false)
	require.NotNil(t, issue)

	_, err = c.PinIssue(repo.Owner.UserName, repo.Name, issue.Index)
	require.NoError(t, err)

	pinned, _, err := c.ListPinnedIssues(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	require.Len(t, pinned, 1)
	assert.Equal(t, issue.Index, pinned[0].Index)

	_, err = c.UnpinIssue(repo.Owner.UserName, repo.Name, issue.Index)
	require.NoError(t, err)

	pinned, _, err = c.ListPinnedIssues(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	assert.Empty(t, pinned)
}
