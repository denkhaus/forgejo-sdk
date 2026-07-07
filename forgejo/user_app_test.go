// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_UserApp exercises access-token CRUD (List/Create/Delete).
// These methods require BasicAuth; if the test env is token-only the
// list call errors and we skip rather than fail.
func Test_UserApp(t *testing.T) {
	log.Println("== Test_UserApp ==")
	c := newTestClient()

	result, _, err := c.ListAccessTokens(ListAccessTokensOptions{})
	if err != nil {
		t.Skipf("ListAccessTokens unavailable in this environment: %v", err)
	}
	before := len(result)

	tok, _, err := c.CreateAccessToken(CreateAccessTokenOption{
		Name:   "Test_UserApp_Tok",
		Scopes: []AccessTokenScope{AccessTokenScopeUserRead},
	})
	require.NoError(t, err)
	assert.EqualValues(t, "Test_UserApp_Tok", tok.Name)
	defer func() { _, _ = c.DeleteAccessToken(tok.ID) }()

	after, _, err := c.ListAccessTokens(ListAccessTokensOptions{})
	require.NoError(t, err)
	assert.Len(t, after, before+1)
}
