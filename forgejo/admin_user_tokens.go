// Copyright 2024 The Forgejo Authors. All rights reserved.
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

// AdminListUserAccessTokens lists all access tokens of the specified user (admin only).
func (c *Client) AdminListUserAccessTokens(username string, opt ListAccessTokensOptions) ([]*models.AccessToken, *Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/admin/users/%s/tokens", username))
	link.RawQuery = opt.getURLQuery().Encode()
	tokens := make([]*models.AccessToken, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &tokens)
	return tokens, resp, err
}

// AdminCreateUserAccessToken creates an access token for the specified user (admin only).
func (c *Client) AdminCreateUserAccessToken(username string, opt CreateAccessTokenOption) (*models.AccessToken, *Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, nil, err
	}
	if len(opt.Name) == 0 {
		return nil, nil, fmt.Errorf("name is empty")
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	token := new(models.AccessToken)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/admin/users/%s/tokens", username), jsonHeader, bytes.NewReader(body), token)
	return token, resp, err
}

// AdminDeleteUserAccessToken deletes an access token (identified by name) for the
// specified user (admin only).
func (c *Client) AdminDeleteUserAccessToken(username, token string) (*Response, error) {
	if err := escapeValidatePathSegments(&username, &token); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/admin/users/%s/tokens/%s", username, token), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user or token not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusUnprocessableEntity:
		return resp, fmt.Errorf("invalid token")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
