// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRepoSubscribersV15 exercises listing repo subscribers.
func TestRepoSubscribersV15(t *testing.T) {
	log.Println("== TestRepoSubscribersV15 ==")
	c := newTestClient()

	repo, err := createTestRepo(t, "v15-subscribers-repo", c)
	require.NoError(t, err)

	_, _, err = c.ListRepoSubscribers(repo.Owner.UserName, repo.Name, ListStargazersOptions{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
}
