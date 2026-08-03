// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2020 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
)

func getForgejoURL() string {
	return os.Getenv("FORGEJO_SDK_TEST_URL")
}

func getForgejoToken() string {
	return os.Getenv("FORGEJO_SDK_TEST_TOKEN")
}

func getForgejoUsername() string {
	return os.Getenv("FORGEJO_SDK_TEST_USERNAME")
}

func getForgejoPassword() string {
	return os.Getenv("FORGEJO_SDK_TEST_PASSWORD")
}

func newTestClient() *Client {
	c, err := NewClient(getForgejoURL(), newTestClientAuth())
	if err != nil {
		// Without a live server NewClient fails its version check and returns
		// (nil, err). A nil client dereferences its RWMutex on the first request
		// and panics (see TestGetActionsRun for the same rationale). Fall back to
		// a minimal, non-nil client so the suite never panics — callers that
		// genuinely need a server get a clean request error instead.
		c = &Client{
			url:    getForgejoURL(),
			client: &http.Client{},
			ctx:    context.Background(),
		}
	}
	return c
}

func newTestClientAuth() ClientOption {
	token := getForgejoToken()
	if token == "" {
		return SetBasicAuth(getForgejoUsername(), getForgejoPassword())
	}
	return SetToken(getForgejoToken())
}

// TestMain runs the suite. The Forgejo test instance is managed externally via
// the Makefile (make test-instance-start / make test), not from within the Go
// test process, so TestMain only emits connection diagnostics and runs.
func TestMain(m *testing.M) {
	log.Printf("testing with %v, %v, %v\n", getForgejoURL(), getForgejoUsername(), getForgejoPassword())
	exitCode := m.Run()
	exit(exitCode)
}

func exit(code int) {
	os.Exit(code)
}
