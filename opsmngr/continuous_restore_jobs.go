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

const continuousRestoreJobsPath = "api/public/v1.0/groups/%s/clusters/%s/restoreJobs"

// ContinuousRestoreJobsServiceOp handles communication with the Continuous Backup Restore Jobs related methods
// of the MongoDB Ops Manager API.
type ContinuousRestoreJobsServiceOp service

type (
	ContinuousJob        = atlas.ContinuousJob
	ContinuousJobs       = atlas.ContinuousJobs
	ContinuousJobRequest = atlas.ContinuousJobRequest
)

var _ atlas.ContinuousRestoreJobsService = &ContinuousRestoreJobsServiceOp{}

// List lists all continuous backup jobs in Atlas
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/restorejobs/get-all-restore-jobs-for-one-cluster/
func (s *ContinuousRestoreJobsServiceOp) List(ctx context.Context, groupID, clusterID string, opts *ListOptions) (*ContinuousJobs, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(continuousRestoreJobsPath, groupID, clusterID)

	path, err := setQueryParams(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousJobs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets a continuous backup job in Atlas
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/restorejobs/get-one-single-restore-job-for-one-cluster/
func (s *ContinuousRestoreJobsServiceOp) Get(ctx context.Context, groupID, clusterID, jobID string) (*ContinuousJob, *Response, error) {
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if jobID == "" {
		return nil, nil, NewArgError("jobID", "must be set")
	}
	defaultPath := fmt.Sprintf(continuousRestoreJobsPath, groupID, clusterID)

	path := fmt.Sprintf("%s/%s", defaultPath, jobID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousJob)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create creates a continuous backup job in Atlas
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/restorejobs/create-one-restore-job-for-one-cluster/
func (s *ContinuousRestoreJobsServiceOp) Create(ctx context.Context, groupID, clusterID string, request *ContinuousJobRequest) (*ContinuousJobs, *Response, error) {
	if request == nil {
		return nil, nil, NewArgError("request", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
	}
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(continuousRestoreJobsPath, groupID, clusterID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, request)

	if err != nil {
		return nil, nil, err
	}

	root := new(ContinuousJobs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
