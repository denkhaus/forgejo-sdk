// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gitea Authors. All rights reserved.
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

func TestGetGlobalSettings(t *testing.T) {
	log.Println("== TestGetGlobalSettings ==")
	c := newTestClient()

	uiSettings, _, err := c.GetGlobalUISettings()
	require.NoError(t, err)
	expectedAllowedReactions := []string{"+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"}
	assert.ElementsMatch(t, expectedAllowedReactions, uiSettings.AllowedReactions)

	repoSettings, _, err := c.GetGlobalRepoSettings()
	require.NoError(t, err)
	assert.EqualValues(t, &models.GeneralRepoSettings{
		HTTPGitDisabled: false,
		MirrorsDisabled: false,
		LFSDisabled:     true,
	}, repoSettings)

	apiSettings, _, err := c.GetGlobalAPISettings()
	require.NoError(t, err)
	assert.EqualValues(t, &models.GeneralAPISettings{
		MaxResponseItems:       50,
		DefaultPagingNum:       30,
		DefaultGitTreesPerPage: 1000,
		DefaultMaxBlobSize:     10485760,
	}, apiSettings)

	attachSettings, _, err := c.GetGlobalAttachmentSettings()
	require.NoError(t, err)
	if assert.NotEmpty(t, attachSettings.AllowedTypes) {
		attachSettings.AllowedTypes = ""
	}
	assert.EqualValues(t, &models.GeneralAttachmentSettings{
		Enabled:  true,
		MaxSize:  2048,
		MaxFiles: 5,
	}, attachSettings)
}
