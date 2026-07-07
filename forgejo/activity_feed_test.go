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

// TestActivityFeedsV15 exercises the repo/user activity-feed endpoints.
func TestActivityFeedsV15(t *testing.T) {
	log.Println("== TestActivityFeedsV15 ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "v15-feeds-repo", c)
	require.NoError(t, err)

	opt := ListActivityFeedsOption{ListOptions: ListOptions{PageSize: 5}}
	feeds, _, err := c.ListRepoActivityFeeds(repo.Owner.UserName, repo.Name, opt)
	require.NoError(t, err)
	assert.NotNil(t, feeds)

	_, _, err = c.ListUserActivityFeeds(repo.Owner.UserName, opt)
	require.NoError(t, err)
}
