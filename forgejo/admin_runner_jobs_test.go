// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAdminRunnerJobsV15 exercises admin action run job listing and search.
func TestAdminRunnerJobsV15(t *testing.T) {
	log.Println("== TestAdminRunnerJobsV15 ==")
	c := newTestClient()

	_, _, err := c.AdminGetActionRunJobs(AdminSearchActionRunJobsOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)

	_, _, err = c.AdminSearchRunJobs(AdminSearchActionRunJobsOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
}
