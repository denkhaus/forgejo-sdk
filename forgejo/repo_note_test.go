// Copyright 2024 The Forgejo Authors. All rights reserved.
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
