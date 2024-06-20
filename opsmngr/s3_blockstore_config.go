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

const backupAdministratorS3BlockstoreBasePath = "api/public/v1.0/admin/backup/snapshot/s3Configs"

// S3BlockstoreConfigService is an interface for using the S3 Blockstore Service
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/s3-blockstore-config/
type S3BlockstoreConfigService interface {
	List(context.Context, *ListOptions) (*S3Blockstores, *Response, error)
	Get(context.Context, string) (*S3Blockstore, *Response, error)
	Create(context.Context, *S3Blockstore) (*S3Blockstore, *Response, error)
	Update(context.Context, string, *S3Blockstore) (*S3Blockstore, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// S3BlockstoreConfigServiceOp provides an implementation of the S3BlockstoreConfigServiceinterface.
type S3BlockstoreConfigServiceOp service

var _ S3BlockstoreConfigService = &S3BlockstoreConfigServiceOp{}

// S3Blockstore represents a S3Blockstore in the MongoDB Ops Manager API.
type S3Blockstore struct {
	BackupStore
	AWSAccessKey           string `json:"awsAccessKey,omitempty"`
	AWSSecretKey           string `json:"awsSecretKey,omitempty"`
	S3AuthMethod           string `json:"s3AuthMethod,omitempty"`
	S3BucketEndpoint       string `json:"s3BucketEndpoint,omitempty"`
	S3BucketName           string `json:"s3BucketName,omitempty"`
	S3MaxConnections       int64  `json:"s3MaxConnections,omitempty"`
	DisableProxyS3         *bool  `json:"disableProxyS3,omitempty"`
	AcceptedTos            *bool  `json:"acceptedTos"`
	SSEEnabled             *bool  `json:"sseEnabled"`
	PathStyleAccessEnabled *bool  `json:"pathStyleAccessEnabled"`
}

// S3Blockstores represents a paginated collection of S3Blockstore.
type S3Blockstores struct {
	Links      []*Link         `json:"links"`
	Results    []*S3Blockstore `json:"results"`
	TotalCount int             `json:"totalCount"`
}

// Get retrieves a Get.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/get-one-s3-blockstore-configuration-by-id/
func (s *S3BlockstoreConfigServiceOp) Get(ctx context.Context, s3BlockstoreID string) (*S3Blockstore, *Response, error) {
	if s3BlockstoreID == "" {
		return nil, nil, NewArgError("s3BlockstoreID", "must be set")
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

// List retrieves all the S3Blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/get-all-s3-blockstore-configurations/
func (s *S3BlockstoreConfigServiceOp) List(ctx context.Context, options *ListOptions) (*S3Blockstores, *Response, error) {
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

// Create create a S3Blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/create-one-s3-blockstore-configuration/
func (s *S3BlockstoreConfigServiceOp) Create(ctx context.Context, blockstore *S3Blockstore) (*S3Blockstore, *Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorS3BlockstoreBasePath, blockstore)
	if err != nil {
		return nil, nil, err
	}

	root := new(S3Blockstore)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates a S3Blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/update-one-s3-blockstore-configuration/
func (s *S3BlockstoreConfigServiceOp) Update(ctx context.Context, s3BlockstoreID string, blockstore *S3Blockstore) (*S3Blockstore, *Response, error) {
	if s3BlockstoreID == "" {
		return nil, nil, NewArgError("s3BlockstoreID", "must be set")
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

// Delete removes a blockstore.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/snapshot/s3Configs/delete-one-s3-blockstore-configuration/
func (s *S3BlockstoreConfigServiceOp) Delete(ctx context.Context, s3BlockstoreID string) (*Response, error) {
	if s3BlockstoreID == "" {
		return nil, NewArgError("s3BlockstoreID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorS3BlockstoreBasePath, s3BlockstoreID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
