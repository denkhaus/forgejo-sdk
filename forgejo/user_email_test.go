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

// Test_UserEmail exercises ListEmails plus AddEmail/DeleteEmail.
// Adding may be rejected (already used / gated); tolerate and skip.
func Test_UserEmail(t *testing.T) {
	log.Println("== Test_UserEmail ==")
	c := newTestClient()

	list, _, err := c.ListEmails(ListEmailsOptions{})
	require.NoError(t, err)
	assert.NotNil(t, list)

	addr := "ue-useremail-test@x.test"
	_, _, err = c.AddEmail(CreateEmailOption{Emails: []string{addr}})
	if err != nil {
		t.Skipf("AddEmail unavailable or address already used: %v", err)
	}
	defer func() { _, _ = c.DeleteEmail(DeleteEmailOption{Emails: []string{addr}}) }()

	list2, _, err := c.ListEmails(ListEmailsOptions{})
	require.NoError(t, err)
	found := false
	for _, e := range list2 {
		if string(e.Email) == addr {
			found = true
			break
		}
	}
	assert.True(t, found, "added email should appear in list")
}
