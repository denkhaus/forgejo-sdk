// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// DeleteActionRun deletes a completed workflow run.
func (c *Client) DeleteActionRun(owner, repo string, runID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/repos/%s/%s/actions/runs/%d", owner, repo, runID), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("repository or run not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusBadRequest:
		return resp, fmt.Errorf("bad request: run may not be completed")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// CancelActionRun cancels a pending or running workflow run.
func (c *Client) CancelActionRun(owner, repo string, runID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("POST", fmt.Sprintf("/repos/%s/%s/actions/runs/%d/cancel", owner, repo, runID), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK, http.StatusAccepted:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("repository or run not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// ListActionRunJobs lists the jobs of a workflow run.
func (c *Client) ListActionRunJobs(owner, repo string, runID int64) ([]*models.ActionRunJob, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	jobs := make([]*models.ActionRunJob, 0)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/actions/runs/%d/jobs", owner, repo, runID), jsonHeader, nil, &jobs)
	return jobs, resp, err
}

// GetActionRunLogs downloads a ZIP archive containing the plaintext logs of
// every job in a workflow run. The returned io.ReadCloser streams the raw
// archive bytes and must be closed by the caller.
func (c *Client) GetActionRunLogs(owner, repo string, runID int64) (io.ReadCloser, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	return c.getResponseReader(fmt.Sprintf("/repos/%s/%s/actions/runs/%d/logs", owner, repo, runID), jsonHeader)
}

// GetActionJobLogsOption options for downloading an action job's logs.
type GetActionJobLogsOption struct {
	// Attempt selects a specific (1-based) run attempt. When zero (the
	// default) the logs of the latest attempt are returned.
	Attempt int64
}

// GetActionJobLogs downloads the plaintext logs of a single action job. The
// returned io.ReadCloser streams the raw log bytes and must be closed by the
// caller.
func (c *Client) GetActionJobLogs(owner, repo string, jobID int64, opt GetActionJobLogsOption) (io.ReadCloser, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/jobs/%d/logs", owner, repo, jobID))
	if opt.Attempt > 0 {
		q := link.Query()
		q.Set("attempt", strconv.FormatInt(opt.Attempt, 10))
		link.RawQuery = q.Encode()
	}
	return c.getResponseReader(link.String(), jsonHeader)
}
