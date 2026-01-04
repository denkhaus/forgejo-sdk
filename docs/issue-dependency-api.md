# Issue Dependency API Documentation

## Overview

This document describes the Forgejo API endpoints for managing issue dependencies. Issue dependencies allow you to express blocking relationships between issues - i.e., "Issue A cannot be completed until Issue B is resolved."

**Target API Version**: Forgejo 13.0.3+

**Key Concepts**:
- **Dependencies** (also called "blocked by"): Issues that block the current issue from being completed
- **Blocks**: Issues that are blocked by the current issue

## API Endpoints

### 1. List Issue Dependencies
**Endpoint**: `GET /repos/{owner}/{repo}/issues/{index}/dependencies`

**Description**: Lists all issues that block the given issue. These are the dependencies that must be resolved before the current issue can be completed.

**Path Parameters**:
- `owner` (string): Repository owner
- `repo` (string): Repository name
- `index` (int64): Issue number (index)

**Query Parameters** (optional):
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: server config, max: server config)

**Response**: `200 OK`
```json
[
  {
    "id": 12345,
    "number": 101,
    "title": "Implement authentication feature",
    "body": "Description of the issue...",
    "state": "open",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-16T14:20:00Z",
    "closed_at": null,
    "user": { ... },
    "assignees": [ ... ],
    "labels": [ ... ],
    "milestone": { ... },
    "html_url": "https://forgejo.example.com/owner/repo/issues/101",
    "url": "https://forgejo.example.com/api/v1/repos/owner/repo/issues/101"
  }
]
```

**SDK Method Signature**:
```go
func (c *Client) ListIssueDependencies(owner, repo string, index int64, opt ListDependenciesOptions) ([]*models.Issue, *Response, error)
```

---

### 2. List Blocked Issues
**Endpoint**: `GET /repos/{owner}/{repo}/issues/{index}/blocks`

**Description**: Lists all issues that are blocked by the given issue. These are the issues that cannot be completed until the current issue is resolved.

**Path Parameters**:
- `owner` (string): Repository owner
- `repo` (string): Repository name
- `index` (int64): Issue number (index)

**Query Parameters** (optional):
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: server config, max: server config)

**Response**: `200 OK`
```json
[
  {
    "id": 12346,
    "number": 102,
    "title": "Add user profile page",
    "body": "Depends on issue #101...",
    "state": "open",
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-16T14:20:00Z",
    "closed_at": null,
    "user": { ... },
    "assignees": [ ... ],
    "labels": [ ... ],
    "milestone": { ... },
    "html_url": "https://forgejo.example.com/owner/repo/issues/102",
    "url": "https://forgejo.example.com/api/v1/repos/owner/repo/issues/102"
  }
]
```

**SDK Method Signature**:
```go
func (c *Client) ListBlockedIssues(owner, repo string, index int64, opt ListBlockedOptions) ([]*models.Issue, *Response, error)
```

---

### 3. Create Issue Dependency
**Endpoint**: `POST /repos/{owner}/{repo}/issues/{index}/dependencies`

**Description**: Creates a dependency relationship, making the issue in the URL depend on the issue specified in the request body. In other words, "the issue at {index} is now blocked by {newDependency}".

**Path Parameters**:
- `owner` (string): Repository owner
- `repo` (string): Repository name
- `index` (int64): Issue number (the issue that will be blocked)

**Request Body**:
```json
{
  "newDependency": "dependency-type"  // This is the TYPE of dependency
}
```

**Correct Request Format** (based on common API patterns):
```json
{
  "issue": 101  // or { "id": 101 } - the issue that blocks
}
```

**Response**: `201 Created`
Returns the created dependency information or the updated issue.

**SDK Method Signature**:
```go
type CreateDependencyOption struct {
    // The issue ID or index that will block the current issue
    Issue int64 `json:"issue"`
}

func (opt CreateDependencyOption) Validate() error {
    if opt.Issue == 0 {
        return fmt.Errorf("issue ID is required")
    }
    return nil
}

func (c *Client) CreateIssueDependency(owner, repo string, index int64, opt CreateDependencyOption) (*Response, error)
```

---

### 4. Remove Issue Dependency
**Endpoint**: `DELETE /repos/{owner}/{repo}/issues/{index}/dependencies`

**Description**: Removes a dependency relationship between issues.

