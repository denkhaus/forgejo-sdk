// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"strings"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserQuotaV15 exercises the current-user quota endpoints (Forgejo v15+).
func TestUserQuotaV15(t *testing.T) {
	log.Println("== TestUserQuotaV15 ==")
	c := newTestClient()

	q, _, err := c.GetMyQuota()
	if err != nil && strings.Contains(err.Error(), "404") {
		t.Skip("quota API not enabled on this instance")
	}
	require.NoError(t, err)
	require.NotNil(t, q)
	assert.NotNil(t, q.Used)

	ok, _, err := c.CheckMyQuota()
	require.NoError(t, err)
	_ = ok // bool: within quota

	opt := ListQuotaUsedOption{ListOptions: ListOptions{PageSize: 1}}
	_, _, err = c.ListMyQuotaArtifacts(opt)
	require.NoError(t, err)
	_, _, err = c.ListMyQuotaAttachments(opt)
	require.NoError(t, err)
	_, _, err = c.ListMyQuotaPackages(opt)
	require.NoError(t, err)
}

// TestUserBlockV15 exercises blocking/unblocking a user (Forgejo v15+).
func TestUserBlockV15(t *testing.T) {
	log.Println("== TestUserBlockV15 ==")
	c := newTestClient()
	u := createTestUser(t, "v15-block-user", c)

	_, err := c.BlockUser(u.UserName)
	require.NoError(t, err)

	blocked, _, err := c.ListBlockedUsers(ListBlockedUsersOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotEmpty(t, blocked)

	_, err = c.UnblockUser(u.UserName)
	require.NoError(t, err)
}

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

// TestActivityFeedsV15 exercises the repo/user activity-feed endpoints.
func TestActivityFeedsV15(t *testing.T) {
	log.Println("== TestActivityFeedsV15 ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "v15-feeds-repo", c)
	require.NoError(t, err)

	opt := ListActivityFeedsOption{ListOptions: ListOptions{PageSize: 5}}
	feeds, _, err := c.ListRepoActivityFeeds(repo.Owner.UserName, repo.Name, opt)
	require.NoError(t, err)
	assert.NotNil(t, feeds)

	_, _, err = c.ListUserActivityFeeds(repo.Owner.UserName, opt)
	require.NoError(t, err)
}

// TestMiscV15 exercises instance-level v15 endpoints.
func TestMiscV15(t *testing.T) {
	log.Println("== TestMiscV15 ==")
	c := newTestClient()

	ni, _, err := c.GetNodeInfo()
	if err == nil { // nodeinfo may be disabled (federation) on the instance
		require.NotNil(t, ni)
		assert.NotEmpty(t, ni.Version)
	}

	// signing keys depend on instance configuration and may be absent/unset
	_, _, _ = c.GetSigningKey()
	_, _, _ = c.GetSSHSigningKey()

	html, _, err := c.RenderMarkdown(models.MarkdownOption{Text: "# hello v15", Mode: "markdown"})
	require.NoError(t, err)
	assert.Contains(t, html, "hello v15")

	raw, _, err := c.RenderMarkdownRaw([]byte("# raw hello"))
	require.NoError(t, err)
	assert.Contains(t, raw, "raw hello")

	git, _, err := c.ListGitignoreTemplates()
	require.NoError(t, err)
	assert.NotEmpty(t, git)

	_, _, err = c.ListLicenseTemplates()
	require.NoError(t, err)

	_, _, err = c.ListLabelTemplates()
	require.NoError(t, err)

	_, _, err = c.SearchTopics(SearchTopicOption{Query: "forgejo", ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
}

// TestRunnersV15 exercises the v15 interactive runner registration flow.
func TestRunnersV15(t *testing.T) {
	log.Println("== TestRunnersV15 ==")
	c := newTestClient()

	// listing works even when no runners exist
	_, _, err := c.ListAdminRunners(ListActionRunnersOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)

	// register an ephemeral runner via the v15 interactive flow
	name := "v15-test-runner"
	rr, _, err := c.RegisterAdminRunner(models.RegisterRunnerOptions{Name: &name, Ephemeral: true})
	require.NoError(t, err)
	require.NotNil(t, rr)
	assert.NotZero(t, rr.ID)
	assert.NotEmpty(t, rr.Token)
	assert.NotEmpty(t, rr.UUID)

	runners, _, err := c.ListAdminRunners(ListActionRunnersOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotEmpty(t, runners)

	got, _, err := c.GetAdminRunner(rr.ID)
	require.NoError(t, err)
	assert.Equal(t, rr.ID, got.ID)

	_, err = c.DeleteAdminRunner(rr.ID)
	require.NoError(t, err)

	_, _, err = c.GetAdminRunner(rr.ID)
	require.Error(t, err)
}
