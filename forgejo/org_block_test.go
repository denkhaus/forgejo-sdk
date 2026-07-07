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

// TestOrgBlockV15 exercises org block/unblock of a user.
func TestOrgBlockV15(t *testing.T) {
	log.Println("== TestOrgBlockV15 ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{Name: "v15-org-block"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg("v15-org-block") }()

	target := createTestUser(t, "v15-block-target", c)

	_, err = c.BlockOrgUser(org.UserName, target.UserName)
	require.NoError(t, err)

	blocked, _, err := c.ListOrgBlockedUsers(org.UserName, ListBlockedUsersOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	assert.NotEmpty(t, blocked)

	_, err = c.UnblockOrgUser(org.UserName, target.UserName)
	require.NoError(t, err)
}
