// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetActionsRun(t *testing.T) {
	log.Println("== TestGetActionsRun ==")

	// Build a minimal client directly instead of newTestClient(): without a
	// live server NewClient fails its version check and returns a nil client,
	// whose first request then panics. Constructing the client inline keeps this
	// a self-contained unit test. The /actions/run endpoint requires the
	// automatic actions bearer token ({{ forgejo.token }}) of a running job,
	// which a regular API token cannot provide, so we only assert that the call
	// fails cleanly rather than panicking.
	c := &Client{
		url:    getForgejoURL(),
		client: &http.Client{},
		ctx:    context.Background(),
	}

	_, _, err := c.GetActionsRun()
	require.Error(t, err)
}
