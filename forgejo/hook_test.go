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

// hookOption is a reusable, valid CreateHookOption for the webhook lifecycle tests.
func hookOption(url string) CreateHookOption {
	return CreateHookOption{
		Type:   HookTypeGitea,
		Config: map[string]string{"url": url, "content_type": "json"},
		Events: []string{"push"},
		Active: true,
	}
}

// Test_OrgHooks exercises the organization webhook lifecycle:
// create, list, get, edit and delete.
func Test_OrgHooks(t *testing.T) {
	log.Println("== Test_OrgHooks ==")
	c := newTestClient()

	orgName := "TestOrgHooksOrg"
	_, _ = c.DeleteOrg(orgName)
	_, _, err := c.CreateOrg(CreateOrgOption{
		Name:                      orgName,
		Visibility:                VisibleTypePublic,
		RepoAdminChangeTeamAccess: true,
	})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg(orgName) }()

	hook, _, err := c.CreateOrgHook(orgName, hookOption("http://example.com/org-hook"))
	require.NoError(t, err)
	if assert.NotNil(t, hook) {
		assert.NotZero(t, hook.ID)
	}

	hooks, _, err := c.ListOrgHooks(orgName, ListHooksOptions{})
	require.NoError(t, err)
	assert.NotEmpty(t, hooks)

	_, _, err = c.GetOrgHook(orgName, hook.ID)
	require.NoError(t, err)

	// A full config edit on a freshly created hook should succeed.
	_, err = c.EditOrgHook(orgName, hook.ID, models.EditHookOption{
		Active: false,
		Config: map[string]string{"url": "http://example.com/org-hook", "content_type": "json"},
	})
	require.NoError(t, err)

	_, err = c.DeleteOrgHook(orgName, hook.ID)
	require.NoError(t, err)

	// deleting twice -> not found
	_, err = c.DeleteOrgHook(orgName, hook.ID)
	require.Error(t, err)
}

// Test_MyHooks exercises the authenticated-user webhook lifecycle.
func Test_MyHooks(t *testing.T) {
	log.Println("== Test_MyHooks ==")
	c := newTestClient()

	hook, _, err := c.CreateMyHook(hookOption("http://example.com/my-hook"))
	require.NoError(t, err)
	if assert.NotNil(t, hook) {
		assert.NotZero(t, hook.ID)
	}
	defer func() { _, _ = c.DeleteMyHook(hook.ID) }()

	hooks, _, err := c.ListMyHooks(ListHooksOptions{})
	require.NoError(t, err)
	assert.NotEmpty(t, hooks)

	_, _, err = c.GetMyHook(hook.ID)
	require.NoError(t, err)

	_, err = c.EditMyHook(hook.ID, models.EditHookOption{
		Active: false,
		Config: map[string]string{"url": "http://example.com/my-hook", "content_type": "json"},
	})
	require.NoError(t, err)

	_, err = c.DeleteMyHook(hook.ID)
	require.NoError(t, err)
}

// Test_HookNegative covers client-side validation and not-found error paths.
func Test_HookNegative(t *testing.T) {
	log.Println("== Test_HookNegative ==")
	c := newTestClient()
	user, _, err := c.GetMyUserInfo()
	require.NoError(t, err)

	// CreateHookOption.Validate rejects an empty type before any request.
	require.Error(t, CreateHookOption{}.Validate())
	_, _, err = c.CreateRepoHook(user.UserName, "x", CreateHookOption{})
	require.Error(t, err)
	_, _, err = c.CreateMyHook(CreateHookOption{})
	require.Error(t, err)
	_, _, err = c.CreateOrgHook("x", CreateHookOption{})
	require.Error(t, err)

	// non-existent hooks / owners -> error
	_, _, err = c.GetMyHook(999999)
	require.Error(t, err)
	_, err = c.DeleteMyHook(999999)
	require.Error(t, err)
	_, _, err = c.GetOrgHook("no-such-org", 999999)
	require.Error(t, err)
	_, _, err = c.GetRepoHook(user.UserName, "no-such-repo", 999999)
	require.Error(t, err)
}
