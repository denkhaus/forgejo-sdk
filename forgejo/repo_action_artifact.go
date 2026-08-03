// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListActionArtifactsOption options for listing action artifacts.
type ListActionArtifactsOption struct {
	ListOptions
	// Name filters the result by artifact name (case-sensitive).
	Name string
}

// getURLQuery extends the pagination query with the optional name filter.
func (opt ListActionArtifactsOption) getURLQuery() url.Values {
	query := opt.ListOptions.getURLQuery()
	if opt.Name != "" {
		query.Add("name", opt.Name)
	}
	return query
}

// ListActionArtifacts lists a repository's artifacts.
func (c *Client) ListActionArtifacts(owner, repo string, opt ListActionArtifactsOption) ([]*models.ActionArtifact, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/artifacts", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	artifacts := make([]*models.ActionArtifact, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &artifacts)
	return artifacts, resp, err
}

// ListActionRunArtifacts lists the artifacts of a single workflow run.
func (c *Client) ListActionRunArtifacts(owner, repo string, runID int64, opt ListActionArtifactsOption) ([]*models.ActionArtifact, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/runs/%d/artifacts", owner, repo, runID))
	link.RawQuery = opt.getURLQuery().Encode()
	artifacts := make([]*models.ActionArtifact, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &artifacts)
	return artifacts, resp, err
}

// GetActionArtifact gets a repository's artifact by ID.
func (c *Client) GetActionArtifact(owner, repo string, artifactID int64) (*models.ActionArtifact, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	artifact := new(models.ActionArtifact)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/actions/artifacts/%d", owner, repo, artifactID), jsonHeader, nil, artifact)
	return artifact, resp, err
}

// DownloadActionArtifact downloads an artifact's ZIP archive. The returned
// io.ReadCloser streams the raw archive bytes and must be closed by the caller.
func (c *Client) DownloadActionArtifact(owner, repo string, artifactID int64) (io.ReadCloser, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	return c.getResponseReader(fmt.Sprintf("/repos/%s/%s/actions/artifacts/%d/zip", owner, repo, artifactID), jsonHeader)
}

// DeleteActionArtifact marks an artifact for deletion.
func (c *Client) DeleteActionArtifact(owner, repo string, artifactID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/repos/%s/%s/actions/artifacts/%d", owner, repo, artifactID), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("repository or artifact not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
