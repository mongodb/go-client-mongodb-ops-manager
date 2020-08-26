// Copyright 2019 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	usersBasePath = "groups/%s/users"
)

// UsersService provides access to the user related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/
type UsersService interface {
	List(context.Context, string, *atlas.ListOptions) ([]*User, *atlas.Response, error)
	Get(context.Context, string, string) (*User, *atlas.Response, error)
	GetByName(context.Context, string) (*User, *atlas.Response, error)
	Create(context.Context, *User) (*User, *atlas.Response, error)
	Delete(context.Context, string) (*atlas.Response, error)
}

// UsersServiceOp provides an implementation of the UsersService interface
type UsersServiceOp service

var _ UsersService = &UsersServiceOp{}

// User wrapper for a user response, augmented with a few extra fields
type User struct {
	Username     string        `json:"username"`
	Password     string        `json:"password,omitempty"`
	FirstName    string        `json:"firstName,omitempty"`
	LastName     string        `json:"lastName,omitempty"`
	EmailAddress string        `json:"emailAddress,omitempty"`
	ID           string        `json:"id,omitempty"`
	Links        []*atlas.Link `json:"links,omitempty"`
	Roles        []*UserRole   `json:"roles,omitempty"`
}

// UserRole denotes a single user role
type UserRole struct {
	RoleName string `json:"roleName"`
	GroupID  string `json:"groupId,omitempty"`
	OrgID    string `json:"orgId,omitempty"`
}

// UsersResponse represents a array of users
type UsersResponse struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*User    `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// List gets all users in a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/get-all-users-in-one-group/
func (s *UsersServiceOp) List(ctx context.Context, projectID string, opts *atlas.ListOptions) ([]*User, *atlas.Response, error) {
	path := fmt.Sprintf(usersBasePath, projectID)

	path, err := setQueryParams(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(UsersResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets a single user.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/user-get-by-id/
func (s *UsersServiceOp) Get(ctx context.Context, projectID, userID string) (*User, *atlas.Response, error) {
	if userID == "" {
		return nil, nil, atlas.NewArgError("userID", "must be set")
	}

	path := fmt.Sprintf("%s/%s/%s", usersBasePath, projectID, userID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(User)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetByName gets a single project by its name.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/get-one-group-by-name/
func (s *UsersServiceOp) GetByName(ctx context.Context, groupName string) (*User, *atlas.Response, error) {
	if groupName == "" {
		return nil, nil, atlas.NewArgError("groupName", "must be set")
	}

	path := fmt.Sprintf("%s/byName/%s", projectBasePath, groupName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(User)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates an IAM user.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/create-one-group/
func (s *UsersServiceOp) Create(ctx context.Context, createRequest *User) (*User, *atlas.Response, error) {
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, projectBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(User)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/delete-one-group/
func (s *UsersServiceOp) Delete(ctx context.Context, projectID string) (*atlas.Response, error) {
	if projectID == "" {
		return nil, atlas.NewArgError("projectID", "must be set")
	}

	basePath := fmt.Sprintf("%s/%s", projectBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
