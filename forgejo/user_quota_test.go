// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserQuotaV15 exercises the current-user quota endpoints (Forgejo v15+).
func TestUserQuotaV15(t *testing.T) {
	log.Println("== TestUserQuotaV15 ==")
	c := newTestClient()

	q, _, err := c.GetMyQuota()
	if err != nil && strings.Contains(err.Error(), "404") {
		t.Skip("quota API not enabled on this instance")
	}
	require.NoError(t, err)
	require.NotNil(t, q)
	assert.NotNil(t, q.Used)

	ok, _, err := c.CheckMyQuota("size:all")
	require.NoError(t, err)
	_ = ok // bool: within quota

	opt := ListQuotaUsedOption{ListOptions: ListOptions{PageSize: 1}}
	_, _, err = c.ListMyQuotaArtifacts(opt)
	require.NoError(t, err)
	_, _, err = c.ListMyQuotaAttachments(opt)
	require.NoError(t, err)
	_, _, err = c.ListMyQuotaPackages(opt)
	require.NoError(t, err)
}
