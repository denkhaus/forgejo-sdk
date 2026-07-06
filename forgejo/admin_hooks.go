// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// AdminListHooksOption options for listing system hooks
type AdminListHooksOption struct {
	ListOptions
}

// AdminListHooks lists all system (admin) hooks
func (c *Client) AdminListHooks(opt AdminListHooksOption) ([]*models.Hook, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/admin/hooks")
	link.RawQuery = opt.getURLQuery().Encode()
	hooks := make([]*models.Hook, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &hooks)
	return hooks, resp, err
}

// AdminGetHook gets a system hook by ID
func (c *Client) AdminGetHook(id int64) (*models.Hook, *Response, error) {
	h := new(models.Hook)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/admin/hooks/%d", id), jsonHeader, nil, h)
	return h, resp, err
}

// AdminCreateHook creates a system hook
func (c *Client) AdminCreateHook(opt models.CreateHookOption) (*models.Hook, *Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	h := new(models.Hook)
	resp, err := c.getParsedResponse("POST", "/admin/hooks", jsonHeader, bytes.NewReader(body), h)
	return h, resp, err
}

// AdminEditHook edits a system hook by ID
func (c *Client) AdminEditHook(id int64, opt models.EditHookOption) (*models.Hook, *Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	h := new(models.Hook)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/admin/hooks/%d", id), jsonHeader, bytes.NewReader(body), h)
	return h, resp, err
}

// AdminDeleteHook deletes a system hook by ID
func (c *Client) AdminDeleteHook(id int64) (*Response, error) {
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/admin/hooks/%d", id), nil, nil)
	return resp, err
}
