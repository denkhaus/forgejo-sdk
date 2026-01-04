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

// TestIssueDependencies tests issue dependency CRUD operations
func TestIssueDependencies(t *testing.T) {
	log.Println("== TestIssueDependencies ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "IssueDepsTestsRepo", c)
	require.NoError(t, err)

	// Create test issues
	issue1, _, err := c.CreateIssue(repo.Owner.UserName, repo.Name, CreateIssueOption{
		Title: "Test Issue 1",
		Body:  "This issue will have dependencies",
	})
	require.NoError(t, err)

	issue2, _, err := c.CreateIssue(repo.Owner.UserName, repo.Name, CreateIssueOption{
		Title: "Test Issue 2",
		Body:  "Blocking issue",
	})
	require.NoError(t, err)

	issue3, _, err := c.CreateIssue(repo.Owner.UserName, repo.Name, CreateIssueOption{
		Title: "Test Issue 3",
		Body:  "Another blocking issue",
	})
	require.NoError(t, err)

	// Test 1: Create dependency
	resp, err := c.CreateIssueDependency(repo.Owner.UserName, repo.Name, issue1.Index, CreateIssueDependencyOption{
		NewDependency: issue2.Index,
	})
	require.NoError(t, err)
	assert.NotNil(t, resp)

	// Test 2: List dependencies
	deps, _, err := c.ListIssueDependencies(repo.Owner.UserName, repo.Name, issue1.Index, ListDependenciesOptions{})
	require.NoError(t, err)
	assert.Len(t, deps, 1)
	assert.EqualValues(t, issue2.Index, deps[0].Index)

	// Test 3: Add second dependency
	resp, err = c.CreateIssueDependency(repo.Owner.UserName, repo.Name, issue1.Index, CreateIssueDependencyOption{
		NewDependency: issue3.Index,
	})
	require.NoError(t, err)

	deps, _, err = c.ListIssueDependencies(repo.Owner.UserName, repo.Name, issue1.Index, ListDependenciesOptions{})
	require.NoError(t, err)
	assert.Len(t, deps, 2)

	// Test 4: List blocked issues (issues blocked BY this issue)
	// Issue 1 doesn't block any other issues, so this should be empty
	blocked, _, err := c.ListBlockedIssues(repo.Owner.UserName, repo.Name, issue1.Index)
	require.NoError(t, err)
	assert.Len(t, blocked, 0)

	// Test 5: List blocking issues (issues that block THIS issue)
	// This is an alias for ListIssueDependencies
	blocking, _, err := c.ListBlockingIssues(repo.Owner.UserName, repo.Name, issue1.Index)
	require.NoError(t, err)
	assert.Len(t, blocking, 2)
	// Issue 2 should be in the list
	assert.EqualValues(t, issue2.Index, blocking[0].Index)

	// Test 6: Remove dependency
	resp, err = c.RemoveIssueDependency(repo.Owner.UserName, repo.Name, issue1.Index, issue2.Index)
	require.NoError(t, err)

	deps, _, err = c.ListIssueDependencies(repo.Owner.UserName, repo.Name, issue1.Index, ListDependenciesOptions{})
	require.NoError(t, err)
	assert.Len(t, deps, 1)
	assert.EqualValues(t, issue3.Index, deps[0].Index)

	// Test 7: Validation - invalid issue number
	opt := CreateIssueDependencyOption{NewDependency: -1}
	err = opt.Validate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "positive")

	// Test 8: Remove non-existent dependency
	resp, err = c.RemoveIssueDependency(repo.Owner.UserName, repo.Name, issue1.Index, 99999)
	require.Error(t, err) // Should return error from server
}

// TestIssueDependenciesEmpty tests listing dependencies for issue with none
func TestIssueDependenciesEmpty(t *testing.T) {
	log.Println("== TestIssueDependenciesEmpty ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "IssueDepsEmptyTestsRepo", c)
	require.NoError(t, err)

	issue, _, err := c.CreateIssue(repo.Owner.UserName, repo.Name, CreateIssueOption{
		Title: "Test Issue",
	})
	require.NoError(t, err)

	// Should return empty list, not error
	deps, _, err := c.ListIssueDependencies(repo.Owner.UserName, repo.Name, issue.Index, ListDependenciesOptions{})
	require.NoError(t, err)
	assert.Empty(t, deps)

	// ListBlockedIssues should also return empty
	blocked, _, err := c.ListBlockedIssues(repo.Owner.UserName, repo.Name, issue.Index)
	require.NoError(t, err)
	assert.Empty(t, blocked)

	// ListBlockingIssues should also return empty
	blocking, _, err := c.ListBlockingIssues(repo.Owner.UserName, repo.Name, issue.Index)
	require.NoError(t, err)
	assert.Empty(t, blocking)
}

// TestIssueDependencyValidation tests validation logic
func TestIssueDependencyValidation(t *testing.T) {
	log.Println("== TestIssueDependencyValidation ==")

	tests := []struct {
		name        string
		option      CreateIssueDependencyOption
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid dependency",
			option:  CreateIssueDependencyOption{NewDependency: 5},
			wantErr: false,
		},
		{
			name:        "zero dependency",
			option:      CreateIssueDependencyOption{NewDependency: 0},
			wantErr:     true,
			errContains: "positive",
		},
		{
			name:        "negative dependency",
			option:      CreateIssueDependencyOption{NewDependency: -1},
			wantErr:     true,
			errContains: "positive",
		},
		{
			name:        "negative large dependency",
			option:      CreateIssueDependencyOption{NewDependency: -999},
			wantErr:     true,
			errContains: "positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.option.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
