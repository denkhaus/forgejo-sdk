// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// PinIssue pins an issue to the top of the issue list
func (c *Client) PinIssue(owner, repo string, index int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/repos/%s/%s/issues/%d/pin", owner, repo, index), nil, nil)
	return resp, err
}

// UnpinIssue unpins an issue
func (c *Client) UnpinIssue(owner, repo string, index int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/issues/%d/pin", owner, repo, index), nil, nil)
	return resp, err
}

// MoveIssuePin moves a pinned issue to a new position (1-based)
func (c *Client) MoveIssuePin(owner, repo string, index int64, position int) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("PATCH", fmt.Sprintf("/repos/%s/%s/issues/%d/pin/%d", owner, repo, index, position), nil, nil)
	return resp, err
}

// ListPinnedIssues lists the pinned issues of a repository
func (c *Client) ListPinnedIssues(owner, repo string) ([]*models.Issue, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	issues := make([]*models.Issue, 0, 10)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/pinned", owner, repo), jsonHeader, nil, &issues)
	return issues, resp, err
}

// ListPinnedPullRequests lists the pinned pull requests of a repository
func (c *Client) ListPinnedPullRequests(owner, repo string) ([]*models.PullRequest, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	prs := make([]*models.PullRequest, 0, 10)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/pulls/pinned", owner, repo), jsonHeader, nil, &prs)
	return prs, resp, err
}
