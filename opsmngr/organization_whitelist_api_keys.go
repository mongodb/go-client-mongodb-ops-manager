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

const organizationWhitelistAPIKeysPath = "api/public/v1.0/orgs/%s/apiKeys/%s/whitelist" //nolint:gosec // This is a path

// WhitelistAPIKeysServiceOp handles communication with the Whitelist API keys related methods of the
// MongoDB Ops Manager API.
type WhitelistAPIKeysServiceOp service

var _ atlas.WhitelistAPIKeysService = &WhitelistAPIKeysServiceOp{} //nolint:staticcheck //we keep whitelist to support OM 4.2 and 4.4

// List gets all Whitelist API keys.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/get-all-org-api-key-whitelist/
func (s *WhitelistAPIKeysServiceOp) List(ctx context.Context, orgID, apiKeyID string, listOptions *atlas.ListOptions) (*atlas.WhitelistAPIKeys, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, atlas.NewArgError("apiKeyID", "must be set")
	}

	path := fmt.Sprintf(organizationWhitelistAPIKeysPath, orgID, apiKeyID)
	path, err := setQueryParams(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.WhitelistAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get gets the Whitelist API keys.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/get-one-org-api-key-whitelist/
func (s *WhitelistAPIKeysServiceOp) Get(ctx context.Context, orgID, apiKeyID, ipAddress string) (*atlas.WhitelistAPIKey, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, atlas.NewArgError("apiKeyID", "must be set")
	}
	if ipAddress == "" {
		return nil, nil, atlas.NewArgError("ipAddress", "must be set")
	}

	path := fmt.Sprintf(organizationWhitelistAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.WhitelistAPIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create a submit a POST request containing ipAddress or cidrBlock values which are not already present in the whitelist,
// Atlas adds those entries to the list of existing entries in the whitelist.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/create-org-api-key-whitelist/
func (s *WhitelistAPIKeysServiceOp) Create(ctx context.Context, orgID, apiKeyID string, createRequest []*atlas.WhitelistAPIKeysReq) (*atlas.WhitelistAPIKeys, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	if apiKeyID == "" {
		return nil, nil, atlas.NewArgError("apiKeyID", "must be set")
	}
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(organizationWhitelistAPIKeysPath, orgID, apiKeyID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.WhitelistAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete deletes the Whitelist API keys.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/delete-one-org-api-key-whitelist/
func (s *WhitelistAPIKeysServiceOp) Delete(ctx context.Context, orgID, apiKeyID, ipAddress string) (*Response, error) {
	if orgID == "" {
		return nil, atlas.NewArgError("groupId", "must be set")
	}
	if apiKeyID == "" {
		return nil, atlas.NewArgError("clusterName", "must be set")
	}
	if ipAddress == "" {
		return nil, atlas.NewArgError("snapshotId", "must be set")
	}

	path := fmt.Sprintf(organizationWhitelistAPIKeysPath+"/%s", orgID, apiKeyID, ipAddress)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
