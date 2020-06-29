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
	projectBasePath = "groups"
)

// ProjectsService provides access to the project related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/
type ProjectsService interface {
	List(context.Context, *atlas.ListOptions) (*Projects, *atlas.Response, error)
	Get(context.Context, string) (*Project, *atlas.Response, error)
	GetByName(context.Context, string) (*Project, *atlas.Response, error)
	Create(context.Context, *Project) (*Project, *atlas.Response, error)
	Delete(context.Context, string) (*atlas.Response, error)
}

// ProjectsServiceOp provides an implementation of the ProjectsService interface
type ProjectsServiceOp service

var _ ProjectsService = &ProjectsServiceOp{}

// HostCount number of processes per project.
type HostCount struct {
	Arbiter   int `json:"arbiter"`
	Config    int `json:"config"`
	Master    int `json:"master"`
	Mongos    int `json:"mongos"`
	Primary   int `json:"primary"`
	Secondary int `json:"secondary"`
	Slave     int `json:"slave"`
}

// Project represents the structure of a project.
type Project struct {
	ActiveAgentCount int           `json:"activeAgentCount,omitempty"`
	HostCounts       *HostCount    `json:"hostCounts,omitempty"`
	ID               string        `json:"id,omitempty"`
	LastActiveAgent  string        `json:"lastActiveAgent,omitempty"`
	Links            []*atlas.Link `json:"links,omitempty"`
	Name             string        `json:"name,omitempty"`
	OrgID            string        `json:"orgId,omitempty"`
	PublicAPIEnabled bool          `json:"publicApiEnabled,omitempty"`
	ReplicaSetCount  int           `json:"replicaSetCount,omitempty"`
	ShardCount       int           `json:"shardCount,omitempty"`
	Tags             []*string     `json:"tags,omitempty"`
}

// Projects represents a array of project
type Projects struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Project    `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// List gets all projects.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/get-all-groups-for-current-user/
func (s *ProjectsServiceOp) List(ctx context.Context, opts *atlas.ListOptions) (*Projects, *atlas.Response, error) {
	path, err := setQueryParams(projectBasePath, opts)
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

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get gets a single project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/get-one-group-by-id/
func (s *ProjectsServiceOp) Get(ctx context.Context, groupID string) (*Project, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", projectBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Project)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetByName gets a single project by its name.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/get-one-group-by-name/
func (s *ProjectsServiceOp) GetByName(ctx context.Context, groupName string) (*Project, *atlas.Response, error) {
	if groupName == "" {
		return nil, nil, atlas.NewArgError("groupName", "must be set")
	}

	path := fmt.Sprintf("%s/byName/%s", projectBasePath, groupName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Project)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create creates a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/create-one-group/
func (s *ProjectsServiceOp) Create(ctx context.Context, createRequest *Project) (*Project, *atlas.Response, error) {
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, projectBasePath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Project)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/delete-one-group/
func (s *ProjectsServiceOp) Delete(ctx context.Context, projectID string) (*atlas.Response, error) {
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
