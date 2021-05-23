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

const backupAdministratorSyncBasePath = "admin/backup/sync/mongoConfigs"

// SyncService is an interface for using the Sync
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync-store-config/
type SyncStoreConfigService interface {
	List(context.Context, *atlas.ListOptions) (*BackupStores, *atlas.Response, error)
	Get(context.Context, string) (*BackupStore, *atlas.Response, error)
	Create(context.Context, *BackupStore) (*BackupStore, *atlas.Response, error)
	Update(context.Context, string, *BackupStore) (*BackupStore, *atlas.Response, error)
	Delete(context.Context, string) (*atlas.Response, error)
}

// SyncStoreConfigServiceOp provides an implementation of the SyncStoreConfigService interface.
type SyncStoreConfigServiceOp service

var _ SyncStoreConfigService = &SyncStoreConfigServiceOp{}

// Get retrieves a Sync.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/get-one-sync-store-configuration-by-id/
func (s *SyncStoreConfigServiceOp) Get(ctx context.Context, syncID string) (*BackupStore, *atlas.Response, error) {
	if syncID == "" {
		return nil, nil, atlas.NewArgError("syncID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorSyncBasePath, syncID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List retrieves all the Syncs.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/get-all-sync-store-configurations/
func (s *SyncStoreConfigServiceOp) List(ctx context.Context, options *atlas.ListOptions) (*BackupStores, *atlas.Response, error) {
	path, err := setQueryParams(backupAdministratorSyncBasePath, options)
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

// Create create a Sync.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/create-one-sync-store-configuration/
func (s *SyncStoreConfigServiceOp) Create(ctx context.Context, sync *BackupStore) (*BackupStore, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorSyncBasePath, sync)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates a Sync.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/update-one-sync-store-configuration/
func (s *SyncStoreConfigServiceOp) Update(ctx context.Context, syncID string, sync *BackupStore) (*BackupStore, *atlas.Response, error) {
	if syncID == "" {
		return nil, nil, atlas.NewArgError("syncID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorSyncBasePath, syncID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, sync)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete removes a Sync.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/delete-one-sync-store-configuration/
func (s *SyncStoreConfigServiceOp) Delete(ctx context.Context, syncID string) (*atlas.Response, error) {
	if syncID == "" {
		return nil, atlas.NewArgError("syncID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorSyncBasePath, syncID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
