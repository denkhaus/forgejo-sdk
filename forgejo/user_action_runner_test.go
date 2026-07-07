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

func TestSearchUserRunnerJobs(t *testing.T) {
	log.Println("== TestSearchUserRunnerJobs ==")
	c := newTestClient()

	// Search user runner jobs (may be empty for new user)
	jobs, _, err := c.SearchUserRunnerJobs(SearchUserRunnerJobsOption{})
	require.NoError(t, err)
	assert.NotNil(t, jobs)
}

// TestUserRunnersV15 exercises user-level v15 interactive runner registration+management.
func TestUserRunnersV15(t *testing.T) {
	log.Println("== TestUserRunnersV15 ==")
	c := newTestClient()

	urr, _, err := c.RegisterUserRunner(models.RegisterRunnerOptions{Name: OptionalString("user-runner"), Ephemeral: true})
	require.NoError(t, err)
	require.NotZero(t, urr.ID)

	_, _, err = c.GetUserRunner(urr.ID)
	require.NoError(t, err)
	_, _, err = c.ListUserRunners(ListActionRunnersOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	_, err = c.DeleteUserRunner(urr.ID)
	require.NoError(t, err)
}
