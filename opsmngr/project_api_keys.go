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

const projectAPIKeysPath = "api/public/v1.0/groups/%s/apiKeys" //nolint:gosec // This is a path

// ProjectAPIKeysOp handles communication with the APIKey related methods
// of the MongoDB Ops Manager API.
type ProjectAPIKeysOp service

var _ atlas.ProjectAPIKeysService = &ProjectAPIKeysOp{}

// List all API-KEY in the organization associated to {GROUP-ID}.
//
// See more: https://docs.cloudmanager.mongodb.com/reference/api/api-keys/project/get-all-apiKeys-in-one-project/
func (s *ProjectAPIKeysOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]atlas.APIKey, *Response, error) {
	path := fmt.Sprintf(projectAPIKeysPath, groupID)

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

// Create an API Key by the {GROUP-ID}.
//
// See more https://docs.cloudmanager.mongodb.com/reference/api/api-keys/project/create-one-apiKey-in-one-project/
func (s *ProjectAPIKeysOp) Create(ctx context.Context, groupID string, createRequest *atlas.APIKeyInput) (*atlas.APIKey, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(projectAPIKeysPath, groupID)

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

// Assign an API-KEY related to {GROUP-ID} to a the project with {API-KEY-ID}.
//
// See more: https://docs.cloudmanager.mongodb.com/reference/api/api-keys/project/assign-one-org-apiKey-to-one-project/
func (s *ProjectAPIKeysOp) Assign(ctx context.Context, groupID, keyID string, assignAPIKeyRequest *atlas.AssignAPIKey) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupID", "must be set")
	}

	if keyID == "" {
		return nil, NewArgError("keyID", "must be set")
	}

	basePath := fmt.Sprintf(projectAPIKeysPath, groupID)

	path := fmt.Sprintf("%s/%s", basePath, keyID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, assignAPIKeyRequest)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// Unassign an API-KEY related to {GROUP-ID} to a the project with {API-KEY-ID}.
//
// See more: https://docs.cloudmanager.mongodb.com/reference/api/api-keys/project/delete-one-apiKey-in-one-project/
func (s *ProjectAPIKeysOp) Unassign(ctx context.Context, groupID, keyID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("apiKeyID", "must be set")
	}

	if keyID == "" {
		return nil, NewArgError("keyID", "must be set")
	}

	basePath := fmt.Sprintf(projectAPIKeysPath, groupID)

	path := fmt.Sprintf("%s/%s", basePath, keyID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
