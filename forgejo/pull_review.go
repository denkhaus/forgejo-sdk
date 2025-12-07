// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ReviewStateType review state type
type ReviewStateType string

const (
	// ReviewStateApproved pr is approved
	ReviewStateApproved ReviewStateType = "APPROVED"
	// ReviewStatePending pr state is pending
	ReviewStatePending ReviewStateType = "PENDING"
	// ReviewStateComment is a comment review
	ReviewStateComment ReviewStateType = "COMMENT"
	// ReviewStateRequestChanges changes for pr are requested
	ReviewStateRequestChanges ReviewStateType = "REQUEST_CHANGES"
	// ReviewStateRequestReview review is requested from user
	ReviewStateRequestReview ReviewStateType = "REQUEST_REVIEW"
	// ReviewStateUnknown state of pr is unknown
	ReviewStateUnknown ReviewStateType = ""
)

// CreatePullReviewOptions are options to create a pull review
type CreatePullReviewOptions struct {
	State    ReviewStateType           `json:"event"`
	Body     string                    `json:"body"`
	CommitID string                    `json:"commit_id"`
	Comments []CreatePullReviewComment `json:"comments"`
}

// CreatePullReviewComment represent a review comment for creation api
type CreatePullReviewComment struct {
	// the tree path
	Path string `json:"path"`
	Body string `json:"body"`
	// if comment to old file line or 0
	OldLineNum int64 `json:"old_position"`
	// if comment to new file line or 0
	NewLineNum int64 `json:"new_position"`
}

// SubmitPullReviewOptions are options to submit a pending pull review
type SubmitPullReviewOptions struct {
	State ReviewStateType `json:"event"`
	Body  string          `json:"body"`
}

// ListPullReviewsOptions options for listing PullReviews
type ListPullReviewsOptions struct {
	ListOptions
}

// Validate the CreatePullReviewOptions struct
func (opt CreatePullReviewOptions) Validate() error {
	if opt.State != ReviewStateApproved && len(opt.Comments) == 0 && len(strings.TrimSpace(opt.Body)) == 0 {
		return fmt.Errorf("body is empty")
	}
	for i := range opt.Comments {
		if err := opt.Comments[i].Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate the SubmitPullReviewOptions struct
func (opt SubmitPullReviewOptions) Validate() error {
	if opt.State != ReviewStateApproved && len(strings.TrimSpace(opt.Body)) == 0 {
		return fmt.Errorf("body is empty")
	}
	return nil
}

// Validate the CreatePullReviewComment struct
func (opt CreatePullReviewComment) Validate() error {
	if len(strings.TrimSpace(opt.Body)) == 0 {
		return fmt.Errorf("body is empty")
	}
	if opt.NewLineNum != 0 && opt.OldLineNum != 0 {
		return fmt.Errorf("old and new line num are set, cant identify the code comment position")
	}
	return nil
}

// ListPullReviews lists all reviews of a pull request
func (c *Client) ListPullReviews(owner, repo string, index int64, opt ListPullReviewsOptions) ([]*models.PullReview, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	rs := make([]*models.PullReview, 0, opt.PageSize)

	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews", owner, repo, index))
	link.RawQuery = opt.ListOptions.getURLQuery().Encode()

	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &rs)
	return rs, resp, err
}

// GetPullReview gets a specific review of a pull request
func (c *Client) GetPullReview(owner, repo string, index, id int64) (*models.PullReview, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}

	r := new(models.PullReview)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d", owner, repo, index, id), jsonHeader, nil, &r)
	return r, resp, err
}

// ListPullReviewComments lists all comments of a pull request review
func (c *Client) ListPullReviewComments(owner, repo string, index, id int64) ([]*models.PullReviewComment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}
	rcl := make([]*models.PullReviewComment, 0, 4)
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d/comments", owner, repo, index, id))

	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &rcl)
	return rcl, resp, err
}

// DeletePullReview delete a specific review from a pull request
func (c *Client) DeletePullReview(owner, repo string, index, id int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, err
	}

	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d", owner, repo, index, id), jsonHeader, nil)
	return resp, err
}

// CreatePullReview create a review to an pull request
func (c *Client) CreatePullReview(owner, repo string, index int64, opt CreatePullReviewOptions) (*models.PullReview, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}
	if err := opt.Validate(); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}

	r := new(models.PullReview)
	resp, err := c.getParsedResponse("POST",
		fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews", owner, repo, index),
		jsonHeader, bytes.NewReader(body), r)
	return r, resp, err
}

// SubmitPullReview submit a pending review to an pull request
func (c *Client) SubmitPullReview(owner, repo string, index, id int64, opt SubmitPullReviewOptions) (*models.PullReview, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}
	if err := opt.Validate(); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}

	r := new(models.PullReview)
	resp, err := c.getParsedResponse("POST",
		fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d", owner, repo, index, id),
		jsonHeader, bytes.NewReader(body), r)
	return r, resp, err
}

// CreateReviewRequests create review requests to an pull request
func (c *Client) CreateReviewRequests(owner, repo string, index int64, opt models.PullReviewRequestOptions) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_14_0); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	_, resp, err := c.getResponse("POST",
		fmt.Sprintf("/repos/%s/%s/pulls/%d/requested_reviewers", owner, repo, index),
		jsonHeader, bytes.NewReader(body))
	return resp, err
}

// DeleteReviewRequests delete review requests to an pull request
func (c *Client) DeleteReviewRequests(owner, repo string, index int64, opt models.PullReviewRequestOptions) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_14_0); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	_, resp, err := c.getResponse("DELETE",
		fmt.Sprintf("/repos/%s/%s/pulls/%d/requested_reviewers", owner, repo, index),
		jsonHeader, bytes.NewReader(body))
	return resp, err
}

// DismissPullReview dismiss a review for a pull request
func (c *Client) DismissPullReview(owner, repo string, index, id int64, opt models.DismissPullReviewOptions) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_14_0); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	_, resp, err := c.getResponse("POST",
		fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d/dismissals", owner, repo, index, id),
		jsonHeader, bytes.NewReader(body))
	return resp, err
}

// UnDismissPullReview cancel to dismiss a review for a pull request
func (c *Client) UnDismissPullReview(owner, repo string, index, id int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_14_0); err != nil {
		return nil, err
	}

	_, resp, err := c.getResponse("POST",
		fmt.Sprintf("/repos/%s/%s/pulls/%d/reviews/%d/undismissals", owner, repo, index, id),
		jsonHeader, nil)
	return resp, err
}
