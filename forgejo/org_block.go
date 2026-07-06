// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/http"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// BlockOrgUser blocks a user on behalf of an organization
func (c *Client) BlockOrgUser(org, username string) (*Response, error) {
	if err := escapeValidatePathSegments(&org, &username); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/orgs/%s/block/%s", org, username), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("organization or user not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// UnblockOrgUser unblocks a user on behalf of an organization
func (c *Client) UnblockOrgUser(org, username string) (*Response, error) {
	if err := escapeValidatePathSegments(&org, &username); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/orgs/%s/unblock/%s", org, username), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("organization or user not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// ListOrgBlockedUsers lists the users blocked by an organization
func (c *Client) ListOrgBlockedUsers(org string, opt ListBlockedUsersOption) ([]*models.BlockedUser, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/list_blocked", org))
	link.RawQuery = opt.getURLQuery().Encode()
	users := make([]*models.BlockedUser, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &users)
	return users, resp, err
}
