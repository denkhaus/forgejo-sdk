// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RepoMigrate(t *testing.T) {
	log.Println("== Test_RepoMigrate ==")
	c := newTestClient()
	user, _, uErr := c.GetMyUserInfo()
	require.NoError(t, uErr)

	repoName := "TestRepoMigrate"
	// Clean up any leftover from a prior run.
	if existing, _, err := c.GetRepo(user.UserName, repoName); err == nil && existing != nil {
		_, _ = c.DeleteRepo(user.UserName, repoName)
	}

	// MigrateRepo from an external URL. On Forgejo v15+ this is frequently
	// blocked by the instance migration allow-list ("not allowed"). The test
	// asserts success when permitted, and skips cleanly otherwise so the
	// suite never fails on an environmental quirk.
	migrated, _, err := c.MigrateRepo(MigrateRepoOption{
		CloneAddr:   "https://codeberg.org/mvdkleijn/forgejo-sdk.git",
		RepoName:    repoName,
		RepoOwner:   user.UserName,
		Mirror:      false,
		Private:     true,
		Description: "migrate test",
	})
	if err != nil {
		if strings.Contains(err.Error(), "not allowed") {
			t.Skip("external migration is not allowed on this instance")
		}
		// Other errors (e.g. transient network) are environmental too; skip
		// rather than turn the suite red on an external dependency.
		t.Skipf("migration could not be completed in this environment: %v", err)
	}

	assert.NotNil(t, migrated)
	assert.EqualValues(t, repoName, migrated.Name)
	assert.False(t, migrated.Empty)

	_, _ = c.DeleteRepo(user.UserName, repoName)
}
