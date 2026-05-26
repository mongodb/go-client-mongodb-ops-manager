package opsmngr

import (
	"context"
	"fmt"
	"net/http"
)

const backupAdministratorS3OplogBasePath = "api/public/v1.0/admin/backup/oplog/s3Configs"

// S3OplogStoreConfigService is an interface for using the S3 Oplog Store
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/oplog-store-config/
type S3OplogStoreConfigService interface {
	List(context.Context, *ListOptions) (*S3Blockstores, *Response, error)
	Get(context.Context, string) (*S3Blockstore, *Response, error)
	Create(context.Context, *S3Blockstore) (*S3Blockstore, *Response, error)
	Update(context.Context, string, *S3Blockstore) (*S3Blockstore, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// S3OplogStoreConfigServiceOp provides an implementation of S3OplogStoreConfigService.
type S3OplogStoreConfigServiceOp service

var _ S3OplogStoreConfigService = &S3OplogStoreConfigServiceOp{}

func (s *S3OplogStoreConfigServiceOp) Get(ctx context.Context, id string) (*S3Blockstore, *Response, error) {
	if id == "" {
		return nil, nil, NewArgError("id", "must be set")
	}
	path := fmt.Sprintf("%s/%s", backupAdministratorS3OplogBasePath, id)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	root := new(S3Blockstore)
	resp, err := s.Client.Do(ctx, req, root)
	return root, resp, err
}

func (s *S3OplogStoreConfigServiceOp) List(ctx context.Context, options *ListOptions) (*S3Blockstores, *Response, error) {
	path, err := setQueryParams(backupAdministratorS3OplogBasePath, options)
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

func (s *S3OplogStoreConfigServiceOp) Create(ctx context.Context, store *S3Blockstore) (*S3Blockstore, *Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodPost, backupAdministratorS3OplogBasePath, store)
	if err != nil {
		return nil, nil, err
	}
	root := new(S3Blockstore)
	resp, err := s.Client.Do(ctx, req, root)
	return root, resp, err
}

func (s *S3OplogStoreConfigServiceOp) Update(ctx context.Context, id string, store *S3Blockstore) (*S3Blockstore, *Response, error) {
	if id == "" {
		return nil, nil, NewArgError("id", "must be set")
	}
	path := fmt.Sprintf("%s/%s", backupAdministratorS3OplogBasePath, id)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, store)
	if err != nil {
		return nil, nil, err
	}
	root := new(S3Blockstore)
	resp, err := s.Client.Do(ctx, req, root)
	return root, resp, err
}

func (s *S3OplogStoreConfigServiceOp) Delete(ctx context.Context, id string) (*Response, error) {
	if id == "" {
		return nil, NewArgError("id", "must be set")
	}
	path := fmt.Sprintf("%s/%s", backupAdministratorS3OplogBasePath, id)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return s.Client.Do(ctx, req, nil)
}
