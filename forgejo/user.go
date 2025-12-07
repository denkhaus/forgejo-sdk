// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"
	"strconv"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetUserInfo get user info by user's name
func (c *Client) GetUserInfo(user string) (*models.User, *Response, error) {
	if err := escapeValidatePathSegments(&user); err != nil {
		return nil, nil, err
	}
	u := new(models.User)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/users/%s", user), nil, nil, u)
	return u, resp, err
}

// GetMyUserInfo get user info of current user
func (c *Client) GetMyUserInfo() (*models.User, *Response, error) {
	u := new(models.User)
	resp, err := c.getParsedResponse("GET", "/user", nil, nil, u)
	return u, resp, err
}

// GetUserByID returns user by a given user ID
func (c *Client) GetUserByID(id int64) (*models.User, *Response, error) {
	if id < 0 {
		return nil, nil, fmt.Errorf("invalid user id %d", id)
	}

	query := make(url.Values)
	query.Add("uid", strconv.FormatInt(id, 10))
	users, resp, err := c.searchUsers(query.Encode())
	if err != nil {
		return nil, resp, err
	}

	if len(users) == 1 {
		return users[0], resp, err
	}

	return nil, resp, fmt.Errorf("user not found with id %d", id)
}
