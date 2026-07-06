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

// ListOrgLabelsOption options for listing an organization's labels
type ListOrgLabelsOption struct {
	ListOptions
}

// ListOrgLabels lists an organization's labels
func (c *Client) ListOrgLabels(org string, opt ListOrgLabelsOption) ([]*models.Label, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/labels", org))
	link.RawQuery = opt.getURLQuery().Encode()
	labels := make([]*models.Label, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &labels)
	return labels, resp, err
}

// GetOrgLabel gets an organization label by ID
func (c *Client) GetOrgLabel(org string, id int64) (*models.Label, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	label := new(models.Label)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/labels/%d", org, id), jsonHeader, nil, label)
	return label, resp, err
}

// CreateOrgLabel creates a label for an organization
func (c *Client) CreateOrgLabel(org string, opt models.CreateLabelOption) (*models.Label, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	label := new(models.Label)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/orgs/%s/labels", org), jsonHeader, bytes.NewReader(body), label)
	return label, resp, err
}

// EditOrgLabel edits an organization label by ID
func (c *Client) EditOrgLabel(org string, id int64, opt models.EditLabelOption) (*models.Label, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	label := new(models.Label)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/orgs/%s/labels/%d", org, id), jsonHeader, bytes.NewReader(body), label)
	return label, resp, err
}

// DeleteOrgLabel deletes an organization label by ID
func (c *Client) DeleteOrgLabel(org string, id int64) (*Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/orgs/%s/labels/%d", org, id), nil, nil)
	return resp, err
}
