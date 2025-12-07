// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// AdminCreateRepo create a repo
func (c *Client) AdminCreateRepo(user string, opt CreateRepoOption) (*models.Repository, *Response, error) {
	if err := escapeValidatePathSegments(&user); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	repo := new(models.Repository)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/admin/users/%s/repos", user), jsonHeader, bytes.NewReader(body), repo)
	return repo, resp, err
}
