// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// AdminListEmailsOption options for listing all emails
type AdminListEmailsOption struct {
	ListOptions
}

// AdminSearchEmailsOption options for searching emails
type AdminSearchEmailsOption struct {
	ListOptions
	// Query is the (partial) email address to search for
	Query string
}

// AdminListEmails lists all email addresses across the instance (admin)
func (c *Client) AdminListEmails(opt AdminListEmailsOption) ([]*models.Email, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/admin/emails")
	link.RawQuery = opt.getURLQuery().Encode()
	emails := make([]*models.Email, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &emails)
	return emails, resp, err
}

// AdminSearchEmails searches email addresses across the instance (admin)
func (c *Client) AdminSearchEmails(opt AdminSearchEmailsOption) ([]*models.Email, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/admin/emails/search")
	query := opt.getURLQuery()
	if opt.Query != "" {
		query.Set("query", opt.Query)
	}
	link.RawQuery = query.Encode()
	emails := make([]*models.Email, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &emails)
	return emails, resp, err
}
