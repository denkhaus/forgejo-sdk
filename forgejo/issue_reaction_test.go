// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_IssueReaction exercises adding, listing and removing a reaction on an issue.
func Test_IssueReaction(t *testing.T) {
	log.Println("== Test_IssueReaction ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "TestIssueReactionRepo", c)
	require.NoError(t, err)
	issue := createTestIssue(t, c, repo.Name, "Reaction Issue", "", nil, nil, 0, nil, false, false)
	owner := repo.Owner.UserName

	reaction, _, err := c.PostIssueReaction(owner, repo.Name, issue.Index, "+1")
	require.NoError(t, err)
	if assert.NotNil(t, reaction) {
		assert.EqualValues(t, "+1", reaction.Reaction)
	}

	reactions, _, err := c.GetIssueReactions(owner, repo.Name, issue.Index)
	require.NoError(t, err)
	assert.NotEmpty(t, reactions)

	_, err = c.DeleteIssueReaction(owner, repo.Name, issue.Index, "+1")
	require.NoError(t, err)

	reactionsAfter, _, err := c.GetIssueReactions(owner, repo.Name, issue.Index)
	require.NoError(t, err)
	assert.Empty(t, reactionsAfter)
}
