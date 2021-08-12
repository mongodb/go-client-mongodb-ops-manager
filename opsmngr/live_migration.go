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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const liveMigrationBasePath = "api/public/v1.0/orgs/%s/liveExport/migrationLink"

// LiveMigrationService is an interface for interfacing with the Live Migration
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/cloud-migration/
type LiveMigrationService interface {
	Create(context.Context, string, *atlas.LinkToken) (*atlas.LiveMigration, *Response, error)
	Delete(context.Context, string) (*Response, error)
	Get(context.Context, string) (*atlas.LiveMigration, *Response, error)
}

// LiveMigrationOp provides an implementation of the LiveMigrationService interface.
type LiveMigrationServiceOp service

var _ LiveMigrationService = &LiveMigrationServiceOp{}

// Get returns the status of the connection between the specified source Ops Manager organization and the target MongoDB Atlas organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/cloud-migration/return-the-status-of-the-organization-link/
func (s *LiveMigrationServiceOp) Get(ctx context.Context, orgID string) (*atlas.LiveMigration, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	basePath := fmt.Sprintf(liveMigrationBasePath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, "status")

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.LiveMigration)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create connects the source Ops Manager organization with a target MongoDB Atlas organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/cloud-migration/link-the-organization-with-atlas/
func (s *LiveMigrationServiceOp) Create(ctx context.Context, orgID string, linkToken *atlas.LinkToken) (*atlas.LiveMigration, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}

	if linkToken == nil {
		return nil, nil, atlas.NewArgError("linkToken", "must be set")
	}

	path := fmt.Sprintf(liveMigrationBasePath, orgID)
	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, linkToken)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.LiveMigration)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete removes the connection between the source Ops Manager organization and the target MongoDB Atlas organization.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-delete-one/
func (s *LiveMigrationServiceOp) Delete(ctx context.Context, orgID string) (*Response, error) {
	if orgID == "" {
		return nil, atlas.NewArgError("orgID", "must be set")
	}

	path := fmt.Sprintf(liveMigrationBasePath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
