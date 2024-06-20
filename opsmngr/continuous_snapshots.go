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

const (
	continuousSnapshotsBasePath = "api/public/v1.0/groups/%s/clusters/%s/snapshots"
)

// ContinuousSnapshotsServiceOp handles communication with the Continuous Snapshots related methods of the
// MongoDB Ops Manager API.
type ContinuousSnapshotsServiceOp service

var _ atlas.ContinuousSnapshotsService = &ContinuousSnapshotsServiceOp{}

type (
	ContinuousSnapshots = atlas.ContinuousSnapshots
	ContinuousSnapshot  = atlas.ContinuousSnapshot
)

// List lists continuous snapshots for the given cluster
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/snapshots/get-all-snapshots-for-one-cluster/
func (s *ContinuousSnapshotsServiceOp) List(ctx context.Context, groupID, clusterID string, listOptions *ListOptions) (*ContinuousSnapshots, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}

	path := fmt.Sprintf(continuousSnapshotsBasePath, groupID, clusterID)

	// Add query params from listOptions
	path, err := setQueryParams(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousSnapshots)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets the continuous snapshot for the given cluster and snapshot ID
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/snapshots/get-one-snapshot-for-one-cluster/
func (s *ContinuousSnapshotsServiceOp) Get(ctx context.Context, groupID, clusterID, snapshotID string) (*ContinuousSnapshot, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if snapshotID == "" {
		return nil, nil, NewArgError("snapshotID", "must be set")
	}

	basePath := fmt.Sprintf(continuousSnapshotsBasePath, groupID, clusterID)
	path := fmt.Sprintf("%s/%s", basePath, snapshotID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousSnapshot)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// ChangeExpiry changes the expiry date for the given cluster and snapshot ID
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/snapshots/change-expiry-for-one-snapshot/
func (s *ContinuousSnapshotsServiceOp) ChangeExpiry(ctx context.Context, groupID, clusterID, snapshotID string, updateRequest *ContinuousSnapshot) (*ContinuousSnapshot, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if snapshotID == "" {
		return nil, nil, NewArgError("snapshotID", "must be set")
	}

	basePath := fmt.Sprintf(continuousSnapshotsBasePath, groupID, clusterID)
	path := fmt.Sprintf("%s/%s", basePath, snapshotID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousSnapshot)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete deletes the given continuous snapshot
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/snapshots/remove-one-snapshot-from-one-cluster/
func (s *ContinuousSnapshotsServiceOp) Delete(ctx context.Context, groupID, clusterID, snapshotID string) (*Response, error) {
	if groupID == "" {
		return nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, NewArgError("clusterID", "must be set")
	}
	if snapshotID == "" {
		return nil, NewArgError("snapshotID", "must be set")
	}

	basePath := fmt.Sprintf(continuousSnapshotsBasePath, groupID, clusterID)
	path := fmt.Sprintf("%s/%s", basePath, snapshotID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	return s.Client.Do(ctx, req, nil)
}
