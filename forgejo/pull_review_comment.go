// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetPullReviewComment gets a single comment of a pull request review
func (c *Client) GetPullReviewComment(owner, repo string, index, reviewID, commentID int64) (*models.PullReviewComment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	comment := new(models.PullReviewComment)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d/comments/%d", owner, repo, index, reviewID, commentID), jsonHeader, nil, comment)
	return comment, resp, err
}

// CreatePullReviewComment creates a comment on a pull request review
func (c *Client) CreatePullReviewComment(owner, repo string, index, reviewID int64, opt models.CreatePullReviewCommentOptions) (*models.PullReviewComment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	comment := new(models.PullReviewComment)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d/comments", owner, repo, index, reviewID), jsonHeader, bytes.NewReader(body), comment)
	return comment, resp, err
}

// DeletePullReviewComment deletes a comment of a pull request review
func (c *Client) DeletePullReviewComment(owner, repo string, index, reviewID, commentID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d/comments/%d", owner, repo, index, reviewID, commentID), nil, nil)
	return resp, err
}
