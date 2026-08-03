// Copyright 2024 The Forgejo Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package forgejo

import (
	"bytes"
	"encoding/json"
	"fmt"

	"codeberg.org/mvdkleijn/forgejo-sdk/forgejo/v2/models"
)

// ---------------------------------------------------------------------------
// Instance actor
// ---------------------------------------------------------------------------

// GetInstanceActor returns the instance's ActivityPub actor.
func (c *Client) GetInstanceActor() (*models.ActivityPub, *Response, error) {
	actor := new(models.ActivityPub)
	resp, err := c.getParsedResponse("GET", "/activitypub/actor", jsonHeader, nil, actor)
	return actor, resp, err
}

// SendToInstanceActorInbox sends a message to the instance actor's inbox.
func (c *Client) SendToInstanceActorInbox() (*Response, error) {
	_, resp, err := c.getResponse("POST", "/activitypub/actor/inbox", jsonHeader, nil)
	return resp, err
}

// PostInstanceActorOutbox posts to the instance actor's outbox.
func (c *Client) PostInstanceActorOutbox() (models.ForgeOutbox, *Response, error) {
	var out models.ForgeOutbox
	resp, err := c.getParsedResponse("POST", "/activitypub/actor/outbox", jsonHeader, nil, &out)
	return out, resp, err
}

// ---------------------------------------------------------------------------
// Repository actor
// ---------------------------------------------------------------------------

// GetRepositoryActor returns a repository's ActivityPub actor.
func (c *Client) GetRepositoryActor(repositoryID int64) (*models.ActivityPub, *Response, error) {
	actor := new(models.ActivityPub)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/activitypub/repository-id/%d", repositoryID), jsonHeader, nil, actor)
	return actor, resp, err
}

// SendToRepositoryInbox sends a ForgeLike (or other activity) to a repository's inbox.
func (c *Client) SendToRepositoryInbox(repositoryID int64, like models.ForgeLike) (*Response, error) {
	body, err := json.Marshal(&like)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/activitypub/repository-id/%d/inbox", repositoryID), jsonHeader, bytes.NewReader(body))
	return resp, err
}

// PostRepositoryOutbox posts to a repository's outbox.
func (c *Client) PostRepositoryOutbox(repositoryID int64) (models.ForgeOutbox, *Response, error) {
	var out models.ForgeOutbox
	resp, err := c.getParsedResponse("POST", fmt.Sprintf("/activitypub/repository-id/%d/outbox", repositoryID), jsonHeader, nil, &out)
	return out, resp, err
}

// ---------------------------------------------------------------------------
// Person (user) actor
// ---------------------------------------------------------------------------

// GetPersonActor returns a user's ActivityPub actor.
func (c *Client) GetPersonActor(userID int64) (*models.ActivityPub, *Response, error) {
	actor := new(models.ActivityPub)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/activitypub/user-id/%d", userID), jsonHeader, nil, actor)
	return actor, resp, err
}

// SendToPersonInbox sends a message to a user's inbox.
func (c *Client) SendToPersonInbox(userID int64) (*Response, error) {
	_, resp, err := c.getResponse("POST", fmt.Sprintf("/activitypub/user-id/%d/inbox", userID), jsonHeader, nil)
	return resp, err
}

// GetPersonFeed returns a user's ActivityPub outbox feed.
func (c *Client) GetPersonFeed(userID int64) (models.ForgeOutbox, *Response, error) {
	var out models.ForgeOutbox
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/activitypub/user-id/%d/outbox", userID), jsonHeader, nil, &out)
	return out, resp, err
}

// GetPersonActivityNote returns the Activity Note object for a user activity.
func (c *Client) GetPersonActivityNote(userID int64, activityID string) (*models.ActivityPub, *Response, error) {
	if err := escapeValidatePathSegments(&activityID); err != nil {
		return nil, nil, err
	}
	actor := new(models.ActivityPub)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/activitypub/user-id/%d/activities/%s", userID, activityID), jsonHeader, nil, actor)
	return actor, resp, err
}

// GetPersonActivity returns the Activity object for a user activity.
func (c *Client) GetPersonActivity(userID int64, activityID string) (*models.ActivityPub, *Response, error) {
	if err := escapeValidatePathSegments(&activityID); err != nil {
		return nil, nil, err
	}
	actor := new(models.ActivityPub)
	resp, err := c.getParsedResponse("GET", fmt.Sprintf("/activitypub/user-id/%d/activities/%s/activity", userID, activityID), jsonHeader, nil, actor)
	return actor, resp, err
}

// ---------------------------------------------------------------------------
// Remote follow
// ---------------------------------------------------------------------------

// FollowActivityPub follows a remote ActivityPub account identified by the
// Target URI in opt (e.g. an actor or repository inbox URL).
func (c *Client) FollowActivityPub(opt models.APRemoteFollowOption) (*Response, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}
	_, resp, err := c.getResponse("POST", "/user/activitypub/follow", jsonHeader, bytes.NewReader(body))
	return resp, err
}
