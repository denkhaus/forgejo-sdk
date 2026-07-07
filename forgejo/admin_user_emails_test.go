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

// Test_AdminUserEmails exercises AdminListUserEmails (read-only).
func Test_AdminUserEmails(t *testing.T) {
	log.Println("== Test_AdminUserEmails ==")
	c := newTestClient()

	emails, _, err := c.AdminListUserEmails("test01")
	require.NoError(t, err)
	assert.NotEmpty(t, emails)
}
