// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ListIssueAttachmentsOptions options for listing issue or comment attachments
type ListIssueAttachmentsOptions struct {
	ListOptions
}

// ---------------------------------------------------------------------------
// Issue attachments
// ---------------------------------------------------------------------------

// ListIssueAttachments lists the attachments of an issue
func (c *Client) ListIssueAttachments(owner, repo string, index int64, opt ListIssueAttachmentsOptions) ([]*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	attachments := make([]*models.Attachment, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/assets?%s", owner, repo, index, opt.getURLQuery().Encode()), nil, nil, &attachments)
	return attachments, resp, err
}

// GetIssueAttachment returns the requested issue attachment
func (c *Client) GetIssueAttachment(owner, repo string, index, attachmentID int64) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	a := new(models.Attachment)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/%d/assets/%d", owner, repo, index, attachmentID), nil, nil, a)
	return a, resp, err
}

// CreateIssueAttachment creates an attachment for the given issue
func (c *Client) CreateIssueAttachment(owner, repo string, index int64, file io.Reader, filename string) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
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
	a := new(models.Attachment)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/issues/%d/assets", owner, repo, index), http.Header{"Content-Type": []string{writer.FormDataContentType()}}, body, a)
	return a, resp, err
}

// EditIssueAttachment updates an issue attachment
func (c *Client) EditIssueAttachment(owner, repo string, index, attachmentID int64, form models.EditAttachmentOptions) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&form)
	if err != nil {
		return nil, nil, err
	}
	a := new(models.Attachment)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/issues/%d/assets/%d", owner, repo, index, attachmentID), jsonHeader, bytes.NewReader(body), a)
	return a, resp, err
}

// DeleteIssueAttachment deletes an issue attachment
func (c *Client) DeleteIssueAttachment(owner, repo string, index, attachmentID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/issues/%d/assets/%d", owner, repo, index, attachmentID), nil, nil)
	return resp, err
}

// ---------------------------------------------------------------------------
// Issue comment attachments
// ---------------------------------------------------------------------------

// ListIssueCommentAttachments lists the attachments of an issue comment
func (c *Client) ListIssueCommentAttachments(owner, repo string, commentID int64, opt ListIssueAttachmentsOptions) ([]*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	opt.setDefaults()
	attachments := make([]*models.Attachment, 0, opt.PageSize)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/comments/%d/assets?%s", owner, repo, commentID, opt.getURLQuery().Encode()), nil, nil, &attachments)
	return attachments, resp, err
}

// GetIssueCommentAttachment returns the requested issue comment attachment
func (c *Client) GetIssueCommentAttachment(owner, repo string, commentID, attachmentID int64) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	a := new(models.Attachment)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/issues/comments/%d/assets/%d", owner, repo, commentID, attachmentID), nil, nil, a)
	return a, resp, err
}

// CreateIssueCommentAttachment creates an attachment for the given issue comment
func (c *Client) CreateIssueCommentAttachment(owner, repo string, commentID int64, file io.Reader, filename string) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
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
	a := new(models.Attachment)
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/issues/comments/%d/assets", owner, repo, commentID), http.Header{"Content-Type": []string{writer.FormDataContentType()}}, body, a)
	return a, resp, err
}

// EditIssueCommentAttachment updates an issue comment attachment
func (c *Client) EditIssueCommentAttachment(owner, repo string, commentID, attachmentID int64, form models.EditAttachmentOptions) (*models.Attachment, *Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, nil, err
	}
	body, err := json.Marshal(&form)
	if err != nil {
		return nil, nil, err
	}
	a := new(models.Attachment)
	resp, err := c.getParsedResponse("PATCH", fmt.Sprintf("/repos/%s/%s/issues/comments/%d/assets/%d", owner, repo, commentID, attachmentID), jsonHeader, bytes.NewReader(body), a)
	return a, resp, err
}

// DeleteIssueCommentAttachment deletes an issue comment attachment
func (c *Client) DeleteIssueCommentAttachment(owner, repo string, commentID, attachmentID int64) (*Response, error) {
	if err := escapeValidatePathSegments(&owner, &repo); err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("DELETE", fmt.Sprintf("/repos/%s/%s/issues/comments/%d/assets/%d", owner, repo, commentID, attachmentID), nil, nil)
	return resp, err
}
