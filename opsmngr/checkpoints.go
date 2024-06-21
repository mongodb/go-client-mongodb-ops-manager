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
	checkpoints = "api/public/v1.0/groups/%s/clusters/%s/checkpoints"
)

// Checkpoint represents MongoDB Checkpoint.
type Checkpoint struct {
	ClusterID  string  `json:"clusterId"`
	Completed  string  `json:"completed,omitempty"`
	GroupID    string  `json:"groupId"`
	ID         string  `json:"id,omitempty"`    // Unique identifier of the checkpoint.
	Links      []*Link `json:"links,omitempty"` // One or more links to sub-resources and/or related resources.
	Parts      []*Part `json:"parts,omitempty"`
	Restorable bool    `json:"restorable"`
	Started    string  `json:"started"`
	Timestamp  string  `json:"timestamp"`
}

// CheckpointPart represents the individual parts that comprise the complete checkpoint.
type CheckpointPart struct {
	ShardName       string            `json:"shardName"`
	TokenDiscovered bool              `json:"tokenDiscovered"`
	TokenTimestamp  SnapshotTimestamp `json:"tokenTimestamp"`
}

// Checkpoints represents all the backup checkpoints related to a cluster.
type Checkpoints struct {
	Results    []*Checkpoint `json:"results,omitempty"`    // Includes one Checkpoint object for each item detailed in the results array section.
	Links      []*Link       `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int           `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

// CheckpointsService provides access to the backup related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/checkpoints/
type CheckpointsService interface {
	List(context.Context, string, string, *ListOptions) (*Checkpoints, *Response, error)
	Get(context.Context, string, string, string) (*Checkpoint, *Response, error)
}

// CheckpointsServiceOp provides an implementation of the CheckpointsService interface.
type CheckpointsServiceOp service

var _ CheckpointsService = &CheckpointsServiceOp{}

// List lists checkpoints.
//
// See https://docs.opsmanager.mongodb.com/current/reference/api/checkpoints/#get-all-checkpoints
func (s *CheckpointsServiceOp) List(ctx context.Context, groupID, clusterName string, listOptions *ListOptions) (*Checkpoints, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterName == "" {
		return nil, nil, NewArgError("clusterName", "must be set")
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

	root := new(Checkpoints)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets a checkpoint.
//
// See https://docs.opsmanager.mongodb.com/current/reference/api/checkpoints/#get-one-checkpoint
func (s *CheckpointsServiceOp) Get(ctx context.Context, groupID, clusterID, checkpointID string) (*Checkpoint, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if checkpointID == "" {
		return nil, nil, NewArgError("checkpointID", "must be set")
	}

	basePath := fmt.Sprintf(checkpoints, groupID, clusterID)
	path := fmt.Sprintf("%s/%s", basePath, checkpointID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(Checkpoint)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
