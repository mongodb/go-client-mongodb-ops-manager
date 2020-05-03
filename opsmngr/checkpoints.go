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
	checkpoints = "groups/%s/clusters/%s/checkpoints"
)

// CheckpointsService provides access to the backup related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/checkpoints/
type CheckpointsService interface {
	List(context.Context, string, string, *atlas.ListOptions) (*atlas.Checkpoints, *atlas.Response, error)
	Get(context.Context, string, string, string) (*atlas.Checkpoint, *atlas.Response, error)
}

// CheckpointsServiceOp provides an implementation of the CheckpointsService interface
type CheckpointsServiceOp service

var _ CheckpointsService = &CheckpointsServiceOp{}

// List lists checkpoints.
//
// See https://docs.opsmanager.mongodb.com/current/reference/api/checkpoints/#get-all-checkpoints
func (s *CheckpointsServiceOp) List(ctx context.Context, groupID, clusterName string, listOptions *atlas.ListOptions) (*atlas.Checkpoints, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, atlas.NewArgError("clusterName", "must be set")
	}

	basePath := fmt.Sprintf(checkpoints, groupID, clusterName)
	path, err := setQueryParams(basePath, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Checkpoints)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets a checkpoint.
//
// See https://docs.opsmanager.mongodb.com/current/reference/api/checkpoints/#get-one-checkpoint
func (s *CheckpointsServiceOp) Get(ctx context.Context, groupID, clusterID, checkpointID string) (*atlas.Checkpoint, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, atlas.NewArgError("clusterID", "must be set")
	}
	if checkpointID == "" {
		return nil, nil, atlas.NewArgError("checkpointID", "must be set")
	}

	basePath := fmt.Sprintf(checkpoints, groupID, clusterID)
	path := fmt.Sprintf("%s/%s", basePath, checkpointID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Checkpoint)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
