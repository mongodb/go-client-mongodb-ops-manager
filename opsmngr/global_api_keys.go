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
	"net/url"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const apiKeysPath = "api/public/v1.0/admin/apiKeys" //nolint:gosec // This is a path

// GlobalAPIKeysService provides access to the global alerts related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/global-api-keys/
type GlobalAPIKeysService interface {
	List(context.Context, *ListOptions) ([]atlas.APIKey, *Response, error)
	Get(context.Context, string) (*atlas.APIKey, *Response, error)
	Create(context.Context, *atlas.APIKeyInput) (*atlas.APIKey, *Response, error)
	Update(context.Context, string, *atlas.APIKeyInput) (*atlas.APIKey, *Response, error)
	Delete(context.Context, string) (*Response, error)
	// Roles(context.Context) ([]atlas.AtlasRole, *Response, error)
}

// GlobalAPIKeysServiceOp provides an implementation of the GlobalAPIKeysService interface.
type GlobalAPIKeysServiceOp service

var _ GlobalAPIKeysService = &GlobalAPIKeysServiceOp{}

// List gets all Global API Keys.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/get-all-global-api-keys/
func (s *GlobalAPIKeysServiceOp) List(ctx context.Context, listOptions *ListOptions) ([]atlas.APIKey, *Response, error) {
	// Add query params from listOptions
	path, err := setQueryParams(apiKeysPath, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(APIKeysResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Get gets one API Key with ID apiKeyID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/get-one-global-api-key/
func (s *GlobalAPIKeysServiceOp) Get(ctx context.Context, apiKeyID string) (*atlas.APIKey, *Response, error) {
	if apiKeyID == "" {
		return nil, nil, atlas.NewArgError("apiKeyID", "must be set")
	}
	escapedEntry := url.PathEscape(apiKeyID)
	path := fmt.Sprintf("%s/%s", apiKeysPath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.APIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create an API Key.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/create-one-global-api-key/
func (s *GlobalAPIKeysServiceOp) Create(ctx context.Context, createRequest *atlas.APIKeyInput) (*atlas.APIKey, *Response, error) {
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, apiKeysPath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.APIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update one API Key with ID apiKeyID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/update-one-global-api-key/
func (s *GlobalAPIKeysServiceOp) Update(ctx context.Context, apiKeyID string, updateRequest *atlas.APIKeyInput) (*atlas.APIKey, *Response, error) {
	if updateRequest == nil {
		return nil, nil, atlas.NewArgError("updateRequest", "cannot be nil")
	}
	escapedEntry := url.PathEscape(apiKeyID)
	path := fmt.Sprintf("%s/%s", apiKeysPath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.APIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the API Key with ID apiKeyID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/delete-one-global-api-key/
func (s *GlobalAPIKeysServiceOp) Delete(ctx context.Context, apiKeyID string) (*Response, error) {
	if apiKeyID == "" {
		return nil, atlas.NewArgError("apiKeyID", "must be set")
	}
	escapedEntry := url.PathEscape(apiKeyID)
	path := fmt.Sprintf("%s/%s", apiKeysPath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
