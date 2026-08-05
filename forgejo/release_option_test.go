// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEditReleaseOptionNoteMarshal guards the fix for forgejo-cli #138:
// EditReleaseOption.Note is *string with omitempty, so an unset note is OMITTED
// from the PATCH body (preserving the release's existing notes) rather than
// serialized as "body":"" (which the API treats as "clear the notes"). An
// explicit OptionalString("") still sends "body":"" for a deliberate clear.
//
// This is a pure marshal test — no server required.
func TestEditReleaseOptionNoteMarshal(t *testing.T) {
	t.Run("nil note omits body (preserves existing notes)", func(t *testing.T) {
		b, err := json.Marshal(EditReleaseOption{})
		require.NoError(t, err)
		assert.NotContains(t, string(b), `"body"`, "unset Note must not be sent (would clobber existing notes)")
	})
	t.Run("explicit empty note sends body (deliberate clear)", func(t *testing.T) {
		b, err := json.Marshal(EditReleaseOption{Note: OptionalString("")})
		require.NoError(t, err)
		assert.Contains(t, string(b), `"body":""`, "explicit empty Note must be sent to clear notes")
	})
	t.Run("non-empty note sends body", func(t *testing.T) {
		b, err := json.Marshal(EditReleaseOption{Note: OptionalString("release notes")})
		require.NoError(t, err)
		assert.Contains(t, string(b), `"body":"release notes"`)
	})
}
