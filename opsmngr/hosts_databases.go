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

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	hostsDatabasesBasePath = "groups/%s/hosts/%s/databases"
)

// HostDatabasesService is an interface for interfacing with databases in MongoDB Ops Manager APIs
// https://docs.opsmanager.mongodb.com/current/reference/api/databases/
type HostDatabasesService interface {
	Get(context.Context, string, string, string) (*atlas.ProcessDatabase, *atlas.Response, error)
	List(context.Context, string, string, *atlas.ListOptions) (*atlas.ProcessDatabasesResponse, *atlas.Response, error)
}

type HostDatabasesServiceOp struct {
	Client atlas.RequestDoer
}

// Get gets the MongoDB databases with the specified host ID and database name.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/disk-get-one/
func (s *HostDatabasesServiceOp) Get(ctx context.Context, groupID, hostID, partitionName string) (*atlas.ProcessDatabase, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if hostID == "" {
		return nil, nil, atlas.NewArgError("hostID", "must be set")
	}
	basePath := fmt.Sprintf(hostsDatabasesBasePath, groupID, hostID)
	path := fmt.Sprintf("%s/%s", basePath, partitionName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.ProcessDatabase)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List lists all MongoDB databases in a host.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/databases-get-all-on-host/
func (s *HostDatabasesServiceOp) List(ctx context.Context, groupID, hostID string, opts *atlas.ListOptions) (*atlas.ProcessDatabasesResponse, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if hostID == "" {
		return nil, nil, atlas.NewArgError("hostID", "must be set")
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

	root := new(atlas.ProcessDatabasesResponse)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
