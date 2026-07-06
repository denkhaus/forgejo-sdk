// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// AdminListUserEmails lists the email addresses of a user (admin)
func (c *Client) AdminListUserEmails(username string) ([]*models.Email, *Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, nil, err
	}
	emails := make([]*models.Email, 0, 10)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/admin/users/%s/emails", username), jsonHeader, nil, &emails)
	return emails, resp, err
}

// AdminDeleteUserEmails deletes email addresses from a user (admin)
func (c *Client) AdminDeleteUserEmails(username string, opt models.DeleteEmailOption) (*Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/admin/users/%s/emails", username), jsonHeader, bytes.NewReader(body))
	return resp, err
}
