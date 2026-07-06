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

// ListActionRunnersOption options for listing action runners
type ListActionRunnersOption struct {
	ListOptions
}

// GetAdminRunnerRegistrationToken gets a global runner registration token
func (c *Client) GetAdminRunnerRegistrationToken() (*models.RegistrationToken, *Response, error) {
	token := new(models.RegistrationToken)
	resp, err := c.getParsedResponse("GET", "/admin/runners/registration-token", jsonHeader, nil, &token)
	return token, resp, err
}

// ListAdminRunners lists all instance-level action runners
func (c *Client) ListAdminRunners(opt ListActionRunnersOption) ([]*models.ActionRunner, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/admin/actions/runners")
	link.RawQuery = opt.getURLQuery().Encode()
	runners := make([]*models.ActionRunner, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &runners)
	return runners, resp, err
}

// GetAdminRunner gets an instance-level action runner by ID
func (c *Client) GetAdminRunner(runnerID int64) (*models.ActionRunner, *Response, error) {
	runner := new(models.ActionRunner)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/admin/actions/runners/%d", runnerID), jsonHeader, nil, runner)
	return runner, resp, err
}

// RegisterAdminRunner registers a new instance-level runner through the
// interactive registration flow (Forgejo v15+). The returned token and UUID
// are used by the runner to connect. Set Ephemeral to register a single-job
// runner whose registration is removed after the job completes.
func (c *Client) RegisterAdminRunner(opt models.RegisterRunnerOptions) (*models.RegisterRunnerResponse, *Response, error) {
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	out := new(models.RegisterRunnerResponse)
	resp, err := c.getParsedResponse("POST", "/admin/actions/runners", jsonHeader, bytes.NewReader(body), out)
	return out, resp, err
}

// DeleteAdminRunner deletes an instance-level action runner by ID
func (c *Client) DeleteAdminRunner(runnerID int64) (*Response, error) {
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/admin/actions/runners/%d", runnerID), jsonHeader, nil)
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
