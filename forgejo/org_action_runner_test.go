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

func TestGetOrgRunnerRegistrationToken(t *testing.T) {
	log.Println("== TestGetOrgRunnerRegistrationToken ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{
		Name:    "test-runner-org",
		FullName: "Test Runner Org",
	})
	require.NoError(t, err)
	assert.NotNil(t, org)

	// Get org runner registration token
	token, resp, err := c.GetOrgRunnerRegistrationToken(org.UserName)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, token.Token)
}
