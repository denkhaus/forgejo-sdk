// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUserHeatmapV15 exercises fetching a user's heatmap data.
func TestUserHeatmapV15(t *testing.T) {
	log.Println("== TestUserHeatmapV15 ==")
	c := newTestClient()

	user, _, err := c.GetMyUserInfo()
	require.NoError(t, err)

	_, _, err = c.GetUserHeatmapData(user.UserName)
	require.NoError(t, err)
}
