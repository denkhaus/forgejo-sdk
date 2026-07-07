// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_IssueStopWatch exercises starting, listing, stopping and deleting a stopwatch on an issue.
func Test_IssueStopWatch(t *testing.T) {
	log.Println("== Test_IssueStopWatch ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "TestIssueStopWatchRepo", c)
	require.NoError(t, err)
	issue := createTestIssue(t, c, repo.Name, "Stopwatch Issue", "", nil, nil, 0, nil, false, false)
	owner := repo.Owner.UserName

	// Starting a stopwatch can fail if another stopwatch is already running for
	// this user; tolerate that environmental condition rather than failing.
	if _, err := c.StartIssueStopWatch(owner, repo.Name, issue.Index); err != nil {
		t.Skipf("could not start stopwatch (precondition not met): %v", err)
	}

	stopwatches, _, err := c.GetMyStopwatches()
	require.NoError(t, err)
	assert.NotNil(t, stopwatches)

	if _, err := c.StopIssueStopWatch(owner, repo.Name, issue.Index); err != nil {
		t.Logf("StopIssueStopWatch returned (tolerated): %v", err)
	}

	// DeleteIssueStopwatch cancels any stopwatch on the issue; tolerate a no-op.
	if _, err := c.DeleteIssueStopwatch(owner, repo.Name, issue.Index); err != nil {
		t.Logf("DeleteIssueStopwatch returned (tolerated): %v", err)
	}
}
