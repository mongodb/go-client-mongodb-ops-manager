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

const (
	hostsDatabasesBasePath = "api/public/v1.0/groups/%s/hosts/%s/databases"
)

// ProcessDatabasesResponse is the response from the ProcessDatabasesService.List.
type ProcessDatabasesResponse struct {
	Links      []*Link            `json:"links"`
	Results    []*ProcessDatabase `json:"results"`
	TotalCount int                `json:"totalCount"`
}

// ProcessDatabase is the database information of a process.
type ProcessDatabase struct {
	Links        []*Link `json:"links"`
	DatabaseName string  `json:"databaseName"`
}

// ListDatabases retrieve all databases running on the specified host.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/databases-get-all-on-host/
func (s *DeploymentsServiceOp) ListDatabases(ctx context.Context, groupID, hostID string, opts *ListOptions) (*ProcessDatabasesResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if hostID == "" {
		return nil, nil, NewArgError("hostID", "must be set")
	}
	basePath := fmt.Sprintf(hostsDatabasesBasePath, groupID, hostID)
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDatabasesResponse)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// GetDatabase retrieve a single database by name.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/disk-get-one/
func (s *DeploymentsServiceOp) GetDatabase(ctx context.Context, groupID, hostID, name string) (*ProcessDatabase, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if hostID == "" {
		return nil, nil, NewArgError("hostID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("name", "must be set")
	}
	basePath := fmt.Sprintf(hostsDatabasesBasePath, groupID, hostID)
	path := fmt.Sprintf("%s/%s", basePath, name)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDatabase)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
