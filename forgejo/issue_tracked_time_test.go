// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_IssueTrackedTime exercises adding and listing tracked time on an issue.
func Test_IssueTrackedTime(t *testing.T) {
	log.Println("== Test_IssueTrackedTime ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "TestIssueTrackedTimeRepo", c)
	require.NoError(t, err)
	issue := createTestIssue(t, c, repo.Name, "Tracked Time Issue", "", nil, nil, 0, nil, false, false)
	owner := repo.Owner.UserName

	tt, _, err := c.AddTime(owner, repo.Name, issue.Index, AddTimeOption{Time: 60})
	require.NoError(t, err)
	if assert.NotNil(t, tt) {
		assert.EqualValues(t, 60, tt.Time)
	}

	times, _, err := c.ListIssueTrackedTimes(owner, repo.Name, issue.Index, ListTrackedTimesOptions{})
	require.NoError(t, err)
	assert.NotEmpty(t, times)

	myTimes, _, err := c.GetMyTrackedTimes()
	require.NoError(t, err)
	assert.NotNil(t, myTimes)

	// cleanup: remove the tracked time entry we just created.
	if tt != nil && tt.ID != 0 {
		_, _ = c.DeleteTime(owner, repo.Name, issue.Index, tt.ID)
	}
}
