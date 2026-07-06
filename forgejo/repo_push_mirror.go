// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListPushMirrorsOption options for listing push mirrors
type ListPushMirrorsOption struct {
	ListOptions
}

// ListPushMirrors lists a repository's push mirrors
func (c *Client) ListPushMirrors(owner, repo string, opt ListPushMirrorsOption) ([]*models.PushMirror, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/push_mirrors", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	mirrors := make([]*models.PushMirror, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &mirrors)
	return mirrors, resp, err
}

// GetPushMirror gets a push mirror by remote name
func (c *Client) GetPushMirror(owner, repo, name string) (*models.PushMirror, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, nil, err
	}
	m := new(models.PushMirror)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/push_mirrors/%s", owner, repo, name), jsonHeader, nil, m)
	return m, resp, err
}

// DeletePushMirror deletes a push mirror by remote name
func (c *Client) DeletePushMirror(owner, repo, name string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &name); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/push_mirrors/%s", owner, repo, name), nil, nil)
	return resp, err
}

// SyncPushMirrors syncs all push mirrors of a repository immediately
func (c *Client) SyncPushMirrors(owner, repo string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/repos/%s/%s/push_mirrors-sync", owner, repo), nil, nil)
	return resp, err
}
