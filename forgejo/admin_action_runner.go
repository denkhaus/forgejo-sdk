// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetAdminRunnerRegistrationToken gets a global runner registration token
func (c *Client) GetAdminRunnerRegistrationToken() (*models.RegistrationToken, *Response, error) {
	token := new(models.RegistrationToken)
	resp, err := c.getParsedResponse("GET", "/admin/runners/registration-token", jsonHeader, nil, &token)
	return token, resp, err
}
