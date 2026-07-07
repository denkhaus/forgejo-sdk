// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/require"
)

// TestIssueDeadlineV15 exercises setting an issue deadline.
func TestIssueDeadlineV15(t *testing.T) {
	log.Println("== TestIssueDeadlineV15 ==")
	c := newTestClient()

	repo, err := createTestRepo(t, "v15-deadline-repo", c)
	require.NoError(t, err)

	issue := createTestIssue(t, c, repo.Name, "deadline-issue", "", nil, nil, 0, nil, false, false)
	require.NotNil(t, issue)

	due := strfmt.DateTime(time.Now().Add(48 * time.Hour))
	_, _, err = c.EditIssueDeadline(repo.Owner.UserName, repo.Name, issue.Index, models.EditDeadlineOption{Deadline: &due})
	require.NoError(t, err)
}
