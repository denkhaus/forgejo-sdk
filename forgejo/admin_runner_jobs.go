// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// AdminSearchActionRunJobsOption options for listing admin action run jobs
type AdminSearchActionRunJobsOption struct {
	ListOptions
}

// AdminGetActionRunJobs lists all action run jobs across the instance (admin)
func (c *Client) AdminGetActionRunJobs(opt AdminSearchActionRunJobsOption) ([]*models.ActionRunJob, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/admin/actions/runners/jobs")
	link.RawQuery = opt.getURLQuery().Encode()
	jobs := make([]*models.ActionRunJob, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &jobs)
	return jobs, resp, err
}

// AdminSearchRunJobs searches for action run jobs across the instance (admin)
func (c *Client) AdminSearchRunJobs(opt AdminSearchActionRunJobsOption) ([]*models.ActionRunJob, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/admin/runners/jobs")
	link.RawQuery = opt.getURLQuery().Encode()
	jobs := make([]*models.ActionRunJob, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &jobs)
	return jobs, resp, err
}
