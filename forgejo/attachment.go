// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo // import "code.codeberg.org/mvdkleijn/forgejo-sdk"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListReleaseAttachmentsOptions options for listing release's attachments
type ListReleaseAttachmentsOptions struct {
	ListOptions
}

// ListReleaseAttachments list release's attachments
func (c *Client) ListReleaseAttachments(user, repo string, release int64, opt ListReleaseAttachmentsOptions) ([]*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&user, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	attachments := make([]*models.Attachment, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET",
		fmt.Sprintf("/repos/%s/%s/releases/%d/assets?%s", user, repo, release, opt.getURLQuery().Encode()),
		nil, nil, &attachments)
	return attachments, resp, err
}

// GetReleaseAttachment returns the requested attachment
func (c *Client) GetReleaseAttachment(user, repo string, release, id int64) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&user, &repo); err != nil {
		return nil, nil, err
	}
	a := new(models.Attachment)
	resp, err := c.getParsedResponse("GET",
		fmt.Sprintf("/repos/%s/%s/releases/%d/assets/%d", user, repo, release, id),
		nil, nil, &a)
	return a, resp, err
}

// CreateReleaseAttachment creates an attachment for the given release
func (c *Client) CreateReleaseAttachment(user, repo string, release int64, file io.Reader, filename string) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&user, &repo); err != nil {
		return nil, nil, err
	}
	// Write file to body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("attachment", filename)
	if err != nil {
		return nil, nil, err
	}

	if _, err = io.Copy(part, file); err != nil {
		return nil, nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, nil, err
	}

	// Send request
	attachment := new(models.Attachment)
	resp, err := c.getParsedResponse("POST",
		fmt.Sprintf("/repos/%s/%s/releases/%d/assets", user, repo, release),
		http.Header{"Content-Type": []string{writer.FormDataContentType()}}, body, &attachment)
	return attachment, resp, err
}

// EditReleaseAttachment updates the given attachment with the given options
func (c *Client) EditReleaseAttachment(user, repo string, release, attachment int64, form models.EditAttachmentOptions) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&user, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&form)
	if err != nil {
		return nil, nil, err
	}
	attach := new(models.Attachment)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/releases/%d/assets/%d", user, repo, release, attachment), jsonHeader, bytes.NewReader(body), attach)
	return attach, resp, err
}

// DeleteReleaseAttachment deletes the given attachment including the uploaded file
func (c *Client) DeleteReleaseAttachment(user, repo string, release, id int64) (*Response, error) {
	if err := escapeValidatePathSegments(&user, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/releases/%d/assets/%d", user, repo, release, id), nil, nil)
	return resp, err
}
