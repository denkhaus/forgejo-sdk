// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2023 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParsedPaging parses a Link header (RFC 5988) into first/prev/next/last
// page numbers. The URL host/path are arbitrary fixture strings — the parser
// only inspects the page query parameter.
func TestParsedPaging(t *testing.T) {
	resp := newResponse(&http.Response{
		Header: http.Header{
			"Link": []string{
				strings.Join(
					[]string{
						`<https://codeberg.org/api/v1/repos/forgejo/forgejo/issues/1/comments?page=3>; rel="next"`,
						`<https://codeberg.org/api/v1/repos/forgejo/forgejo/issues/1/comments?page=4>; rel="last"`,
						`<https://codeberg.org/api/v1/repos/forgejo/forgejo/issues/1/comments?page=1>; rel="first"`,
						`<https://codeberg.org/api/v1/repos/forgejo/forgejo/issues/1/comments?page=1>; rel="prev"`,
					}, ",",
				),
			},
		},
	})

	assert.Equal(t, 1, resp.FirstPage)
	assert.Equal(t, 1, resp.PrevPage)
	assert.Equal(t, 3, resp.NextPage)
	assert.Equal(t, 4, resp.LastPage)
}
