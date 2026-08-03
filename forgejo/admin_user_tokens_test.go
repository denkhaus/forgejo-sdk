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

func TestAdminUserAccessTokens(t *testing.T) {
	log.Println("== TestAdminUserAccessTokens ==")
	c := newTestClient()

	user := createTestUser(t, "admin_tokens_user", c)

	tokenName := "admin-test-token"
	// pre-cleanup in case a prior run left a token behind
	_, _ = c.AdminDeleteUserAccessToken(user.UserName, tokenName)

	// create (Forgejo requires at least one scope on a token)
	tok, resp, err := c.AdminCreateUserAccessToken(user.UserName, CreateAccessTokenOption{
		Name:   tokenName,
		Scopes: []AccessTokenScope{AccessTokenScopeAll},
	})
	require.NoError(t, err)
	require.NotNil(t, tok)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, tokenName, tok.Name)
	assert.NotEmpty(t, tok.Token)

	// list contains the new token
	tokens, _, err := c.AdminListUserAccessTokens(user.UserName, ListAccessTokensOptions{ListOptions: ListOptions{PageSize: 50}})
	require.NoError(t, err)
	assert.NotEmpty(t, tokens)
	assert.True(t, containsAccessTokenName(tokens, tokenName))

	// delete
	_, err = c.AdminDeleteUserAccessToken(user.UserName, tokenName)
	require.NoError(t, err)

	// verify it is gone
	tokens, _, err = c.AdminListUserAccessTokens(user.UserName, ListAccessTokensOptions{})
	require.NoError(t, err)
	assert.False(t, containsAccessTokenName(tokens, tokenName))
}

func TestAdminCreateUserAccessTokenValidation(t *testing.T) {
	log.Println("== TestAdminCreateUserAccessTokenValidation ==")
	c := newTestClient()

	user := createTestUser(t, "admin_tokens_user", c)

	// empty name must fail client-side before any request is made
	_, _, err := c.AdminCreateUserAccessToken(user.UserName, CreateAccessTokenOption{Name: ""})
	require.Error(t, err)
}

func containsAccessTokenName(tokens []*models.AccessToken, name string) bool {
	for _, tk := range tokens {
		if tk.Name == name {
			return true
		}
	}
	return false
}
