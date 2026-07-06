// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListQuotaUsedOption options for listing quota-used resources (artifacts, attachments, packages)
type ListQuotaUsedOption struct {
	ListOptions
}

// ---------------------------------------------------------------------------
// Quota groups (admin)
// ---------------------------------------------------------------------------

// ListQuotaGroups lists all quota groups
func (c *Client) ListQuotaGroups() ([]*models.QuotaGroup, *Response, error) {
	groups := make([]*models.QuotaGroup, 0, 10)
	resp, err := c.getParsedResponse("GET", "/admin/quota/groups", jsonHeader, nil, &groups)
	return groups, resp, err
}

// GetQuotaGroup gets a quota group by name
func (c *Client) GetQuotaGroup(name string) (*models.QuotaGroup, *Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, nil, err
	}
	group := new(models.QuotaGroup)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/admin/quota/groups/%s", name), jsonHeader, nil, group)
	return group, resp, err
}

// CreateQuotaGroup creates a quota group
func (c *Client) CreateQuotaGroup(opt models.CreateQuotaGroupOptions) (*models.QuotaGroup, *Response, error) {
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	group := new(models.QuotaGroup)
	resp, err := c.getParsedResponse("POST", "/admin/quota/groups", jsonHeader, bytes.NewReader(body), group)
	return group, resp, err
}

// DeleteQuotaGroup deletes a quota group by name
func (c *Client) DeleteQuotaGroup(name string) (*Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/admin/quota/groups/%s", name), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("quota group not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// ListQuotaGroupUsers lists the users assigned to a quota group
func (c *Client) ListQuotaGroupUsers(name string) ([]*models.User, *Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, nil, err
	}
	users := make([]*models.User, 0, 10)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/admin/quota/groups/%s/users", name), jsonHeader, nil, &users)
	return users, resp, err
}

// AddUserToQuotaGroup adds a user to a quota group
func (c *Client) AddUserToQuotaGroup(group, username string) (*Response, error) {
	if err := escapeValidatePathSegments(&group, &username); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/admin/quota/groups/%s/users/%s", group, username), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("quota group or user not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// RemoveUserFromQuotaGroup removes a user from a quota group
func (c *Client) RemoveUserFromQuotaGroup(group, username string) (*Response, error) {
	if err := escapeValidatePathSegments(&group, &username); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/admin/quota/groups/%s/users/%s", group, username), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("quota group or user not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// AddRuleToQuotaGroup adds a rule to a quota group
func (c *Client) AddRuleToQuotaGroup(group, rule string) (*Response, error) {
	if err := escapeValidatePathSegments(&group, &rule); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("PUT", fmt.Sprintf("/admin/quota/groups/%s/rules/%s", group, rule), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("quota group or rule not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// RemoveRuleFromQuotaGroup removes a rule from a quota group
func (c *Client) RemoveRuleFromQuotaGroup(group, rule string) (*Response, error) {
	if err := escapeValidatePathSegments(&group, &rule); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/admin/quota/groups/%s/rules/%s", group, rule), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("quota group or rule not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// ---------------------------------------------------------------------------
// Quota rules (admin)
// ---------------------------------------------------------------------------

// ListQuotaRules lists all quota rules
func (c *Client) ListQuotaRules() ([]*models.QuotaRuleInfo, *Response, error) {
	rules := make([]*models.QuotaRuleInfo, 0, 10)
	resp, err := c.getParsedResponse("GET", "/admin/quota/rules", jsonHeader, nil, &rules)
	return rules, resp, err
}

// GetQuotaRule gets a quota rule by name
func (c *Client) GetQuotaRule(name string) (*models.QuotaRuleInfo, *Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, nil, err
	}
	rule := new(models.QuotaRuleInfo)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/admin/quota/rules/%s", name), jsonHeader, nil, rule)
	return rule, resp, err
}

// CreateQuotaRule creates a quota rule
func (c *Client) CreateQuotaRule(opt models.CreateQuotaRuleOptions) (*models.QuotaRuleInfo, *Response, error) {
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	rule := new(models.QuotaRuleInfo)
	resp, err := c.getParsedResponse("POST", "/admin/quota/rules", jsonHeader, bytes.NewReader(body), rule)
	return rule, resp, err
}

// EditQuotaRule edits a quota rule by name
func (c *Client) EditQuotaRule(name string, opt models.EditQuotaRuleOptions) (*models.QuotaRuleInfo, *Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	rule := new(models.QuotaRuleInfo)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/admin/quota/rules/%s", name), jsonHeader, bytes.NewReader(body), rule)
	return rule, resp, err
}

// DeleteQuotaRule deletes a quota rule by name
func (c *Client) DeleteQuotaRule(name string) (*Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/admin/quota/rules/%s", name), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("quota rule not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}

// ---------------------------------------------------------------------------
// User quota (admin)
// ---------------------------------------------------------------------------

// GetUserQuota gets the quota for a user
func (c *Client) GetUserQuota(username string) (*models.QuotaInfo, *Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, nil, err
	}
	info := new(models.QuotaInfo)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/admin/users/%s/quota", username), jsonHeader, nil, info)
	return info, resp, err
}

// SetUserQuotaGroups sets the quota groups assigned to a user
func (c *Client) SetUserQuotaGroups(username string, opt models.SetUserQuotaGroupsOptions) (*Response, error) {
	if err := escapeValidatePathSegments(&username); err != nil {
		return nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("POST", fmt.Sprintf("/admin/users/%s/quota/groups", username), jsonHeader, bytes.NewReader(body))
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("user not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	case http.StatusBadRequest:
		return resp, fmt.Errorf("bad request: invalid quota groups")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
