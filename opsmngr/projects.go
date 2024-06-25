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
)

const (
	projectBasePath = "api/public/v1.0/groups"
)

// ProjectsListOptions filtering options for projects.
type ProjectsListOptions struct {
	Name string `url:"name,omitempty"`
	ListOptions
}

type CreateProjectOptions struct {
	ProjectOwnerID string `url:"projectOwnerId,omitempty"` // Unique 24-hexadecimal digit string that identifies the Atlas user account to be granted the Project Owner role on the specified project.
}

// ProjectTeam represents the kind of role that has the team.
type ProjectTeam struct {
	TeamID    string   `json:"teamId,omitempty"`
	RoleNames []string `json:"roleNames,omitempty"`
}

// TeamsAssigned represents the one team assigned to the project.
type TeamsAssigned struct {
	Links      []*Link   `json:"links"`
	Results    []*Result `json:"results"`
	TotalCount int       `json:"totalCount"`
}

// Result is part og TeamsAssigned structure.
type Result struct {
	Links     []*Link  `json:"links"`
	RoleNames []string `json:"roleNames"`
	TeamID    string   `json:"teamId"`
}

// ProjectsService provides access to the project related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/
type ProjectsService interface {
	List(context.Context, *ListOptions) (*Projects, *Response, error)
	ListUsers(context.Context, string, *ListOptions) ([]*User, *Response, error)
	Get(context.Context, string) (*Project, *Response, error)
	GetByName(context.Context, string) (*Project, *Response, error)
	Create(context.Context, *Project, *CreateProjectOptions) (*Project, *Response, error)
	Delete(context.Context, string) (*Response, error)
	RemoveUser(context.Context, string, string) (*Response, error)
	AddTeamsToProject(context.Context, string, []*ProjectTeam) (*TeamsAssigned, *Response, error)
	GetTeams(context.Context, string, *ListOptions) (*TeamsAssigned, *Response, error)
	Invitations(context.Context, string, *InvitationOptions) ([]*Invitation, *Response, error)
	Invitation(context.Context, string, string) (*Invitation, *Response, error)
	InviteUser(context.Context, string, *Invitation) (*Invitation, *Response, error)
	UpdateInvitation(context.Context, string, *Invitation) (*Invitation, *Response, error)
	UpdateInvitationByID(context.Context, string, string, *Invitation) (*Invitation, *Response, error)
	DeleteInvitation(context.Context, string, string) (*Response, error)
}

// ProjectsServiceOp provides an implementation of the ProjectsService interface.
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

// LDAPGroupMapping for LDAP-backed Ops Manager,
// the mappings of LDAP groups to Ops Manager project roles.
// Only present for LDAP-backed Ops Manager.
type LDAPGroupMapping struct {
	RoleName   string   `json:"roleName"`
	LDAPGroups []string `json:"ldapGroups"`
}

// Project represents the structure of a project.
type Project struct {
	HostCounts                *HostCount          `json:"hostCounts,omitempty"`
	LDAPGroupMappings         []*LDAPGroupMapping `json:"ldapGroupMappings,omitempty"`
	Links                     []*Link             `json:"links,omitempty"`
	Name                      string              `json:"name,omitempty"`
	OrgID                     string              `json:"orgId,omitempty"`
	ID                        string              `json:"id,omitempty"`
	AgentAPIKey               string              `json:"agentApiKey,omitempty"`
	LastActiveAgent           string              `json:"lastActiveAgent,omitempty"`
	Tags                      []*string           `json:"tags,omitempty"`
	PublicAPIEnabled          bool                `json:"publicApiEnabled,omitempty"`
	WithDefaultAlertsSettings *bool               `json:"withDefaultAlertsSettings,omitempty"`
	ReplicaSetCount           int                 `json:"replicaSetCount,omitempty"`
	ShardCount                int                 `json:"shardCount,omitempty"`
	ActiveAgentCount          int                 `json:"activeAgentCount,omitempty"`
}

// Projects represents a array of project.
type Projects struct {
	Links      []*Link    `json:"links"`
	Results    []*Project `json:"results"`
	TotalCount int        `json:"totalCount"`
}

// List gets all projects.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/get-all-groups-for-current-user/
func (s *ProjectsServiceOp) List(ctx context.Context, opts *ListOptions) (*Projects, *Response, error) {
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

// ListUsers gets all users in a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/get-all-users-in-one-group/
func (s *ProjectsServiceOp) ListUsers(ctx context.Context, projectID string, opts *ListOptions) ([]*User, *Response, error) {
	path := fmt.Sprintf("%s/%s/users", projectBasePath, projectID)

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

// Get gets a single project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/get-one-group-by-id/
func (s *ProjectsServiceOp) Get(ctx context.Context, groupID string) (*Project, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
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
func (s *ProjectsServiceOp) GetByName(ctx context.Context, groupName string) (*Project, *Response, error) {
	if groupName == "" {
		return nil, nil, NewArgError("groupName", "must be set")
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
func (s *ProjectsServiceOp) Create(ctx context.Context, createRequest *Project, opts *CreateProjectOptions) (*Project, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path, err := setQueryParams(projectBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
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
func (s *ProjectsServiceOp) Delete(ctx context.Context, projectID string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}

	basePath := fmt.Sprintf("%s/%s", projectBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// RemoveUser removes a user from a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/remove-one-user-from-one-group/
func (s *ProjectsServiceOp) RemoveUser(ctx context.Context, projectID, userID string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}

	if userID == "" {
		return nil, NewArgError("userID", "must be set")
	}

	basePath := fmt.Sprintf("%s/%s/users/%s", projectBasePath, projectID, userID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, basePath, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// AddTeamsToProject adds teams to a project
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/project-add-team/
func (s *ProjectsServiceOp) AddTeamsToProject(ctx context.Context, projectID string, createRequest []*ProjectTeam) (*TeamsAssigned, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s/teams", projectBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(TeamsAssigned)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetTeams gets all teams in a project
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/groups/project-get-teams/
func (s *ProjectsServiceOp) GetTeams(ctx context.Context, projectID string, opts *ListOptions) (*TeamsAssigned, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s/teams", projectBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, opts)
	if err != nil {
		return nil, nil, err
	}

	root := new(TeamsAssigned)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
