// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetOrgQuota gets the quota of an organization
func (c *Client) GetOrgQuota(org string) (*models.QuotaInfo, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	info := new(models.QuotaInfo)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/quota", org), jsonHeader, nil, info)
	return info, resp, err
}

// CheckOrgQuota reports whether an organization is within its quota
func (c *Client) CheckOrgQuota(org string) (bool, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return false, nil, err
	}
	var ok bool
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/quota/check", org), jsonHeader, nil, &ok)
	return ok, resp, err
}

// ListOrgQuotaArtifacts lists the artifacts contributing to an organization's quota usage
func (c *Client) ListOrgQuotaArtifacts(org string, opt ListQuotaUsedOption) ([]*models.QuotaUsedArtifact, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/quota/artifacts", org))
	link.RawQuery = opt.getURLQuery().Encode()
	artifacts := make([]*models.QuotaUsedArtifact, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &artifacts)
	return artifacts, resp, err
}

// ListOrgQuotaAttachments lists the attachments contributing to an organization's quota usage
func (c *Client) ListOrgQuotaAttachments(org string, opt ListQuotaUsedOption) ([]*models.QuotaUsedAttachment, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/quota/attachments", org))
	link.RawQuery = opt.getURLQuery().Encode()
	attachments := make([]*models.QuotaUsedAttachment, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &attachments)
	return attachments, resp, err
}

// ListOrgQuotaPackages lists the packages contributing to an organization's quota usage
func (c *Client) ListOrgQuotaPackages(org string, opt ListQuotaUsedOption) ([]*models.QuotaUsedPackage, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/quota/packages", org))
	link.RawQuery = opt.getURLQuery().Encode()
	packages := make([]*models.QuotaUsedPackage, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &packages)
	return packages, resp, err
}
