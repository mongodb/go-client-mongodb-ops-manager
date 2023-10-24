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
	"net/url"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const apiKeysOrgPath = "api/public/v1.0/orgs/%s/apiKeys" //nolint:gosec // This is a path

// APIKeysServiceOp handles communication with the APIKey related methods
// of the MongoDB Ops Manager API.
type APIKeysServiceOp service

var _ atlas.APIKeysService = &APIKeysServiceOp{}

// APIKeysResponse is the response from the APIKeysService.List.
type APIKeysResponse struct {
	Links      []*atlas.Link  `json:"links,omitempty"`
	Results    []atlas.APIKey `json:"results,omitempty"`
	TotalCount int            `json:"totalCount,omitempty"`
}

// List all API-KEY in the organization associated to {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/get-all-org-api-keys/
func (s *APIKeysServiceOp) List(ctx context.Context, orgID string, listOptions *ListOptions) ([]APIKey, *Response, error) {
	path := fmt.Sprintf(apiKeysOrgPath, orgID)

	// Add query params from listOptions
	path, err := setQueryParams(path, listOptions)
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

// Get gets the APIKey specified to {API-KEY-ID} from the organization associated to {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/get-one-org-api-key/
func (s *APIKeysServiceOp) Get(ctx context.Context, orgID, apiKeyID string) (*APIKey, *Response, error) {
	if apiKeyID == "" {
		return nil, nil, atlas.NewArgError("name", "must be set")
	}

	basePath := fmt.Sprintf(apiKeysOrgPath, orgID)
	escapedEntry := url.PathEscape(apiKeyID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

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

// Create an API Key by the {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/create-one-org-api-key/
func (s *APIKeysServiceOp) Create(ctx context.Context, orgID string, createRequest *APIKeyInput) (*APIKey, *Response, error) {
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(apiKeysOrgPath, orgID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createRequest)
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

// Update a API Key in the organization associated to {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/update-one-org-api-key/
func (s *APIKeysServiceOp) Update(ctx context.Context, orgID, apiKeyID string, updateRequest *APIKeyInput) (*APIKey, *Response, error) {
	if updateRequest == nil {
		return nil, nil, atlas.NewArgError("updateRequest", "cannot be nil")
	}

	basePath := fmt.Sprintf(apiKeysOrgPath, orgID)
	path := fmt.Sprintf("%s/%s", basePath, apiKeyID)

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

// Delete the API Key specified to {API-KEY-ID} from the organization associated to {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/org/delete-one-api-key/
func (s *APIKeysServiceOp) Delete(ctx context.Context, orgID, apiKeyID string) (*Response, error) {
	if apiKeyID == "" {
		return nil, atlas.NewArgError("apiKeyID", "must be set")
	}

	basePath := fmt.Sprintf(apiKeysOrgPath, orgID)
	escapedEntry := url.PathEscape(apiKeyID)
	path := fmt.Sprintf("%s/%s", basePath, escapedEntry)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
