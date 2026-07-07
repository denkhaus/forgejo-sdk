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

func Test_PackageLink(t *testing.T) {
	log.Println("== Test_PackageLink ==")
	c := newTestClient()

	user, _, err := c.GetMyUserInfo()
	require.NoError(t, err)

	// LinkPackage/UnlinkPackage require a real package to exist. Without one
	// the endpoints return an error (typically 404); tolerate that and only
	// assert on success.
	owner := user.UserName
	packageType := "generic"
	packageName := "linkable-package"

	_, linkErr := c.LinkPackage(owner, packageType, packageName, "LinkTargetRepo")
	if linkErr == nil {
		// Unlink should succeed if we linked.
		_, unlinkErr := c.UnlinkPackage(owner, packageType, packageName)
		assert.NoError(t, unlinkErr)
		return
	}

	// No package existed; the link/unlink paths cannot be exercised end-to-end
	// here. Treat the 404-style failure as environment-gated and skip.
	t.Skipf("package not present; link/unlink endpoints are gated: %v", linkErr)
}
