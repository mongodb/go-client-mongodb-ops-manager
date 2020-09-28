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
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	backupAdministratorBlockstoreBasePath                    = "admin/backup/snapshot/mongoConfigs"
	backupAdministratorFileSystemStoreConfigurationsBasePath = "admin/backup/snapshot/fileSystemConfigs"
	backupAdministratorS3BlockstoreBasePath                  = "admin/backup/snapshot/s3Configs"
)

// BackupAdministratorService is an interface for using the Backup Administrator
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/nav/administration-backup/
type BackupAdministratorService interface {
	BlockstoreService
	FileSystemStoreConfigurationsService
	S3BlockstoreService
}

// BlockstoreService is an interface for using the Blockstore
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/blockstore-config/
type BlockstoreService interface {
	ListBlockstores(context.Context, *atlas.ListOptions) (*Blockstores, *atlas.Response, error)
	GetBlockstore(context.Context, string) (*Blockstore, *atlas.Response, error)
	CreateBlockstore(context.Context, *Blockstore) (*Blockstore, *atlas.Response, error)
	UpdateBlockstore(context.Context, string, *Blockstore) (*Blockstore, *atlas.Response, error)
	DeleteBlockstore(context.Context, string) (*atlas.Response, error)
}

// FileSystemStoreConfigurationsService is an interface for using the File System Store Configuration
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/file-system-store-config/
type FileSystemStoreConfigurationsService interface {
	ListFileSystemStoreConfigurations(context.Context, *atlas.ListOptions) (*FileSystemStoreConfigurations, *atlas.Response, error)
	GetFileSystemStoreConfiguration(context.Context, string) (*FileSystemStoreConfiguration, *atlas.Response, error)
	CreateFileSystemStoreConfiguration(context.Context, *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *atlas.Response, error)
	UpdateFileSystemStoreConfiguration(context.Context, string, *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *atlas.Response, error)
	DeleteFileSystemStoreConfiguration(context.Context, string) (*atlas.Response, error)
}

// S3BlockstoreService is an interface for using the S3BlockstoreService
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/s3-blockstore-config/
type S3BlockstoreService interface {
	ListS3Blockstores(context.Context, *atlas.ListOptions) (*S3Blockstores, *atlas.Response, error)
	GetS3Blockstore(context.Context, string) (*S3Blockstore, *atlas.Response, error)
	CreateS3Blockstore(context.Context, *S3Blockstore) (*S3Blockstore, *atlas.Response, error)
	UpdateS3Blockstore(context.Context, string, *S3Blockstore) (*S3Blockstore, *atlas.Response, error)
	DeleteS3Blockstore(context.Context, string) (*atlas.Response, error)
}

// BackupConfigsServiceOp provides an implementation of the BackupConfigsService interface
type BackupAdministratorServiceOp service

var _ BackupAdministratorService = &BackupAdministratorServiceOp{}

// S3Blockstore represents a S3Blockstore in the MongoDB Ops Manager API
type S3Blockstore struct {
	Blockstore
	AWSAccessKey           string `json:"awsAccessKey,omitempty"`
	AWSSecretKey           string `json:"awsSecretKey,omitempty"`
	DisableProxyS3         string `json:"disableProxyS3,omitempty"`
	PathStyleAccessEnabled bool   `json:"pathStyleAccessEnabled,omitempty"`
	S3AuthMethod           string `json:"s3AuthMethod,omitempty"`
	S3BucketEndpoint       string `json:"s3BucketEndpoint,omitempty"`
	S3BucketName           string `json:"s3BucketName,omitempty"`
	S3MaxConnections       int64  `json:"s3MaxConnections,omitempty"`
	AcceptedTos            bool   `json:"acceptedTos,omitempty"`
	SSEEnabled             bool   `json:"sseEnabled,omitempty"`
}

// S3Blockstores represents an array of S3Blockstore
type S3Blockstores struct {
	Links      []*atlas.Link   `json:"links"`
	Results    []*S3Blockstore `json:"results"`
	TotalCount int             `json:"totalCount"`
}

// Blockstore represents a Blockstore in the MongoDB Ops Manager API
type Blockstore struct {
	ID                   string        `json:"id,omitempty"`
	AssignmentEnabled    bool          `json:"assignmentEnabled,omitempty"`
	EncryptedCredentials bool          `json:"encryptedCredentials,omitempty"`
	LoadFactor           int64         `json:"loadFactor,omitempty"`
	MaxCapacityGB        int64         `json:"maxCapacityGB,omitempty"`
	URI                  string        `json:"uri,omitempty"`
	Labels               []string      `json:"labels,omitempty"`
	SSL                  bool          `json:"ssl,omitempty"`
	UsedSize             int64         `json:"usedSize,omitempty"`
	WriteConcern         string        `json:"writeConcern,omitempty"`
	Provisioned          bool          `json:"provisioned,omitempty"`
	SyncSource           string        `json:"syncSource,omitempty"`
	Username             string        `json:"username,omitempty"`
	Links                []*atlas.Link `json:"links"`
}

