// Copyright 2029 MongoDB Inc
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

const backupAdministratorBlockstoreBasePath = "admin/backup/snapshot/mongoConfigs"

// BlockstoreConfigService is an interface for using the Blockstore Configuration
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/blockstore-config/
type BlockstoreConfigService interface {
	List(context.Context, *atlas.ListOptions) (*BackupStores, *atlas.Response, error)
	Get(context.Context, string) (*BackupStore, *atlas.Response, error)
	Create(context.Context, *BackupStore) (*BackupStore, *atlas.Response, error)
	Update(context.Context, string, *BackupStore) (*BackupStore, *atlas.Response, error)
	Delete(context.Context, string) (*atlas.Response, error)
}

// BackupConfigsServiceOp provides an implementation of the BackupConfigsService interface
type BlockstoreConfigServiceOp service

var _ BlockstoreConfigService = &BlockstoreConfigServiceOp{}

// BackupStore represents a Balockstore, Oplog and Sync in the MongoDB Ops Manager API
type BackupStore struct {
	AdminBackupConfig
	LoadFactor    int64  `json:"loadFactor,omitempty"`
	MaxCapacityGB int64  `json:"maxCapacityGB,omitempty"`
	Provisioned   *bool  `json:"provisioned,omitempty"`
	SyncSource    string `json:"syncSource,omitempty"`
	Username      string `json:"username,omitempty"`
}

// BackupStores represents a paginated collection of BackupStore
type BackupStores struct {
	Links      []*atlas.Link  `json:"links"`
	Results    []*BackupStore `json:"results"`
	TotalCount int            `json:"totalCount"`
}

// Get retrieves a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/get-one-blockstore-configuration-by-id/
func (s *BlockstoreConfigServiceOp) Get(ctx context.Context, blockstoreID string) (*BackupStore, *atlas.Response, error) {
	if blockstoreID == "" {
		return nil, nil, atlas.NewArgError("blockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorBlockstoreBasePath, blockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List retrieves all the blockstores.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/get-all-blockstore-configurations/
func (s *BlockstoreConfigServiceOp) List(ctx context.Context, options *atlas.ListOptions) (*BackupStores, *atlas.Response, error) {
	path, err := setQueryParams(backupAdministratorBlockstoreBasePath, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStores)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create create a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/create-one-blockstore-configuration/
func (s *BlockstoreConfigServiceOp) Create(ctx context.Context, blockstore *BackupStore) (*BackupStore, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorBlockstoreBasePath, blockstore)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/update-one-blockstore-configuration/
func (s *BlockstoreConfigServiceOp) Update(ctx context.Context, blockstoreID string, blockstore *BackupStore) (*BackupStore, *atlas.Response, error) {
	if blockstoreID == "" {
		return nil, nil, atlas.NewArgError("blockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorBlockstoreBasePath, blockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, blockstore)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete removes a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/delete-one-blockstore-configuration/
func (s *BlockstoreConfigServiceOp) Delete(ctx context.Context, blockstoreID string) (*atlas.Response, error) {
	if blockstoreID == "" {
		return nil, atlas.NewArgError("blockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorBlockstoreBasePath, blockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
