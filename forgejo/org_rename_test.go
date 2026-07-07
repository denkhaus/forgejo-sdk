// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/require"
)

// TestOrgRenameV15 exercises renaming an organization.
func TestOrgRenameV15(t *testing.T) {
	log.Println("== TestOrgRenameV15 ==")
	c := newTestClient()

	_, _, err := c.CreateOrg(CreateOrgOption{Name: "v15-org-rename"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg("v15-org-rename-2") }()

	_, err = c.RenameOrg("v15-org-rename", models.RenameOrgOption{NewName: OptionalString("v15-org-rename-2")})
	require.NoError(t, err)
}
