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

// quotaEnabled skips the test when the quota feature is disabled on the test
// instance (the /admin/quota endpoints then return 404). It is a safety net:
// the test instance app.ini sets [quota] ENABLED = true, but other environments
// may run without quota.
func quotaEnabled(t *testing.T, c *Client) {
	t.Helper()
	if _, _, err := c.ListQuotaGroups(); err != nil {
		t.Skipf("quota disabled on this instance: %v", err)
	}
}

func Test_AdminQuotaGroups(t *testing.T) {
	log.Println("== Test_AdminQuotaGroups ==")
	c := newTestClient()
	quotaEnabled(t, c)

	ruleName := "test-size-all-rule"
	groupName := "test-quota-group"
	// pre-clean any leftovers from a previous run
	_, _ = c.DeleteQuotaGroup(groupName)
	_, _ = c.DeleteQuotaRule(ruleName)

	// a rule to attach to the group
	rule, _, err := c.CreateQuotaRule(models.CreateQuotaRuleOptions{
		Name:     ruleName,
		Limit:    1024,
		Subjects: []string{"size:all"},
	})
	require.NoError(t, err)
	require.NotNil(t, rule)
	defer func() { _, _ = c.DeleteQuotaRule(ruleName) }()

	// create + get + list
	group, _, err := c.CreateQuotaGroup(models.CreateQuotaGroupOptions{Name: groupName})
	require.NoError(t, err)
	if assert.NotNil(t, group) {
		assert.Equal(t, groupName, group.Name)
	}
	defer func() { _, _ = c.DeleteQuotaGroup(groupName) }()

	got, _, err := c.GetQuotaGroup(groupName)
	require.NoError(t, err)
	assert.Equal(t, groupName, got.Name)

	groups, _, err := c.ListQuotaGroups()
	require.NoError(t, err)
	assert.NotEmpty(t, groups)

	// attach the rule and verify it shows up on the group
	_, err = c.AddRuleToQuotaGroup(groupName, ruleName)
	require.NoError(t, err)
	got, _, err = c.GetQuotaGroup(groupName)
	require.NoError(t, err)
	if assert.Len(t, got.Rules, 1) {
		assert.Equal(t, ruleName, got.Rules[0].Name)
	}

	// user membership add/list/remove
	user := createTestUser(t, "quota_group_user", c)
	_, err = c.AddUserToQuotaGroup(groupName, user.UserName)
	require.NoError(t, err)
	users, _, err := c.ListQuotaGroupUsers(groupName)
	require.NoError(t, err)
	assert.NotEmpty(t, users)

	_, err = c.RemoveUserFromQuotaGroup(groupName, user.UserName)
	require.NoError(t, err)
	users, _, err = c.ListQuotaGroupUsers(groupName)
	require.NoError(t, err)
	assert.Empty(t, users)

	// detach the rule again
	_, err = c.RemoveRuleFromQuotaGroup(groupName, ruleName)
	require.NoError(t, err)
	got, _, err = c.GetQuotaGroup(groupName)
	require.NoError(t, err)
	assert.Empty(t, got.Rules)

	// delete the group; a follow-up get errors
	_, err = c.DeleteQuotaGroup(groupName)
	require.NoError(t, err)
	_, _, err = c.GetQuotaGroup(groupName)
	require.Error(t, err)
}

func Test_AdminQuotaRules(t *testing.T) {
	log.Println("== Test_AdminQuotaRules ==")
	c := newTestClient()
	quotaEnabled(t, c)

	ruleName := "test-crud-rule"
	_, _ = c.DeleteQuotaRule(ruleName)

	// create + get + list
	rule, _, err := c.CreateQuotaRule(models.CreateQuotaRuleOptions{
		Name:     ruleName,
		Limit:    2048,
		Subjects: []string{"size:all"},
	})
	require.NoError(t, err)
	if assert.NotNil(t, rule) {
		assert.Equal(t, ruleName, rule.Name)
		assert.EqualValues(t, 2048, rule.Limit)
	}
	defer func() { _, _ = c.DeleteQuotaRule(ruleName) }()

	got, _, err := c.GetQuotaRule(ruleName)
	require.NoError(t, err)
	assert.Equal(t, ruleName, got.Name)

	rules, _, err := c.ListQuotaRules()
	require.NoError(t, err)
	assert.NotEmpty(t, rules)

	// edit the rule's limit and verify
	edited, _, err := c.EditQuotaRule(ruleName, models.EditQuotaRuleOptions{
		Limit:    4096,
		Subjects: []string{"size:all"},
	})
	require.NoError(t, err)
	if assert.NotNil(t, edited) {
		assert.EqualValues(t, 4096, edited.Limit)
	}

	// delete the rule; a follow-up get errors
	_, err = c.DeleteQuotaRule(ruleName)
	require.NoError(t, err)
	_, _, err = c.GetQuotaRule(ruleName)
	require.Error(t, err)
}

func Test_AdminQuotaUser(t *testing.T) {
	log.Println("== Test_AdminQuotaUser ==")
	c := newTestClient()
	quotaEnabled(t, c)

	user := createTestUser(t, "quota_user_user", c)
	groupName := "test-user-quota-group"
	_, _ = c.DeleteQuotaGroup(groupName)
	_, _, err := c.CreateQuotaGroup(models.CreateQuotaGroupOptions{Name: groupName})
	require.NoError(t, err)
	defer func() { _, _ = c.DeleteQuotaGroup(groupName) }()

	// GetUserQuota returns quota info (no groups assigned yet)
	info, _, err := c.GetUserQuota(user.UserName)
	require.NoError(t, err)
	require.NotNil(t, info)

	// assign the group, then the user's quota lists it
	_, err = c.SetUserQuotaGroups(user.UserName, models.SetUserQuotaGroupsOptions{
		Groups: []string{groupName},
	})
	require.NoError(t, err)
	info, _, err = c.GetUserQuota(user.UserName)
	require.NoError(t, err)
	require.Len(t, info.Groups, 1)
	assert.Equal(t, groupName, info.Groups[0].Name)

	// clear the assignment
	_, err = c.SetUserQuotaGroups(user.UserName, models.SetUserQuotaGroupsOptions{
		Groups: []string{},
	})
	require.NoError(t, err)
	info, _, err = c.GetUserQuota(user.UserName)
	require.NoError(t, err)
	assert.Empty(t, info.Groups)
}

func Test_AdminQuotaNegative(t *testing.T) {
	log.Println("== Test_AdminQuotaNegative ==")
	c := newTestClient()
	quotaEnabled(t, c)

	// non-existent group -> get/delete error
	_, _, err := c.GetQuotaGroup("no-such-quota-group")
	require.Error(t, err)
	_, err = c.DeleteQuotaGroup("no-such-quota-group")
	require.Error(t, err)

	// non-existent rule -> get/delete error
	_, _, err = c.GetQuotaRule("no-such-quota-rule")
	require.Error(t, err)
	_, err = c.DeleteQuotaRule("no-such-quota-rule")
	require.Error(t, err)

	// attaching a user/rule to a non-existent group -> error
	me, _, err := c.GetMyUserInfo()
	require.NoError(t, err)
	_, err = c.AddUserToQuotaGroup("no-such-quota-group", me.UserName)
	require.Error(t, err)
	_, err = c.AddRuleToQuotaGroup("no-such-quota-group", "no-such-quota-rule")
	require.Error(t, err)

	// quota of a non-existent user -> error
	_, _, err = c.GetUserQuota("no-such-quota-user")
	require.Error(t, err)
}
