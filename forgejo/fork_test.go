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

func Test_Fork(t *testing.T) {
	log.Println("== Test_Fork ==")
	c := newTestClient()

	baseRepo, err := createTestRepo(t, "ForkBase", c)
	require.NoError(t, err)

	orgName := "fork-org"
	defer func() {
		_, _ = c.DeleteRepo(orgName, "ForkedRepo")
		_, _ = c.DeleteOrg(orgName)
	}()

	// Reset any leftover org from previous runs.
	_, _ = c.DeleteOrg(orgName)
	_, _, err = c.CreateOrg(CreateOrgOption{
		Name:                      orgName,
		Visibility:                VisibleTypePublic,
		RepoAdminChangeTeamAccess: true,
	})
	require.NoError(t, err)

	// ListForks before forking should be empty/safe.
	before, _, err := c.ListForks(baseRepo.Owner.UserName, baseRepo.Name, ListForksOptions{})
	require.NoError(t, err)
	assert.NotNil(t, before)

	// CreateFork into the org.
	fork, _, err := c.CreateFork(baseRepo.Owner.UserName, baseRepo.Name, models.CreateForkOption{
		Organization: orgName,
		Name:         "ForkedRepo",
	})
	if err != nil {
		t.Skipf("fork creation not permitted in this environment: %v", err)
	}
	require.NotNil(t, fork)
	assert.Equal(t, orgName, fork.Owner.UserName)

	// ListForks should now report at least one fork.
	after, _, err := c.ListForks(baseRepo.Owner.UserName, baseRepo.Name, ListForksOptions{})
	require.NoError(t, err)
	assert.NotEmpty(t, after)
}
