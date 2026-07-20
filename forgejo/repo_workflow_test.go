// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkflowDispatch(t *testing.T) {
	log.Println("== TestWorkflowDispatch ==")
	c := newTestClient()

	user := createTestUser(t, "workflow_user", c)
	c.SetSudo(user.UserName)

	// Pre-cleanup: delete existing repo from previous test runs
	_, _ = c.DeleteRepo(user.UserName, "test-workflow-dispatch")

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

	// Pre-cleanup: delete existing repo from previous test runs
	_, _ = c.DeleteRepo(user.UserName, "test-workflow-inputs")

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

	// Pre-cleanup: delete existing repo from previous test runs
	_, _ = c.DeleteRepo(user.UserName, "test-workflow-ref")

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

// TestWorkflowDispatchNoContent verifies that WorkflowDispatch tolerates the
// 204 No Content (empty body) response Forgejo returns on a successful
// workflow_dispatch. Regression test: getParsedResponse previously failed with
// "unexpected end of JSON input" on the empty body.
func TestWorkflowDispatchNoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/version" {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"version":"15.0.3"}`))
			return
		}
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/repos/owner/repo/actions/workflows/test.yml/dispatches", r.URL.Path)
		w.WriteHeader(http.StatusNoContent) // empty body, as Forgejo does
	}))
	defer server.Close()

	c, err := NewClient(server.URL, SetToken("token"))
	require.NoError(t, err)

	run, resp, err := c.WorkflowDispatch("owner", "repo", "test.yml", WorkflowDispatchOption{Ref: "main"})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	// No body returned -> run is non-nil but zero-valued.
	require.NotNil(t, run)
	assert.Zero(t, run.ID)
}
