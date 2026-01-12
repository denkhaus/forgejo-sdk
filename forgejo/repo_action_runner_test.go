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

func TestListActionRuns(t *testing.T) {
	log.Println("== TestListActionRuns ==")
	c := newTestClient()

	user := createTestUser(t, "action_runner_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-action-runs",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// List action runs (may be empty for new repo)
	runs, _, err := c.ListActionRuns(newRepo.Owner.UserName, newRepo.Name, ListActionRunsOption{})
	require.NoError(t, err)
	assert.NotNil(t, runs)
}

func TestGetActionRun(t *testing.T) {
	log.Println("== TestGetActionRun ==")
	c := newTestClient()

	user := createTestUser(t, "action_runner_user2", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-action-run",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// Try to get a non-existent run (should return 404 or error)
	_, resp, err := c.GetActionRun(newRepo.Owner.UserName, newRepo.Name, 999999)
	require.Error(t, err)
	assert.NotNil(t, resp)
}

func TestListActionTasks(t *testing.T) {
	log.Println("== TestListActionTasks ==")
	c := newTestClient()

	user := createTestUser(t, "action_runner_user3", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-action-tasks",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// List action tasks (may be empty for new repo)
	tasks, _, err := c.ListActionTasks(newRepo.Owner.UserName, newRepo.Name, ListActionTasksOption{})
	require.NoError(t, err)
	assert.NotNil(t, tasks)
}

func TestSearchRunnerJobs(t *testing.T) {
	log.Println("== TestSearchRunnerJobs ==")
	c := newTestClient()

	user := createTestUser(t, "action_runner_user4", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-runner-jobs",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// Search runner jobs (may be empty for new repo)
	jobs, _, err := c.SearchRunnerJobs(newRepo.Owner.UserName, newRepo.Name, SearchRunnerJobsOption{})
	require.NoError(t, err)
	assert.NotNil(t, jobs)
}

func TestListRunnerJobs(t *testing.T) {
	log.Println("== TestListRunnerJobs ==")
	c := newTestClient()

	user := createTestUser(t, "action_runner_user_list", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-list-runner-jobs",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// List runner jobs (may be empty for new repo)
	jobs, _, err := c.ListRunnerJobs(newRepo.Owner.UserName, newRepo.Name, ListRunnerJobsOption{})
	require.NoError(t, err)
	assert.NotNil(t, jobs)
}

func TestGetRepoRunnerRegistrationToken(t *testing.T) {
	log.Println("== TestGetRepoRunnerRegistrationToken ==")
	c := newTestClient()

	user := createTestUser(t, "action_runner_user5", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-runner-token",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// Get runner registration token
	token, resp, err := c.GetRepoRunnerRegistrationToken(newRepo.Owner.UserName, newRepo.Name)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, token.Token)
}
