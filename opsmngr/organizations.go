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

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	orgsBasePath = "orgs"
)

// OrganizationsService provides access to the organization related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/
type OrganizationsService interface {
	List(context.Context, *atlas.ListOptions) (*Organizations, *atlas.Response, error)
	Get(context.Context, string) (*Organization, *atlas.Response, error)
	GetProjects(context.Context, string, *atlas.ListOptions) (*Projects, *atlas.Response, error)
	Create(context.Context, *Organization) (*Organization, *atlas.Response, error)
	Delete(context.Context, string) (*atlas.Response, error)
}

// OrganizationsServiceOp provides an implementation of the OrganizationsService interface
type OrganizationsServiceOp service

var _ OrganizationsService = &OrganizationsServiceOp{}

// Organization represents the structure of an organization.
type Organization struct {
	ID    string        `json:"id,omitempty"`
	Links []*atlas.Link `json:"links,omitempty"`
	Name  string        `json:"name,omitempty"`
}

// Organizations represents an array of organization
type Organizations struct {
	Links      []*atlas.Link   `json:"links"`
	Results    []*Organization `json:"results"`
	TotalCount int             `json:"totalCount"`
}

// List gets all organizations.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-get-all/
func (s *OrganizationsServiceOp) List(ctx context.Context, opts *atlas.ListOptions) (*Organizations, *atlas.Response, error) {
	path, err := setQueryParams(orgsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organizations)
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
func (s *OrganizationsServiceOp) Get(ctx context.Context, orgID string) (*Organization, *atlas.Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", orgsBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organization)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetProjects gets all projects for the given organization ID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-get-all-projects/
func (s *OrganizationsServiceOp) GetProjects(ctx context.Context, orgID string, opts *atlas.ListOptions) (*Projects, *atlas.Response, error) {
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
func (s *OrganizationsServiceOp) Create(ctx context.Context, createRequest *Organization) (*Organization, *atlas.Response, error) {
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, orgsBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organization)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes an organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/organizations/organization-delete-one/
func (s *OrganizationsServiceOp) Delete(ctx context.Context, orgID string) (*atlas.Response, error) {
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
