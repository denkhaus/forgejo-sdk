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

// ListActionRunsOption options for listing action runs
type ListActionRunsOption struct {
	ListOptions
	// Filter by workflow event (e.g., push, pull_request, workflow_dispatch)
	Events []string
	// Filter by run status (unknown, waiting, running, success, failure, cancelled, skipped, blocked)
	Status []string
	// Filter by run number
	RunNumber int64
	// Filter by head commit SHA
	HeadSHA string
}

// ListActionRuns lists a repository's action runs
func (c *Client) ListActionRuns(owner, repo string, opt ListActionRunsOption) ([]*models.ActionRun, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/runs", owner, repo))
	query := opt.getURLQuery()

	if len(opt.Events) > 0 {
		for _, event := range opt.Events {
			query.Add("event", event)
		}
	}
	if len(opt.Status) > 0 {
		for _, status := range opt.Status {
			query.Add("status", status)
		}
	}
	if opt.RunNumber != 0 {
		query.Add("run_number", fmt.Sprintf("%d", opt.RunNumber))
	}
	if opt.HeadSHA != "" {
		query.Add("head_sha", opt.HeadSHA)
	}

	link.RawQuery = query.Encode()
	resp := &models.ListActionRunResponse{}
	httpResp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &resp)
	return resp.Entries, httpResp, err
}

// GetActionRun gets an action run by ID
func (c *Client) GetActionRun(owner, repo string, runID int64) (*models.ActionRun, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	run := new(models.ActionRun)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/actions/runs/%d", owner, repo, runID), jsonHeader, nil, &run)
	return run, resp, err
}

// ListActionTasksOption options for listing action tasks
type ListActionTasksOption struct {
	ListOptions
}

// ListActionTasks lists a repository's action tasks
func (c *Client) ListActionTasks(owner, repo string, opt ListActionTasksOption) ([]*models.ActionTask, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/tasks", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	resp := &models.ActionTaskResponse{}
	httpResp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &resp)
	return resp.Entries, httpResp, err
}

// SearchRunnerJobsOption options for searching runner jobs
type SearchRunnerJobsOption struct {
	// Filter by job labels (comma-separated list)
	Labels []string
}

// SearchRunnerJobs searches for repository's action jobs according filter conditions
func (c *Client) SearchRunnerJobs(owner, repo string, opt SearchRunnerJobsOption) ([]*models.ActionRunJob, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/runners/jobs", owner, repo))
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

// ListRunnerJobsOption options for listing runner jobs
type ListRunnerJobsOption struct {
	ListOptions
}

// ListRunnerJobs lists a repository's runner jobs
func (c *Client) ListRunnerJobs(owner, repo string, opt ListRunnerJobsOption) ([]*models.ActionRunJob, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/runners/jobs", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	jobs := make([]*models.ActionRunJob, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &jobs)
	if jobs == nil {
		jobs = make([]*models.ActionRunJob, 0, opt.PageSize)
	}
	return jobs, resp, err
}

// GetRepoRunnerRegistrationToken gets a repository's runner registration token
func (c *Client) GetRepoRunnerRegistrationToken(owner, repo string) (*models.RegistrationToken, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	token := new(models.RegistrationToken)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/actions/runners/registration-token", owner, repo), jsonHeader, nil, &token)
	return token, resp, err
}

// ListRepoRunners lists a repository's action runners
func (c *Client) ListRepoRunners(owner, repo string, opt ListActionRunnersOption) ([]*models.ActionRunner, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/runners", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	runners := make([]*models.ActionRunner, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &runners)
	return runners, resp, err
}

// GetRepoRunner gets a repository's action runner by ID
func (c *Client) GetRepoRunner(owner, repo string, runnerID int64) (*models.ActionRunner, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	runner := new(models.ActionRunner)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/actions/runners/%d", owner, repo, runnerID), jsonHeader, nil, runner)
	return runner, resp, err
}

// RegisterRepoRunner registers a new repository-level runner through the
// interactive registration flow (Forgejo v15+). Set Ephemeral to register a
// single-job runner.
func (c *Client) RegisterRepoRunner(owner, repo string, opt models.RegisterRunnerOptions) (*models.RegisterRunnerResponse, *Response, error) {
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
	out := new(models.RegisterRunnerResponse)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/actions/runners", owner, repo), jsonHeader, bytes.NewReader(body), out)
	return out, resp, err
}

// DeleteRepoRunner deletes a repository's action runner by ID
func (c *Client) DeleteRepoRunner(owner, repo string, runnerID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/repos/%s/%s/actions/runners/%d", owner, repo, runnerID), jsonHeader, nil)
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
