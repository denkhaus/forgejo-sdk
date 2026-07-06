// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"log"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// minimal 1x1 transparent PNG, base64-encoded
const onePxPNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="

// TestOrgLabelsV15 exercises organization label CRUD (Forgejo v15+).
func TestOrgLabelsV15(t *testing.T) {
	log.Println("== TestOrgLabelsV15 ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{Name: "v15-org-labels"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg(org.UserName) }()

	l1, _, err := c.CreateOrgLabel(org.UserName, models.CreateLabelOption{Name: OptionalString("org-label-1"), Color: OptionalString("#aabbcc")})
	require.NoError(t, err)
	assert.NotZero(t, l1.ID)

	list, _, err := c.ListOrgLabels(org.UserName, ListOrgLabelsOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotEmpty(t, list)

	got, _, err := c.GetOrgLabel(org.UserName, l1.ID)
	require.NoError(t, err)
	assert.Equal(t, l1.ID, got.ID)

	_, _, err = c.EditOrgLabel(org.UserName, l1.ID, models.EditLabelOption{Description: "updated"})
	require.NoError(t, err)

	_, err = c.DeleteOrgLabel(org.UserName, l1.ID)
	require.NoError(t, err)

	list, _, err = c.ListOrgLabels(org.UserName, ListOrgLabelsOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.Empty(t, list)
}

// TestTagProtectionsV15 exercises repository tag-protection rules.
func TestTagProtectionsV15(t *testing.T) {
	log.Println("== TestTagProtectionsV15 ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "v15-tag-prot-repo", c)
	require.NoError(t, err)

	tp, _, err := c.CreateRepoTagProtection(repo.Owner.UserName, repo.Name, models.CreateTagProtectionOption{NamePattern: "v*", WhitelistUsernames: []string{repo.Owner.UserName}})
	require.NoError(t, err)
	assert.NotZero(t, tp.ID)

	list, _, err := c.ListRepoTagProtections(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	require.Len(t, list, 1)

	got, _, err := c.GetRepoTagProtection(repo.Owner.UserName, repo.Name, tp.ID)
	require.NoError(t, err)
	assert.Equal(t, tp.ID, got.ID)

	_, _, err = c.EditRepoTagProtection(repo.Owner.UserName, repo.Name, tp.ID, models.EditTagProtectionOption{NamePattern: "release-*"})
	require.NoError(t, err)

	_, err = c.DeleteRepoTagProtection(repo.Owner.UserName, repo.Name, tp.ID)
	require.NoError(t, err)

	list, _, err = c.ListRepoTagProtections(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	assert.Empty(t, list)
}

// TestGitNotesV15 exercises get/set/remove of git notes on a commit.
func TestGitNotesV15(t *testing.T) {
	log.Println("== TestGitNotesV15 ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "v15-notes-repo", c)
	require.NoError(t, err)

	commits, _, err := c.ListRepoCommits(repo.Owner.UserName, repo.Name, ListCommitOptions{ListOptions: ListOptions{PageSize: 1}})
	require.NoError(t, err)
	require.NotEmpty(t, commits)
	sha := commits[0].SHA

	_, _, err = c.SetNote(repo.Owner.UserName, repo.Name, sha, models.NoteOptions{Message: "a v15 note"})
	require.NoError(t, err)

	note, _, err := c.GetNote(repo.Owner.UserName, repo.Name, sha)
	require.NoError(t, err)
	assert.Contains(t, note.Message, "a v15 note")

	_, err = c.RemoveNote(repo.Owner.UserName, repo.Name, sha)
	require.NoError(t, err)
}

// TestIssueAttachmentsV15 exercises issue-attachment CRUD (multipart upload).
func TestIssueAttachmentsV15(t *testing.T) {
	log.Println("== TestIssueAttachmentsV15 ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "v15-issue-att-repo", c)
	require.NoError(t, err)

	issue := createTestIssue(t, c, repo.Name, "attachment-issue", "", nil, nil, 0, nil, false, false)
	require.NotNil(t, issue)

	a, _, err := c.CreateIssueAttachment(repo.Owner.UserName, repo.Name, issue.Index, bytes.NewReader([]byte("hello attachment")), "note.txt")
	require.NoError(t, err)
	assert.NotZero(t, a.ID)

	list, _, err := c.ListIssueAttachments(repo.Owner.UserName, repo.Name, issue.Index, ListIssueAttachmentsOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	require.Len(t, list, 1)

	got, _, err := c.GetIssueAttachment(repo.Owner.UserName, repo.Name, issue.Index, a.ID)
	require.NoError(t, err)
	assert.Equal(t, a.ID, got.ID)

	_, _, err = c.EditIssueAttachment(repo.Owner.UserName, repo.Name, issue.Index, a.ID, models.EditAttachmentOptions{Name: "renamed.txt"})
	require.NoError(t, err)

	_, err = c.DeleteIssueAttachment(repo.Owner.UserName, repo.Name, issue.Index, a.ID)
	require.NoError(t, err)

	list, _, err = c.ListIssueAttachments(repo.Owner.UserName, repo.Name, issue.Index, ListIssueAttachmentsOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.Empty(t, list)
}

// TestAvatarsV15 exercises user and repository avatar update/delete.
func TestAvatarsV15(t *testing.T) {
	log.Println("== TestAvatarsV15 ==")
	c := newTestClient()

	_, err := c.UpdateUserAvatar(models.UpdateUserAvatarOption{Image: onePxPNG})
	require.NoError(t, err)
	_, err = c.DeleteUserAvatar()
	require.NoError(t, err)

	repo, err := createTestRepo(t, "v15-avatar-repo", c)
	require.NoError(t, err)

	_, err = c.UpdateRepoAvatar(repo.Owner.UserName, repo.Name, models.UpdateRepoAvatarOption{Image: onePxPNG})
	require.NoError(t, err)
	_, err = c.DeleteRepoAvatar(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
}
