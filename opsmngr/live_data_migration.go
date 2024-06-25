// Copyright 2021 MongoDB Inc
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

const liveMigrationBasePath = "api/public/v1.0/orgs/%s/liveExport/migrationLink"

type LinkToken struct {
	LinkToken string `json:"linkToken,omitempty"` // Ops Manager-generated token that links the source (Cloud Manager or Ops Manager) and destination (Atlas) clusters for migration.
}

// LiveDataMigrationService is an interface for interfacing with the Live Migration
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/cloud-migration/
type LiveDataMigrationService interface {
	ConnectOrganizations(context.Context, string, *LinkToken) (*ConnectionStatus, *Response, error)
	DeleteConnection(context.Context, string) (*Response, error)
	ConnectionStatus(context.Context, string) (*ConnectionStatus, *Response, error)
}

// ConnectionStatus represents the response of LiveDataMigrationService.ConnectOrganizations and LiveDataMigrationService.ConnectionStatus.
type ConnectionStatus struct {
	Status string `json:"status,omitempty"` // Status represents the state of the connection that exists between this organization and the target cluster in the MongoDB Ops Manager organization.
}

// LiveDataMigrationServiceOp provides an implementation of the LiveDataMigrationService interface.
type LiveDataMigrationServiceOp service

var _ LiveDataMigrationService = &LiveDataMigrationServiceOp{}

// ConnectionStatus returns the status of the connection between the specified source Ops Manager organization and the target MongoDB Ops Manager organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/cloud-migration/return-the-status-of-the-organization-link/
func (s *LiveDataMigrationServiceOp) ConnectionStatus(ctx context.Context, orgID string) (*ConnectionStatus, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(liveMigrationBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, "status")

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ConnectionStatus)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// ConnectOrganizations connects the source Ops Manager organization with a target MongoDB Ops Manager organization.
func (s *LiveDataMigrationServiceOp) ConnectOrganizations(ctx context.Context, orgID string, linkToken *LinkToken) (*ConnectionStatus, *Response, error) {
	if orgID == "" {
		return nil, nil, NewArgError("orgID", "must be set")
	}

	if linkToken == nil {
		return nil, nil, NewArgError("linkToken", "must be set")
	}

	path := fmt.Sprintf(liveMigrationBasePath, orgID)
	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, linkToken)
	if err != nil {
		return nil, nil, err
	}

	root := new(ConnectionStatus)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// DeleteConnection removes the connection between the source Ops Manager organization and the target MongoDB Ops Manager organization.
func (s *LiveDataMigrationServiceOp) DeleteConnection(ctx context.Context, orgID string) (*Response, error) {
	if orgID == "" {
		return nil, NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf(liveMigrationBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
