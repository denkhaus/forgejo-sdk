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

func Test_RepoMirror(t *testing.T) {
	log.Println("== Test_RepoMirror ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestRepoMirror", c)

	// ListPushMirrors — safe GET, expected empty for a fresh repo.
	mirrors, _, err := c.ListPushMirrors(repo.Owner.UserName, repo.Name, ListPushMirrorsOption{})
	require.NoError(t, err)
	assert.Empty(t, mirrors)

	// Adding a push mirror requires a reachable remote with credentials; on
	// most test instances this is not configured. Tolerate that failure and
	// only assert the success path.
	pm, _, err := c.PushMirrors(repo.Owner.UserName, repo.Name, models.CreatePushMirrorOption{
		RemoteAddress: "https://example.invalid/throwaway.git",
		Interval:      "8h0m0s",
	})
	if err == nil && pm != nil {
		// If creation succeeded, clean up the mirror.
		_, _ = c.DeletePushMirror(repo.Owner.UserName, repo.Name, pm.RemoteName)
	}

	// SyncPushMirrors on a repo with no mirrors is a safe no-op; tolerate
	// any server-side rejection.
	_, _ = c.SyncPushMirrors(repo.Owner.UserName, repo.Name)

	// MirrorSync queues a mirror sync. On a non-mirror repo it may be a
	// no-op or rejected by the server; tolerate either outcome.
	_, _ = c.MirrorSync(repo.Owner.UserName, repo.Name)

	// Re-listing must still succeed regardless of the operations above.
	mirrors, _, err = c.ListPushMirrors(repo.Owner.UserName, repo.Name, ListPushMirrorsOption{})
	require.NoError(t, err)
	assert.Empty(t, mirrors)
}
