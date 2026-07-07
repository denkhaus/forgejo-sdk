// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOrgLabelsV15 exercises organization label CRUD (Forgejo v15+).
func TestOrgLabelsV15(t *testing.T) {
	log.Println("== TestOrgLabelsV15 ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{Name: "v15-org-labels"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg(org.UserName) }()

	l1, _, err := c.CreateOrgLabel(org.UserName, models.CreateLabelOption{Name: OptionalString("org-label-1"), Color: OptionalString("#aabbcc")})
	require.NoError(t, err)
	assert.NotZero(t, l1.ID)

	list, _, err := c.ListOrgLabels(org.UserName, ListOrgLabelsOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotEmpty(t, list)

	got, _, err := c.GetOrgLabel(org.UserName, l1.ID)
	require.NoError(t, err)
	assert.Equal(t, l1.ID, got.ID)

	_, _, err = c.EditOrgLabel(org.UserName, l1.ID, models.EditLabelOption{Description: "updated"})
	require.NoError(t, err)

	_, err = c.DeleteOrgLabel(org.UserName, l1.ID)
	require.NoError(t, err)

	list, _, err = c.ListOrgLabels(org.UserName, ListOrgLabelsOption{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.Empty(t, list)
}
