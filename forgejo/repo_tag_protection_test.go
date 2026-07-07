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
