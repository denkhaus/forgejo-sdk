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

// Test_AdminUser exercises AdminListUsers (read-only; at least test01 exists).
func Test_AdminUser(t *testing.T) {
	log.Println("== Test_AdminUser ==")
	c := newTestClient()

	users, _, err := c.AdminListUsers(AdminListUsersOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotEmpty(t, users)
}
