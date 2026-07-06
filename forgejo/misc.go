// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ---------------------------------------------------------------------------
// Instance metadata
// ---------------------------------------------------------------------------

// GetSigningKey returns the instance's GPG signing key (ASCII-armored).
func (c *Client) GetSigningKey() (string, *Response, error) {
	data, resp, err := c.getResponse("GET", "/signing-key.gpg", jsonHeader, nil)
	return string(data), resp, err
}

// GetSSHSigningKey returns the instance's SSH signing key.
func (c *Client) GetSSHSigningKey() (string, *Response, error) {
	data, resp, err := c.getResponse("GET", "/signing-key.ssh", jsonHeader, nil)
	return string(data), resp, err
}

// GetNodeInfo returns the instance's nodeinfo document.
func (c *Client) GetNodeInfo() (*models.NodeInfo, *Response, error) {
	info := new(models.NodeInfo)
	resp, err := c.getParsedResponse("GET", "/nodeinfo", jsonHeader, nil, info)
	return info, resp, err
}

// ---------------------------------------------------------------------------
// Markdown / markup rendering
// ---------------------------------------------------------------------------

// RenderMarkdown renders a MarkdownOption to HTML.
func (c *Client) RenderMarkdown(opt models.MarkdownOption) (string, *Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return "", nil, err
	}
	data, resp, err := c.getResponse("POST", "/markdown", jsonHeader, bytes.NewReader(body))
	return string(data), resp, err
}

// RenderMarkdownRaw renders raw markdown bytes to HTML.
func (c *Client) RenderMarkdownRaw(body []byte) (string, *Response, error) {
	header := http.Header{"Content-Type": []string{"text/plain"}}
	data, resp, err := c.getResponse("POST", "/markdown/raw", header, bytes.NewReader(body))
	return string(data), resp, err
}

// RenderMarkup renders a MarkupOption (e.g. asciidoc, rst) to HTML.
func (c *Client) RenderMarkup(opt models.MarkupOption) (string, *Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return "", nil, err
	}
	data, resp, err := c.getResponse("POST", "/markup", jsonHeader, bytes.NewReader(body))
	return string(data), resp, err
}

// ---------------------------------------------------------------------------
// Repository/file templates
// ---------------------------------------------------------------------------

// ListGitignoreTemplates lists the available .gitignore template names.
func (c *Client) ListGitignoreTemplates() ([]string, *Response, error) {
	templates := make([]string, 0, 10)
	resp, err := c.getParsedResponse("GET", "/gitignore/templates", jsonHeader, nil, &templates)
	return templates, resp, err
}

// GetGitignoreTemplateInfo returns the content of a .gitignore template.
func (c *Client) GetGitignoreTemplateInfo(name string) (*models.GitignoreTemplateInfo, *Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, nil, err
	}
	info := new(models.GitignoreTemplateInfo)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/gitignore/templates/%s", name), jsonHeader, nil, info)
	return info, resp, err
}

// ListLabelTemplates lists the available label template names.
func (c *Client) ListLabelTemplates() ([]string, *Response, error) {
	templates := make([]string, 0, 10)
	resp, err := c.getParsedResponse("GET", "/label/templates", jsonHeader, nil, &templates)
	return templates, resp, err
}

// GetLabelTemplateInfo returns the labels defined by a label template.
func (c *Client) GetLabelTemplateInfo(name string) (*models.LabelTemplate, *Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, nil, err
	}
	t := new(models.LabelTemplate)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/label/templates/%s", name), jsonHeader, nil, t)
	return t, resp, err
}

// ListLicenseTemplates lists the available license templates.
func (c *Client) ListLicenseTemplates() ([]*models.LicensesTemplateListEntry, *Response, error) {
	templates := make([]*models.LicensesTemplateListEntry, 0, 10)
	resp, err := c.getParsedResponse("GET", "/licenses", jsonHeader, nil, &templates)
	return templates, resp, err
}

// GetLicenseTemplateInfo returns the content of a license template.
func (c *Client) GetLicenseTemplateInfo(name string) (*models.LicenseTemplateInfo, *Response, error) {
	if err := escapeValidatePathSegments(&name); err != nil {
		return nil, nil, err
	}
	info := new(models.LicenseTemplateInfo)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/licenses/%s", name), jsonHeader, nil, info)
	return info, resp, err
}
