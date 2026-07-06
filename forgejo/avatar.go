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

// UpdateUserAvatar updates the avatar of the current authenticated user
func (c *Client) UpdateUserAvatar(opt models.UpdateUserAvatarOption) (*Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", "/user/avatar", jsonHeader, bytes.NewReader(body))
	return resp, err
}

// DeleteUserAvatar deletes the avatar of the current authenticated user
func (c *Client) DeleteUserAvatar() (*Response, error) {
	_, resp, err := c.getResponse("DELETE", "/user/avatar", nil, nil)
	return resp, err
}

// UpdateOrgAvatar updates the avatar of an organization
func (c *Client) UpdateOrgAvatar(org string, opt models.UpdateUserAvatarOption) (*Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/orgs/%s/avatar", org), jsonHeader, bytes.NewReader(body))
	return resp, err
}

// DeleteOrgAvatar deletes the avatar of an organization
func (c *Client) DeleteOrgAvatar(org string) (*Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/orgs/%s/avatar", org), nil, nil)
	return resp, err
}

// UpdateRepoAvatar updates the avatar of a repository
func (c *Client) UpdateRepoAvatar(owner, repo string, opt models.UpdateRepoAvatarOption) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/repos/%s/%s/avatar", owner, repo), jsonHeader, bytes.NewReader(body))
	return resp, err
}

// DeleteRepoAvatar deletes the avatar of a repository
func (c *Client) DeleteRepoAvatar(owner, repo string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/avatar", owner, repo), nil, nil)
	return resp, err
}
