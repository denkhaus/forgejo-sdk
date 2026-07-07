// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ActivityPub(t *testing.T) {
	log.Println("== Test_ActivityPub ==")
	c := newTestClient()

	_, _, err := c.GetInstanceActor()
	if err != nil && (strings.Contains(err.Error(), "404") || strings.Contains(strings.ToLower(err.Error()), "not found")) {
		t.Skip("activitypub/federation disabled on this instance")
	}
	require.NoError(t, err)
}
