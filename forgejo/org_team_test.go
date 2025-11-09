// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestOrgTeams(t *testing.T, c *Client, org, name string, accessMode AccessMode, units []RepoUnitType) (*Team, error) {
	team, _, e := c.CreateTeam(org, CreateTeamOption{
		Name:                    name,
		Description:             name + "'s team desc",
		Permission:              accessMode,
		CanCreateOrgRepo:        false,
		IncludesAllRepositories: false,
		Units:                   units,
	})
	require.NoError(t, e)
	assert.NotNil(t, team)
	return team, e
}

func TestTeamSearch(t *testing.T) {
	log.Println("== TestTeamSearch ==")
	c := newTestClient()

	orgName := "TestTeamsOrg"
	// prepare for test
	_, _, err := c.CreateOrg(CreateOrgOption{
		Name:                      orgName,
		Visibility:                VisibleTypePublic,
		RepoAdminChangeTeamAccess: true,
	})
	defer func() {
		_, _ = c.DeleteOrg(orgName)
	}()

	require.NoError(t, err)

	if _, err = createTestOrgTeams(t, c, orgName, "Admins", AccessModeAdmin, []RepoUnitType{RepoUnitCode, RepoUnitIssues, RepoUnitPulls, RepoUnitReleases}); err != nil {
		return
	}

	teams, _, err := c.SearchOrgTeams(orgName, &SearchTeamsOptions{
		Query: "Admins",
	})
	require.NoError(t, err)
	if assert.Len(t, teams, 1) {
		assert.Equal(t, "Admins", teams[0].Name)
	}
}

func TestCreateTeamWithUnitsMap(t *testing.T) {
	log.Println("== TestCreateTeamWithUnitsMap ==")
	c := newTestClient()

	orgName := "TestUnitsMapOrg"
	// Create test organization
	_, _, err := c.CreateOrg(CreateOrgOption{
		Name:                      orgName,
		Visibility:                VisibleTypePublic,
		RepoAdminChangeTeamAccess: true,
	})
	defer func() {
		_, _ = c.DeleteOrg(orgName)
	}()

	require.NoError(t, err)

	// Create team with per-unit permissions using UnitsMap
	unitsMap := map[string]string{
		"repo.code":    "write",
		"repo.issues":  "write",
		"repo.pulls":   "read",
		"repo.wiki":    "read",
		"repo.actions": "none",
	}

	team, _, err := c.CreateTeam(orgName, CreateTeamOption{
		Name:                    "developers",
		Description:             "Developers team with per-unit permissions",
		Permission:              AccessModeWrite,
		CanCreateOrgRepo:        false,
		IncludesAllRepositories: false,
		UnitsMap:                unitsMap,
	})

	require.NoError(t, err)
	assert.NotNil(t, team)
	assert.Equal(t, "developers", team.Name)
	assert.Equal(t, unitsMap, team.UnitsMap)
}

func TestEditTeamWithUnitsMap(t *testing.T) {
	log.Println("== TestEditTeamWithUnitsMap ==")
	c := newTestClient()

	orgName := "TestEditUnitsMapOrg"
	// Create test organization
	_, _, err := c.CreateOrg(CreateOrgOption{
		Name:                      orgName,
		Visibility:                VisibleTypePublic,
		RepoAdminChangeTeamAccess: true,
	})
	defer func() {
		_, _ = c.DeleteOrg(orgName)
	}()

	require.NoError(t, err)

	// Create team first with UnitsMap for consistency
	initialTeam, _, err := c.CreateTeam(orgName, CreateTeamOption{
		Name:                    "readers",
		Description:             "Initial readers team",
		Permission:              AccessModeRead,
		CanCreateOrgRepo:        false,
		IncludesAllRepositories: false,
		UnitsMap: map[string]string{
			"repo.code": "read",
		},
	})

	require.NoError(t, err)
	assert.NotNil(t, initialTeam)

	// Edit team to add per-unit permissions
	updatedUnitsMap := map[string]string{
		"repo.code":   "read",
		"repo.issues": "write",
		"repo.pulls":  "none",
	}

	_, err = c.EditTeam(initialTeam.ID, EditTeamOption{
		Name:       "readers",
		Permission: AccessModeRead,
		UnitsMap:   updatedUnitsMap,
	})

	require.NoError(t, err)

	// Verify the team was updated
	updatedTeam, _, err := c.GetTeam(initialTeam.ID)
	require.NoError(t, err)
	assert.NotNil(t, updatedTeam)
	assert.Equal(t, updatedUnitsMap, updatedTeam.UnitsMap)
}

func TestUnitsMapSerialization(t *testing.T) {
	log.Println("== TestUnitsMapSerialization ==")

	// Test CreateTeamOption serialization
	opt := CreateTeamOption{
		Name:        "test-team",
		Description: "Test team",
		Permission:  AccessModeWrite,
		UnitsMap: map[string]string{
			"repo.code":   "write",
			"repo.issues": "read",
		},
	}

	// Verify UnitsMap field is set correctly
	assert.NotNil(t, opt.UnitsMap)
	assert.Len(t, opt.UnitsMap, 2)
	assert.Equal(t, "write", opt.UnitsMap["repo.code"])
	assert.Equal(t, "read", opt.UnitsMap["repo.issues"])

	// Test EditTeamOption serialization
	editOpt := EditTeamOption{
		Name: "test-team",
		UnitsMap: map[string]string{
			"repo.code": "none",
		},
	}

	assert.NotNil(t, editOpt.UnitsMap)
	assert.Len(t, editOpt.UnitsMap, 1)
	assert.Equal(t, "none", editOpt.UnitsMap["repo.code"])

	// Test Team struct
	team := &Team{
		ID:   1,
		Name: "test-team",
		UnitsMap: map[string]string{
			"repo.wiki":     "read",
			"repo.releases": "write",
		},
	}

	assert.NotNil(t, team.UnitsMap)
	assert.Len(t, team.UnitsMap, 2)
	assert.Equal(t, "read", team.UnitsMap["repo.wiki"])
	assert.Equal(t, "write", team.UnitsMap["repo.releases"])
}

func TestUnitsMapBackwardCompatibility(t *testing.T) {
	log.Println("== TestUnitsMapBackwardCompatibility ==")

	// Test that teams without UnitsMap still work
	opt := CreateTeamOption{
		Name:                    "legacy-team",
		Permission:              AccessModeRead,
		CanCreateOrgRepo:        true,
		IncludesAllRepositories: true,
		Units:                   []RepoUnitType{RepoUnitCode, RepoUnitIssues},
	}

	// UnitsMap should be nil (not set)
	assert.Nil(t, opt.UnitsMap)

	// Verify other fields are still accessible
	assert.Equal(t, "legacy-team", opt.Name)
	assert.Equal(t, AccessModeRead, opt.Permission)
	assert.True(t, opt.CanCreateOrgRepo)
	assert.True(t, opt.IncludesAllRepositories)
	assert.Len(t, opt.Units, 2)

	// Test EditTeamOption without UnitsMap
	editOpt := EditTeamOption{
		Name:  "legacy-team",
		Units: []RepoUnitType{RepoUnitCode},
	}

	assert.Nil(t, editOpt.UnitsMap)
	assert.Len(t, editOpt.Units, 1)
}
