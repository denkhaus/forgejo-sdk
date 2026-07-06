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

// ListBlockedUsersOption options for listing blocked users
type ListBlockedUsersOption struct {
	ListOptions
}

// BlockUser blocks a user on behalf of the current authenticated user
func (c *Client) BlockUser(username string) (*Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/user/block/%s", username), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// UnblockUser unblocks a user on behalf of the current authenticated user
func (c *Client) UnblockUser(username string) (*Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/user/unblock/%s", username), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// ListBlockedUsers lists the users blocked by the current authenticated user
func (c *Client) ListBlockedUsers(opt ListBlockedUsersOption) ([]*models.BlockedUser, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/user/list_blocked")
	link.RawQuery = opt.getURLQuery().Encode()
	users := make([]*models.BlockedUser, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &users)
	return users, resp, err
}
