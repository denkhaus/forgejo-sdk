// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRepoGapsV15 exercises the remaining gap endpoints.
func TestRepoGapsV15(t *testing.T) {
	log.Println("== TestRepoGapsV15 ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "v15-gaps-repo", c)
	require.NoError(t, err)

	// editorconfig (tolerant: depends on a .editorconfig being present)
	if ec, _, err := c.GetRepoEditorConfig(repo.Owner.UserName, repo.Name, "README.md"); err == nil {
		assert.NotNil(t, ec)
	}

	cfg, _, err := c.GetRepoIssueConfig(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	v, _, err := c.ValidateRepoIssueConfig(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	assert.NotNil(t, v)

	subs, _, err := c.ListRepoSubscribers(repo.Owner.UserName, repo.Name, ListStargazersOptions{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	_ = subs

	issue := createTestIssue(t, c, repo.Name, "timeline-issue", "", nil, nil, 0, nil, false, false)
	require.NotNil(t, issue)

	timeline, _, err := c.GetIssueTimeline(repo.Owner.UserName, repo.Name, issue.Index, ListIssueTimelineOptions{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	_ = timeline

	tok, _, err := c.GetMyRunnerRegistrationToken()
	require.NoError(t, err)
	assert.NotNil(t, tok)
}

// TestRunnersAllLevelsV15 exercises org/repo/user runner registration+management.
func TestRunnersAllLevelsV15(t *testing.T) {
	log.Println("== TestRunnersAllLevelsV15 ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{Name: "v15-runners-org"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg(org.UserName) }()

	repo, err := createTestRepo(t, "v15-runners-repo", c)
	require.NoError(t, err)

	// ---- org level ----
	orr, _, err := c.RegisterOrgRunner(org.UserName, models.RegisterRunnerOptions{Name: OptionalString("org-runner"), Ephemeral: true})
	require.NoError(t, err)
	require.NotZero(t, orr.ID)
	_, _, err = c.GetOrgRunner(org.UserName, orr.ID)
	require.NoError(t, err)
	_, _, err = c.ListOrgRunners(org.UserName, ListActionRunnersOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	_, err = c.DeleteOrgRunner(org.UserName, orr.ID)
	require.NoError(t, err)

	// ---- repo level ----
	rrr, _, err := c.RegisterRepoRunner(repo.Owner.UserName, repo.Name, models.RegisterRunnerOptions{Name: OptionalString("repo-runner"), Ephemeral: true})
	require.NoError(t, err)
	require.NotZero(t, rrr.ID)
	_, _, err = c.GetRepoRunner(repo.Owner.UserName, repo.Name, rrr.ID)
	require.NoError(t, err)
	_, _, err = c.ListRepoRunners(repo.Owner.UserName, repo.Name, ListActionRunnersOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	_, err = c.DeleteRepoRunner(repo.Owner.UserName, repo.Name, rrr.ID)
	require.NoError(t, err)

	// ---- user level ----
	urr, _, err := c.RegisterUserRunner(models.RegisterRunnerOptions{Name: OptionalString("user-runner"), Ephemeral: true})
	require.NoError(t, err)
	require.NotZero(t, urr.ID)
	_, _, err = c.GetUserRunner(urr.ID)
	require.NoError(t, err)
	_, _, err = c.ListUserRunners(ListActionRunnersOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	_, err = c.DeleteUserRunner(urr.ID)
	require.NoError(t, err)
}

// TestAdminBatchV15 exercises admin emails, runner jobs, unadopted repos, rename.
func TestAdminBatchV15(t *testing.T) {
	log.Println("== TestAdminBatchV15 ==")
	c := newTestClient()

	_, _, err := c.AdminListEmails(AdminListEmailsOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)

	_, _, err = c.AdminSearchEmails(AdminSearchEmailsOption{Query: "test01", ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)

	_, _, err = c.AdminGetActionRunJobs(AdminSearchActionRunJobsOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)

	_, _, err = c.AdminSearchRunJobs(AdminSearchActionRunJobsOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)

	_, _, err = c.AdminListUnadoptedRepositories(AdminListUnadoptedOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)

	u := createTestUser(t, "v15-rename-src", c)
	_, err = c.AdminRenameUser(u.UserName, models.RenameUserOption{NewName: OptionalString("v15-rename-dst")})
	require.NoError(t, err)
}

// TestOrgBlockRenameV15 exercises org block/unblock and org rename.
func TestOrgBlockRenameV15(t *testing.T) {
	log.Println("== TestOrgBlockRenameV15 ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{Name: "v15-org-block"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg("v15-org-block-renamed") }() // renamed below

	target := createTestUser(t, "v15-block-target", c)

	_, err = c.BlockOrgUser(org.UserName, target.UserName)
	require.NoError(t, err)

	blocked, _, err := c.ListOrgBlockedUsers(org.UserName, ListBlockedUsersOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	assert.NotEmpty(t, blocked)

	_, err = c.UnblockOrgUser(org.UserName, target.UserName)
	require.NoError(t, err)

	_, err = c.RenameOrg(org.UserName, models.RenameOrgOption{NewName: OptionalString("v15-org-block-renamed")})
	require.NoError(t, err)
}

// TestUserExtrasV15 exercises heatmap, GPG verify token and issue deadline.
func TestUserExtrasV15(t *testing.T) {
	log.Println("== TestUserExtrasV15 ==")
	c := newTestClient()
	user, _, err := c.GetMyUserInfo()
	require.NoError(t, err)

	_, _, err = c.GetUserHeatmapData(user.UserName)
	require.NoError(t, err)

	_, _, err = c.GetMyGPGKeyVerificationToken()
	require.NoError(t, err)

	repo, err := createTestRepo(t, "v15-deadline-repo", c)
	require.NoError(t, err)
	issue := createTestIssue(t, c, repo.Name, "deadline-issue", "", nil, nil, 0, nil, false, false)
	require.NotNil(t, issue)

	due := strfmt.DateTime(time.Now().Add(48 * time.Hour))
	_, _, err = c.EditIssueDeadline(repo.Owner.UserName, repo.Name, issue.Index, models.EditDeadlineOption{Deadline: &due})
	require.NoError(t, err)
}

// TestMiscMoreV15 exercises markup rendering and template lookups.
func TestMiscMoreV15(t *testing.T) {
	log.Println("== TestMiscMoreV15 ==")
	c := newTestClient()

	html, _, err := c.RenderMarkup(models.MarkupOption{Text: "<b>hi</b>", Mode: "markdown"})
	require.NoError(t, err)
	_ = html

	git, _, err := c.ListGitignoreTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, git)
	info, _, err := c.GetGitignoreTemplateInfo(git[0])
	require.NoError(t, err)
	assert.NotNil(t, info)

	labels, _, err := c.ListLabelTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, labels)
	li, _, err := c.GetLabelTemplateInfo(labels[0])
	require.NoError(t, err)
	assert.NotEmpty(t, li)

	lics, _, err := c.ListLicenseTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, lics)
	lci, _, err := c.GetLicenseTemplateInfo(lics[0].Name)
	require.NoError(t, err)
	assert.NotNil(t, lci)
}
