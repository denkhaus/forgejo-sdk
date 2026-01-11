// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetOrgRunnerRegistrationToken gets an organization's runner registration token
func (c *Client) GetOrgRunnerRegistrationToken(org string) (*models.RegistrationToken, *Response, error) {
	if err := escapeValidatePathSegments(&org); err != nil {
		return nil, nil, err
	}
	token := new(models.RegistrationToken)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/orgs/%s/actions/runners/registration-token", org), jsonHeader, nil, &token)
	return token, resp, err
}
