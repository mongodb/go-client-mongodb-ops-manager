// Copyright 2020 MongoDB Inc
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
	usersBasePath    = "api/public/v1.0/users"
	orgUsersBasePath = "api/public/v1.0/orgs/%s/users"
)

// UsersService provides access to the user related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/users/
type UsersService interface {
	Get(context.Context, string) (*User, *Response, error)
	GetByName(context.Context, string) (*User, *Response, error)
	Create(context.Context, *User) (*User, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// UsersServiceOp provides an implementation of the UsersService interface.
type UsersServiceOp service

var _ UsersService = &UsersServiceOp{}

// User wrapper for a user response, augmented with a few extra fields.
type User struct {
	ID           string        `json:"id,omitempty"`
	EmailAddress string        `json:"emailAddress,omitempty"`
	FirstName    string        `json:"firstName,omitempty"`
	LastName     string        `json:"lastName,omitempty"`
	Links        []*atlas.Link `json:"links,omitempty"`
	MobileNumber string        `json:"mobileNumber,omitempty"`
	Password     string        `json:"password,omitempty"`
	TeamIds      []string      `json:"teamIds,omitempty"` //nolint
	Roles        []*UserRole   `json:"roles,omitempty"`
	Username     string        `json:"username"`
	Country      string        `json:"country,omitempty"`
}

// UserRole denotes a single user role.
type UserRole struct {
	RoleName string `json:"roleName"`
	GroupID  string `json:"groupId,omitempty"`
	OrgID    string `json:"orgId,omitempty"`
}

// UsersResponse represents a array of users.
type UsersResponse struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*User       `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// Get gets a single user by ID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/user-get-by-id/
func (s *UsersServiceOp) Get(ctx context.Context, userID string) (*User, *Response, error) {
	if userID == "" {
		return nil, nil, NewArgError("userID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", usersBasePath, userID)

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

// GetByName gets a single user by name.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/user-get-by-name/
func (s *UsersServiceOp) GetByName(ctx context.Context, username string) (*User, *Response, error) {
	if username == "" {
		return nil, nil, NewArgError("username", "must be set")
	}

	path := fmt.Sprintf("%s/byName/%s", usersBasePath, username)

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

// Create creates a new IAM user.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/user-create/
func (s *UsersServiceOp) Create(ctx context.Context, createRequest *User) (*User, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, usersBasePath, createRequest)
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

// Delete deletes a user.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/users/delete-one-user/
func (s *UsersServiceOp) Delete(ctx context.Context, userID string) (*Response, error) {
	if userID == "" {
		return nil, NewArgError("userID", "must be set")
	}

	basePath := fmt.Sprintf("%s/%s", usersBasePath, userID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
