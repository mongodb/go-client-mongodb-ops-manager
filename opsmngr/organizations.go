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
	orgsBasePath = "api/public/v1.0/orgs"
)

// OrganizationsService provides access to the organization related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/
type OrganizationsService interface {
	List(context.Context, *atlas.OrganizationsListOptions) (*atlas.Organizations, *Response, error)
	ListUsers(context.Context, string, *atlas.ListOptions) (*UsersResponse, *Response, error)
	Get(context.Context, string) (*atlas.Organization, *Response, error)
	Projects(context.Context, string, *atlas.ProjectsListOptions) (*Projects, *Response, error)
	Create(context.Context, *atlas.Organization) (*atlas.Organization, *Response, error)
	Delete(context.Context, string) (*Response, error)
	Invitations(context.Context, string, *atlas.InvitationOptions) ([]*atlas.Invitation, *Response, error)
	Invitation(context.Context, string, string) (*atlas.Invitation, *Response, error)
	InviteUser(context.Context, string, *atlas.Invitation) (*atlas.Invitation, *Response, error)
	UpdateInvitation(context.Context, string, *atlas.Invitation) (*atlas.Invitation, *Response, error)
	UpdateInvitationByID(context.Context, string, string, *atlas.Invitation) (*atlas.Invitation, *Response, error)
	DeleteInvitation(context.Context, string, string) (*Response, error)
}

// OrganizationsServiceOp provides an implementation of the OrganizationsService interface.
type OrganizationsServiceOp service

var _ OrganizationsService = &OrganizationsServiceOp{}

// List gets all organizations.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-get-all/
func (s *OrganizationsServiceOp) List(ctx context.Context, opts *atlas.OrganizationsListOptions) (*atlas.Organizations, *Response, error) {
	path, err := setQueryParams(orgsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Organizations)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// ListUsers gets all users in an organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-get-all-users/
func (s *OrganizationsServiceOp) ListUsers(ctx context.Context, orgID string, opts *atlas.ListOptions) (*UsersResponse, *Response, error) {
	path := fmt.Sprintf(orgUsersBasePath, orgID)

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

	return root, resp, nil
}

// Get gets a single organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-get-one/
func (s *OrganizationsServiceOp) Get(ctx context.Context, orgID string) (*atlas.Organization, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", orgsBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Organization)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Projects gets all projects for the given organization ID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-get-all-projects/
func (s *OrganizationsServiceOp) Projects(ctx context.Context, orgID string, opts *atlas.ProjectsListOptions) (*Projects, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	basePath := fmt.Sprintf("%s/%s/groups", orgsBasePath, orgID)

	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Projects)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates an organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-create-one/
func (s *OrganizationsServiceOp) Create(ctx context.Context, createRequest *atlas.Organization) (*atlas.Organization, *Response, error) {
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, orgsBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Organization)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes an organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-delete-one/
func (s *OrganizationsServiceOp) Delete(ctx context.Context, orgID string) (*Response, error) {
	if orgID == "" {
		return nil, atlas.NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf("%s/%s", orgsBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
