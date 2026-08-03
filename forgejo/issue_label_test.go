// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gitea Authors. All rights reserved.
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

// TestLabels test label related func
func TestLabels(t *testing.T) {
	log.Println("== TestLabels ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "LabelTestsRepo", c)
	require.NoError(t, err)

	createOpts := CreateLabelOption{
		Name:        " ",
		Description: "",
		Color:       "",
	}
	err = createOpts.Validate()
	require.Error(t, err)
	assert.EqualValues(t, "invalid color format", err.Error())
	createOpts.Color = "12345f"
	err = createOpts.Validate()
	require.Error(t, err)
	assert.EqualValues(t, "empty name not allowed", err.Error())
	createOpts.Name = "label one"

	labelOne, _, err := c.CreateLabel(repo.Owner.UserName, repo.Name, createOpts)
	require.NoError(t, err)
	assert.EqualValues(t, createOpts.Name, labelOne.Name)
	assert.EqualValues(t, createOpts.Color, labelOne.Color)

	labelTwo, _, err := c.CreateLabel(repo.Owner.UserName, repo.Name, CreateLabelOption{
		Name:        "blue",
		Color:       "#0000FF",
		Description: "CMYB(100%, 100%, 0%, 0%)",
	})
	require.NoError(t, err)
	_, _, err = c.CreateLabel(repo.Owner.UserName, repo.Name, CreateLabelOption{
		Name:        "gray",
		Color:       "808080",
		Description: "CMYB(0%, 0%, 0%, 50%)",
	})
	require.NoError(t, err)
	_, _, err = c.CreateLabel(repo.Owner.UserName, repo.Name, CreateLabelOption{
		Name:        "green",
		Color:       "#98F76C",
		Description: "CMYB(38%, 0%, 56%, 3%)",
	})
	require.NoError(t, err)

	labels, resp, err := c.ListRepoLabels(repo.Owner.UserName, repo.Name, ListLabelsOptions{ListOptions: ListOptions{PageSize: 3}})
	require.NoError(t, err)
	assert.Len(t, labels, 3)
	assert.NotNil(t, resp)
	assert.Contains(t, labels, labelTwo)
	assert.NotContains(t, labels, labelOne)

	label, _, err := c.GetRepoLabel(repo.Owner.UserName, repo.Name, labelTwo.ID)
	require.NoError(t, err)
	assert.EqualValues(t, labelTwo, label)

	label, _, err = c.EditLabel(repo.Owner.UserName, repo.Name, labelTwo.ID, EditLabelOption{
		Color:       OptionalString("#0e0175"),
		Description: OptionalString("blueish"),
	})
	require.NoError(t, err)
	assert.EqualValues(t, &models.Label{
		ID:          labelTwo.ID,
		Name:        labelTwo.Name,
		Color:       "0e0175",
		Description: "blueish",
		URL:         labelTwo.URL,
	}, label)
	labels, _, _ = c.ListRepoLabels(repo.Owner.UserName, repo.Name, ListLabelsOptions{ListOptions: ListOptions{PageSize: 3}})

	createTestIssue(t, c, repo.Name, "test-issue", "", nil, nil, 0, []int64{label.ID}, false, false)
	issueIndex := int64(1)

	issueLabels, _, err := c.GetIssueLabels(repo.Owner.UserName, repo.Name, issueIndex, ListLabelsOptions{})
	require.NoError(t, err)
	assert.Len(t, issueLabels, 1)
	assert.EqualValues(t, label, issueLabels[0])

	_, _, err = c.AddIssueLabels(repo.Owner.UserName, repo.Name, issueIndex, models.IssueLabelsOption{Labels: []any{labels[0].ID}})
	require.NoError(t, err)

	issueLabels, _, err = c.AddIssueLabels(repo.Owner.UserName, repo.Name, issueIndex, models.IssueLabelsOption{Labels: []any{labels[1].ID, labels[2].ID}})
	require.NoError(t, err)
	assert.Len(t, issueLabels, 3)
	assert.EqualValues(t, labels, issueLabels)

	labels, _, _ = c.ListRepoLabels(repo.Owner.UserName, repo.Name, ListLabelsOptions{})
	assert.Len(t, labels, 11)

	issueLabels, _, err = c.ReplaceIssueLabels(repo.Owner.UserName, repo.Name, issueIndex, models.IssueLabelsOption{Labels: []any{labels[0].ID, labels[1].ID}})
	require.NoError(t, err)
	assert.Len(t, issueLabels, 2)

	_, err = c.DeleteIssueLabel(repo.Owner.UserName, repo.Name, issueIndex, labels[0].ID)
	require.NoError(t, err)
	issueLabels, _, _ = c.GetIssueLabels(repo.Owner.UserName, repo.Name, issueIndex, ListLabelsOptions{})
	assert.Len(t, issueLabels, 1)

	_, err = c.ClearIssueLabels(repo.Owner.UserName, repo.Name, issueIndex)
	require.NoError(t, err)
	issueLabels, _, _ = c.GetIssueLabels(repo.Owner.UserName, repo.Name, issueIndex, ListLabelsOptions{})
	assert.Empty(t, issueLabels)

	_, err = c.DeleteLabel(repo.Owner.UserName, repo.Name, labelTwo.ID)
	require.NoError(t, err)
	labels, _, _ = c.ListRepoLabels(repo.Owner.UserName, repo.Name, ListLabelsOptions{})
	assert.Len(t, labels, 10)
}

// Test_Labels_Negative covers client-side validation and server-side not-found
// error paths for both repo labels and issue labels.
func Test_Labels_Negative(t *testing.T) {
	log.Println("== Test_Labels_Negative ==")
	c := newTestClient()
	repo, err := createTestRepo(t, "LabelNegativeRepo", c)
	require.NoError(t, err)
	owner, name := repo.Owner.UserName, repo.Name

	// --- client-side validation (no server round-trip) ---
	t.Run("CreateLabelOption.Validate", func(t *testing.T) {
		require.Error(t, CreateLabelOption{Name: "x", Color: "nope"}.Validate())    // bad color
		require.Error(t, CreateLabelOption{Name: " ", Color: "#00aabb"}.Validate()) // empty name
		require.NoError(t, CreateLabelOption{Name: "x", Color: "#00aabb"}.Validate())
		require.NoError(t, CreateLabelOption{Name: "x", Color: "00aabb"}.Validate()) // '#' optional
	})
	t.Run("EditLabelOption.Validate", func(t *testing.T) {
		badColor := "nope"
		emptyName := "  "
		require.Error(t, EditLabelOption{Color: &badColor}.Validate())
		require.Error(t, EditLabelOption{Name: &emptyName}.Validate())
		require.NoError(t, EditLabelOption{}.Validate()) // nil fields = no-op, valid
	})

	// --- server-side not-found paths ---
	const ghostID int64 = 999999

	// repo label operations on a non-existent label id -> error (get/edit)
	_, _, err = c.GetRepoLabel(owner, name, ghostID)
	require.Error(t, err)
	_, _, err = c.EditLabel(owner, name, ghostID, EditLabelOption{})
	require.Error(t, err)
	// DeleteLabel is idempotent on a valid repo: a non-existent label id is a
	// no-op and returns no error.
	_, err = c.DeleteLabel(owner, name, ghostID)
	require.NoError(t, err)

	// non-existent repo -> error for get and delete
	_, _, err = c.GetRepoLabel(owner, "no-such-repo", 1)
	require.Error(t, err)
	_, err = c.DeleteLabel(owner, "no-such-repo", ghostID)
	require.Error(t, err)

	// issue-label operations on a non-existent issue -> error
	_, _, err = c.GetIssueLabels(owner, name, ghostID, ListLabelsOptions{})
	require.Error(t, err)
	_, _, err = c.AddIssueLabels(owner, name, ghostID, models.IssueLabelsOption{Labels: []any{int64(1)}})
	require.Error(t, err)
	_, _, err = c.ReplaceIssueLabels(owner, name, ghostID, models.IssueLabelsOption{Labels: []any{int64(1)}})
	require.Error(t, err)
	_, err = c.DeleteIssueLabel(owner, name, ghostID, 1)
	require.Error(t, err)
	_, err = c.ClearIssueLabels(owner, name, ghostID)
	require.Error(t, err)
}
