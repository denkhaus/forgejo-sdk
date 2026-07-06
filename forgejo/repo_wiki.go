// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListWikiPagesOption options for listing wiki pages or revisions
type ListWikiPagesOption struct {
	ListOptions
}

// ListRepoWikiPages lists a repository's wiki pages
func (c *Client) ListRepoWikiPages(owner, repo string, opt ListWikiPagesOption) ([]*models.WikiPageMetaData, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/wiki/pages", owner, repo))
	link.RawQuery = opt.getURLQuery().Encode()
	pages := make([]*models.WikiPageMetaData, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &pages)
	return pages, resp, err
}

// GetRepoWikiPage gets a single wiki page
func (c *Client) GetRepoWikiPage(owner, repo, pageName string) (*models.WikiPage, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &pageName); err != nil {
		return nil, nil, err
	}
	page := new(models.WikiPage)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/wiki/page/%s", owner, repo, pageName), jsonHeader, nil, page)
	return page, resp, err
}

// GetRepoWikiPageRevisions lists the revisions of a wiki page
func (c *Client) GetRepoWikiPageRevisions(owner, repo, pageName string, opt ListWikiPagesOption) ([]*models.WikiCommit, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &pageName); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	link, _ := url.Parse(fmt.Sprintf("/repos/%s/%s/wiki/revisions/%s", owner, repo, pageName))
	link.RawQuery = opt.getURLQuery().Encode()
	commits := make([]*models.WikiCommit, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", link.String(), jsonHeader, nil, &commits)
	return commits, resp, err
}

// CreateRepoWikiPage creates a new wiki page in a repository
func (c *Client) CreateRepoWikiPage(owner, repo string, opt models.CreateWikiPageOptions) (*models.WikiPage, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	page := new(models.WikiPage)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/wiki/new", owner, repo), jsonHeader, bytes.NewReader(body), page)
	return page, resp, err
}

// EditRepoWikiPage edits a wiki page in a repository
func (c *Client) EditRepoWikiPage(owner, repo, pageName string, opt models.CreateWikiPageOptions) (*models.WikiPage, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &pageName); err != nil {
		return nil, nil, err
	}
	if err := (&opt).Validate(nil); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, nil, err
	}
	page := new(models.WikiPage)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/wiki/page/%s", owner, repo, pageName), jsonHeader, bytes.NewReader(body), page)
	return page, resp, err
}

// DeleteRepoWikiPage deletes a wiki page in a repository
func (c *Client) DeleteRepoWikiPage(owner, repo, pageName string) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo, &pageName); err != nil {
		return nil, err
	}
	status, resp, err := c.getStatusCode("DELETE", fmt.Sprintf("/repos/%s/%s/wiki/page/%s", owner, repo, pageName), jsonHeader, nil)
	if err != nil {
		return resp, err
	}
	switch status {
	case http.StatusNoContent, http.StatusOK:
		return resp, nil
	case http.StatusNotFound:
		return resp, fmt.Errorf("wiki page not found")
	case http.StatusForbidden:
		return resp, fmt.Errorf("forbidden: permission denied")
	default:
		return resp, fmt.Errorf("unexpected Status: %d", status)
	}
}
