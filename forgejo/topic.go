// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// SearchTopicOption options for searching topics
type SearchTopicOption struct {
	ListOptions
	// Query is the (partial) topic name to search for
	Query string
}

// SearchTopics searches for topics matching the query
func (c *Client) SearchTopics(opt SearchTopicOption) ([]*models.TopicResponse, *Response, error) {
	opt.setDefaults()
	link, _ := url.Parse("/topics/search")
	query := opt.getURLQuery()
	if opt.Query != "" {
		query.Set("q", opt.Query)
	}
	link.RawQuery = query.Encode()
	topics := make([]*models.TopicResponse, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &topics)
	return topics, resp, err
}
