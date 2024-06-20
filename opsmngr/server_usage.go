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
	"io"
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	serverUsageBasePath  = "api/public/v1.0/usage"
	serverUsageGroupPath = serverUsageBasePath + "/groups/%s"
	serverUsageOrgPath   = serverUsageBasePath + "/organizations/%s"
)

// ServerUsageService is an interface for using the Server Usage Service
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/
type ServerUsageService interface {
	GenerateDailyUsageSnapshot(context.Context) (*Response, error)
	ListAllHostAssignment(context.Context, *ServerTypeOptions) (*HostAssignments, *Response, error)
	ProjectHostAssignments(context.Context, string, *ServerTypeOptions) (*HostAssignments, *Response, error)
	OrganizationHostAssignments(context.Context, string, *ServerTypeOptions) (*HostAssignments, *Response, error)
	GetServerTypeProject(context.Context, string) (*ServerType, *Response, error)
	GetServerTypeOrganization(context.Context, string) (*ServerType, *Response, error)
	UpdateProjectServerType(context.Context, string, *ServerTypeRequest) (*Response, error)
	UpdateOrganizationServerType(context.Context, string, *ServerTypeRequest) (*Response, error)
}

// ServerUsageServiceOp provides an implementation of the ServerUsageService.
type ServerUsageServiceOp service

var _ ServerUsageService = &ServerUsageServiceOp{}

// ServerUsageReportService interface is an interface for downloading the service usage report.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/create-one-report/
type ServerUsageReportService interface {
	Download(context.Context, *ServerTypeOptions, io.Writer) (*Response, error)
}

// ServerUsageReportServiceOp handles communication with the Log Collection Jobs download method of the
// MongoDB Ops Manager API.
type ServerUsageReportServiceOp struct {
	Client atlas.GZipRequestDoer
}

var _ ServerUsageReportService = &ServerUsageReportServiceOp{}

// ServerTypeOptions specifies the optional parameters to List methods that
// support pagination.
type ServerTypeOptions struct {
	ListOptions
	StartDate  string `url:"startDate,omitempty"`
	EndDate    string `url:"endDate,omitempty"`
	FileFormat string `url:"fileFormat,omitempty"`
	Redact     *bool  `url:"redact,omitempty"`
}

// HostAssignments represents a paginated collection of HostAssignment.
type HostAssignments struct {
	Links      []*atlas.Link     `json:"links"`
	Results    []*HostAssignment `json:"results"`
	TotalCount int               `json:"totalCount"`
}

// HostAssignment represents a HostAssignment in the MongoDB Ops Manager API.
type HostAssignment struct {
	GroupID      string                   `json:"groupId,omitempty"`
	Hostname     string                   `json:"hostname,omitempty"`
	Processes    []*HostAssignmentProcess `json:"processes,omitempty"`
	ServerType   *ServerType              `json:"serverType,omitempty"`
	MemSizeMB    int64                    `json:"memSizeMB,omitempty"` //nolint:tagliatelle // Bytes vs bits
	IsChargeable *bool                    `json:"isChargeable,omitempty"`
}

// HostAssignmentProcess represents a HostAssignmentProcess in the MongoDB Ops Manager API.
type HostAssignmentProcess struct {
	Cluster                  string `json:"cluster,omitempty"`
	GroupName                string `json:"groupName,omitempty"`
	OrgName                  string `json:"orgName,omitempty"`
	GroupID                  string `json:"groupId,omitempty"`
	Name                     string `json:"name,omitempty"`
	HasConflictingServerType bool   `json:"hasConflictingServerType"`
	ProcessType              int64  `json:"processType"`
}

// ServerType represents a ServerType in the MongoDB Ops Manager API.
type ServerType struct {
	Name  string `json:"name,omitempty"`
	Label string `json:"label,omitempty"`
}

// ServerTypeRequest contains request body parameters for Server Usage Service.
type ServerTypeRequest struct {
	ServerType *ServerType `json:"serverType,omitempty"`
}

// GenerateDailyUsageSnapshot generates snapshot of usage for the processes Ops Manager manages..
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/generate-daily-usage-snapshot/
func (s *ServerUsageServiceOp) GenerateDailyUsageSnapshot(ctx context.Context) (*Response, error) {
	path := fmt.Sprintf("%s/%s", serverUsageBasePath, "dailyCapture")
	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// ListAllHostAssignment retrieves all host assignments.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/list-all-host-assignments/
func (s *ServerUsageServiceOp) ListAllHostAssignment(ctx context.Context, options *ServerTypeOptions) (*HostAssignments, *Response, error) {
	basePath := fmt.Sprintf("%s/%s", serverUsageBasePath, "assignments")
	path, err := setQueryParams(basePath, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(HostAssignments)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// ProjectHostAssignments retrieves all host assignments in a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/list-all-host-assignments-in-one-project/
func (s *ServerUsageServiceOp) ProjectHostAssignments(ctx context.Context, groupID string, options *ServerTypeOptions) (*HostAssignments, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(serverUsageGroupPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "hosts")
	path, err := setQueryParams(path, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(HostAssignments)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// OrganizationHostAssignments retrieves all host assignments in a organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/list-all-host-assignments-in-one-organization/
func (s *ServerUsageServiceOp) OrganizationHostAssignments(ctx context.Context, orgID string, options *ServerTypeOptions) (*HostAssignments, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(serverUsageOrgPath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, "hosts")
	path, err := setQueryParams(path, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(HostAssignments)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// GetServerTypeProject retrieves one default server type for one project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/get-default-server-type-for-one-project/
func (s *ServerUsageServiceOp) GetServerTypeProject(ctx context.Context, groupID string) (*ServerType, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(serverUsageGroupPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "defaultServerType")

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ServerType)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// GetServerTypeOrganization retrieves one default server type for one organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/get-default-server-type-for-one-organization/
func (s *ServerUsageServiceOp) GetServerTypeOrganization(ctx context.Context, orgID string) (*ServerType, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(serverUsageOrgPath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, "defaultServerType")

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ServerType)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// UpdateProjectServerType update the default server type for one project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/update-default-server-type-for-one-project/
func (s *ServerUsageServiceOp) UpdateProjectServerType(ctx context.Context, groupID string, serverType *ServerTypeRequest) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(serverUsageGroupPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "defaultServerType")

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, serverType)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// UpdateOrganizationServerType update the default server type for one organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/update-default-server-type-for-one-organization/
func (s *ServerUsageServiceOp) UpdateOrganizationServerType(ctx context.Context, groupID string, serverType *ServerTypeRequest) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(serverUsageOrgPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "defaultServerType")

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, serverType)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// Download downloads a compressed report of server usage in a given timeframe.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/usage/create-one-report/
func (s *ServerUsageReportServiceOp) Download(ctx context.Context, options *ServerTypeOptions, out io.Writer) (*Response, error) {
	path := fmt.Sprintf("%s/%s", serverUsageBasePath, "report")
	path, err := setQueryParams(path, options)
	if err != nil {
		return nil, err
	}
	req, err := s.Client.NewGZipRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, out)

	return resp, err
}
