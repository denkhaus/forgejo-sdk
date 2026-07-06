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

// ListRepoFlags lists the flags set on a repository
func (c *Client) ListRepoFlags(owner, repo string) ([]string, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	flags := make([]string, 0, 10)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/flags", owner, repo), jsonHeader, nil, &flags)
	return flags, resp, err
}

// CheckRepoFlag reports whether a flag is set on a repository. A nil error
// means the flag is set; a 404 means it is not.
func (c *Client) CheckRepoFlag(owner, repo, flag string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &flag); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("GET", fmt.Sprintf("/repos/%s/%s/flags/%s", owner, repo, flag), nil, nil)
	return resp, err
}

// ReplaceRepoFlags replaces all flags on a repository
func (c *Client) ReplaceRepoFlags(owner, repo string, opt models.ReplaceFlagsOption) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("PUT", fmt.Sprintf("/repos/%s/%s/flags", owner, repo), jsonHeader, bytes.NewReader(body))
	return resp, err
}

// AddRepoFlag adds a flag to a repository
func (c *Client) AddRepoFlag(owner, repo, flag string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &flag); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("PUT", fmt.Sprintf("/repos/%s/%s/flags/%s", owner, repo, flag), nil, nil)
	return resp, err
}

// DeleteAllRepoFlags removes all flags from a repository
func (c *Client) DeleteAllRepoFlags(owner, repo string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/flags", owner, repo), nil, nil)
	return resp, err
}

// DeleteRepoFlag removes a single flag from a repository
func (c *Client) DeleteRepoFlag(owner, repo, flag string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &flag); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/flags/%s", owner, repo, flag), nil, nil)
	return resp, err
}
