// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUserGPGKeyVerificationV15 exercises fetching the GPG key verification token.
func TestUserGPGKeyVerificationV15(t *testing.T) {
	log.Println("== TestUserGPGKeyVerificationV15 ==")
	c := newTestClient()

	_, _, err := c.GetMyGPGKeyVerificationToken()
	require.NoError(t, err)
}
