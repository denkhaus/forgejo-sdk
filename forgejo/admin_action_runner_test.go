// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"net/http"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
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

// TestRunnersV15 exercises the v15 interactive runner registration flow.
func TestRunnersV15(t *testing.T) {
	log.Println("== TestRunnersV15 ==")
	c := newTestClient()

	// listing works even when no runners exist
	_, _, err := c.ListAdminRunners(ListActionRunnersOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)

	// register an ephemeral runner via the v15 interactive flow
	name := "v15-test-runner"
	rr, _, err := c.RegisterAdminRunner(models.RegisterRunnerOptions{Name: &name, Ephemeral: true})
	require.NoError(t, err)
	require.NotNil(t, rr)
	assert.NotZero(t, rr.ID)
	assert.NotEmpty(t, rr.Token)
	assert.NotEmpty(t, rr.UUID)

	runners, _, err := c.ListAdminRunners(ListActionRunnersOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotEmpty(t, runners)

	got, _, err := c.GetAdminRunner(rr.ID)
	require.NoError(t, err)
	assert.Equal(t, rr.ID, got.ID)

	_, err = c.DeleteAdminRunner(rr.ID)
	require.NoError(t, err)

	_, _, err = c.GetAdminRunner(rr.ID)
	require.Error(t, err)
}
