// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"log"
	"testing"
)

func Test_HTTPSign(t *testing.T) {
	log.Println("== Test_HTTPSign ==")
	// httpsign.go implements internal HTTP request signing used by the client
	// when ssh-key auth is configured; exercise the client setup path.
	_ = newTestClient()
}
