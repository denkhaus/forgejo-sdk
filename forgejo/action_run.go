// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetActionsRun returns the workflow run associated with the request's actions
// token.
//
// This endpoint is intended to be called from within a running action job using
// the automatic actions token (sent as "Authorization: Bearer ${{ forgejo.token }}").
// The client must therefore be configured with that bearer token; a regular API
// access token (sent as "Authorization: token …") is not accepted.
func (c *Client) GetActionsRun() (*models.ActionRun, *Response, error) {
	run := new(models.ActionRun)
	resp, err := c.getParsedResponse("GET", "/actions/run", jsonHeader, nil, run)
	return run, resp, err
}
