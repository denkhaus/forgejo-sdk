// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
)

// GetRepoEditorConfig returns the .editorconfig settings applicable to a file
// path in the repository, as a map of config keys to values (e.g.
// {"indent_style": "space", "indent_size": 4}).
func (c *Client) GetRepoEditorConfig(owner, repo, filepath string) (map[string]any, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &filepath); err != nil {
		return nil, nil, err
	}
	out := make(map[string]any)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/editorconfig/%s", owner, repo, filepath), jsonHeader, nil, &out)
	return out, resp, err
}
