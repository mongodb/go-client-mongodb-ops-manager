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
)

const (
	teamsBasePath = "api/public/v1.0/orgs/%s/teams"
)

// TeamsService provides access to the team related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/
type TeamsService interface {
	List(context.Context, string, *ListOptions) ([]Team, *Response, error)
	Get(context.Context, string, string) (*Team, *Response, error)
	GetOneTeamByName(context.Context, string, string) (*Team, *Response, error)
	GetTeamUsersAssigned(context.Context, string, string) ([]*User, *Response, error)
	Create(context.Context, string, *Team) (*Team, *Response, error)
	Rename(context.Context, string, string, string) (*Team, *Response, error)
	UpdateTeamRoles(context.Context, string, string, *TeamUpdateRoles) ([]TeamRoles, *Response, error)
	AddUsersToTeam(context.Context, string, string, []string) ([]*User, *Response, error)
	RemoveUserToTeam(context.Context, string, string, string) (*Response, error)
	RemoveTeamFromOrganization(context.Context, string, string) (*Response, error)
	RemoveTeamFromProject(context.Context, string, string) (*Response, error)
}

// TeamsServiceOp provides an implementation of the TeamsService interface.
type TeamsServiceOp service

var _ TeamsService = &TeamsServiceOp{}

// TeamsResponse represents a array of project.
type TeamsResponse struct {
	Links      []*Link `json:"links"`
	Results    []Team  `json:"results"`
	TotalCount int     `json:"totalCount"`
}

// Team defines an Ops Manager team structure.
type Team struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name"`
	Usernames []string `json:"usernames,omitempty"`
}

// OMUserAssigned represents the user assigned to the project.
type OMUserAssigned struct {
	Links      []*Link `json:"links"`
	Results    []User  `json:"results"`
	TotalCount int     `json:"totalCount"`
}

// TeamUpdateRoles update request body.
type TeamUpdateRoles struct {
	RoleNames []string `json:"roleNames"`
}

// TeamUpdateRolesResponse update roles response.
type TeamUpdateRolesResponse struct {
	Links      []*Link     `json:"links"`
	Results    []TeamRoles `json:"results"`
	TotalCount int         `json:"totalCount"`
}

// TeamRoles List of roles for a team.
type TeamRoles struct {
	Links     []*Link  `json:"links"`
	RoleNames []string `json:"roleNames"`
	TeamID    string   `json:"teamId"`
}

// List gets all teams.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-get-all/
func (s *TeamsServiceOp) List(ctx context.Context, orgID string, listOptions *ListOptions) ([]Team, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	path := fmt.Sprintf(teamsBasePath, orgID)

	// Add query params from listOptions
	path, err := setQueryParams(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(TeamsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets a single team in the organization by team ID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-get-one-by-id/
func (s *TeamsServiceOp) Get(ctx context.Context, orgID, teamID string) (*Team, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if teamID == "" {
		return nil, nil, NewArgError("teamID", "must be set")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, teamID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Team)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetOneTeamByName gets a single project by its name.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-get-one-by-name/
func (s *TeamsServiceOp) GetOneTeamByName(ctx context.Context, orgID, teamName string) (*Team, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if teamName == "" {
		return nil, nil, NewArgError("teamName", "must be set")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/byName/%s", basePath, teamName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Team)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetTeamUsersAssigned gets all the users assigned to a team.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-get-all-users/
func (s *TeamsServiceOp) GetTeamUsersAssigned(ctx context.Context, orgID, teamID string) ([]*User, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if teamID == "" {
		return nil, nil, NewArgError("teamID", "must be set")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s/users", basePath, teamID)

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

// Create creates a team.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-create-one/
func (s *TeamsServiceOp) Create(ctx context.Context, orgID string, createRequest *Team) (*Team, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, fmt.Sprintf(teamsBasePath, orgID), createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(Team)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Rename renames a team.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-rename-one/
func (s *TeamsServiceOp) Rename(ctx context.Context, orgID, teamID, teamName string) (*Team, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if teamID == "" {
		return nil, nil, NewArgError("teamID", "must be set")
	}
	if teamName == "" {
		return nil, nil, NewArgError("teamName", "cannot be nil")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, teamID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, map[string]interface{}{
		"name": teamName,
	})
	if err != nil {
		return nil, nil, err
	}

	root := new(Team)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// UpdateTeamRoles Update the roles of a team in an Ops Manager project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-update-roles/
func (s *TeamsServiceOp) UpdateTeamRoles(ctx context.Context, orgID, teamID string, updateTeamRolesRequest *TeamUpdateRoles) ([]TeamRoles, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if teamID == "" {
		return nil, nil, NewArgError("teamID", "must be set")
	}
	if updateTeamRolesRequest == nil {
		return nil, nil, NewArgError("updateTeamRolesRequest", "cannot be nil")
	}

	path := fmt.Sprintf("api/public/v1.0/groups/%s/teams/%s", orgID, teamID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateTeamRolesRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(TeamUpdateRolesResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// AddUsersToTeam adds a users from the organization associated with {ORG-ID} to the team with ID {TEAM-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-add-user/
func (s *TeamsServiceOp) AddUsersToTeam(ctx context.Context, orgID, teamID string, usersID []string) ([]*User, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}
	if teamID == "" {
		return nil, nil, NewArgError("teamID", "must be set")
	}
	if len(usersID) < 1 {
		return nil, nil, NewArgError("usersID", "cannot empty at leas one userID must be set")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s/users", basePath, teamID)

	users := make([]map[string]interface{}, len(usersID))
	for i, id := range usersID {
		users[i] = map[string]interface{}{"id": id}
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, users)

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

// RemoveUserToTeam removes the specified user from the specified team.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-remove-user/
func (s *TeamsServiceOp) RemoveUserToTeam(ctx context.Context, orgID, teamID, userID string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("orgID", "must be set")
	}
	if teamID == "" {
		return nil, NewArgError("teamID", "must be set")
	}
	if userID == "" {
		return nil, NewArgError("userID", "cannot be nil")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s/users/%s", basePath, teamID, userID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveTeamFromOrganization deletes the team with ID {TEAM-ID} from the organization specified to {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-delete-one/
func (s *TeamsServiceOp) RemoveTeamFromOrganization(ctx context.Context, orgID, teamID string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("orgID", "must be set")
	}
	if teamID == "" {
		return nil, NewArgError("teamID", "cannot be nil")
	}

	basePath := fmt.Sprintf(teamsBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, teamID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RemoveTeamFromProject removes the specified team from the specified project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/teams/teams-remove-from-project/
func (s *TeamsServiceOp) RemoveTeamFromProject(ctx context.Context, groupID, teamID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}
	if teamID == "" {
		return nil, NewArgError("teamID", "cannot be nil")
	}

	path := fmt.Sprintf("api/public/v1.0/groups/%s/teams/%s", groupID, teamID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
