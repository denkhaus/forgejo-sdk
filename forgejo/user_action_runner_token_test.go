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

// TestUserRunnerRegistrationTokenV15 exercises fetching the user-level runner registration token.
func TestUserRunnerRegistrationTokenV15(t *testing.T) {
	log.Println("== TestUserRunnerRegistrationTokenV15 ==")
	c := newTestClient()

	tok, _, err := c.GetMyRunnerRegistrationToken()
	require.NoError(t, err)
	assert.NotNil(t, tok)
}
