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

const backupAdministratorProjectJobBasePath = "api/public/v1.0/admin/backup/groups"

// ProjectJobConfigService is an interface for using the Project Job
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/backup-group-config/
type ProjectJobConfigService interface {
	List(context.Context, *ListOptions) (*ProjectJobs, *Response, error)
	Get(context.Context, string) (*ProjectJob, *Response, error)
	Update(context.Context, string, *ProjectJob) (*ProjectJob, *Response, error)
}

// ProjectJobConfigServiceOp provides an implementation of the ProjectJobConfigService interface.
type ProjectJobConfigServiceOp service

var _ ProjectJobConfigService = &ProjectJobConfigServiceOp{}

// AdminBackupConfig contains the common fields of backup administrator structs.
type AdminBackupConfig struct {
	ID                   string   `json:"id,omitempty"`
	URI                  string   `json:"uri,omitempty"`
	WriteConcern         string   `json:"writeConcern,omitempty"`
	Labels               []string `json:"labels,omitempty"`
	SSL                  *bool    `json:"ssl,omitempty"`
	AssignmentEnabled    *bool    `json:"assignmentEnabled,omitempty"`
	EncryptedCredentials *bool    `json:"encryptedCredentials,omitempty"`
	UsedSize             int64    `json:"usedSize,omitempty"`
}

// ProjectJob represents a Backup Project Configuration Job in the MongoDB Ops Manager API.
type ProjectJob struct {
	AdminBackupConfig
	KMIPClientCertPassword string         `json:"kmipClientCertPassword,omitempty"`
	KMIPClientCertPath     string         `json:"kmipClientCertPath,omitempty"`
	LabelFilter            []string       `json:"labelFilter,omitempty"`
	SyncStoreFilter        []string       `json:"syncStoreFilter,omitempty"`
	DaemonFilter           []*Machine     `json:"daemonFilter,omitempty"`
	OplogStoreFilter       []*StoreFilter `json:"oplogStoreFilter,omitempty"`
	SnapshotStoreFilter    []*StoreFilter `json:"snapshotStoreFilter,omitempty"`
}

// StoreFilter represents a StoreFilter in the MongoDB Ops Manager API.
type StoreFilter struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// ProjectJobs represents a paginated collection of ProjectJob.
type ProjectJobs struct {
	Links      []*Link       `json:"links"`
	Results    []*ProjectJob `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// List retrieves the configurations of all project’s backup jobs.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/get-all-sync-store-configurations/
func (s *ProjectJobConfigServiceOp) List(ctx context.Context, options *ListOptions) (*ProjectJobs, *Response, error) {
	path, err := setQueryParams(backupAdministratorProjectJobBasePath, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectJobs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get retrieves the configuration of one project’s backup jobs.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/groups/get-one-backup-group-configuration-by-id/
func (s *ProjectJobConfigServiceOp) Get(ctx context.Context, projectJobID string) (*ProjectJob, *Response, error) {
	if projectJobID == "" {
		return nil, nil, NewArgError("projectJobID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorProjectJobBasePath, projectJobID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectJob)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates the configuration of one project’s backup jobs.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/groups/update-one-backup-group-configuration/
func (s *ProjectJobConfigServiceOp) Update(ctx context.Context, projectJobID string, projectJob *ProjectJob) (*ProjectJob, *Response, error) {
	if projectJobID == "" {
		return nil, nil, NewArgError("projectJobID", "must be set")
	}
	path := fmt.Sprintf("%s/%s", backupAdministratorProjectJobBasePath, projectJobID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, projectJob)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProjectJob)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
