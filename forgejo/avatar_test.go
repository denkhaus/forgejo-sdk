// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/require"
)

// minimal 1x1 transparent PNG, base64-encoded
const onePxPNG = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNkYAAAAAYAAjCB0C8AAAAASUVORK5CYII="

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
