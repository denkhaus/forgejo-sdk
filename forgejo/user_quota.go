// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetMyQuota gets the quota of the current authenticated user
func (c *Client) GetMyQuota() (*models.QuotaInfo, *Response, error) {
	info := new(models.QuotaInfo)
	resp, err := c.getParsedResponse("GET", "/user/quota", jsonHeader, nil, info)
	return info, resp, err
}

// CheckMyQuota reports whether the current user may perform an action that
// would consume quota for the given subject (e.g. "size:all").
func (c *Client) CheckMyQuota(subject string) (bool, *Response, error) {
	var ok bool
	link, _ := url.Parse("/user/quota/check")
	q := link.Query()
	q.Set("subject", subject)
	link.RawQuery = q.Encode()
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &ok)
	return ok, resp, err
}

// ListMyQuotaArtifacts lists the artifacts contributing to the current user's quota usage
func (c *Client) ListMyQuotaArtifacts(opt ListQuotaUsedOption) ([]*models.QuotaUsedArtifact, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/user/quota/artifacts")
	link.RawQuery = opt.getURLQuery().Encode()
	artifacts := make([]*models.QuotaUsedArtifact, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &artifacts)
	return artifacts, resp, err
}

// ListMyQuotaAttachments lists the attachments contributing to the current user's quota usage
func (c *Client) ListMyQuotaAttachments(opt ListQuotaUsedOption) ([]*models.QuotaUsedAttachment, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/user/quota/attachments")
	link.RawQuery = opt.getURLQuery().Encode()
	attachments := make([]*models.QuotaUsedAttachment, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &attachments)
	return attachments, resp, err
}

// ListMyQuotaPackages lists the packages contributing to the current user's quota usage
func (c *Client) ListMyQuotaPackages(opt ListQuotaUsedOption) ([]*models.QuotaUsedPackage, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/user/quota/packages")
	link.RawQuery = opt.getURLQuery().Encode()
	packages := make([]*models.QuotaUsedPackage, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &packages)
	return packages, resp, err
}
