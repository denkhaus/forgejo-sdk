// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2024 The Gitea Authors. All rights reserved.
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

func TestCreateOrgActionSecret(t *testing.T) {
	log.Println("== TestCreateOrgActionSecret ==")
	c := newTestClient()

	user := createTestUser(t, "org_action_user", c)
	c.SetSudo(user.UserName)
	newOrg, _, err := c.CreateOrg(CreateOrgOption{Name: "ActionOrg"})
	require.NoError(t, err)
	assert.NotNil(t, newOrg)

	// create secret
	resp, err := c.CreateOrgActionSecret(newOrg.UserName, CreateSecretOption{Name: "test", Data: "test"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// update secret
	resp, err = c.CreateOrgActionSecret(newOrg.UserName, CreateSecretOption{Name: "test", Data: "test2"})
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// list secrets
	secrets, _, err := c.ListOrgActionSecret(newOrg.UserName, ListOrgActionSecretOption{})
	require.NoError(t, err)
	assert.Len(t, secrets, 1)
}
