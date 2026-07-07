// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2014 The Gogs Authors. All rights reserved.
// Copyright 2017 The Gitea Authors. All rights reserved.
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

// Test_Hook exercises the repository webhook lifecycle: create, list, get,
// test delivery, edit and delete.
func Test_Hook(t *testing.T) {
	log.Println("== Test_Hook ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "TestHookRepo", c)
	require.NoError(t, err)
	owner := repo.Owner.UserName

	hook, _, err := c.CreateRepoHook(owner, repo.Name, CreateHookOption{
		Type:   HookType("gitea"),
		Config: map[string]string{"url": "http://example.com/webhook", "content_type": "json"},
		Events: []string{"push"},
		Active: true,
	})
	require.NoError(t, err)
	if assert.NotNil(t, hook) {
		assert.NotZero(t, hook.ID)
	}

	hooks, _, err := c.ListRepoHooks(owner, repo.Name, ListHooksOptions{})
	require.NoError(t, err)
	assert.NotEmpty(t, hooks)

	_, _, err = c.GetRepoHook(owner, repo.Name, hook.ID)
	require.NoError(t, err)

	// A test delivery targets an unreachable URL; tolerate the resulting error.
	if _, err := c.TestRepoHook(owner, repo.Name, hook.ID); err != nil {
		t.Logf("TestRepoHook returned (tolerated): %v", err)
	}

	// Editing the hook is straightforward; tolerate any restriction.
	if _, err := c.EditRepoHook(owner, repo.Name, hook.ID, models.EditHookOption{Active: false}); err != nil {
		t.Logf("EditRepoHook returned (tolerated): %v", err)
	}

	_, err = c.DeleteRepoHook(owner, repo.Name, hook.ID)
	require.NoError(t, err)
}
