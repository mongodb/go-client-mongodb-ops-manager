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
	hostsDisksBasePath = "api/public/v1.0/groups/%s/hosts/%s/disks"
)

// ProcessDisksResponse is the response from the ProcessDisksService.List.
type ProcessDisksResponse struct {
	Links      []*Link        `json:"links"`
	Results    []*ProcessDisk `json:"results"`
	TotalCount int            `json:"totalCount"`
}

// ProcessDisk is the partition information of a process.
type ProcessDisk struct {
	Links         []*Link `json:"links"`
	PartitionName string  `json:"partitionName"`
}

// ListPartitions retrieves all disk partitions on the specified host.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/disks-get-all/
func (s *DeploymentsServiceOp) ListPartitions(ctx context.Context, groupID, hostID string, opts *ListOptions) (*ProcessDisksResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if hostID == "" {
		return nil, nil, NewArgError("hostID", "must be set")
	}
	basePath := fmt.Sprintf(hostsDisksBasePath, groupID, hostID)
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDisksResponse)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// GetPartition retrieves a disk partition.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/disk-get-one/
func (s *DeploymentsServiceOp) GetPartition(ctx context.Context, groupID, hostID, name string) (*ProcessDisk, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if hostID == "" {
		return nil, nil, NewArgError("hostID", "must be set")
	}
	if name == "" {
		return nil, nil, NewArgError("hostID", "must be set")
	}
	basePath := fmt.Sprintf(hostsDisksBasePath, groupID, hostID)
	path := fmt.Sprintf("%s/%s", basePath, name)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDisk)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
