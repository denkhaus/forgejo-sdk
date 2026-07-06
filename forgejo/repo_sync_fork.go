// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/http"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetRepoSyncForkInfo reports whether a fork's default branch can be synced with its base
func (c *Client) GetRepoSyncForkInfo(owner, repo string) (*models.SyncForkInfo, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	info := new(models.SyncForkInfo)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/sync_fork", owner, repo), jsonHeader, nil, info)
	return info, resp, err
}

// GetRepoSyncForkBranchInfo reports whether a fork's given branch can be synced with its base
func (c *Client) GetRepoSyncForkBranchInfo(owner, repo, branch string) (*models.SyncForkInfo, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &branch); err != nil {
		return nil, nil, err
	}
	info := new(models.SyncForkInfo)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/sync_fork/%s", owner, repo, branch), jsonHeader, nil, info)
	return info, resp, err
}

// SyncRepoForkDefault syncs a fork's default branch with its base
func (c *Client) SyncRepoForkDefault(owner, repo string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	return c.syncRepoFork("POST", fmt.Sprintf("/repos/%s/%s/sync_fork", owner, repo))
}

// SyncRepoForkBranch syncs a fork's given branch with its base
func (c *Client) SyncRepoForkBranch(owner, repo, branch string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &branch); err != nil {
		return nil, err
	}
	return c.syncRepoFork("POST", fmt.Sprintf("/repos/%s/%s/sync_fork/%s", owner, repo, branch))
}

func (c *Client) syncRepoFork(method, path string) (*Response, error) {
	status, resp, err := c.getStatusCode(method, path, jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusOK, http.StatusAccepted, http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("repository or branch not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusConflict:
		return resp, fmt.Errorf("conflict: merge conflict or non-fast-forward")
	case http.StatusServiceUnavailable:
		return resp, fmt.Errorf("service unavailable: sync could not be performed")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
