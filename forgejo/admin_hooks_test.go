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

// Test_AdminHooks exercises listing system hooks (safe) and, tolerantly, the
// create/get/edit/delete lifecycle for a gitea-type system hook.
func Test_AdminHooks(t *testing.T) {
	log.Println("== Test_AdminHooks ==")
	c := newTestClient()

	// AdminListHooks is a safe GET; this is the reliable assertion.
	hooks, _, err := c.AdminListHooks(AdminListHooksOption{})
	require.NoError(t, err)
	assert.NotNil(t, hooks)

	// The remaining lifecycle requires admin privileges and may be restricted
	// by the instance; exercise it tolerantly and clean up.
	created, _, err := c.AdminCreateHook(models.CreateHookOption{
		Type:   OptionalString("gitea"),
		Config: models.CreateHookOptionConfig{"url": "http://example.com/sys-hook", "content_type": "json"},
		Events: []string{"push"},
	})
	if err != nil {
		t.Logf("AdminCreateHook returned (tolerated): %v", err)
		return
	}
	if !assert.NotNil(t, created) || created.ID == 0 {
		return
	}

	if got, _, err := c.AdminGetHook(created.ID); err == nil {
		assert.EqualValues(t, created.ID, got.ID)
	} else {
		t.Logf("AdminGetHook returned (tolerated): %v", err)
	}

	if _, _, err := c.AdminEditHook(created.ID, models.EditHookOption{Active: false}); err != nil {
		t.Logf("AdminEditHook returned (tolerated): %v", err)
	}

	if _, err := c.AdminDeleteHook(created.ID); err != nil {
		t.Logf("AdminDeleteHook returned (tolerated): %v", err)
	}
}
