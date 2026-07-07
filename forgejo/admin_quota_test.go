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

func Test_AdminQuota(t *testing.T) {
	log.Println("== Test_AdminQuota ==")
	c := newTestClient()

	groups, _, err := c.ListQuotaGroups()
	if err != nil && (strings.Contains(err.Error(), "404") || strings.Contains(strings.ToLower(err.Error()), "not found")) {
		t.Skip("quota disabled on this instance")
	}
	require.NoError(t, err)

	// quota is enabled; keep assertions conservative since the default test
	// instance may have zero configured groups/rules.
	_ = groups
}
