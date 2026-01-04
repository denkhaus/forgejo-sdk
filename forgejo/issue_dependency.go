// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListDependenciesOptions options for listing issue dependencies
type ListDependenciesOptions struct {
	ListOptions
}

// ListIssueDependencies list all dependencies of a given issue
// Dependencies are issues that block the current issue from being completed
func (c *Client) ListIssueDependencies(owner, repo string, index int64, opt ListDependenciesOptions) ([]*models.Issue, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	dependencies := make([]*models.Issue, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/dependencies?%s", owner, repo, index, opt.getURLQuery().Encode()), nil, nil, &dependencies)
	return dependencies, resp, err
}

// ListBlockedIssues list all issues that block the given issue
// Blocked issues are issues that must be completed before this issue can be worked on
func (c *Client) ListBlockedIssues(owner, repo string, index int64) ([]*models.Issue, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	issues := make([]*models.Issue, 0, 5)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/blocked", owner, repo, index), nil, nil, &issues)
	return issues, resp, err
}

// ListBlockingIssues list all issues that are blocked by the given issue
// Blocking issues are issues that cannot be worked on until this issue is completed
func (c *Client) ListBlockingIssues(owner, repo string, index int64) ([]*models.Issue, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	issues := make([]*models.Issue, 0, 5)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/blocking", owner, repo, index), nil, nil, &issues)
	return issues, resp, err
}
