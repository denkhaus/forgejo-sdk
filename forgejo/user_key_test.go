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

// Test_UserKey exercises public-key listing (safe; adding needs valid key material).
func Test_UserKey(t *testing.T) {
	log.Println("== Test_UserKey ==")
	c := newTestClient()

	myKeys, _, err := c.ListMyPublicKeys(ListPublicKeysOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotNil(t, myKeys)

	_, _, err = c.ListPublicKeys("test01", ListPublicKeysOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
}
