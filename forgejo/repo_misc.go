// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// TestRepoHook triggers a test delivery of a repository web hook
func (c *Client) TestRepoHook(owner, repo string, id int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/repos/%s/%s/hooks/%d/tests", owner, repo, id), nil, nil)
	return resp, err
}

// UpdatePullRequest merges the base branch into the pull request branch.
// style is the merge style (merge, rebase, rebase-merge, squash, commit); empty uses the default.
func (c *Client) UpdatePullRequest(owner, repo string, index int64, style string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/pulls/%d/update", owner, repo, index))
	if style != "" {
		q := link.Query()
		q.Set("style", style)
		link.RawQuery = q.Encode()
	}
	_, resp, err := c.getResponse("POST", link.String(), nil, nil)
	return resp, err
}

// NewIssuePinsAllowed reports whether new issue pins are allowed on a repository
func (c *Client) NewIssuePinsAllowed(owner, repo string) (*models.NewIssuePinsAllowed, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	out := new(models.NewIssuePinsAllowed)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/new_pin_allowed", owner, repo), jsonHeader, nil, out)
	return out, resp, err
}

// GetRepoIssueConfig returns the issue configuration of a repository
func (c *Client) GetRepoIssueConfig(owner, repo string) (*models.IssueConfig, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	cfg := new(models.IssueConfig)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issue_config", owner, repo), jsonHeader, nil, cfg)
	return cfg, resp, err
}

// ValidateRepoIssueConfig validates the issue configuration of a repository
func (c *Client) ValidateRepoIssueConfig(owner, repo string) (*models.IssueConfigValidation, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	out := new(models.IssueConfigValidation)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issue_config/validate", owner, repo), jsonHeader, nil, out)
	return out, resp, err
}
