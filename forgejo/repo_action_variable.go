// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2024 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/http"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListRepoActionVariablesOption list action variables options
type ListRepoActionVariablesOption struct {
	ListOptions
}

// ListRepoActionVariables lists a repository's action variables
func (c *Client) ListRepoActionVariables(owner, repo string, opt ListRepoActionVariablesOption) ([]*models.ActionVariable, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	variables := make([]*models.ActionVariable, 0, opt.PageSize)

	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/actions/variables", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &variables)
	return variables, resp, err
}

// CreateRepoActionVariable creates a repository action variable
func (c *Client) CreateRepoActionVariable(owner, repo, name string, opt models.CreateVariableOption) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, err
	}
	return c.submitActionVariable("POST", fmt.Sprintf("/repos/%s/%s/actions/variables/%s", owner, repo, name), &opt, http.StatusCreated, "repository not found or variable name invalid")
}

// GetRepoActionVariable gets a repository action variable by name
func (c *Client) GetRepoActionVariable(owner, repo, name string) (*models.ActionVariable, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, nil, err
	}

	variable := new(models.ActionVariable)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/actions/variables/%s", owner, repo, name), jsonHeader, nil, &variable)
	return variable, resp, err
}

// UpdateRepoActionVariable updates a repository action variable
func (c *Client) UpdateRepoActionVariable(owner, repo, name string, opt models.UpdateVariableOption) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, err
	}
	return c.submitActionVariable("PUT", fmt.Sprintf("/repos/%s/%s/actions/variables/%s", owner, repo, name), &opt, http.StatusOK, "repository or variable not found")
}

// DeleteRepoActionVariable deletes a repository action variable
func (c *Client) DeleteRepoActionVariable(owner, repo, name string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, err
	}

	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/repos/%s/%s/actions/variables/%s", owner, repo, name), jsonHeader, nil)
	if err != nil {
		return nil, err
	}

	switch status {
	case http.StatusOK:
		return resp, nil
	case http.StatusNoContent:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("repository or variable not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
