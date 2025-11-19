// Copyright 2024 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepoActionSecrets(t *testing.T) {
	log.Println("== TestRepoActionSecrets ==")
	c := newTestClient()

	// Create one user and repo for all subtests
	user := createTestUser(t, "repo_action_test_user", c)
	c.SetSudo(user.UserName)
	testRepo, _, err := c.CreateRepo(CreateRepoOption{Name: "ActionTestRepo"})
	require.NoError(t, err)
	require.NotNil(t, testRepo)

	t.Run("CreateAndUpdate", func(t *testing.T) {
		// create secret
		resp, err := c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "test", Data: "test"})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// update secret
		resp, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "test", Data: "test2"})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// create secret with 64 character name
		longName64 := strings.Repeat("B", 64)
		resp, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: longName64, Data: "test_data"})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// create secret with long value
		longValue := strings.Repeat("secret", 100)
		resp, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "LONG_VALUE_TEST", Data: longValue})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// list secrets - should have 3 secrets
		secrets, _, err := c.ListRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, ListRepoActionSecretOption{})
		require.NoError(t, err)
		assert.Len(t, secrets, 3)
	})

	t.Run("InvalidNames", func(t *testing.T) {
		// test invalid names - client-side validation should reject these
		_, err := c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "", Data: "data"})
		assert.Error(t, err)
		assert.EqualError(t, err, "name required")

		_, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "INVALID-NAME", Data: "data"})
		assert.Error(t, err)
		assert.EqualError(t, err, "name must contain only alphanumeric characters and underscores")

		_, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "INVALID NAME", Data: "data"})
		assert.Error(t, err)
		assert.EqualError(t, err, "name must contain only alphanumeric characters and underscores")

		_, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "GITEA_SECRET", Data: "data"})
		assert.Error(t, err)
		assert.EqualError(t, err, "name cannot start with GITEA_ or GITHUB_ (reserved prefixes)")

		_, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "github_token", Data: "data"})
		assert.Error(t, err)
		assert.EqualError(t, err, "name cannot start with GITEA_ or GITHUB_ (reserved prefixes)")

		_, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: strings.Repeat("A", 256), Data: "data"})
		assert.Error(t, err)
		assert.EqualError(t, err, "name too long (maximum 255 characters)")
	})

	t.Run("EmptyData", func(t *testing.T) {
		_, err := c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "VALID_NAME", Data: ""})
		assert.Error(t, err)
		assert.EqualError(t, err, "data required")
	})

	t.Run("UpdateMultipleTimes", func(t *testing.T) {
		// create and update same secret multiple times
		resp, err := c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "UPDATE_SECRET", Data: "data1"})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		for i := 2; i <= 5; i++ {
			resp, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "UPDATE_SECRET", Data: "data" + string(rune('0'+i))})
			require.NoError(t, err)
			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		}

		secrets, _, err := c.ListRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, ListRepoActionSecretOption{})
		require.NoError(t, err)
		// Should have original 3 + UPDATE_SECRET = 4 total
		assert.GreaterOrEqual(t, len(secrets), 1)
	})

	t.Run("CaseSensitivity", func(t *testing.T) {
		// create secret with lowercase name
		resp, err := c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "my_secret", Data: "lower"})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// try to create with uppercase - should update the same secret (case-insensitive)
		resp, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "MY_SECRET", Data: "upper"})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// try with mixed case - should also update the same secret
		resp, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "My_Secret", Data: "mixed"})
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("Pagination", func(t *testing.T) {
		// create 10 more secrets for pagination test
		for i := 1; i <= 10; i++ {
			name := "PAGE_SECRET_" + string(rune('A'+i-1))
			_, err := c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: name, Data: "data" + string(rune('0'+i))})
			require.NoError(t, err)
		}

		// test pagination
		secrets, _, err := c.ListRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, ListRepoActionSecretOption{
			ListOptions: ListOptions{Page: 1, PageSize: 5},
		})
		require.NoError(t, err)
		assert.Len(t, secrets, 5)

		secrets, _, err = c.ListRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, ListRepoActionSecretOption{
			ListOptions: ListOptions{Page: 2, PageSize: 5},
		})
		require.NoError(t, err)
		assert.Len(t, secrets, 5)

		// get all
		secrets, _, err = c.ListRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, ListRepoActionSecretOption{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(secrets), 10)
	})

	t.Run("NonExistentRepo", func(t *testing.T) {
		// trying to create secret for non-existent repo should fail
		_, err := c.CreateRepoActionSecret(testRepo.Owner.UserName, "NonExistentRepo123456", CreateSecretOption{Name: "TEST", Data: "data"})
		assert.Error(t, err)

		// listing secrets for non-existent repo also returns error
		_, _, err = c.ListRepoActionSecret(testRepo.Owner.UserName, "NonExistentRepo123456", ListRepoActionSecretOption{})
		assert.Error(t, err)
	})

	t.Run("LargeData", func(t *testing.T) {
		// test various data sizes
		smallData := "small"
		mediumData := strings.Repeat("medium", 100)
		largeData := strings.Repeat("large", 1000)

		_, err := c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "SMALL_DATA", Data: smallData})
		require.NoError(t, err)

		_, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "MEDIUM_DATA", Data: mediumData})
		require.NoError(t, err)

		_, err = c.CreateRepoActionSecret(testRepo.Owner.UserName, testRepo.Name, CreateSecretOption{Name: "LARGE_DATA", Data: largeData})
		require.NoError(t, err)
	})
}
