// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAdminRunnerRegistrationToken(t *testing.T) {
	log.Println("== TestGetAdminRunnerRegistrationToken ==")
	c := newTestClient()

	// Get admin runner registration token
	token, resp, err := c.GetAdminRunnerRegistrationToken()
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, token.Token)
}
