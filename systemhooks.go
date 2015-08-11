package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// See: https://gitlab.com/gitlab-org/gitlab-ce/blob/master/doc/system_hooks/system_hooks.md

type SystemHookHeader struct {
	CreatedAt time.Time `json:"created_at"`
	EventName string    `json:"event_name"`
}

type SystemHook interface{}

func ParseSystemHook(r *http.Request) (SystemHook, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var header SystemHookHeader
	if err = json.Unmarshal(b, &header); err != nil {
		return nil, err
	}

	printf("system hook received: %s\n", header.EventName)

	switch header.EventName {
	case "project_create", "project_destroy":
		hook := ProjectHook{}
		err = json.Unmarshal(b, &hook)
		return &hook, err
	}

	return nil, nil
}

type ProjectHook struct {
	CreatedAt         time.Time `json:"created_at"`
	EventName         string    `json:"event_name"`
	Name              string    `json:"name"`
	OwnerEmail        string    `json:"owner_email"`
	OwnerName         string    `json:"owner_name"`
	Path              string    `json:"path"`
	PathWithNamespace string    `json:"path_with_namespace"`
	ProjectID         int       `json:"project_id"`
	ProjectVisibility string    `json:"project_visibility"`
}

type TeamMemberHook struct {
	CreatedAt         time.Time `json:"created_at"`
	EventName         string    `json:"event_name"`
	ProjectAccess     string    `json:"project_access"`
	ProjectID         int       `json:"project_id"`
	ProjectName       string    `json:"project_name"`
	ProjectPath       string    `json:"project_path"`
	UserEmail         string    `json:"user_email"`
	UserName          string    `json:"user_name"`
	UserID            int       `json:"user_id"`
	ProjectVisibility string    `json:"project_visibility"`
}

type UserHook struct {
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	EventName string    `json:"event_name"`
	Name      string    `json:"name"`
	UserID    int       `json:"user_id"`
}

type KeyHook struct {
	EventName string `json:"event_name"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	Key       string `json:"key"`
	ID        int    `json:"id"`
}

type GroupHook struct {
	CreatedAt  time.Time `json:"created_at"`
	EventName  string    `json:"event_name"`
	Name       string    `json:"name"`
	OwnerEmail string    `json:"owner_email"`
	OwnerName  string    `json:"owner_name"`
	Path       string    `json:"path"`
	GroupID    int       `json:"group_id"`
}

type GroupMemberHook struct {
	CreatedAt   time.Time `json:"created_at"`
	EventName   string    `json:"event_name"`
	GroupAccess string    `json:"group_access"`
	GroupID     int       `json:"group_id"`
	GroupName   string    `json:"group_name"`
	GroupPath   string    `json:"group_path"`
	UserEmail   string    `json:"user_email"`
	UserName    string    `json:"user_name"`
	UserID      int       `json:"user_id"`
}
