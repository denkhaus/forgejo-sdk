// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIssueTimelineV15 exercises fetching the timeline of an issue.
func TestIssueTimelineV15(t *testing.T) {
	log.Println("== TestIssueTimelineV15 ==")
	c := newTestClient()

	repo, err := createTestRepo(t, "v15-timeline-repo", c)
	require.NoError(t, err)

	issue := createTestIssue(t, c, repo.Name, "timeline-issue", "", nil, nil, 0, nil, false, false)
	require.NotNil(t, issue)

	_, _, err = c.GetIssueTimeline(repo.Owner.UserName, repo.Name, issue.Index, ListIssueTimelineOptions{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
}
