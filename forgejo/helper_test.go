// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Helper(t *testing.T) {
	b := OptionalBool(true)
	assert.NotNil(t, b)
	assert.True(t, *b)

	s := OptionalString("forgejo")
	assert.NotNil(t, s)
	assert.Equal(t, "forgejo", *s)

	i := OptionalInt64(42)
	assert.NotNil(t, i)
	assert.EqualValues(t, 42, *i)
}
