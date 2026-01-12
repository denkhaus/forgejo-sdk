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

// ListUserActionVariablesOption list action variables options
type ListUserActionVariablesOption struct {
	ListOptions
}

// ListUserActionVariables lists the current user's action variables
func (c *Client) ListUserActionVariables(opt ListUserActionVariablesOption) ([]*models.ActionVariable, *Response, error) {
	opt.setDefaults()
	variables := make([]*models.ActionVariable, 0, opt.PageSize)

	link, _ := url.Parse("/user/actions/variables")
	link.RawQuery = opt.getURLQuery().Encode()
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &variables)
	return variables, resp, err
}

// CreateUserActionVariable creates a user action variable
func (c *Client) CreateUserActionVariable(name string, opt models.CreateVariableOption) (*Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("POST", fmt.Sprintf("/user/actions/variables/%s", name), jsonHeader, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusCreated:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user not found or variable name invalid")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusBadRequest:
		return resp, fmt.Errorf("bad request: invalid variable data")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// GetUserActionVariable gets a user action variable by name
func (c *Client) GetUserActionVariable(name string) (*models.ActionVariable, *Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, nil, err
	}

	variable := new(models.ActionVariable)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/user/actions/variables/%s", name), jsonHeader, nil, &variable)
	return variable, resp, err
}

// UpdateUserActionVariable updates a user action variable
func (c *Client) UpdateUserActionVariable(name string, opt models.UpdateVariableOption) (*Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/user/actions/variables/%s", name), jsonHeader, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user or variable not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusBadRequest:
		return resp, fmt.Errorf("bad request: invalid variable data")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// DeleteUserActionVariable deletes a user action variable
func (c *Client) DeleteUserActionVariable(name string) (*Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/user/actions/variables/%s", name), jsonHeader, nil)
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user or variable not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
