// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListEmailsOptions options for listing current's user emails
type ListEmailsOptions struct {
	ListOptions
}

// ListEmails all the email addresses of user
func (c *Client) ListEmails(opt ListEmailsOptions) ([]*models.Email, *Response, error) {
	opt.setDefaults()
	emails := make([]*models.Email, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/user/emails?%s", opt.getURLQuery().Encode()), nil, nil, &emails)
	return emails, resp, err
}

// CreateEmailOption options when creating email addresses
type CreateEmailOption struct {
	// email addresses to add
	Emails []string `json:"emails"`
}

// AddEmail add one email to current user with options
func (c *Client) AddEmail(opt CreateEmailOption) ([]*models.Email, *Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	emails := make([]*models.Email, 0, 3)
	resp, err := c.getParsedResponse("POST", "/user/emails", jsonHeader, bytes.NewReader(body), &emails)
	return emails, resp, err
}

// DeleteEmailOption options when deleting email addresses
type DeleteEmailOption struct {
	// email addresses to delete
	Emails []string `json:"emails"`
}

// DeleteEmail delete one email of current users'
func (c *Client) DeleteEmail(opt DeleteEmailOption) (*Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", "/user/emails", jsonHeader, bytes.NewReader(body))
	return resp, err
}
