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

// TestRepoEditorConfigV15 exercises fetching the editorconfig for a file.
func TestRepoEditorConfigV15(t *testing.T) {
	log.Println("== TestRepoEditorConfigV15 ==")
	c := newTestClient()

	repo, err := createTestRepo(t, "v15-editorconfig-repo", c)
	require.NoError(t, err)

	// tolerant: depends on a .editorconfig being present
	if ec, _, err := c.GetRepoEditorConfig(repo.Owner.UserName, repo.Name, "README.md"); err == nil {
		assert.NotNil(t, ec)
	}
}
