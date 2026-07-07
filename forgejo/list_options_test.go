// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ListOptions(t *testing.T) {
	// a zero Page should default to 1
	o := ListOptions{}
	o.setDefaults()
	assert.Equal(t, 1, o.Page)

	// a negative Page disables pagination (Page and PageSize -> 0)
	o2 := ListOptions{Page: -1, PageSize: 10}
	o2.setDefaults()
	assert.Equal(t, 0, o2.Page)
	assert.Equal(t, 0, o2.PageSize)

	// getURLQuery reflects the current values
	q := ListOptions{Page: 2, PageSize: 5}.getURLQuery()
	assert.Equal(t, "2", q.Get("page"))
	assert.Equal(t, "5", q.Get("limit"))
}
