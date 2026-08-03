// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2021 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"strings"
	"testing"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestOrgTeams(t *testing.T, c *Client, org, name string, accessMode AccessMode, units []RepoUnitType) (*models.Team, error) {
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

// createTestOrgWithTeam sets up a fresh org + team and returns them, cleaning
// up the org (which cascades to the team) on test completion.
func createTestOrgWithTeam(t *testing.T, c *Client, orgName, teamName string) (*models.Team, func()) {
	t.Helper()
	_, _ = c.DeleteOrg(orgName)
	_, _, err := c.CreateOrg(CreateOrgOption{
		Name:                      orgName,
		Visibility:                VisibleTypePublic,
		RepoAdminChangeTeamAccess: true,
	})
	require.NoError(t, err)
	team, err := createTestOrgTeams(t, c, orgName, teamName, AccessModeAdmin,
		[]RepoUnitType{RepoUnitCode, RepoUnitIssues, RepoUnitPulls, RepoUnitReleases})
	require.NoError(t, err)
	return team, func() { _, _ = c.DeleteOrg(orgName) }
}

// Test_TeamLifecycle exercises create (via helper), list, get, edit and delete.
func Test_TeamLifecycle(t *testing.T) {
	log.Println("== Test_TeamLifecycle ==")
	c := newTestClient()

	orgName := "TestTeamLifecycleOrg"
	team, cleanup := createTestOrgWithTeam(t, c, orgName, "LifecycleTeam")
	defer cleanup()

	// ListOrgTeams finds the freshly created team.
	teams, _, err := c.ListOrgTeams(orgName, ListTeamsOptions{})
	require.NoError(t, err)
	require.NotEmpty(t, teams)

	// GetTeam returns the same team by id.
	got, _, err := c.GetTeam(team.ID)
	require.NoError(t, err)
	assert.Equal(t, team.ID, got.ID)

	// EditTeam renames the team; GetTeam reflects the change.
	desc := "updated team description"
	_, err = c.EditTeam(team.ID, EditTeamOption{
		Name:                    "LifecycleTeamRenamed",
		Description:             &desc,
		Permission:              AccessModeAdmin,
		CanCreateOrgRepo:        OptionalBool(false),
		IncludesAllRepositories: OptionalBool(false),
		Units:                   []RepoUnitType{RepoUnitCode, RepoUnitIssues},
	})
	require.NoError(t, err)
	got, _, err = c.GetTeam(team.ID)
	require.NoError(t, err)
	assert.Equal(t, "LifecycleTeamRenamed", got.Name)

	// ListMyTeams returns the authenticated user's teams (the owner belongs to
	// at least the Owners team). Assert it works and is non-empty.
	myTeams, _, err := c.ListMyTeams(&ListTeamsOptions{})
	require.NoError(t, err)
	assert.NotEmpty(t, myTeams)

	// DeleteTeam removes the team; a second delete is an error.
	_, err = c.DeleteTeam(team.ID)
	require.NoError(t, err)
	_, err = c.DeleteTeam(team.ID)
	require.Error(t, err)
}

// Test_TeamMembers exercises add/get/list/remove of team members. The org owner
// (the authenticated admin user) is already an org member, so it can be added.
func Test_TeamMembers(t *testing.T) {
	log.Println("== Test_TeamMembers ==")
	c := newTestClient()
	user, _, err := c.GetMyUserInfo()
	require.NoError(t, err)

	orgName := "TestTeamMembersOrg"
	team, cleanup := createTestOrgWithTeam(t, c, orgName, "MembersTeam")
	defer cleanup()

	_, err = c.AddTeamMember(team.ID, user.UserName)
	require.NoError(t, err)

	member, _, err := c.GetTeamMember(team.ID, user.UserName)
	require.NoError(t, err)
	assert.Equal(t, user.UserName, member.UserName)

	members, _, err := c.ListTeamMembers(team.ID, ListTeamMembersOptions{})
	require.NoError(t, err)
	assert.NotEmpty(t, members)

	_, err = c.RemoveTeamMember(team.ID, user.UserName)
	require.NoError(t, err)
	members, _, err = c.ListTeamMembers(team.ID, ListTeamMembersOptions{})
	require.NoError(t, err)
	assert.Empty(t, members)
}

// Test_TeamRepositories exercises add/list/remove of team repositories. The
// repository must belong to the org, so it is created via CreateOrgRepo.
func Test_TeamRepositories(t *testing.T) {
	log.Println("== Test_TeamRepositories ==")
	c := newTestClient()

	orgName := "TestTeamReposOrg"
	repoName := "team-repo-test"
	// thorough pre-clean: a repo under the org blocks org deletion, so remove
	// the repo first, then the org, before re-creating.
	_, _ = c.DeleteRepo(orgName, repoName)
	team, cleanup := createTestOrgWithTeam(t, c, orgName, "ReposTeam")
	// on cleanup, delete the repo before the org (an org must be empty to delete).
	defer func() {
		_, _ = c.DeleteRepo(orgName, repoName)
		cleanup()
	}()

	_, _, err := c.CreateOrgRepo(orgName, CreateRepoOption{Name: repoName})
	require.NoError(t, err)

	_, err = c.AddTeamRepository(team.ID, orgName, repoName)
	require.NoError(t, err)

	repos, _, err := c.ListTeamRepositories(team.ID, ListTeamRepositoriesOptions{})
	require.NoError(t, err)
	assert.NotEmpty(t, repos)

	_, err = c.RemoveTeamRepository(team.ID, orgName, repoName)
	require.NoError(t, err)
	repos, _, err = c.ListTeamRepositories(team.ID, ListTeamRepositoriesOptions{})
	require.NoError(t, err)
	assert.Empty(t, repos)
}

// Test_TeamNegative covers not-found error paths and client-side validation.
func Test_TeamNegative(t *testing.T) {
	log.Println("== Test_TeamNegative ==")
	c := newTestClient()

	// non-existent team ids -> error
	_, _, err := c.GetTeam(999999)
	require.Error(t, err)
	_, _, err = c.ListTeamMembers(999999, ListTeamMembersOptions{})
	require.Error(t, err)

	// CreateTeamOption.Validate (table-driven)
	t.Run("CreateTeamOption.Validate", func(t *testing.T) {
		longName := strings.Repeat("a", 256)
		cases := []struct {
			name    string
			opt     CreateTeamOption
			wantErr bool
		}{
			{"empty name", CreateTeamOption{Permission: AccessModeRead}, true},
			{"invalid permission", CreateTeamOption{Name: "x", Permission: AccessMode("bogus")}, true},
			{"name too long", CreateTeamOption{Name: longName, Permission: AccessModeAdmin}, true},
			{"owner normalized to admin", CreateTeamOption{Name: "x", Permission: AccessModeOwner}, false},
			{"valid admin", CreateTeamOption{Name: "x", Permission: AccessModeAdmin}, false},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.opt.Validate()
				if tc.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})

	// EditTeamOption.Validate (table-driven)
	t.Run("EditTeamOption.Validate", func(t *testing.T) {
		longName := strings.Repeat("a", 31) // limit is 30
		tooLongDesc := strings.Repeat("d", 256)
		cases := []struct {
			name    string
			opt     EditTeamOption
			wantErr bool
		}{
			{"empty name", EditTeamOption{Permission: AccessModeRead}, true},
			{"name too long (>30)", EditTeamOption{Name: longName, Permission: AccessModeAdmin}, true},
			{"description too long", EditTeamOption{Name: "x", Permission: AccessModeAdmin, Description: &tooLongDesc}, true},
			{"valid", EditTeamOption{Name: "x", Permission: AccessModeWrite}, false},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.opt.Validate()
				if tc.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			})
		}
	})
}
