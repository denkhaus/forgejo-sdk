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

// Test_UserSearch exercises SearchUsers (read-only query).
func Test_UserSearch(t *testing.T) {
	log.Println("== Test_UserSearch ==")
	c := newTestClient()

	users, _, err := c.SearchUsers(SearchUsersOption{
		ListOptions: ListOptions{PageSize: 10},
		KeyWord:     "test",
	})
	require.NoError(t, err)
	assert.NotNil(t, users)
}
