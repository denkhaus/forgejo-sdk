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

// Test_UserGpgKey exercises GPG-key listing (safe, no key material required).
func Test_UserGpgKey(t *testing.T) {
	log.Println("== Test_UserGpgKey ==")
	c := newTestClient()

	myKeys, _, err := c.ListMyGPGKeys(&ListGPGKeysOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotNil(t, myKeys)

	_, _, err = c.ListGPGKeys("test01", ListGPGKeysOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
}
