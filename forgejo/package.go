// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2023 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListPackagesOptions options for listing packages
type ListPackagesOptions struct {
	ListOptions
}

// ListPackages lists all the packages owned by a given owner (user, organisation)
func (c *Client) ListPackages(owner string, opt ListPackagesOptions) ([]*models.Package, *Response, error) {
	if err := escapeValidatePathSegments(&owner); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	packages := make([]*models.Package, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/packages/%s?%s", owner, opt.getURLQuery().Encode()), nil, nil, &packages)
	return packages, resp, err
}

// GetPackage gets the details of a specific package version
func (c *Client) GetPackage(owner, packageType, name, version string) (*models.Package, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &packageType, &name, &version); err != nil {
		return nil, nil, err
	}
	foundPackage := new(models.Package)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/packages/%s/%s/%s/%s", owner, packageType, name, version), nil, nil, foundPackage)
	return foundPackage, resp, err
}

// DeletePackage deletes a specific package version
func (c *Client) DeletePackage(owner, packageType, name, version string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &packageType, &name, &version); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/packages/%s/%s/%s/%s", owner, packageType, name, version), nil, nil)
	return resp, err
}

// ListPackageFiles lists the files within a package
func (c *Client) ListPackageFiles(owner, packageType, name, version string) ([]*models.PackageFile, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &packageType, &name, &version); err != nil {
		return nil, nil, err
	}
	packageFiles := make([]*models.PackageFile, 0)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/packages/%s/%s/%s/%s/files", owner, packageType, name, version), nil, nil, &packageFiles)
	return packageFiles, resp, err
}
