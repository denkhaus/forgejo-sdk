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

// EditIssueDeadline sets or updates the deadline of an issue
func (c *Client) EditIssueDeadline(owner, repo string, index int64, opt models.EditDeadlineOption) (*models.IssueDeadline, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	deadline := new(models.IssueDeadline)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/issues/%d/deadline", owner, repo, index), jsonHeader, bytes.NewReader(body), deadline)
	return deadline, resp, err
}
