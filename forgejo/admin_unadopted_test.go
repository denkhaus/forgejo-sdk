// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAdminUnadoptedV15 exercises listing unadopted repositories.
func TestAdminUnadoptedV15(t *testing.T) {
	log.Println("== TestAdminUnadoptedV15 ==")
	c := newTestClient()

	_, _, err := c.AdminListUnadoptedRepositories(AdminListUnadoptedOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
}
