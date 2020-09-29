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
	backupAdministratorOplogBasePath                         = "admin/backup/oplog/mongoConfigs"
	backupAdministratorSyncBasePath                         = "admin/backup/sync/mongoConfigs"
)

// BackupAdministratorService is an interface for using the Backup Administrator
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/nav/administration-backup/
type BackupAdministratorService interface {
	BlockstoreService
	FileSystemStoreService
	S3BlockstoreService
	OplogService
	SyncService
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

// FileSystemStoreService is an interface for using the File System Store Configuration
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/file-system-store-config/
type FileSystemStoreService interface {
	ListFileSystemStores(context.Context, *atlas.ListOptions) (*FileSystemStoreConfigurations, *atlas.Response, error)
	GetFileSystemStore(context.Context, string) (*FileSystemStoreConfiguration, *atlas.Response, error)
	CreateFileSystemStore(context.Context, *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *atlas.Response, error)
	UpdateFileSystemStore(context.Context, string, *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *atlas.Response, error)
	DeleteFileSystemStore(context.Context, string) (*atlas.Response, error)
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

// OplogService is an interface for using the Oplog
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog-store-config/
type OplogService interface {
	ListOplog(context.Context, *atlas.ListOptions) (*Oplogs, *atlas.Response, error)
	GetOplog(context.Context, string) (*Oplog, *atlas.Response, error)
	CreateOplog(context.Context, *Oplog) (*Oplog, *atlas.Response, error)
	UpdateOplog(context.Context, string, *Oplog) (*Oplog, *atlas.Response, error)
	DeleteOplog(context.Context, string) (*atlas.Response, error)
}

// SyncService is an interface for using the Sync
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync-store-config/
type SyncService interface {
	ListSyncs(context.Context, *atlas.ListOptions) (*Syncs, *atlas.Response, error)
	GetSync(context.Context, string) (*Sync, *atlas.Response, error)
	CreateSync(context.Context, *Sync) (*Sync, *atlas.Response, error)
	UpdateSync(context.Context, string, *Sync) (*Sync, *atlas.Response, error)
	DeleteSync(context.Context, string) (*atlas.Response, error)
}

// BackupConfigsServiceOp provides an implementation of the BackupConfigsService interface
type BackupAdministratorServiceOp service

var _ BackupAdministratorService = &BackupAdministratorServiceOp{}

// AdminConfig contains the common fields of backup administrator structs
type AdminConfig struct {
	ID                   string   `json:"id,omitempty"`
	AssignmentEnabled    bool     `json:"assignmentEnabled,omitempty"`
	EncryptedCredentials bool     `json:"encryptedCredentials,omitempty"`
	URI                  string   `json:"uri,omitempty"`
	Labels               []string `json:"labels,omitempty"`
	SSL                  bool     `json:"ssl,omitempty"`
	WriteConcern         string   `json:"writeConcern,omitempty"`
	UsedSize             int64    `json:"usedSize,omitempty"`
}

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

// S3Blockstores represents a paginated collection of S3Blockstore
type S3Blockstores struct {
	Links      []*atlas.Link   `json:"links"`
	Results    []*S3Blockstore `json:"results"`
	TotalCount int             `json:"totalCount"`
}

// Blockstore represents a Blockstore in the MongoDB Ops Manager API
type Blockstore struct {
	AdminConfig
	LoadFactor    int64  `json:"loadFactor,omitempty"`
	MaxCapacityGB int64  `json:"maxCapacityGB,omitempty"`
	Provisioned   bool   `json:"provisioned,omitempty"`
	SyncSource    string `json:"syncSource,omitempty"`
	Username      string `json:"username,omitempty"`
}

// Blockstores represents a paginated collection of Blockstore
type Blockstores struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Blockstore `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// FileSystemStoreConfiguration represents a File System Store Configuration in the MongoDB Ops Manager API
type FileSystemStoreConfiguration struct {
	AdminConfig
	LoadFactor               int64  `json:"loadFactor,omitempty"`
	MMAPV1CompressionSetting string `json:"mmapv1CompressionSetting,omitempty"`
	StorePath                string `json:"storePath,omitempty"`
	WTCompressionSetting     string `json:"wtCompressionSetting,omitempty"`
	AssignmentEnabled        bool   `json:"assignmentEnabled,omitempty"`
}

// FileSystemStoreConfigurations represents a paginated collection of FileSystemStoreConfiguration
type FileSystemStoreConfigurations struct {
	Links      []*atlas.Link                   `json:"links"`
	Results    []*FileSystemStoreConfiguration `json:"results"`
	TotalCount int                             `json:"totalCount"`
}

// Oplog represents a Oplog Configuration in the MongoDB Ops Manager API
type Oplog struct {
	Blockstore
}

// Oplogs represents a paginated collection of Oplog
type Oplogs struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Oplog      `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// Sync represents a Sync Configuration in the MongoDB Ops Manager API
type Sync struct {
	Blockstore
}

// Syncs represents a paginated collection of Oplog
type Syncs struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Sync      `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// GetFileSystemStore retrieves a File System Store Configuration.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/get-one-file-system-store-configuration-by-id/
func (s *BackupAdministratorServiceOp) GetFileSystemStore(ctx context.Context, fileSystemID string) (*FileSystemStoreConfiguration, *atlas.Response, error) {
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

// ListFileSystemStores retrieves the configurations of all file system stores.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/get-all-file-system-store-configurations/
func (s *BackupAdministratorServiceOp) ListFileSystemStores(ctx context.Context, opts *atlas.ListOptions) (*FileSystemStoreConfigurations, *atlas.Response, error) {
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

// CreateFileSystemStore configures one new file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/create-one-file-system-store-configuration/
func (s *BackupAdministratorServiceOp) CreateFileSystemStore(ctx context.Context, fileSystem *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorFileSystemStoreConfigurationsBasePath, fileSystem)
	if err != nil {
		return nil, nil, err
	}

	root := new(FileSystemStoreConfiguration)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// UpdateFileSystemStore updates the configuration of one file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/update-one-file-system-store-configuration/
func (s *BackupAdministratorServiceOp) UpdateFileSystemStore(ctx context.Context, fileSystemID string, fileSystem *FileSystemStoreConfiguration) (*FileSystemStoreConfiguration, *atlas.Response, error) {
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

// DeleteFileSystemStore deletes the configuration of one file system store.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/fileSystemConfigs/delete-one-file-system-store-configuration/
func (s *BackupAdministratorServiceOp) DeleteFileSystemStore(ctx context.Context, fileSystemID string) (*atlas.Response, error) {
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

// GetOplog retrieves a Oplog.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog/mongoConfigs/get-one-oplog-configuration-by-id/
func (s *BackupAdministratorServiceOp) GetOplog(ctx context.Context, oplog string) (*Oplog, *atlas.Response, error) {
	if oplog == "" {
		return nil, nil, atlas.NewArgError("oplog", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorOplogBasePath, oplog)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Oplog)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// ListOplog retrieves all the Oplogs.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog/mongoConfigs/get-all-oplog-configurations/
func (s *BackupAdministratorServiceOp) ListOplog(ctx context.Context, options *atlas.ListOptions) (*Oplogs, *atlas.Response, error) {
	path, err := setQueryParams(backupAdministratorOplogBasePath, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Oplogs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// CreateOplog create a Oplog.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog/mongoConfigs/create-one-oplog-configuration/
func (s *BackupAdministratorServiceOp) CreateOplog(ctx context.Context, oplog *Oplog) (*Oplog, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorOplogBasePath, oplog)
	if err != nil {
		return nil, nil, err
	}

	root := new(Oplog)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// UpdateOplog updates a Oplog.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog/mongoConfigs/update-one-oplog-configuration/
func (s *BackupAdministratorServiceOp) UpdateOplog(ctx context.Context, oplogID string, oplog *Oplog) (*Oplog, *atlas.Response, error) {
	if oplogID == "" {
		return nil, nil, atlas.NewArgError("oplogID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorOplogBasePath, oplogID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, oplog)
	if err != nil {
		return nil, nil, err
	}

	root := new(Oplog)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// DeleteOplog removes a Oplog.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/delete-one-s3-blockstore-configuration/
func (s *BackupAdministratorServiceOp) DeleteOplog(ctx context.Context, oplogID string) (*atlas.Response, error) {
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

// GetSync retrieves a Sync.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/get-one-sync-store-configuration-by-id/
func (s *BackupAdministratorServiceOp) GetSync(ctx context.Context, syncID string) (*Sync, *atlas.Response, error) {
	if syncID == "" {
		return nil, nil, atlas.NewArgError("syncID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorSyncBasePath, syncID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Sync)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// ListSyncs retrieves all the Syncs.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/get-all-sync-store-configurations/
func (s *BackupAdministratorServiceOp) ListSyncs(ctx context.Context, options *atlas.ListOptions) (*Syncs, *atlas.Response, error) {
	path, err := setQueryParams(backupAdministratorSyncBasePath, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Syncs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// CreateSync create a Sync.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/create-one-sync-store-configuration/
func (s *BackupAdministratorServiceOp) CreateSync(ctx context.Context, sync *Sync) (*Sync, *atlas.Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorSyncBasePath, sync)
	if err != nil {
		return nil, nil, err
	}

	root := new(Sync)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// UpdateSync updates a Sync.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/update-one-sync-store-configuration/
func (s *BackupAdministratorServiceOp) UpdateSync(ctx context.Context, syncID string, sync *Sync) (*Sync, *atlas.Response, error) {
	if syncID == "" {
		return nil, nil, atlas.NewArgError("syncID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorSyncBasePath, syncID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, sync)
	if err != nil {
		return nil, nil, err
	}

	root := new(Sync)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// DeleteSync removes a Sync.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/sync/mongoConfigs/delete-one-sync-store-configuration/
func (s *BackupAdministratorServiceOp) DeleteSync(ctx context.Context, syncID string) (*atlas.Response, error) {
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