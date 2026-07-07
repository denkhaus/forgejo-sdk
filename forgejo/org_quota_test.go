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

func Test_OrgQuota(t *testing.T) {
	log.Println("== Test_OrgQuota ==")
	c := newTestClient()

	orgName := "OrgQuotaFileTest"
	_, _ = c.DeleteOrg(orgName)
	_, _, err := c.CreateOrg(CreateOrgOption{
		Name:       orgName,
		Visibility: VisibleTypePublic,
	})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg(orgName) }()

	info, _, err := c.GetOrgQuota(orgName)
	if err != nil && (strings.Contains(err.Error(), "404") || strings.Contains(strings.ToLower(err.Error()), "not found")) {
		t.Skip("quota disabled on this instance")
	}
	require.NoError(t, err)
	require.NotNil(t, info)
}
