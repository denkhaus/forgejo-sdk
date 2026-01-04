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

// ListBlockedIssues list all issues that are blocked by this issue
// These issues cannot be worked on until this issue is completed
func (c *Client) ListBlockedIssues(owner, repo string, index int64) ([]*models.Issue, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	issues := make([]*models.Issue, 0, 5)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/blocks", owner, repo, index), nil, nil, &issues)
	return issues, resp, err
}

// ListBlockingIssues list all issues that block this issue
// These issues must be completed before this issue can be worked on
// This is an alias for ListIssueDependencies with default options
func (c *Client) ListBlockingIssues(owner, repo string, index int64) ([]*models.Issue, *Response, error) {
	return c.ListIssueDependencies(owner, repo, index, ListDependenciesOptions{})
}

// CreateIssueDependencyOption options for creating an issue dependency
type CreateIssueDependencyOption struct {
	// NewDependency is the issue index (number) of the blocking issue
	NewDependency int64 `json:"newDependency"`
}

// Validate the CreateIssueDependencyOption struct
func (opt CreateIssueDependencyOption) Validate() error {
	if opt.NewDependency <= 0 {
		return fmt.Errorf("newDependency must be a positive issue number")
	}
	return nil
}

// CreateIssueDependency adds a dependency relationship to an issue
// The dependency issue (specified by NewDependency) must be completed
// before this issue can be worked on
func (c *Client) CreateIssueDependency(owner, repo string, index int64, opt CreateIssueDependencyOption) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	if err := opt.Validate(); err != nil {
		return nil, err
	}
	// API expects IssueMeta with index, owner, and repo fields
	meta := models.IssueMeta{
		Index: opt.NewDependency,
		Owner: owner,
		Name:  repo,
	}
	body, err := json.Marshal(&meta)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST",
		fmt.Sprintf("/repos/%s/%s/issues/%d/dependencies", owner, repo, index),
		jsonHeader, bytes.NewReader(body))
	return resp, err
}

// RemoveIssueDependency removes a dependency relationship from an issue
// The dependency issue is specified by its index (number)
func (c *Client) RemoveIssueDependency(owner, repo string, index, dependency int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	// API expects IssueMeta with index, owner, and repo fields in body
	meta := models.IssueMeta{
		Index: dependency,
		Owner: owner,
		Name:  repo,
	}
	body, err := json.Marshal(&meta)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE",
		fmt.Sprintf("/repos/%s/%s/issues/%d/dependencies", owner, repo, index),
		jsonHeader, bytes.NewReader(body))
	return resp, err
}
