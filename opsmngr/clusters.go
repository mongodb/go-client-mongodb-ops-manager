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
	clustersBasePath = "groups/%s/clusters"
)

// ClustersService is an interface for interfacing with Clusters in MongoDB Ops Manager APIs
// https://docs.opsmanager.mongodb.com/current/reference/api/clusters/
type ClustersService interface {
	Get(context.Context, string, string) (*Cluster, *atlas.Response, error)
	List(context.Context, string, *atlas.ListOptions) (*Clusters, *atlas.Response, error)
}

type ClustersServiceOp service

// Cluster a cluster details
type Cluster struct {
	ClusterName    string        `json:"clusterName,omitempty"`
	GroupID        string        `json:"groupId,omitempty"`
	ID             string        `json:"id,omitempty"`
	LastHeartbeat  string        `json:"lastHeartbeat,omitempty"`
	Links          []*atlas.Link `json:"links,omitempty"`
	ReplicaSetName string        `json:"replicaSetName,omitempty"`
	ShardName      string        `json:"shardName,omitempty"`
	TypeName       string        `json:"typeName,omitempty"`
}

// Clusters is a list of clusters
type Clusters struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Cluster    `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// List all clusters for a project
func (s *ClustersServiceOp) List(ctx context.Context, groupID string, listOptions *atlas.ListOptions) (*Clusters, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupId", "must be set")
	}
	basePath := fmt.Sprintf(clustersBasePath, groupID)
	path, err := setQueryParams(basePath, listOptions)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Clusters)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

func (s *ClustersServiceOp) Get(ctx context.Context, groupID, clusterID string) (*Cluster, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, atlas.NewArgError("clusterID", "must be set")
	}
	basePath := fmt.Sprintf(clustersBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, clusterID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)

	if err != nil {
		return nil, nil, err
	}

	root := new(Cluster)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
