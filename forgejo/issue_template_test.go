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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_IssueTemplate lists the issue templates of a repository (a safe GET).
func Test_IssueTemplate(t *testing.T) {
	log.Println("== Test_IssueTemplate ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "TestIssueTemplateRepo", c)
	require.NoError(t, err)

	templates, _, err := c.GetIssueTemplates(repo.Owner.UserName, repo.Name)
	require.NoError(t, err)
	// A freshly auto-initialized repository ships no issue templates by
	// default, so an empty list is a valid result; the call returning
	// without error (above) is sufficient.
	for _, tpl := range templates {
		assert.NotEmpty(t, tpl.Filename)
	}
}
