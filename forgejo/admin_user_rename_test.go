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

// TestAdminRenameUserV15 exercises renaming a user via the admin API.
func TestAdminRenameUserV15(t *testing.T) {
	log.Println("== TestAdminRenameUserV15 ==")
	c := newTestClient()

	u := createTestUser(t, "v15-rename-src", c)
	_, err := c.AdminRenameUser(u.UserName, models.RenameUserOption{NewName: OptionalString("v15-rename-dst")})
	require.NoError(t, err)
}
