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

// Test_AdminCron exercises ListCronTasks (read-only; deliberately does NOT run tasks).
func Test_AdminCron(t *testing.T) {
	log.Println("== Test_AdminCron ==")
	c := newTestClient()

	tasks, _, err := c.ListCronTasks(ListCronTaskOptions{ListOptions: ListOptions{PageSize: 10}})
	require.NoError(t, err)
	assert.NotNil(t, tasks)
}
