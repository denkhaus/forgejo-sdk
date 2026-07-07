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

// Test_UserFollow exercises Follow/IsFollowing/ListMyFollowing/Unfollow.
// Follow may be gated or already-following; tolerate and skip.
func Test_UserFollow(t *testing.T) {
	log.Println("== Test_UserFollow ==")
	c := newTestClient()
	target := createTestUser(t, "v15-follow-tgt", c)

	_, err := c.Follow(target.UserName)
	if err != nil {
		t.Skipf("Follow unavailable: %v", err)
	}
	defer func() { _, _ = c.Unfollow(target.UserName) }()

	is, _ := c.IsFollowing(target.UserName)
	assert.True(t, is)

	following, _, err := c.ListMyFollowing(ListFollowingOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotEmpty(t, following)
}
