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

// OrganizationsService is an interface for interfacing with the Organizations
// endpoints of the MongoDB Atlas API.
// See more: https://docs.cloudmanager.mongodb.com/reference/api/organizations/
type OrganizationsService interface {
	GetAllOrganizations(context.Context) (*Organizations, *atlas.Response, error)
	GetOneOrganization(context.Context, string) (*Organization, *atlas.Response, error)
	GetProjects(context.Context, string) (*Projects, *atlas.Response, error)
	Create(context.Context, *Organization) (*Organization, *atlas.Response, error)
	Delete(context.Context, string) (*atlas.Response, error)
}

// OrganizationsServiceOp handles communication with the Projects related methods of the
// MongoDB Atlas API
type OrganizationsServiceOp struct {
	client *Client
}

var _ OrganizationsService = &OrganizationsServiceOp{}

// Organization represents the structure of an organization.
type Organization struct {
	ID    string        `json:"id,omitempty"`
	Links []*atlas.Link `json:"links,omitempty"`
	Name  string        `json:"name,omitempty"`
}

// Organizations represents a array of organization
type Organizations struct {
	Links      []*atlas.Link   `json:"links"`
	Results    []*Organization `json:"results"`
	TotalCount int             `json:"totalCount"`
}

// GetAllOrganizations gets all organizations.
// See more: https://docs.cloudmanager.mongodb.com/reference/api/organizations/organization-get-all/
func (s *OrganizationsServiceOp) GetAllOrganizations(ctx context.Context) (*Organizations, *atlas.Response, error) {

	req, err := s.client.NewRequest(ctx, http.MethodGet, orgsBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organizations)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// GetOneOrganization gets a single organization.
// See more: https://docs.cloudmanager.mongodb.com/reference/api/organizations/organization-get-one/
func (s *OrganizationsServiceOp) GetOneOrganization(ctx context.Context, orgID string) (*Organization, *atlas.Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", orgsBasePath, orgID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organization)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetProjects gets all projects for the given organization ID
// See more: https://docs.cloudmanager.mongodb.com/reference/api/organizations/organization-get-all-projects/
func (s *OrganizationsServiceOp) GetProjects(ctx context.Context, orgID string) (*Projects, *atlas.Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf("%s/%s/groups", orgsBasePath, orgID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Projects)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates an organization.
// See more: https://docs.cloudmanager.mongodb.com/reference/api/organizations/organization-create-one/
func (s *OrganizationsServiceOp) Create(ctx context.Context, createRequest *Organization) (*Organization, *atlas.Response, error) {
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, orgsBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Organization)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes an organization.
// See more: https://docs.cloudmanager.mongodb.com/reference/api/organizations/organization-delete-one/
func (s *OrganizationsServiceOp) Delete(ctx context.Context, orgID string) (*atlas.Response, error) {
	if orgID == "" {
		return nil, atlas.NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf("%s/%s", orgsBasePath, orgID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}
