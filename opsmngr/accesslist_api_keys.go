// Copyright 2021 MongoDB Inc
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

type (
	AccessListAPIKey     = atlas.AccessListAPIKey
	AccessListAPIKeys    = atlas.AccessListAPIKeys
	AccessListAPIKeysReq = atlas.AccessListAPIKeysReq
)

const accessListAPIKeysPath = "api/public/v1.0/orgs/%s/apiKeys/%s/accessList" //nolint:gosec // This is a path

// AccessListAPIKeysServiceOp handles communication with the AccessList API keys related methods of the MongoDB Ops Manager API.
type AccessListAPIKeysServiceOp service

var _ atlas.AccessListAPIKeysService = &AccessListAPIKeysServiceOp{}

// List gets all AccessList API keys.
func (s *AccessListAPIKeysServiceOp) List(ctx context.Context, orgID, apiKeyID string, listOptions *ListOptions) (*AccessListAPIKeys, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, atlas.NewArgError("apiKeyID", "must be set")
	}

	path := fmt.Sprintf(accessListAPIKeysPath, orgID, apiKeyID)
	path, err := setQueryParams(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AccessListAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get retrieve information on a single API Key access list entry using the unique identifier for the API Key and desired permitted address.
func (s *AccessListAPIKeysServiceOp) Get(ctx context.Context, orgID, apiKeyID, ipAddress string) (*AccessListAPIKey, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, atlas.NewArgError("apiKeyID", "must be set")
	}
	if ipAddress == "" {
		return nil, nil, atlas.NewArgError("ipAddress", "must be set")
	}

	path := fmt.Sprintf(accessListAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AccessListAPIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create one or more new access list entries for the specified API Key.
func (s *AccessListAPIKeysServiceOp) Create(ctx context.Context, orgID, apiKeyID string, createRequest []*AccessListAPIKeysReq) (*AccessListAPIKeys, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, atlas.NewArgError("apiKeyID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(accessListAPIKeysPath, orgID, apiKeyID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(AccessListAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes the AccessList API keys.
func (s *AccessListAPIKeysServiceOp) Delete(ctx context.Context, orgID, apiKeyID, ipAddress string) (*Response, error) {
	if orgID == "" {
		return nil, atlas.NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, atlas.NewArgError("apiKeyID", "must be set")
	}
	if ipAddress == "" {
		return nil, atlas.NewArgError("ipAddress", "must be set")
	}

	path := fmt.Sprintf(accessListAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
