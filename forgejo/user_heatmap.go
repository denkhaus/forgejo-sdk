// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetUserHeatmapData returns a user's contribution heatmap
func (c *Client) GetUserHeatmapData(username string) ([]*models.UserHeatmapData, *Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, nil, err
	}
	data := make([]*models.UserHeatmapData, 0, 10)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/users/%s/heatmap", username), jsonHeader, nil, &data)
	return data, resp, err
}
