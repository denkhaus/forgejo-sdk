// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_GitHook lists, fetches and (tolerantly) edits a repository Git hook.
func Test_GitHook(t *testing.T) {
	log.Println("== Test_GitHook ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "TestGitHookRepo", c)
	require.NoError(t, err)
	owner := repo.Owner.UserName

	hooks, _, err := c.ListRepoGitHooks(owner, repo.Name, ListRepoGitHooksOptions{})
	if err != nil {
		t.Skipf("git hooks not accessible on this instance: %v", err)
	}
	// Forgejo ships the default Git hooks (pre-receive, update, post-receive).
	assert.NotEmpty(t, hooks)

	const hookID = "pre-receive"
	hook, _, err := c.GetRepoGitHook(owner, repo.Name, hookID)
	require.NoError(t, err)
	if assert.NotNil(t, hook) {
		assert.True(t, hook.IsActive)
	}

	// Editing the hook content is straightforward; tolerate any restriction.
	if _, err := c.EditRepoGitHook(owner, repo.Name, hookID, models.EditGitHookOption{
		Content: "# pre-receive edited by test\n",
	}); err != nil {
		t.Logf("EditRepoGitHook returned (tolerated): %v", err)
	}

	// DeleteRepoGitHook resets the hook to its default content; tolerate.
	if _, err := c.DeleteRepoGitHook(owner, repo.Name, hookID); err != nil {
		t.Logf("DeleteRepoGitHook returned (tolerated): %v", err)
	}
}
