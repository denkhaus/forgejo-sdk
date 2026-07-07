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

// ApplyRepoDiffPatch applies a diff patch to a repository file and returns the result.
func (c *Client) ApplyRepoDiffPatch(owner, repo string, opt UpdateFileOptions) (*models.FileResponse, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	fr := new(models.FileResponse)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/diffpatch", owner, repo), jsonHeader, bytes.NewReader(body), fr)
	return fr, resp, err
}
