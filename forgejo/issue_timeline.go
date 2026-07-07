// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListIssueTimelineOptions options for listing an issue's timeline
type ListIssueTimelineOptions struct {
	ListOptions
}

// GetIssueTimeline returns the comments and timeline events of an issue.
func (c *Client) GetIssueTimeline(owner, repo string, index int64, opt ListIssueTimelineOptions) ([]*models.TimelineComment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/issues/%d/timeline", owner, repo, index))
	link.RawQuery = opt.getURLQuery().Encode()
	timeline := make([]*models.TimelineComment, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &timeline)
	return timeline, resp, err
}
