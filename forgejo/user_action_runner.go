// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// SearchUserRunnerJobsOption options for searching user runner jobs
type SearchUserRunnerJobsOption struct {
	// Filter by job labels (comma-separated list)
	Labels []string
}

// SearchUserRunnerJobs searches for the current user's action jobs according to filter conditions
func (c *Client) SearchUserRunnerJobs(opt SearchUserRunnerJobsOption) ([]*models.ActionRunJob, *Response, error) {
	link, _ := url.Parse("/user/actions/runners/jobs")
	query := make(url.Values)

	if len(opt.Labels) > 0 {
		// Labels is a comma-separated list
		labelsStr := ""
		for i, label := range opt.Labels {
			if i > 0 {
				labelsStr += ","
			}
			labelsStr += label
		}
		query.Add("labels", labelsStr)
	}

	link.RawQuery = query.Encode()
	jobs := make([]*models.ActionRunJob, 0, 10)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &jobs)
	if jobs == nil {
		jobs = make([]*models.ActionRunJob, 0, 10)
	}
	return jobs, resp, err
}
