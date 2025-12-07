// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2018 The Gogs Authors. All rights reserved.
// Copyright 2019 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"
	"strconv"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetSingleCommit returns a single commit
func (c *Client) GetSingleCommit(user, repo, commitID string) (*models.Commit, *Response, error) {
	if err := escapeValidatePathSegments(&user, &repo, &commitID); err != nil {
		return nil, nil, err
	}
	commit := new(models.Commit)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/git/commits/%s", user, repo, commitID), nil, nil, &commit)
	return commit, resp, err
}

// ListCommitOptions list commit options
type ListCommitOptions struct {
	ListOptions
	// SHA or branch to start listing commits from (usually 'master')
	SHA string
	// Path indicates that only commits that include the path's file/dir should be returned.
	Path string
	// Stat includes diff stats for every commit (disable for speedup)
	Stat bool
	// Verification includes verification for every commit (disable for speedup)
	Verification bool
	// Files includes a list of affected files for every commit (disable for speedup)
	Files bool
	// Not is a string used such that commits that match the given specifier will not be listed.
	Not string
}

// QueryEncode turns options into querystring argument
func (opt *ListCommitOptions) QueryEncode() string {
	query := opt.getURLQuery()
	if opt.SHA != "" {
		query.Add("sha", opt.SHA)
	}
	if opt.Path != "" {
		query.Add("path", opt.Path)
	}
	query.Add("stat", strconv.FormatBool(opt.Stat))
	query.Add("verification", strconv.FormatBool(opt.Stat))
	query.Add("files", strconv.FormatBool(opt.Stat))
	if opt.Not != "" {
		query.Add("not", opt.Not)
	}
	return query.Encode()
}

// ListRepoCommits return list of commits from a repo
func (c *Client) ListRepoCommits(user, repo string, opt ListCommitOptions) ([]*models.Commit, *Response, error) {
	if err := escapeValidatePathSegments(&user, &repo); err != nil {
		return nil, nil, err
	}
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/commits", user, repo))
	opt.setDefaults()
	commits := make([]*models.Commit, 0, opt.PageSize)
	link.RawQuery = opt.QueryEncode()
	resp, err := c.getParsedResponse("GET", link.String(), nil, nil, &commits)
	return commits, resp, err
}

// GetCommitDiff returns the commit's raw diff.
func (c *Client) GetCommitDiff(user, repo, commitID string) ([]byte, *Response, error) {
	if err := c.checkServerVersionGreaterThanOrEqual(version1_16_0); err != nil {
		return nil, nil, err
	}

	if err := escapeValidatePathSegments(&user, &repo); err != nil {
		return nil, nil, err
	}

	return c.getResponse("GET", fmt.Sprintf("/repos/%s/%s/git/commits/%s.%s", user, repo, commitID, pullRequestDiffTypeDiff), nil, nil)
}

// GetCommitPatch returns the commit's raw patch.
func (c *Client) GetCommitPatch(user, repo, commitID string) ([]byte, *Response, error) {
	if err := c.checkServerVersionGreaterThanOrEqual(version1_16_0); err != nil {
		return nil, nil, err
	}

	if err := escapeValidatePathSegments(&user, &repo); err != nil {
		return nil, nil, err
	}

	return c.getResponse("GET", fmt.Sprintf("/repos/%s/%s/git/commits/%s.%s", user, repo, commitID, pullRequestDiffTypePatch), nil, nil)
}
