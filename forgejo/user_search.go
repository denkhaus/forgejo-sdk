// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

type searchUsersResponse struct {
	Users []*models.User `json:"data"`
}

// SearchUsersOption options for SearchUsers
type SearchUsersOption struct {
	ListOptions
	KeyWord string
}

// QueryEncode turns options into querystring argument
func (opt *SearchUsersOption) QueryEncode() string {
	query := make(url.Values)
	if opt.Page > 0 {
		query.Add("page", fmt.Sprintf("%d", opt.Page))
	}
	if opt.PageSize > 0 {
		query.Add("limit", fmt.Sprintf("%d", opt.PageSize))
	}
	if len(opt.KeyWord) > 0 {
		query.Add("q", opt.KeyWord)
	}
	return query.Encode()
}

func (c *Client) searchUsers(rawQuery string) ([]*models.User, *Response, error) {
	link, _ := url.Parse("/users/search")
	link.RawQuery = rawQuery
	userResp := new(searchUsersResponse)
	resp, err := c.getParsedResponse("GET", link.String(), nil, nil, &userResp)
	return userResp.Users, resp, err
}

// SearchUsers finds users by query
func (c *Client) SearchUsers(opt SearchUsersOption) ([]*models.User, *Response, error) {
	return c.searchUsers(opt.QueryEncode())
}
