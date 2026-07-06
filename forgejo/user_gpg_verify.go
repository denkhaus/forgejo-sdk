// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// GetMyGPGKeyVerificationToken returns the token to sign for GPG key verification
func (c *Client) GetMyGPGKeyVerificationToken() (string, *Response, error) {
	data, resp, err := c.getResponse("GET", "/user/gpg_key_token", nil, nil)
	return string(data), resp, err
}

// VerifyMyGPGKey verifies a GPG key with the signed token and returns the verified key
func (c *Client) VerifyMyGPGKey(opt models.VerifyGPGKeyOption) (*models.GPGKey, *Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	key := new(models.GPGKey)
	resp, err := c.getParsedResponse("POST", "/user/gpg_key_verify", jsonHeader, bytes.NewReader(body), key)
	return key, resp, err
}
