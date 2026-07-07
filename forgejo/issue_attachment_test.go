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
