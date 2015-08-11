/*
 * GitLab Hook Server (C) 2015  Ryan Armstrong <ryan@cavaliercoder.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// See: https://gitlab.com/gitlab-org/gitlab-ce/blob/master/doc/web_hooks/web_hooks.md
package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	// "time"
)

type HookRequest struct {
	RequestID   string `json:"_"`
	EventHeader string `json:"-"`
	RemoteAddr  string `json:"-"`

	After             string     `json:"after"`
	Before            string     `json:"before"`
	Commits           []Commit   `json:"commits"`
	ObjectKind        string     `json:"object_kind"`
	ProjectID         int        `json:"project_id"`
	Ref               string     `json:"ref"`
	Repository        Repository `json:"repository"`
	TotalCommitsCount int        `json:"total_commits_count"`
	UserEmail         string     `json:"user_email"`
	UserID            int        `json:"user_id"`
	UserName          string     `json:"user_name"`

	User             User             `json:"user"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
}

type Commit struct {
	Author struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"author"`
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	URL       string `json:"url"`
}

func (c *Commit) String() string {
	return c.ID[0:7]
}

type Repository struct {
	Description     string `json:"description"`
	GitHTTPURL      string `json:"git_http_url"`
	GitSSHURL       string `json:"git_ssh_url"`
	Homepage        string `json:"homepage"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	VisibilityLevel int    `json:"visibility_level"`
}

func (c *Repository) String() string {
	return c.Name
}

type User struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type ObjectAttributes struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	AssigneeID int    `json:"assignee_id"`
	AuthorID   int    `json:"author_id"`
	ProjectID  int    `json:"project_id"`
	// CreatedAt   time.Time   `json:"created_at"`
	// UpdatedAt   time.Time   `json:"updated_at"`
	Position    int         `json:"position"`
	BranchName  interface{} `json:"branch_name"`
	Description string      `json:"description"`
	MilestoneID interface{} `json:"milestone_id"`
	State       string      `json:"state"`
	Iid         int         `json:"iid"`
	URL         string      `json:"url"`
	Action      string      `json:"action"`
}

func NewHookRequest(r *http.Request) (*HookRequest, error) {
	hookRequest := HookRequest{
		RequestID: uuid.New()[:7],
	}

	if r != nil {
		hookRequest.EventHeader = r.Header.Get(EVENT_HEADER)
		hookRequest.RemoteAddr = r.RemoteAddr

		// check for known event header
		switch hookRequest.EventHeader {
		case PUSH_EVENT, TAG_EVENT, ISSUE_EVENT:
		default:
			return nil, errors.New(fmt.Sprintf("Unknown event header: %s", hookRequest.EventHeader))
		}

		// parse JSON request body
		if r.Body != nil {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&hookRequest)
			if err != nil {
				return nil, err
			}
		}
	}

	return &hookRequest, nil
}

func (c *HookRequest) String() string {
	switch c.ObjectKind {
	case "push":
		return fmt.Sprintf("Push: %s (ID: %d) - %s -> %s", c.Repository.Name, c.ProjectID, c.Ref, c.After[0:7])
	case "tag_push":
		return fmt.Sprintf("Tag: %s (ID: %d) - %s -> %s", c.Repository.Name, c.ProjectID, c.Ref, c.After[0:7])
	case "issue":
		return fmt.Sprintf("Issue: %s (ID: %d)", c.ObjectAttributes.Title, c.ObjectAttributes.Iid)

	default:
		return c.RequestID
	}
}
