// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"
)

// AdminListUnadoptedOption options for listing unadopted repositories
type AdminListUnadoptedOption struct {
	ListOptions
	// Query filters the returned repositories (owner/name substring)
	Query string
}

// AdminListUnadoptedRepositories lists unadopted repositories (admin)
func (c *Client) AdminListUnadoptedRepositories(opt AdminListUnadoptedOption) ([]string, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/admin/unadopted")
	query := opt.getURLQuery()
	if opt.Query != "" {
		query.Set("query", opt.Query)
	}
	link.RawQuery = query.Encode()
	repos := make([]string, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &repos)
	return repos, resp, err
}

// AdminAdoptRepository adopts an unadopted repository (admin)
func (c *Client) AdminAdoptRepository(owner, repo string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/admin/unadopted/%s/%s", owner, repo), nil, nil)
	return resp, err
}

// AdminDeleteUnadoptedRepository deletes an unadopted repository (admin)
func (c *Client) AdminDeleteUnadoptedRepository(owner, repo string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/admin/unadopted/%s/%s", owner, repo), nil, nil)
	return resp, err
}
