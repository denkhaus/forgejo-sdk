// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListRepoTagProtections lists a repository's tag protections
func (c *Client) ListRepoTagProtections(owner, repo string) ([]*models.TagProtection, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	protections := make([]*models.TagProtection, 0, 10)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/tag_protections", owner, repo), jsonHeader, nil, &protections)
	return protections, resp, err
}

// GetRepoTagProtection gets a repository's tag protection by ID
func (c *Client) GetRepoTagProtection(owner, repo string, id int64) (*models.TagProtection, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	p := new(models.TagProtection)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/tag_protections/%d", owner, repo, id), jsonHeader, nil, p)
	return p, resp, err
}

// CreateRepoTagProtection creates a tag protection for a repository
func (c *Client) CreateRepoTagProtection(owner, repo string, opt models.CreateTagProtectionOption) (*models.TagProtection, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	p := new(models.TagProtection)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/tag_protections", owner, repo), jsonHeader, bytes.NewReader(body), p)
	return p, resp, err
}

// EditRepoTagProtection edits a repository's tag protection by ID
func (c *Client) EditRepoTagProtection(owner, repo string, id int64, opt models.EditTagProtectionOption) (*models.TagProtection, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	p := new(models.TagProtection)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/tag_protections/%d", owner, repo, id), jsonHeader, bytes.NewReader(body), p)
	return p, resp, err
}

// DeleteRepoTagProtection deletes a repository's tag protection by ID
func (c *Client) DeleteRepoTagProtection(owner, repo string, id int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/repos/%s/%s/tag_protections/%d", owner, repo, id), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("tag protection not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
