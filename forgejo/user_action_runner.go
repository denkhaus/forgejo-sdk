// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// SearchUserRunnerJobsOption options for searching user runner jobs
type SearchUserRunnerJobsOption struct {
	// Filter by job labels (comma-separated list)
	Labels []string
}

// SearchUserRunnerJobs searches for the current user's action jobs according to filter conditions
func (c *Client) SearchUserRunnerJobs(opt SearchUserRunnerJobsOption) ([]*models.ActionRunJob, *Response, error) {
	link, _ := url.Parse("/user/actions/runners/jobs")
	query := make(url.Values)

	if len(opt.Labels) > 0 {
		// Labels is a comma-separated list
		labelsStr := ""
		for i, label := range opt.Labels {
			if i > 0 {
				labelsStr += ","
			}
			labelsStr += label
		}
		query.Add("labels", labelsStr)
	}

	link.RawQuery = query.Encode()
	jobs := make([]*models.ActionRunJob, 0, 10)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &jobs)
	if jobs == nil {
		jobs = make([]*models.ActionRunJob, 0, 10)
	}
	return jobs, resp, err
}

// ListUserRunners lists the current user's action runners
func (c *Client) ListUserRunners(opt ListActionRunnersOption) ([]*models.ActionRunner, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/user/actions/runners")
	link.RawQuery = opt.getURLQuery().Encode()
	runners := make([]*models.ActionRunner, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &runners)
	return runners, resp, err
}

// GetUserRunner gets the current user's action runner by ID
func (c *Client) GetUserRunner(runnerID int64) (*models.ActionRunner, *Response, error) {
	runner := new(models.ActionRunner)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/user/actions/runners/%d", runnerID), jsonHeader, nil, runner)
	return runner, resp, err
}

// RegisterUserRunner registers a new user-level runner through the interactive
// registration flow (Forgejo v15+). Set Ephemeral to register a single-job
// runner.
func (c *Client) RegisterUserRunner(opt models.RegisterRunnerOptions) (*models.RegisterRunnerResponse, *Response, error) {
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	out := new(models.RegisterRunnerResponse)
	resp, err := c.getParsedResponse("POST", "/user/actions/runners", jsonHeader, bytes.NewReader(body), out)
	return out, resp, err
}

// DeleteUserRunner deletes the current user's action runner by ID
func (c *Client) DeleteUserRunner(runnerID int64) (*Response, error) {
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/user/actions/runners/%d", runnerID), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("runner not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
