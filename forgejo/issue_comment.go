// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListIssueCommentOptions list comment options
type ListIssueCommentOptions struct {
	ListOptions
	Since  time.Time
	Before time.Time
}

// QueryEncode turns options into querystring argument
func (opt *ListIssueCommentOptions) QueryEncode() string {
	query := opt.getURLQuery()
	if !opt.Since.IsZero() {
		query.Add("since", opt.Since.Format(time.RFC3339))
	}
	if !opt.Before.IsZero() {
		query.Add("before", opt.Before.Format(time.RFC3339))
	}
	return query.Encode()
}

// ListIssueComments list comments on an issue.
func (c *Client) ListIssueComments(owner, repo string, index int64, opt ListIssueCommentOptions) ([]*models.Comment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, index))
	link.RawQuery = opt.QueryEncode()
	comments := make([]*models.Comment, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), nil, nil, &comments)
	return comments, resp, err
}

// ListRepoIssueComments list comments for a given repo.
func (c *Client) ListRepoIssueComments(owner, repo string, opt ListIssueCommentOptions) ([]*models.Comment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/issues/comments", owner, repo))
	link.RawQuery = opt.QueryEncode()
	comments := make([]*models.Comment, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), nil, nil, &comments)
	return comments, resp, err
}

// GetIssueComment get a comment for a given repo by id.
func (c *Client) GetIssueComment(owner, repo string, id int64) (*models.Comment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	comment := new(models.Comment)
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return comment, nil, err
	}
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/comments/%d", owner, repo, id), nil, nil, &comment)
	return comment, resp, err
}

// CreateIssueCommentOption options for creating a comment on an issue
type CreateIssueCommentOption struct {
	Body string `json:"body"`
}

// Validate the CreateIssueCommentOption struct
func (opt CreateIssueCommentOption) Validate() error {
	if len(opt.Body) == 0 {
		return fmt.Errorf("body is empty")
	}
	return nil
}

// CreateIssueComment create comment on an issue.
func (c *Client) CreateIssueComment(owner, repo string, index int64, opt CreateIssueCommentOption) (*models.Comment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := opt.Validate(); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	comment := new(models.Comment)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/issues/%d/comments", owner, repo, index), jsonHeader, bytes.NewReader(body), comment)
	return comment, resp, err
}

// EditIssueCommentOption options for editing a comment
type EditIssueCommentOption struct {
	Body string `json:"body"`
}

// Validate the EditIssueCommentOption struct
func (opt EditIssueCommentOption) Validate() error {
	if len(opt.Body) == 0 {
		return fmt.Errorf("body is empty")
	}
	return nil
}

// EditIssueComment edits an issue comment.
func (c *Client) EditIssueComment(owner, repo string, commentID int64, opt EditIssueCommentOption) (*models.Comment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := opt.Validate(); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	comment := new(models.Comment)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/issues/comments/%d", owner, repo, commentID), jsonHeader, bytes.NewReader(body), comment)
	return comment, resp, err
}

// DeleteIssueComment deletes an issue comment.
func (c *Client) DeleteIssueComment(owner, repo string, commentID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/issues/comments/%d", owner, repo, commentID), nil, nil)
	return resp, err
}
