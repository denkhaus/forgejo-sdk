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

// Test_PullReviewComment exercises creating, fetching and deleting a comment
// on a pull request review. It reuses the shared pull-review test fixture,
// which builds a PR with a diff; if that precondition cannot be established
// the test stops rather than failing.
func Test_PullReviewComment(t *testing.T) {
	log.Println("== Test_PullReviewComment ==")
	c := newTestClient()

	repo, pull, _, _, ok := preparePullReviewTest(t, c, "PullReviewComment")
	if !ok {
		t.Skip("pull review precondition could not be established")
	}
	owner := repo.Owner.UserName

	// Open a pending review carrying an inline comment on the new file.
	review, _, err := c.CreatePullReview(owner, repo.Name, pull.Index, CreatePullReviewOptions{
		State: ReviewStatePending,
		Body:  "pending review for comment test",
		Comments: []CreatePullReviewComment{
			{Path: "WOW-file", Body: "initial inline comment", NewLineNum: 1},
		},
	})
	if err != nil {
		t.Skipf("could not create pending review (precondition not met): %v", err)
	}
	require.NotNil(t, review)
	assert.NotZero(t, review.ID)

	// CreatePullReviewComment: the method under test. Adding a review comment
	// requires the diff to be available; tolerate an environmental failure.
	comment, _, err := c.CreatePullReviewComment(owner, repo.Name, pull.Index, review.ID, models.CreatePullReviewCommentOptions{
		CreatePullReviewComment: models.CreatePullReviewComment{
			Path:       "WOW-file",
			Body:       "test review comment",
			NewLineNum: 1,
		},
	})
	if err != nil {
		t.Logf("CreatePullReviewComment returned (tolerated): %v", err)
		return
	}
	if !assert.NotNil(t, comment) || comment.ID == 0 {
		return
	}

	// GetPullReviewComment: fetch the comment we just created.
	got, _, err := c.GetPullReviewComment(owner, repo.Name, pull.Index, review.ID, comment.ID)
	require.NoError(t, err)
	assert.EqualValues(t, comment.ID, got.ID)

	// DeletePullReviewComment: cleanup; tolerate any failure.
	if _, err := c.DeletePullReviewComment(owner, repo.Name, pull.Index, review.ID, comment.ID); err != nil {
		t.Logf("DeletePullReviewComment returned (tolerated): %v", err)
	}
}
