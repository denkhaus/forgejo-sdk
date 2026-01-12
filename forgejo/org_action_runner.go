// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// SearchOrgRunnerJobsOption options for searching org runner jobs
type SearchOrgRunnerJobsOption struct {
	// Filter by job labels (comma-separated list)
	Labels []string
}

// SearchOrgRunnerJobs searches for an organization's action jobs according to filter conditions
func (c *Client) SearchOrgRunnerJobs(org string, opt SearchOrgRunnerJobsOption) ([]*models.ActionRunJob, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/actions/runners/jobs", org))
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

// GetOrgRunnerRegistrationToken gets an organization's runner registration token
func (c *Client) GetOrgRunnerRegistrationToken(org string) (*models.RegistrationToken, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	token := new(models.RegistrationToken)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/actions/runners/registration-token", org), jsonHeader, nil, &token)
	return token, resp, err
}
