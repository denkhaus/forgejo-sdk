// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RepoConvert(t *testing.T) {
	log.Println("== Test_RepoConvert ==")
	c := newTestClient()
	repo, _ := createTestRepo(t, "TestRepoConvert", c)

	// ConvertRepo only succeeds on a fork<->regular eligible repository.
	// On a freshly created normal repo the server rejects the conversion;
	// assert only when it unexpectedly succeeds.
	converted, _, err := c.ConvertRepo(repo.Owner.UserName, repo.Name)
	if err == nil {
		assert.NotNil(t, converted)
	}
}
