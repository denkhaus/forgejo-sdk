// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestAdminEmailsV15 exercises admin email listing and search.
func TestAdminEmailsV15(t *testing.T) {
	log.Println("== TestAdminEmailsV15 ==")
	c := newTestClient()

	_, _, err := c.AdminListEmails(AdminListEmailsOption{ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)

	_, _, err = c.AdminSearchEmails(AdminSearchEmailsOption{Query: "test01", ListOptions: ListOptions{PageSize: 5}})
	require.NoError(t, err)
}