**Path Parameters**:
- `owner` (string): Repository owner
- `repo` (string): Repository name
- `index` (int64): Issue number (the issue that is blocked)

**Query Parameters**:
- `dependency` (string, required): The dependency to remove (format TBD - likely issue ID or type)

**Response**: `204 No Content`

**SDK Method Signature**:
```go
func (c *Client) RemoveIssueDependency(owner, repo string, index, dependency int64) (*Response, error)
```

---

### 5. Create Blocking Relationship
**Endpoint**: `POST /repos/{owner}/{repo}/issues/{index}/blocks`

**Description**: Creates a blocking relationship, making the issue specified in the request body be blocked by the issue in the URL. This is the inverse of creating a dependency.

**Path Parameters**:
- `owner` (string): Repository owner
- `repo` (string): Repository name
- `index` (int64): Issue number (the blocking issue)

**Request Body**:
```json
{
  "issue": 102  // The issue that will be blocked
}
```

**Response**: `201 Created`

**SDK Method Signature**:
```go
type CreateBlockingOption struct {
    Issue int64 `json:"issue"`
}

func (opt CreateBlockingOption) Validate() error {
    if opt.Issue == 0 {
        return fmt.Errorf("issue ID is required")
    }
    return nil
}

func (c *Client) CreateIssueBlocking(owner, repo string, index int64, opt CreateBlockingOption) (*Response, error)
```

---

### 6. Remove Blocking Relationship
**Endpoint**: `DELETE /repos/{owner}/{repo}/issues/{index}/blocks`

**Description**: Removes a blocking relationship. This is the inverse of removing a dependency.

**Path Parameters**:
- `owner` (string): Repository owner
- `repo` (string): Repository name
- `index` (int64): Issue number (the blocking issue)

**Request Body** (likely required):
```json
{
  "issue": 102  // The issue to unblock
}
```

**Response**: `204 No Content`

**SDK Method Signature**:
```go
func (c *Client) RemoveIssueBlocking(owner, repo string, index, blockedIssue int64) (*Response, error)
```

---

## Required Models

### models.Issue

The main model used for API responses. Key fields include:

```go
type Issue struct {
    ID           int64        `json:"id"`
    Index        int64        `json:"number"`      // Issue number (public-facing)
    Title        string       `json:"title"`
    Body         string       `json:"body,omitempty"`
    State        StateType    `json:"state,omitempty"`   // "open" or "closed"
    HTMLURL      string       `json:"html_url,omitempty"`
    URL          string       `json:"url,omitempty"`

    // Timestamps (strfmt.DateTime = RFC3339)
    Created      strfmt.DateTime `json:"created_at,omitempty"`
    Updated      strfmt.DateTime `json:"updated_at,omitempty"`
    Closed       strfmt.DateTime `json:"closed_at,omitempty"`
    Deadline     strfmt.DateTime `json:"due_date,omitempty"`

    // References
    User         *User              `json:"user,omitempty"`
    Assignees    []*User            `json:"assignees"`
    Labels       []*Label           `json:"labels"`
    Milestone    *Milestone         `json:"milestone,omitempty"`
    PullRequest  *PullRequestMeta   `json:"pull_request,omitempty"`
    Repository   *RepositoryMeta    `json:"repository,omitempty"`
    Assignee     *User              `json:"assignee,omitempty"`

    // Additional fields
    IsLocked     bool               `json:"is_locked,omitempty"`
    Comments     int64              `json:"comments,omitempty"`
    PinOrder     int64              `json:"pin_order,omitempty"`
    Ref          string             `json:"ref,omitempty"`

    // Original author (for migrated issues)
    OriginalAuthor     string `json:"original_author,omitempty"`
    OriginalAuthorID   int64  `json:"original_author_id,omitempty"`

    // Attachments (using "assets" as JSON key)
    Attachments        []*Attachment `json:"assets"`
}
```

---

## Implementation Notes

### Pagination

All list endpoints support standard pagination:

```go
type ListDependenciesOptions struct {
    ListOptions  // Embeds Page and PageSize fields
}
```

- Setting `Page = -1` disables pagination
- Default `Page` is 1 if not set
- `PageSize` defaults to server's `DEFAULT_PAGING_NUM`
- Maximum `PageSize` is server's `MAX_RESPONSE_ITEMS`

### Path Validation

All endpoints must validate path segments to prevent path injection:

