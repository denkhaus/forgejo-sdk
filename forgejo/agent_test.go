// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"
)

func Test_Agent(t *testing.T) {
	log.Println("== Test_Agent ==")
	// agent.go wires up the ssh-agent transport used for HTTP signing.
	// Exercise the client setup path without asserting on agent availability.
	_ = newTestClient()
}
