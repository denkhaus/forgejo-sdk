// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2024 The Gitea Authors. All rights reserved.
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

// ListOrgActionVariablesOption list action variables options
type ListOrgActionVariablesOption struct {
	ListOptions
}

// ListOrgActionVariables lists an organization's action variables
func (c *Client) ListOrgActionVariables(org string, opt ListOrgActionVariablesOption) ([]*models.ActionVariable, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	variables := make([]*models.ActionVariable, 0, opt.PageSize)

	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/actions/variables", org))
	link.RawQuery = opt.getURLQuery().Encode()
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &variables)
	return variables, resp, err
}

// CreateOrgActionVariable creates an organization action variable
func (c *Client) CreateOrgActionVariable(org, name string, opt models.CreateVariableOption) (*Response, error) {
	if err := escapeValidatePathSegments(&org, &name); err != nil {
		return nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("POST", fmt.Sprintf("/orgs/%s/actions/variables/%s", org, name), jsonHeader, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusCreated:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("organization not found or variable name invalid")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusBadRequest:
		return resp, fmt.Errorf("bad request: invalid variable data")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// GetOrgActionVariable gets an organization action variable by name
func (c *Client) GetOrgActionVariable(org, name string) (*models.ActionVariable, *Response, error) {
	if err := escapeValidatePathSegments(&org, &name); err != nil {
		return nil, nil, err
	}

	variable := new(models.ActionVariable)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/actions/variables/%s", org, name), jsonHeader, nil, &variable)
	return variable, resp, err
}

// UpdateOrgActionVariable updates an organization action variable
func (c *Client) UpdateOrgActionVariable(org, name string, opt models.UpdateVariableOption) (*Response, error) {
	if err := escapeValidatePathSegments(&org, &name); err != nil {
		return nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/orgs/%s/actions/variables/%s", org, name), jsonHeader, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("organization or variable not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusBadRequest:
		return resp, fmt.Errorf("bad request: invalid variable data")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// DeleteOrgActionVariable deletes an organization action variable
func (c *Client) DeleteOrgActionVariable(org, name string) (*Response, error) {
	if err := escapeValidatePathSegments(&org, &name); err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/orgs/%s/actions/variables/%s", org, name), jsonHeader, nil)
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("organization or variable not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
