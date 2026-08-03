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

func TestRepoActionRuns(t *testing.T) {
	log.Println("== TestRepoActionRuns ==")
	c := newTestClient()

	user := createTestUser(t, "repo_action_run_user", c)
	c.SetSudo(user.UserName)

	// pre-cleanup then create a fresh repo
	_, _ = c.DeleteRepo(user.UserName, "runs-test")
	repo, _, err := c.CreateRepo(CreateRepoOption{Name: "runs-test"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteRepo(user.UserName, repo.Name) }()

	owner, name := user.UserName, repo.Name

	// jobs of a non-existent run -> error
	_, _, err = c.ListActionRunJobs(owner, name, 999999)
	require.Error(t, err)

	// delete / cancel a non-existent run -> error
	_, err = c.DeleteActionRun(owner, name, 999999)
	require.Error(t, err)

	_, err = c.CancelActionRun(owner, name, 999999)
	require.Error(t, err)

	// logs of a non-existent run / job -> error (close any returned reader)
	runLogs, _, err := c.GetActionRunLogs(owner, name, 999999)
	if runLogs != nil {
		runLogs.Close()
	}
	require.Error(t, err)

	jobLogs, _, err := c.GetActionJobLogs(owner, name, 999999, GetActionJobLogsOption{Attempt: 1})
	if jobLogs != nil {
		jobLogs.Close()
	}
	require.Error(t, err)

	// without an attempt filter the request is still well-formed and fails
	jobLogs2, _, err := c.GetActionJobLogs(owner, name, 999999, GetActionJobLogsOption{})
	if jobLogs2 != nil {
		jobLogs2.Close()
	}
	assert.Error(t, err)
}
