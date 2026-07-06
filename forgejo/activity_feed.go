// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListActivityFeedsOption options for listing activity feeds
type ListActivityFeedsOption struct {
	ListOptions
}

// ListRepoActivityFeeds lists a repository's activity feeds
func (c *Client) ListRepoActivityFeeds(owner, repo string, opt ListActivityFeedsOption) ([]*models.Activity, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/activities/feeds", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	feeds := make([]*models.Activity, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &feeds)
	return feeds, resp, err
}

// ListOrgActivityFeeds lists an organization's activity feeds
func (c *Client) ListOrgActivityFeeds(org string, opt ListActivityFeedsOption) ([]*models.Activity, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/orgs/%s/activities/feeds", org))
	link.RawQuery = opt.getURLQuery().Encode()
	feeds := make([]*models.Activity, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &feeds)
	return feeds, resp, err
}

// ListUserActivityFeeds lists a user's activity feeds
func (c *Client) ListUserActivityFeeds(username string, opt ListActivityFeedsOption) ([]*models.Activity, *Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/users/%s/activities/feeds", username))
	link.RawQuery = opt.getURLQuery().Encode()
	feeds := make([]*models.Activity, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &feeds)
	return feeds, resp, err
}

// ListTeamActivityFeeds lists a team's activity feeds
func (c *Client) ListTeamActivityFeeds(teamID int64, opt ListActivityFeedsOption) ([]*models.Activity, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/teams/%d/activities/feeds", teamID))
	link.RawQuery = opt.getURLQuery().Encode()
	feeds := make([]*models.Activity, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &feeds)
	return feeds, resp, err
}
