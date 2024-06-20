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
	clustersBasePath = "api/public/v1.0/groups/%s/clusters"
)

// ClustersService provides access to the cluster related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/clusters/
type ClustersService interface {
	List(context.Context, string, *ListOptions) (*Clusters, *Response, error)
	Get(context.Context, string, string) (*Cluster, *Response, error)
	ListAll(ctx context.Context) (*AllClustersProjects, *Response, error)
}

// ClustersServiceOp provides an implementation of the ClustersService interface.
type ClustersServiceOp service

// List all clusters for a project
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/clusters/#get-all-clusters
func (s *ClustersServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) (*Clusters, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
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

// Get get a single cluster by ID.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/clusters/#get-a-cluster
func (s *ClustersServiceOp) Get(ctx context.Context, groupID, clusterID string) (*Cluster, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupId", "must be set")
	}
	if clusterID == "" {
		return nil, nil, NewArgError("clusterID", "must be set")
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

// Cluster represents a cluster in Ops Manager.
type Cluster struct {
	ClusterName    string  `json:"clusterName,omitempty"`
	GroupID        string  `json:"groupId,omitempty"`
	ID             string  `json:"id,omitempty"`
	LastHeartbeat  string  `json:"lastHeartbeat,omitempty"`
	Links          []*Link `json:"links,omitempty"`
	ReplicaSetName string  `json:"replicaSetName,omitempty"`
	ShardName      string  `json:"shardName,omitempty"`
	TypeName       string  `json:"typeName,omitempty"`
}

// Clusters is a list of clusters.
type Clusters struct {
	Links      []*Link    `json:"links"`
	Results    []*Cluster `json:"results"`
	TotalCount int        `json:"totalCount"`
}

// ListAll list all clusters available to the user.
func (s *ClustersServiceOp) ListAll(ctx context.Context) (*AllClustersProjects, *Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodGet, "api/public/v1.0/clusters", nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AllClustersProjects)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

type AllClustersProject struct {
	GroupName string               `json:"groupName"`
	OrgName   string               `json:"orgName"`
	PlanType  string               `json:"planType,omitempty"`
	GroupID   string               `json:"groupId"`
	OrgID     string               `json:"orgId"`
	Tags      []string             `json:"tags"`
	Clusters  []AllClustersCluster `json:"clusters"`
}

// AllClustersCluster represents MongoDB cluster.
type AllClustersCluster struct {
	ClusterID     string   `json:"clusterId"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Availability  string   `json:"availability"`
	Versions      []string `json:"versions"`
	BackupEnabled bool     `json:"backupEnabled"`
	AuthEnabled   bool     `json:"authEnabled"`
	SSLEnabled    bool     `json:"sslEnabled"`
	AlertCount    int64    `json:"alertCount"`
	DataSizeBytes int64    `json:"dataSizeBytes"`
	NodeCount     int64    `json:"nodeCount"`
}

type AllClustersProjects struct {
	Links      []*Link               `json:"links"`
	Results    []*AllClustersProject `json:"results"`
	TotalCount int                   `json:"totalCount"`
}
