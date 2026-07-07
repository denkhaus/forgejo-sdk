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

func Test_AdminOrg(t *testing.T) {
	log.Println("== Test_AdminOrg ==")
	c := newTestClient()
	user, _, err := c.GetMyUserInfo()
	require.NoError(t, err)

	orgName := "AdminOrgFileTest"
	// clean up any leftover org from a previous run
	_, _ = c.DeleteOrg(orgName)

	newOrg, _, err := c.AdminCreateOrg(user.UserName, CreateOrgOption{
		Name:        orgName,
		FullName:    orgName + " FullName",
		Description: "admin_org_test",
		Visibility:  VisibleTypePublic,
	})
	require.NoError(t, err)
	assert.EqualValues(t, orgName, newOrg.UserName)

	orgs, _, err := c.AdminListOrgs(AdminListOrgsOptions{ListOptions: ListOptions{Page: 1}})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(orgs), 1)

	_, err = c.DeleteOrg(orgName)
	require.NoError(t, err)
}
