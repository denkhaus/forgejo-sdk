// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ConvertRepo converts a fork into a regular repository, or a regular
// repository into a fork. The repository must be eligible for conversion.
func (c *Client) ConvertRepo(owner, repo string) (*models.Repository, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	r := new(models.Repository)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/convert", owner, repo), jsonHeader, nil, r)
	return r, resp, err
}
