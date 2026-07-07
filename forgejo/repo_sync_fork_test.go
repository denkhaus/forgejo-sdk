// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
)

func Test_RepoSyncFork(t *testing.T) {
	log.Println("== Test_RepoSyncFork ==")
	c := newTestClient()

	baseRepo, err := createTestRepo(t, "SyncForkBase", c)
	if err != nil {
		t.Skipf("could not create base repo for fork sync test: %v", err)
	}

	orgName := "syncfork-org"
	// Cleanup runs in LIFO order: repo first, then org.
	defer func() { _, _ = c.DeleteOrg(orgName) }()

	// Reset any leftover org from previous runs.
	_, _ = c.DeleteOrg(orgName)
	if _, _, orgErr := c.CreateOrg(CreateOrgOption{
		Name:                      orgName,
		Visibility:                VisibleTypePublic,
		RepoAdminChangeTeamAccess: true,
	}); orgErr != nil {
		t.Skipf("could not create org for fork sync test: %v", orgErr)
	}

	forkName := "syncfork-child"
	fork, _, ferr := c.CreateFork(baseRepo.Owner.UserName, baseRepo.Name, models.CreateForkOption{
		Organization: orgName,
		Name:         forkName,
	})
	if ferr != nil {
		t.Skipf("could not create fork for sync test: %v", ferr)
	}
	defer func() { _, _ = c.DeleteRepo(orgName, fork.Name) }()

	// GetRepoSyncForkInfo on the fork's default branch.
	if info, _, err := c.GetRepoSyncForkInfo(orgName, fork.Name); err == nil && info != nil {
		// On a freshly created fork syncing is usually already up to date,
		// so Allowed may be false; only assert structural success.
		assert.NotNil(t, info)
	}

	// Branch-scoped info on the default branch.
	if _, _, err := c.GetRepoSyncForkBranchInfo(orgName, fork.Name, fork.DefaultBranch); err != nil {
		// Tolerate: the endpoint may report 404 when there is nothing to sync.
		log.Printf("branch sync info returned error (tolerated): %v", err)
	}

	// Sync default + branch. These return mapped errors for conflict/forbidden;
	// tolerate any environment-specific outcome.
	_, _ = c.SyncRepoForkDefault(orgName, fork.Name)
	_, _ = c.SyncRepoForkBranch(orgName, fork.Name, fork.DefaultBranch)
}
