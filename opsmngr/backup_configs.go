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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	backupConfigsBasePath = "groups/%s/backupConfigs"
)

// BackupConfigsService is an interface for using the Backup Configurations
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/backup-configurations/
type BackupConfigsService interface {
	List(context.Context, string, *atlas.ListOptions) (*BackupConfigs, *atlas.Response, error)
	Get(context.Context, string, string) (*BackupConfig, *atlas.Response, error)
	Update(context.Context, string, string, *BackupConfig) (*BackupConfig, *atlas.Response, error)
}

// BackupConfigsServiceOp provides an implementation of the BackupConfigsService interface
type BackupConfigsServiceOp service

var _ BackupConfigsService = &BackupConfigsServiceOp{}

// BackupConfigs represents a Backup configuration in the MongoDB Ops Manager API
type BackupConfig struct {
	GroupID            string    `json:"groupId,omitempty"`
	ClusterID          string    `json:"clusterId,omitempty"`
	StatusName         string    `json:"statusName,omitempty"`
	StorageEngineName  string    `json:"storageEngineName,omitempty"`
	EncryptionEnabled  bool      `json:"encryptionEnabled,omitempty"`
	SSLEnabled         bool      `json:"sslEnabled,omitempty"`
	ExcludedNamespaces []*string `json:"excludedNamespaces,omitempty"`
	IncludedNamespaces []*string `json:"includedNamespaces,omitempty"`
}

// BackupConfigs represents an array of BackupConfig
type BackupConfigs struct {
	Links      []*atlas.Link   `json:"links"`
	Results    []*BackupConfig `json:"results"`
	TotalCount int             `json:"totalCount"`
}

// List gets all backup configurations.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/backup/get-all-backup-configs-for-group/
func (s *BackupConfigsServiceOp) List(ctx context.Context, groupID string, opts *atlas.ListOptions) (*BackupConfigs, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(backupConfigsBasePath, groupID)
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupConfigs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get retrieves a single backup configuration by cluster ID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/backup/get-one-backup-config-by-cluster-id/
func (s *BackupConfigsServiceOp) Get(ctx context.Context, groupID, clusterID string) (*BackupConfig, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if clusterID == "" {
		return nil, nil, atlas.NewArgError("clusterID", "must be set")
	}

	basePath := fmt.Sprintf(backupConfigsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, clusterID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupConfig)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates a single backup configuration.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/backup/update-backup-config/
func (s *BackupConfigsServiceOp) Update(ctx context.Context, groupID, clusterID string, backupConfig *BackupConfig) (*BackupConfig, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if clusterID == "" {
		return nil, nil, atlas.NewArgError("clusterID", "must be set")
	}

	basePath := fmt.Sprintf(backupConfigsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, clusterID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, backupConfig)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupConfig)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