// Blockstores represents an array of Blockstore
type Blockstores struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Blockstore `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// FileSystemStoreConfiguration represents a File System Store Configuration in the MongoDB Ops Manager API
type FileSystemStoreConfiguration struct {
	ID                       string        `json:"id,omitempty"`
	Labels                   []string      `json:"labels,omitempty"`
	Links                    []*atlas.Link `json:"links"`
	LoadFactor               int64         `json:"loadFactor,omitempty"`
	MMAPV1CompressionSetting string        `json:"mmapv1CompressionSetting,omitempty"`
	StorePath                string        `json:"storePath,omitempty"`
	WTCompressionSetting     string        `json:"wtCompressionSetting,omitempty"`
	AssignmentEnabled        bool          `json:"assignmentEnabled,omitempty"`
}

// FileSystemStoreConfigurations represents an array of FileSystemStoreConfiguration
type FileSystemStoreConfigurations struct {
	Links      []*atlas.Link                   `json:"links"`
	Results    []*FileSystemStoreConfiguration `json:"results"`
	TotalCount int                             `json:"totalCount"`
}

// GetFileSystemStoreConfiguration retrieves a File System Store Configuration.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/get-one-file-system-store-configuration-by-id/
func (s *BackupAdministratorServiceOp) GetFileSystemStoreConfiguration(ctx context.Context, fileSystemID string) (*FileSystemStoreConfiguration, *atlas.Response, error) {
	if fileSystemID == "" {
		return nil, nil, atlas.NewArgError("fileSystemID", "must be set")
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

// ListFileSystemStoreConfigurations retrieves the configurations of all file system stores.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/get-all-file-system-store-configurations/
func (s *BackupAdministratorServiceOp) ListFileSystemStoreConfigurations(ctx context.Context, opts *atlas.ListOptions) (*FileSystemStoreConfigurations, *atlas.Response, error) {
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

// CreateFileSystemStoreConfiguration configures one new file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/create-one-file-system-store-configuration/
func (s *BackupAdministratorServiceOp) CreateFileSystemStoreConfiguration(ctx context.Context, fileSystem *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorFileSystemStoreConfigurationsBasePath, fileSystem)
	if err != nil {
		return nil, nil, err
	}

	root := new(FileSystemStoreConfiguration)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// UpdateFileSystemStoreConfiguration updates the configuration of one file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/update-one-file-system-store-configuration/
func (s *BackupAdministratorServiceOp) UpdateFileSystemStoreConfiguration(ctx context.Context, fileSystemID string, fileSystem *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *atlas.Response, error) {
	if fileSystemID == "" {
		return nil, nil, atlas.NewArgError("fileSystemID", "must be set")
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

// DeleteFileSystemStoreConfiguration deletes the configuration of one file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/delete-one-file-system-store-configuration/
func (s *BackupAdministratorServiceOp) DeleteFileSystemStoreConfiguration(ctx context.Context, fileSystemID string) (*atlas.Response, error) {
	if fileSystemID == "" {
		return nil, atlas.NewArgError("fileSystemID", "must be set")
	}
	path := fmt.Sprintf("%s/%s", backupAdministratorFileSystemStoreConfigurationsBasePath, fileSystemID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// GetBlockstore retrieves a blockstore.
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
	if err != nil {
		return nil, nil, err
	}
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

// GetS3Blockstore retrieves a GetS3Blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/get-one-s3-blockstore-configuration-by-id/
func (s *BackupAdministratorServiceOp) GetS3Blockstore(ctx context.Context, s3BlockstoreID string) (*S3Blockstore, *atlas.Response, error) {
	if s3BlockstoreID == "" {
		return nil, nil, atlas.NewArgError("s3BlockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorS3BlockstoreBasePath, s3BlockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(S3Blockstore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// ListS3Blockstores retrieves all the S3Blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/get-all-s3-blockstore-configurations/
func (s *BackupAdministratorServiceOp) ListS3Blockstores(ctx context.Context, options *atlas.ListOptions) (*S3Blockstores, *atlas.Response, error) {
	path, err := setQueryParams(backupAdministratorS3BlockstoreBasePath, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(S3Blockstores)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// CreateS3Blockstore create a S3Blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/create-one-s3-blockstore-configuration/
func (s *BackupAdministratorServiceOp) CreateS3Blockstore(ctx context.Context, blockstore *S3Blockstore) (*S3Blockstore, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorS3BlockstoreBasePath, blockstore)
	if err != nil {
		return nil, nil, err
	}

	root := new(S3Blockstore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// UpdateS3Blockstore updates a S3Blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/update-one-s3-blockstore-configuration/
func (s *BackupAdministratorServiceOp) UpdateS3Blockstore(ctx context.Context, s3BlockstoreID string, blockstore *S3Blockstore) (*S3Blockstore, *atlas.Response, error) {
	if s3BlockstoreID == "" {
		return nil, nil, atlas.NewArgError("s3BlockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorS3BlockstoreBasePath, s3BlockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, blockstore)
	if err != nil {
		return nil, nil, err
	}

	root := new(S3Blockstore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// DeleteS3Blockstore removes a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/delete-one-s3-blockstore-configuration/
func (s *BackupAdministratorServiceOp) DeleteS3Blockstore(ctx context.Context, s3BlockstoreID string) (*atlas.Response, error) {
	if s3BlockstoreID == "" {
		return nil, atlas.NewArgError("s3BlockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorS3BlockstoreBasePath, s3BlockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
