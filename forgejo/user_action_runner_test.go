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

func TestSearchUserRunnerJobs(t *testing.T) {
	log.Println("== TestSearchUserRunnerJobs ==")
	c := newTestClient()

	// Search user runner jobs (may be empty for new user)
	jobs, _, err := c.SearchUserRunnerJobs(SearchUserRunnerJobsOption{})
	require.NoError(t, err)
	assert.NotNil(t, jobs)
}