```go
if err := escapeValidatePathSegments(&owner, &repo); err != nil {
    return nil, nil, err
}
```

### HTTP Methods

- **GET**: Use `getParsedResponse()` for parsed JSON responses
- **POST**: Use `getParsedResponse()` with `jsonHeader` and request body
- **DELETE**: Use `getResponse()` (no response body expected)

### Response Header

```go
var jsonHeader = http.Header{"Content-Type": []string{"application/json"}}
```

---

## Edge Cases and Validation Rules

### 1. Circular Dependencies

**Server Behavior**: The API should prevent creating circular dependencies (e.g., A depends on B, B depends on A).

**Expected Response**: `400 Bad Request` or `422 Unprocessable Entity`

**SDK Action**: Propagate the error to the caller; do not attempt to detect client-side.

### 2. Self-Dependency

**Validation**: An issue cannot depend on itself.

**Expected Response**: `400 Bad Request`

**SDK Action**: Consider adding client-side validation for better UX:

```go
func (opt CreateDependencyOption) Validate() error {
    if opt.Issue == 0 {
        return fmt.Errorf("issue ID is required")
    }
    return nil
}

// Caller should check: if dependencyIssue == index { return error }
```

### 3. Non-Existent Issues

**Validation**: Referenced issues must exist.

**Expected Response**: `404 Not Found`

**SDK Action**: Propagate server response to caller.

### 4. Cross-Repository Dependencies

**Configuration**: Controlled by server config `ALLOW_CROSS_REPOSITORY_DEPENDENCIES` (default: true).

**Behavior**: When enabled, dependencies can reference issues in other repositories where the user has access.

**SDK Implication**: This SDK implementation will focus on same-repository dependencies first. Cross-repo support would require additional parameters and models.

### 5. Duplicate Dependencies

**Validation**: Creating the same dependency twice should be idempotent.

**Expected Response**: `200 OK` or `409 Conflict` (depending on server implementation)

### 6. Permission Checks

**Required Permissions**: User must have:
- Read access to both issues
- Write access to the issue being modified (to add/remove dependencies)

**Expected Response**: `403 Forbidden` if insufficient permissions

### 7. Closed vs Open Issues

**Validation**: Dependencies can be created regardless of issue state (open or closed).

**Note**: It's common practice to allow dependencies on closed issues (e.g., "re-open this bug fix").

---

## Configuration Options

### Server-Side Settings

These Forgejo configuration options affect dependency behavior:

- `DEFAULT_ENABLE_DEPENDENCIES`: Enable dependencies by default (default: true)
- `ALLOW_CROSS_REPOSITORY_DEPENDENCIES`: Allow cross-repo dependencies (default: true)

---

## Examples

### Example 1: List all dependencies of an issue

```go
client, _ := forgejo.NewClient("https://forgejo.example.com", forgejo.SetToken("..."))

deps, resp, err := client.ListIssueDependencies("owner", "repo", 105, forgejo.ListDependenciesOptions{
    ListOptions: forgejo.ListOptions{Page: 1, PageSize: 50},
})
if err != nil {
    log.Fatalf("Error: %v", err)
}

for _, dep := range deps {
    fmt.Printf("Issue #%d: %s (state: %s)\n", dep.Index, dep.Title, dep.State)
}
```

### Example 2: Create a dependency

```go
// Make issue #105 depend on (be blocked by) issue #101
opt := forgejo.CreateDependencyOption{Issue: 101}
resp, err := client.CreateIssueDependency("owner", "repo", 105, opt)
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

### Example 3: List issues blocked by current issue

```go
// Find all issues that are blocked by issue #101
blocked, resp, err := client.ListBlockedIssues("owner", "repo", 101, forgejo.ListBlockedOptions{
    ListOptions: forgejo.ListOptions{},
})
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Printf("Issue #101 blocks %d other issues\n", len(blocked))
```

### Example 4: Remove a dependency

```go
// Remove the dependency on issue #101 from issue #105
resp, err := client.RemoveIssueDependency("owner", "repo", 105, 101)
if err != nil {
    log.Fatalf("Error: %v", err)
}
```

---

## References

- Forgejo Java SDK: https://github.com/client-api/forgejo-java
- Forgejo Documentation: https://forgejo.org/docs/
- Forgejo Issue #446: Feature request for dependency graph visualization
- Forgejo Config: `ALLOW_CROSS_REPOSITORY_DEPENDENCIES`
