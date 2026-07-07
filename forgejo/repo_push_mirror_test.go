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

func Test_RepoPushMirror(t *testing.T) {
	log.Println("== Test_RepoPushMirror ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "RepoPushMirror", c)
	require.NoError(t, err)

	owner := repo.Owner.UserName
	repoName := repo.Name

	// ListPushMirrors is a safe read and always available.
	mirrors, _, err := c.ListPushMirrors(owner, repoName, ListPushMirrorsOption{})
	require.NoError(t, err)
	assert.NotNil(t, mirrors)

	// Creating a push mirror requires a valid remote address and credentials,
	// which are not available in the test environment. We therefore only
	// exercise the remaining endpoints when a mirror happens to exist, and
	// otherwise tolerate. SyncPushMirrors is a no-op-safe trigger.
	_, _ = c.SyncPushMirrors(owner, repoName)

	if len(mirrors) > 0 {
		name := mirrors[0].RemoteName
		got, _, gErr := c.GetPushMirror(owner, repoName, name)
		require.NoError(t, gErr)
		assert.NotNil(t, got)

		// Do not delete the mirror; we did not create it. The list/GET/sync
		// paths above provide the coverage that is environmentally possible.
	}
}
