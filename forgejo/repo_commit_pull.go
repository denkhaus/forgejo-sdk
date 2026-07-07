// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetCommitPullRequest returns the pull request that includes the given commit,
// or an error if no pull request is associated with it.
func (c *Client) GetCommitPullRequest(owner, repo, sha string) (*models.PullRequest, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &sha); err != nil {
		return nil, nil, err
	}
	pr := new(models.PullRequest)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/commits/%s/pull", owner, repo, sha), jsonHeader, nil, pr)
	return pr, resp, err
}
