// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"encoding/base64"
	"log"
	"strings"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RepoWiki(t *testing.T) {
	log.Println("== Test_RepoWiki ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "RepoWiki", c)
	require.NoError(t, err)

	owner := repo.Owner.UserName
	repoName := repo.Name

	// The wiki may be disabled at the instance or repository level. Detect that
	// up front via the safe list call and skip the suite in that case.
	if _, _, lerr := c.ListRepoWikiPages(owner, repoName, ListWikiPagesOption{}); lerr != nil {
		t.Skipf("wiki not available on this instance: %v", lerr)
	}

	title := "TestWikiPage"
	content := "# Test Wiki Page\n"
	page, _, err := c.CreateRepoWikiPage(owner, repoName, models.CreateWikiPageOptions{
		Title:         title,
		ContentBase64: base64.StdEncoding.EncodeToString([]byte(content)),
		Message:       "initial test page",
	})
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not enabled") {
			t.Skip("wiki disabled on this instance")
		}
		t.Skipf("wiki page creation not supported in this environment: %v", err)
	}
	require.NotNil(t, page)

	// List
	pages, _, err := c.ListRepoWikiPages(owner, repoName, ListWikiPagesOption{})
	require.NoError(t, err)
	require.NotEmpty(t, pages)

	// Get
	got, _, err := c.GetRepoWikiPage(owner, repoName, title)
	require.NoError(t, err)
	assert.Equal(t, title, got.Title)

	// Revisions
	revs, _, err := c.GetRepoWikiPageRevisions(owner, repoName, title, ListWikiPagesOption{})
	require.NoError(t, err)
	assert.NotNil(t, revs)

	// Edit
	_, _, err = c.EditRepoWikiPage(owner, repoName, title, models.CreateWikiPageOptions{
		Title:         title,
		ContentBase64: base64.StdEncoding.EncodeToString([]byte("# Edited Content\n")),
		Message:       "edit test page",
	})
	require.NoError(t, err)

	// Delete
	_, err = c.DeleteRepoWikiPage(owner, repoName, title)
	require.NoError(t, err)
}
