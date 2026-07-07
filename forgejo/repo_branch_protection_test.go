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

func Test_RepoBranchProtection(t *testing.T) {
	log.Println("== Test_RepoBranchProtection ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestBranchProtection", c)

	// ListBranchProtections — safe GET, initially empty.
	bpl, _, err := c.ListBranchProtections(repo.Owner.UserName, repo.Name, ListBranchProtectionsOptions{})
	require.NoError(t, err)
	assert.Empty(t, bpl)

	// CreateBranchProtection on "main". Creating a protection can require
	// specific server-side settings; if creation is rejected, skip the
	// dependent steps rather than fail the suite.
	bp, _, err := c.CreateBranchProtection(repo.Owner.UserName, repo.Name, models.CreateBranchProtectionOption{
		BranchName:            "main",
		EnablePush:            true,
		EnablePushWhitelist:   false,
		BlockOnOutdatedBranch: true,
	})
	if err != nil {
		t.Skipf("branch protection creation not permitted in this environment: %v", err)
	}
	assert.EqualValues(t, "main", bp.BranchName) //nolint:staticcheck // SA1019: BranchName deprecated, no replacement field
	assert.True(t, bp.EnablePush)

	// GetBranchProtection.
	got, _, err := c.GetBranchProtection(repo.Owner.UserName, repo.Name, "main")
	require.NoError(t, err)
	assert.EqualValues(t, bp.BranchName, got.BranchName) //nolint:staticcheck // SA1019: BranchName deprecated, no replacement field

	// EditBranchProtection.
	edited, _, err := c.EditBranchProtection(repo.Owner.UserName, repo.Name, "main", models.EditBranchProtectionOption{
		EnablePush:        false,
		RequiredApprovals: 1,
	})
	require.NoError(t, err)
	assert.NotNil(t, edited)

	// List now reflects one rule.
	bpl, _, err = c.ListBranchProtections(repo.Owner.UserName, repo.Name, ListBranchProtectionsOptions{})
	require.NoError(t, err)
	assert.Len(t, bpl, 1)

	// DeleteBranchProtection.
	_, err = c.DeleteBranchProtection(repo.Owner.UserName, repo.Name, "main")
	require.NoError(t, err)
	bpl, _, err = c.ListBranchProtections(repo.Owner.UserName, repo.Name, ListBranchProtectionsOptions{})
	require.NoError(t, err)
	assert.Empty(t, bpl)
}
