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

const backupAdministratorOplogBasePath = "admin/backup/oplog/mongoConfigs"

// OplogStoreConfigService is an interface for using the Oplog
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog-store-config/
type OplogStoreConfigService interface {
	List(context.Context, *atlas.ListOptions) (*BackupStores, *atlas.Response, error)
	Get(context.Context, string) (*BackupStore, *atlas.Response, error)
	Create(context.Context, *BackupStore) (*BackupStore, *atlas.Response, error)
	Update(context.Context, string, *BackupStore) (*BackupStore, *atlas.Response, error)
	Delete(context.Context, string) (*atlas.Response, error)
}

// BackupConfigsServiceOp provides an implementation of the BackupConfigsService interface
type OplogStoreConfigServiceOp service

var _ OplogStoreConfigService = &OplogStoreConfigServiceOp{}

// Get retrieves a Oplog.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog/mongoConfigs/get-one-oplog-configuration-by-id/
func (s *OplogStoreConfigServiceOp) Get(ctx context.Context, oplogID string) (*BackupStore, *atlas.Response, error) {
	if oplogID == "" {
		return nil, nil, atlas.NewArgError("oplogID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorOplogBasePath, oplogID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List retrieves all the Oplogs.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog/mongoConfigs/get-all-oplog-configurations/
func (s *OplogStoreConfigServiceOp) List(ctx context.Context, options *atlas.ListOptions) (*BackupStores, *atlas.Response, error) {
	path, err := setQueryParams(backupAdministratorOplogBasePath, options)
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

// Create create a Oplog.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog/mongoConfigs/create-one-oplog-configuration/
func (s *OplogStoreConfigServiceOp) Create(ctx context.Context, oplog *BackupStore) (*BackupStore, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorOplogBasePath, oplog)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates a Oplog.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog/mongoConfigs/update-one-oplog-configuration/
func (s *OplogStoreConfigServiceOp) Update(ctx context.Context, oplogID string, oplog *BackupStore) (*BackupStore, *atlas.Response, error) {
	if oplogID == "" {
		return nil, nil, atlas.NewArgError("oplogID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorOplogBasePath, oplogID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, oplog)
	if err != nil {
		return nil, nil, err
	}

	root := new(BackupStore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete removes a Oplog.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/delete-one-s3-blockstore-configuration/
func (s *OplogStoreConfigServiceOp) Delete(ctx context.Context, oplogID string) (*atlas.Response, error) {
	if oplogID == "" {
		return nil, atlas.NewArgError("oplogID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorOplogBasePath, oplogID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
