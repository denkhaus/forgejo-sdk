// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// RenameOrg renames an organization (admin)
func (c *Client) RenameOrg(org string, opt models.RenameOrgOption) (*Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/orgs/%s/rename", org), jsonHeader, bytes.NewReader(body))
	return resp, err
}
