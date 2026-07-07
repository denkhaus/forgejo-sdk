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

func Test_GitBlob(t *testing.T) {
	log.Println("== Test_GitBlob ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "GitBlob", c)
	require.NoError(t, err)

	owner := repo.Owner.UserName
	repoName := repo.Name

	// Resolve a real blob SHA by listing the repo tree (a commit SHA is not a
	// blob, so traverse the tree and pick the first blob entry). The auto-init
	// repo contains README/LICENSE/.gitignore, so at least one blob exists.
	tree, _, err := c.GetTrees(owner, repoName, repo.DefaultBranch, GetTreesOptions{
		Recursive:   true,
		ListOptions: ListOptions{Page: 1, PageSize: 100},
	})
	require.NoError(t, err)
	require.NotNil(t, tree)

	var blobSHA string
	for _, e := range tree.Entries {
		if e.Type == "blob" && e.SHA != "" {
			blobSHA = e.SHA
			break
		}
	}
	if blobSHA == "" {
		t.Skip("no blob entry found in tree")
	}

	blob, _, err := c.GetBlob(owner, repoName, blobSHA)
	require.NoError(t, err)
	require.NotNil(t, blob)
	assert.Equal(t, blobSHA, blob.SHA)
	assert.NotEmpty(t, blob.Content)
}
