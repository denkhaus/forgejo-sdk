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

func Test_Attachment(t *testing.T) {
	log.Println("== Test_Attachment ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "ReleaseAttachment", c)
	require.NoError(t, err)

	owner := repo.Owner.UserName
	repoName := repo.Name

	// Release attachments require an existing release, which in turn needs a
	// tag target. Tolerate release-creation failure and skip the rest.
	release, _, err := c.CreateRelease(owner, repoName, CreateReleaseOption{
		TagName: "att-tag",
		Target:  repo.DefaultBranch,
		Title:   "Attachment Release",
		Note:    "release for attachment test",
	})
	if err != nil {
		t.Skipf("could not create release for attachment test: %v", err)
	}
	defer func() { _, _ = c.DeleteRelease(owner, repoName, release.ID) }()

	// ListReleaseAttachments should be safe and empty initially.
	list, _, err := c.ListReleaseAttachments(owner, repoName, release.ID, ListReleaseAttachmentsOptions{})
	require.NoError(t, err)
	assert.Empty(t, list)

	body := []byte("attachment payload")
	created, _, err := c.CreateReleaseAttachment(owner, repoName, release.ID, bytes.NewReader(body), "release-file.txt")
	require.NoError(t, err)
	assert.NotZero(t, created.ID)

	// GetReleaseAttachment
	got, _, err := c.GetReleaseAttachment(owner, repoName, release.ID, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, got.ID)

	// EditReleaseAttachment
	edited, _, err := c.EditReleaseAttachment(owner, repoName, release.ID, created.ID, models.EditAttachmentOptions{
		Name: "renamed-file.txt",
	})
	require.NoError(t, err)
	assert.Equal(t, "renamed-file.txt", edited.Name)

	// ListReleaseAttachments should now contain one entry.
	list, _, err = c.ListReleaseAttachments(owner, repoName, release.ID, ListReleaseAttachmentsOptions{})
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// DeleteReleaseAttachment
	_, err = c.DeleteReleaseAttachment(owner, repoName, release.ID, created.ID)
	require.NoError(t, err)

	list, _, err = c.ListReleaseAttachments(owner, repoName, release.ID, ListReleaseAttachmentsOptions{})
	require.NoError(t, err)
	assert.Empty(t, list)
}
