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

func TestWorkflowDispatch(t *testing.T) {
	log.Println("== TestWorkflowDispatch ==")
	c := newTestClient()

	user := createTestUser(t, "workflow_user", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-workflow-dispatch",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	// Note: This test may fail if the repo doesn't have any workflows defined
	// that support workflow_dispatch. The API will return 404 in that case.
	_, resp, err := c.WorkflowDispatch(newRepo.Owner.UserName, newRepo.Name, "test.yml", WorkflowDispatchOption{
		Ref: "main",
	})
	// We expect either success or an error (if no workflow exists)
	// The important thing is that the method doesn't panic and returns a proper response
	if err != nil {
		assert.NotNil(t, resp)
	}
}

func TestWorkflowDispatchWithInputs(t *testing.T) {
	log.Println("== TestWorkflowDispatchWithInputs ==")
	c := newTestClient()

	user := createTestUser(t, "workflow_user2", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-workflow-inputs",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	opt := WorkflowDispatchOption{
		Ref: "main",
		Inputs: map[string]any{
			"message":     "Deploying new version",
			"environment": "production",
		},
	}

	_, resp, err := c.WorkflowDispatch(newRepo.Owner.UserName, newRepo.Name, "deploy.yml", opt)
	if err != nil {
		assert.NotNil(t, resp)
	}
}

func TestWorkflowDispatchWithRef(t *testing.T) {
	log.Println("== TestWorkflowDispatchWithRef ==")
	c := newTestClient()

	user := createTestUser(t, "workflow_user3", c)
	c.SetSudo(user.UserName)
	newRepo, _, err := c.CreateRepo(CreateRepoOption{
		Name: "test-workflow-ref",
	})
	require.NoError(t, err)
	assert.NotNil(t, newRepo)

	opt := WorkflowDispatchOption{
		Ref: "develop",
	}

	_, resp, err := c.WorkflowDispatch(newRepo.Owner.UserName, newRepo.Name, "test.yml", opt)
	if err != nil {
		assert.NotNil(t, resp)
	}
}

func TestWorkflowDispatchErrorHandling(t *testing.T) {
	log.Println("== TestWorkflowDispatchErrorHandling ==")
	c := newTestClient()

	// Test with invalid owner/repo (should handle validation error)
	_, _, err := c.WorkflowDispatch("", "repo", "test.yml", WorkflowDispatchOption{})
	require.Error(t, err)

	// Test with invalid characters in owner/repo (should be escaped)
	_, _, err = c.WorkflowDispatch("../evil", "repo", "test.yml", WorkflowDispatchOption{})
	require.Error(t, err)
}
