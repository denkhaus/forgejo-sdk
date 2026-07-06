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

// SearchOrgRunnerJobsOption options for searching org runner jobs
type SearchOrgRunnerJobsOption struct {
	// Filter by job labels (comma-separated list)
	Labels []string
}

// SearchOrgRunnerJobs searches for an organization's action jobs according to filter conditions
func (c *Client) SearchOrgRunnerJobs(org string, opt SearchOrgRunnerJobsOption) ([]*models.ActionRunJob, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/actions/runners/jobs", org))
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

// GetOrgRunnerRegistrationToken gets an organization's runner registration token
func (c *Client) GetOrgRunnerRegistrationToken(org string) (*models.RegistrationToken, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	token := new(models.RegistrationToken)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/actions/runners/registration-token", org), jsonHeader, nil, &token)
	return token, resp, err
}

// ListOrgRunners lists an organization's action runners
func (c *Client) ListOrgRunners(org string, opt ListActionRunnersOption) ([]*models.ActionRunner, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/actions/runners", org))
	link.RawQuery = opt.getURLQuery().Encode()
	runners := make([]*models.ActionRunner, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &runners)
	return runners, resp, err
}

// GetOrgRunner gets an organization's action runner by ID
func (c *Client) GetOrgRunner(org string, runnerID int64) (*models.ActionRunner, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	runner := new(models.ActionRunner)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/actions/runners/%d", org, runnerID), jsonHeader, nil, runner)
	return runner, resp, err
}

// RegisterOrgRunner registers a new organization-level runner through the
// interactive registration flow (Forgejo v15+). Set Ephemeral to register a
// single-job runner.
func (c *Client) RegisterOrgRunner(org string, opt models.RegisterRunnerOptions) (*models.RegisterRunnerResponse, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
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
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/orgs/%s/actions/runners", org), jsonHeader, bytes.NewReader(body), out)
	return out, resp, err
}

// DeleteOrgRunner deletes an organization's action runner by ID
func (c *Client) DeleteOrgRunner(org string, runnerID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/orgs/%s/actions/runners/%d", org, runnerID), jsonHeader, nil)
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
