// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"
)

// LinkPackage links a package to a repository (auto-link, Forgejo v15+).
// Used to associate container packages with a repository.
func (c *Client) LinkPackage(owner, packageType, name, repoName string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &packageType, &name, &repoName); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/packages/%s/%s/%s/-/link/%s", owner, packageType, name, repoName), jsonHeader, nil)
	return resp, err
}

// UnlinkPackage removes the repository link from a package (Forgejo v15+).
func (c *Client) UnlinkPackage(owner, packageType, name string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &packageType, &name); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/packages/%s/%s/%s/-/unlink", owner, packageType, name), jsonHeader, nil)
	return resp, err
}
