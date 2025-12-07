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

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListBranchProtectionsOptions list branch protection options
type ListBranchProtectionsOptions struct {
	ListOptions
}

// ListBranchProtections list branch protections for a repo
func (c *Client) ListBranchProtections(owner, repo string, opt ListBranchProtectionsOptions) ([]*models.BranchProtection, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}
	bps := make([]*models.BranchProtection, 0, opt.PageSize)
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/branch_protections", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &bps)
	return bps, resp, err
}

// GetBranchProtection gets a branch protection
func (c *Client) GetBranchProtection(owner, repo, name string) (*models.BranchProtection, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}
	bp := new(models.BranchProtection)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/branch_protections/%s", owner, repo, name), jsonHeader, nil, bp)
	return bp, resp, err
}

// CreateBranchProtection creates a branch protection for a repo
func (c *Client) CreateBranchProtection(owner, repo string, opt models.CreateBranchProtectionOption) (*models.BranchProtection, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}
	bp := new(models.BranchProtection)
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/branch_protections", owner, repo), jsonHeader, bytes.NewReader(body), bp)
	return bp, resp, err
}

// EditBranchProtection edits a branch protection for a repo
func (c *Client) EditBranchProtection(owner, repo, name string, opt models.EditBranchProtectionOption) (*models.BranchProtection, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, nil, err
	}
	bp := new(models.BranchProtection)
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/branch_protections/%s", owner, repo, name), jsonHeader, bytes.NewReader(body), bp)
	return bp, resp, err
}

// DeleteBranchProtection deletes a branch protection for a repo
func (c *Client) DeleteBranchProtection(owner, repo, name string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, err
	}
	if err := c.checkServerVersionGreaterThanOrEqual(version1_12_0); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/branch_protections/%s", owner, repo, name), jsonHeader, nil)
	return resp, err
}
