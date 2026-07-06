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

// GetNote returns the note for a given commit (by SHA) in a repository
func (c *Client) GetNote(owner, repo, sha string) (*models.Note, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &sha); err != nil {
		return nil, nil, err
	}
	note := new(models.Note)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/git/notes/%s", owner, repo, sha), jsonHeader, nil, note)
	return note, resp, err
}

// SetNote sets or updates the note for a given commit (by SHA) in a repository
func (c *Client) SetNote(owner, repo, sha string, opt models.NoteOptions) (*models.Note, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &sha); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	note := new(models.Note)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/git/notes/%s", owner, repo, sha), jsonHeader, bytes.NewReader(body), note)
	return note, resp, err
}

// RemoveNote removes the note for a given commit (by SHA) in a repository
func (c *Client) RemoveNote(owner, repo, sha string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &sha); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/git/notes/%s", owner, repo, sha), nil, nil)
	return resp, err
}
