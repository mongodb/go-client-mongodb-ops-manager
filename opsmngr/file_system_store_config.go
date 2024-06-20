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

const backupAdministratorFileSystemStoreConfigurationsBasePath = "api/public/v1.0/admin/backup/snapshot/fileSystemConfigs"

// FileSystemStoreConfigService is an interface for using the File System Store Configuration
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/file-system-store-config/
type FileSystemStoreConfigService interface {
	List(context.Context, *ListOptions) (*FileSystemStoreConfigurations, *Response, error)
	Get(context.Context, string) (*FileSystemStoreConfiguration, *Response, error)
	Create(context.Context, *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *Response, error)
	Update(context.Context, string, *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// FileSystemStoreConfigServiceOp provides an implementation of the FileSystemStoreConfigService interface.
type FileSystemStoreConfigServiceOp service

var _ FileSystemStoreConfigService = &FileSystemStoreConfigServiceOp{}

// FileSystemStoreConfiguration represents a File System Store Configuration in the MongoDB Ops Manager API.
type FileSystemStoreConfiguration struct {
	BackupStore
	MMAPV1CompressionSetting string `json:"mmapv1CompressionSetting,omitempty"`
	StorePath                string `json:"storePath,omitempty"`
	WTCompressionSetting     string `json:"wtCompressionSetting,omitempty"`
}

// FileSystemStoreConfigurations represents a paginated collection of FileSystemStoreConfiguration.
type FileSystemStoreConfigurations struct {
	Links      []*Link                         `json:"links"`
	Results    []*FileSystemStoreConfiguration `json:"results"`
	TotalCount int                             `json:"totalCount"`
}

// Get retrieves a File System Store Configuration.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/get-one-file-system-store-configuration-by-id/
func (s *FileSystemStoreConfigServiceOp) Get(ctx context.Context, fileSystemID string) (*FileSystemStoreConfiguration, *Response, error) {
	if fileSystemID == "" {
		return nil, nil, NewArgError("fileSystemID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorFileSystemStoreConfigurationsBasePath, fileSystemID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FileSystemStoreConfiguration)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List retrieves the configurations of all file system stores.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/get-all-file-system-store-configurations/
func (s *FileSystemStoreConfigServiceOp) List(ctx context.Context, opts *ListOptions) (*FileSystemStoreConfigurations, *Response, error) {
	path, err := setQueryParams(backupAdministratorFileSystemStoreConfigurationsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FileSystemStoreConfigurations)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create configures one new file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/create-one-file-system-store-configuration/
func (s *FileSystemStoreConfigServiceOp) Create(ctx context.Context, fileSystem *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorFileSystemStoreConfigurationsBasePath, fileSystem)
	if err != nil {
		return nil, nil, err
	}

	root := new(FileSystemStoreConfiguration)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates the configuration of one file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/update-one-file-system-store-configuration/
func (s *FileSystemStoreConfigServiceOp) Update(ctx context.Context, fileSystemID string, fileSystem *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *Response, error) {
	if fileSystemID == "" {
		return nil, nil, NewArgError("fileSystemID", "must be set")
	}
	path := fmt.Sprintf("%s/%s", backupAdministratorFileSystemStoreConfigurationsBasePath, fileSystemID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, fileSystem)
	if err != nil {
		return nil, nil, err
	}

	root := new(FileSystemStoreConfiguration)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete deletes the configuration of one file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/delete-one-file-system-store-configuration/
func (s *FileSystemStoreConfigServiceOp) Delete(ctx context.Context, fileSystemID string) (*Response, error) {
	if fileSystemID == "" {
		return nil, NewArgError("fileSystemID", "must be set")
	}
	path := fmt.Sprintf("%s/%s", backupAdministratorFileSystemStoreConfigurationsBasePath, fileSystemID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
