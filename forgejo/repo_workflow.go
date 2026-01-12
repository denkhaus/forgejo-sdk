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

// WorkflowDispatchOption options for triggering a workflow dispatch
type WorkflowDispatchOption struct {
	// Branch or tag to run the workflow on (default: default branch)
	Ref string `json:"ref"`
	// Input parameters for the workflow (depends on workflow definition)
	// Keys are the parameter names defined in the workflow file
	Inputs map[string]any `json:"inputs"`
}

// WorkflowDispatch triggers a workflow run via workflow_dispatch event
func (c *Client) WorkflowDispatch(owner, repo, workflowFilename string, opt WorkflowDispatchOption) (*models.ActionRun, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}

	run := new(models.ActionRun)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/actions/workflows/%s/dispatches", owner, repo, workflowFilename), jsonHeader, bytes.NewReader(body), &run)
	return run, resp, err
}
