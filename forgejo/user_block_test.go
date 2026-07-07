// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserBlockV15 exercises blocking/unblocking a user (Forgejo v15+).
func TestUserBlockV15(t *testing.T) {
	log.Println("== TestUserBlockV15 ==")
	c := newTestClient()
	u := createTestUser(t, "v15-block-user", c)

	_, err := c.BlockUser(u.UserName)
	require.NoError(t, err)

	blocked, _, err := c.ListBlockedUsers(ListBlockedUsersOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotEmpty(t, blocked)

	_, err = c.UnblockUser(u.UserName)
	require.NoError(t, err)
}
