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
	atlas "go.mongodb.org/atlas/mongodbatlas"
	"net/http"
)

const (
	backupAdministratorBlockstoreBasePath = "admin/backup/snapshot/mongoConfigs"
)


// BackupAdministratorService is an interface for using the Backup Administrator
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/nav/administration-backup/
type BackupAdministratorService interface {
	ListBlockstores(context.Context, *atlas.ListOptions) (*Blockstores, *atlas.Response, error)
	GetBlockstore(context.Context, string) (*Blockstore, *atlas.Response, error)
	CreateBlockstore(context.Context, *Blockstore) (*Blockstore, *atlas.Response, error)
	UpdateBlockstore(context.Context, string, *Blockstore) (*Blockstore, *atlas.Response, error)
	DeleteBlockstore(context.Context, string) (*atlas.Response, error)
}

// BackupConfigsServiceOp provides an implementation of the BackupConfigsService interface
type BackupAdministratorServiceOp service

var _ BackupAdministratorService = &BackupAdministratorServiceOp{}

// Blockstore represents a Blockstore in the MongoDB Ops Manager API
type Blockstore struct {
	ID            string   `json:"id,omitempty"`
	AssignmentEnabled          bool   `json:"assignmentEnabled,omitempty"`
	EncryptedCredentials         bool   `json:"encryptedCredentials,omitempty"`
	LoadFactor  int64   `json:"loadFactor,omitempty"`
	MaxCapacityGB  int64     `json:"maxCapacityGB,omitempty"`
	URI         string     `json:"uri,omitempty"`
	Labels []string `json:"labels,omitempty"`
	SSL bool `json:"ssl,omitempty"`
	UsedSize  int64   `json:"usedSize,omitempty"`
	WriteConcern           string   `json:"writeConcern,omitempty"`
	Provisioned        bool     `json:"provisioned,omitempty"`
	SyncSource         string   `json:"syncSource,omitempty"`
	Username           string   `json:"username,omitempty"`
	Links      []*atlas.Link   `json:"links"`
}

// Blockstores represents an array of Blockstore
type Blockstores struct {
	Links      []*atlas.Link   `json:"links"`
	Results    []*Blockstore `json:"results"`
	TotalCount int             `json:"totalCount"`
}

// Get retrieves a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/get-one-blockstore-configuration-by-id/
func (s *BackupAdministratorServiceOp) GetBlockstore(ctx context.Context, blockstoreID string) (*Blockstore, *atlas.Response, error) {
	if blockstoreID == "" {
		return nil, nil, atlas.NewArgError("blockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorBlockstoreBasePath, blockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Blockstore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// ListBlockstores retrieves all the blockstores.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/get-all-blockstore-configurations/
func (s *BackupAdministratorServiceOp) ListBlockstores(ctx context.Context, options *atlas.ListOptions) (*Blockstores, *atlas.Response, error) {
	path, err := setQueryParams(backupAdministratorBlockstoreBasePath, options)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Blockstores)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// CreateBlockstore create a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/create-one-blockstore-configuration/
func (s *BackupAdministratorServiceOp) CreateBlockstore(ctx context.Context, blockstore *Blockstore) (*Blockstore, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorBlockstoreBasePath, blockstore)
	if err != nil {
		return nil, nil, err
	}

	root := new(Blockstore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// UpdateBlockstore updates a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/update-one-blockstore-configuration/
func (s *BackupAdministratorServiceOp) UpdateBlockstore(ctx context.Context, blockstoreID string, blockstore *Blockstore) (*Blockstore, *atlas.Response, error) {
	if blockstoreID == "" {
		return nil, nil, atlas.NewArgError("blockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorBlockstoreBasePath, blockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, blockstore)
	if err != nil {
		return nil, nil, err
	}

	root := new(Blockstore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// DeleteBlockstore removes a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/mongoConfigs/delete-one-blockstore-configuration/
func (s *BackupAdministratorServiceOp) DeleteBlockstore(ctx context.Context, blockstoreID string) (*atlas.Response, error) {
	if blockstoreID == "" {
		return  nil, atlas.NewArgError("blockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorBlockstoreBasePath, blockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return  nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return  resp, err
}