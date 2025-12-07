// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import "codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"

// GetGlobalUISettings get global ui settings witch are exposed by API
func (c *Client) GetGlobalUISettings() (*models.GeneralUISettings, *Response, error) {
	if err := c.checkServerVersionGreaterThanOrEqual(version1_13_0); err != nil {
		return nil, nil, err
	}
	conf := new(models.GeneralUISettings)
	resp, err := c.getParsedResponse("GET", "/settings/ui", jsonHeader, nil, &conf)
	return conf, resp, err
}

// GetGlobalRepoSettings get global repository settings witch are exposed by API
func (c *Client) GetGlobalRepoSettings() (*models.GeneralRepoSettings, *Response, error) {
	if err := c.checkServerVersionGreaterThanOrEqual(version1_13_0); err != nil {
		return nil, nil, err
	}
	conf := new(models.GeneralRepoSettings)
	resp, err := c.getParsedResponse("GET", "/settings/repository", jsonHeader, nil, &conf)
	return conf, resp, err
}

// GetGlobalAPISettings get global api settings witch are exposed by it
func (c *Client) GetGlobalAPISettings() (*models.GeneralAPISettings, *Response, error) {
	if err := c.checkServerVersionGreaterThanOrEqual(version1_13_0); err != nil {
		return nil, nil, err
	}
	conf := new(models.GeneralAPISettings)
	resp, err := c.getParsedResponse("GET", "/settings/api", jsonHeader, nil, &conf)
	return conf, resp, err
}

// GetGlobalAttachmentSettings get global repository settings witch are exposed by API
func (c *Client) GetGlobalAttachmentSettings() (*models.GeneralAttachmentSettings, *Response, error) {
	if err := c.checkServerVersionGreaterThanOrEqual(version1_13_0); err != nil {
		return nil, nil, err
	}
	conf := new(models.GeneralAttachmentSettings)
	resp, err := c.getParsedResponse("GET", "/settings/attachment", jsonHeader, nil, &conf)
	return conf, resp, err
}
