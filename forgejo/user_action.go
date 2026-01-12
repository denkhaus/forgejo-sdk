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

// ListUserActionSecretOption list UserActionSecret options
type ListUserActionSecretOption struct {
	ListOptions
}

// ListUserActionSecret lists the current user's action secrets
func (c *Client) ListUserActionSecret(opt ListUserActionSecretOption) ([]*models.Secret, *Response, error) {
	opt.setDefaults()
	secrets := make([]*models.Secret, 0, opt.PageSize)

	link, _ := url.Parse("/user/actions/secrets")
	link.RawQuery = opt.getURLQuery().Encode()
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &secrets)
	return secrets, resp, err
}

// CreateUserActionSecret creates a user action secret
func (c *Client) CreateUserActionSecret(opt CreateSecretOption) (*Response, error) {
	if err := (&opt).Validate(); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/user/actions/secrets/%s", opt.Name), jsonHeader, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusCreated:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user not found or secret name invalid")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusBadRequest:
		return resp, fmt.Errorf("bad request: invalid secret data")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// UpdateUserActionSecretOption options for updating a user secret
type UpdateUserActionSecretOption struct {
	Name  string `json:"name"`
	Value string `json:"data"` // Use "data" for consistency with CreateSecretOption
}

// UpdateUserActionSecret updates a user action secret
func (c *Client) UpdateUserActionSecret(name string, opt UpdateUserActionSecretOption) (*Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, err
	}
	if len(opt.Name) == 0 {
		opt.Name = name
	}
	// Reuse CreateSecretOption for validation
	createOpt := CreateSecretOption{Name: opt.Name, Data: opt.Value}
	if err := (&createOpt).Validate(); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&createOpt)
	if err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/user/actions/secrets/%s", name), jsonHeader, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return resp, nil
	case http.StatusCreated:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user or secret not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusBadRequest:
		return resp, fmt.Errorf("bad request: invalid secret data")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// DeleteUserActionSecret deletes a user action secret
func (c *Client) DeleteUserActionSecret(name string) (*Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/user/actions/secrets/%s", name), jsonHeader, nil)
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user or secret not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
