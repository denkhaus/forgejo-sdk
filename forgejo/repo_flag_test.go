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

func Test_RepoFlag(t *testing.T) {
	log.Println("== Test_RepoFlag ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "RepoFlag", c)
	require.NoError(t, err)

	owner := repo.Owner.UserName
	repoName := repo.Name

	// ListRepoFlags is a safe read and always available.
	flags, _, err := c.ListRepoFlags(owner, repoName)
	if err != nil {
		t.Skipf("repo flags not available on this instance: %v", err)
	}
	assert.NotNil(t, flags)

	// Flags require admin preconfiguration of allowed flag names. Adding a
	// flag may therefore be rejected; tolerate that and only assert the
	// read paths that depend on success.
	flag := "test-flag"
	if _, addErr := c.AddRepoFlag(owner, repoName, flag); addErr == nil {
		// CheckRepoFlag returns nil error when the flag is set.
		_, checkErr := c.CheckRepoFlag(owner, repoName, flag)
		assert.NoError(t, checkErr)

		flagsAfter, _, listErr := c.ListRepoFlags(owner, repoName)
		require.NoError(t, listErr)
		assert.Contains(t, flagsAfter, flag)

		// ReplaceRepoFlags
		_, replaceErr := c.ReplaceRepoFlags(owner, repoName, models.ReplaceFlagsOption{
			Flags: []string{flag, "second"},
		})
		require.NoError(t, replaceErr)

		// DeleteRepoFlag
		_, delErr := c.DeleteRepoFlag(owner, repoName, flag)
		assert.NoError(t, delErr)
	}

	// DeleteAllRepoFlags cleans up regardless of prior state.
	_, _ = c.DeleteAllRepoFlags(owner, repoName)

	// Final list should be clean.
	flags, _, err = c.ListRepoFlags(owner, repoName)
	require.NoError(t, err)
	assert.Empty(t, flags)
}
