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

// TestMiscV15 exercises instance-level v15 endpoints.
func TestMiscV15(t *testing.T) {
	log.Println("== TestMiscV15 ==")
	c := newTestClient()

	ni, _, err := c.GetNodeInfo()
	if err == nil { // nodeinfo may be disabled (federation) on the instance
		require.NotNil(t, ni)
		assert.NotEmpty(t, ni.Version)
	}

	// signing keys depend on instance configuration and may be absent/unset
	_, _, _ = c.GetSigningKey()
	_, _, _ = c.GetSSHSigningKey()

	html, _, err := c.RenderMarkdown(models.MarkdownOption{Text: "# hello v15", Mode: "markdown"})
	require.NoError(t, err)
	assert.Contains(t, html, "hello v15")

	raw, _, err := c.RenderMarkdownRaw([]byte("# raw hello"))
	require.NoError(t, err)
	assert.Contains(t, raw, "raw hello")

	git, _, err := c.ListGitignoreTemplates()
	require.NoError(t, err)
	assert.NotEmpty(t, git)

	_, _, err = c.ListLicenseTemplates()
	require.NoError(t, err)

	_, _, err = c.ListLabelTemplates()
	require.NoError(t, err)

	_, _, err = c.SearchTopics(SearchTopicOption{Query: "forgejo", ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
}

// TestMiscMoreV15 exercises markup rendering and template lookups.
func TestMiscMoreV15(t *testing.T) {
	log.Println("== TestMiscMoreV15 ==")
	c := newTestClient()

	html, _, err := c.RenderMarkup(models.MarkupOption{Text: "<b>hi</b>", Mode: "markdown"})
	require.NoError(t, err)
	_ = html

	git, _, err := c.ListGitignoreTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, git)
	info, _, err := c.GetGitignoreTemplateInfo(git[0])
	require.NoError(t, err)
	assert.NotNil(t, info)

	labels, _, err := c.ListLabelTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, labels)
	li, _, err := c.GetLabelTemplateInfo(labels[0])
	require.NoError(t, err)
	assert.NotEmpty(t, li)

	lics, _, err := c.ListLicenseTemplates()
	require.NoError(t, err)
	require.NotEmpty(t, lics)
	lci, _, err := c.GetLicenseTemplateInfo(lics[0].Name)
	require.NoError(t, err)
	assert.NotNil(t, lci)
}
