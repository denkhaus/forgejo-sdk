// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"net/http"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOrgRunnerRegistrationToken(t *testing.T) {
	log.Println("== TestGetOrgRunnerRegistrationToken ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{
		Name:     "test-runner-org",
		FullName: "Test Runner Org",
	})
	require.NoError(t, err)
	assert.NotNil(t, org)

	// Get org runner registration token
	token, resp, err := c.GetOrgRunnerRegistrationToken(org.UserName)
	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NotEmpty(t, token.Token)
}

func TestSearchOrgRunnerJobs(t *testing.T) {
	log.Println("== TestSearchOrgRunnerJobs ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{
		Name:     "test-runner-org-search",
		FullName: "Test Runner Org Search",
	})
	require.NoError(t, err)
	assert.NotNil(t, org)

	// Search org runner jobs (may be empty for new org)
	jobs, _, err := c.SearchOrgRunnerJobs(org.UserName, SearchOrgRunnerJobsOption{})
	require.NoError(t, err)
	assert.NotNil(t, jobs)
}

// TestOrgRunnersV15 exercises org-level v15 interactive runner registration+management.
func TestOrgRunnersV15(t *testing.T) {
	log.Println("== TestOrgRunnersV15 ==")
	c := newTestClient()

	org, _, err := c.CreateOrg(CreateOrgOption{Name: "v15-org-runners"})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteOrg("v15-org-runners") }()

	orr, _, err := c.RegisterOrgRunner(org.UserName, models.RegisterRunnerOptions{Name: OptionalString("org-runner"), Ephemeral: true})
	require.NoError(t, err)
	require.NotZero(t, orr.ID)
	assert.NotEmpty(t, orr.Token)
	assert.NotEmpty(t, orr.UUID)

	_, _, err = c.GetOrgRunner(org.UserName, orr.ID)
	require.NoError(t, err)
	_, _, err = c.ListOrgRunners(org.UserName, ListActionRunnersOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
	_, err = c.DeleteOrgRunner(org.UserName, orr.ID)
	require.NoError(t, err)
}
