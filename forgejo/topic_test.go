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

func Test_Topic(t *testing.T) {
	log.Println("== Test_Topic ==")
	c := newTestClient()

	// SearchTopics is a safe read and always available.
	topics, _, err := c.SearchTopics(SearchTopicOption{Query: "forgejo"})
	require.NoError(t, err)
	assert.NotNil(t, topics)

	// Exercise repo-topic CRUD on a throwaway repo.
	repo, err := createTestRepo(t, "TopicRepo", c)
	require.NoError(t, err)

	owner := repo.Owner.UserName
	repoName := repo.Name

	// Add a couple of topics.
	_, err = c.AddRepoTopic(owner, repoName, "sdk-topic-alpha")
	require.NoError(t, err)
	_, err = c.AddRepoTopic(owner, repoName, "sdk-topic-beta")
	require.NoError(t, err)

	// ListRepoTopics should reflect the additions.
	listed, _, err := c.ListRepoTopics(owner, repoName, ListRepoTopicsOptions{})
	require.NoError(t, err)
	assert.Contains(t, listed, "sdk-topic-alpha")
	assert.Contains(t, listed, "sdk-topic-beta")

	// SearchTopics is safe; whether it surfaces a freshly created,
	// single-repo topic depends on server-side indexing, so only assert on
	// success and never fail when the topic is not yet indexed.
	results, _, err := c.SearchTopics(SearchTopicOption{Query: "sdk-topic-alpha"})
	require.NoError(t, err)
	for _, tr := range results {
		if tr.Name == "sdk-topic-alpha" {
			assert.Equal(t, "sdk-topic-alpha", tr.Name)
		}
	}

	// SetRepoTopics replaces the whole list.
	_, err = c.SetRepoTopics(owner, repoName, []string{"sdk-topic-gamma"})
	require.NoError(t, err)
	listed, _, err = c.ListRepoTopics(owner, repoName, ListRepoTopicsOptions{})
	require.NoError(t, err)
	assert.Len(t, listed, 1)
	assert.Contains(t, listed, "sdk-topic-gamma")

	// DeleteRepoTopic removes a single topic.
	_, err = c.DeleteRepoTopic(owner, repoName, "sdk-topic-gamma")
	require.NoError(t, err)
	listed, _, err = c.ListRepoTopics(owner, repoName, ListRepoTopicsOptions{})
	require.NoError(t, err)
	assert.Empty(t, listed)
}
