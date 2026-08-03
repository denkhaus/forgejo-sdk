// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2021 The Gitea Authors. All rights reserved.
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

func TestUserSettings(t *testing.T) {
	log.Println("== TestUserSettings ==")
	c := newTestClient()

	userConf, _, err := c.GetUserSettings()
	require.NoError(t, err)
	assert.NotNil(t, userConf)
	assert.EqualValues(t, models.UserSettings{
		Theme:               "forgejo-auto",
		HideEmail:           false,
		HideActivity:        false,
		EnableRepoUnitHints: true, // default flipped to true in Forgejo 16
	}, *userConf)

	userConf, _, err = c.UpdateUserSettings(UserSettingsOptions{
		FullName:  OptionalString("Admin User on Test"),
		Language:  OptionalString("de_de"),
		HideEmail: OptionalBool(true),
	})
	require.NoError(t, err)
	assert.NotNil(t, userConf)
	assert.EqualValues(t, models.UserSettings{
		FullName:            "Admin User on Test",
		Theme:               "forgejo-auto",
		Language:            "de_de",
		HideEmail:           true,
		HideActivity:        false,
		EnableRepoUnitHints: true, // default flipped to true in Forgejo 16
	}, *userConf)

	_, _, err = c.UpdateUserSettings(UserSettingsOptions{
		FullName:  OptionalString(""),
		Language:  OptionalString(""),
		HideEmail: OptionalBool(false),
	})
	require.NoError(t, err)
}
