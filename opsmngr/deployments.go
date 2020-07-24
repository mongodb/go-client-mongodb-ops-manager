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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

// DeploymentsService provides access to the deployment related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/nav/deployments/
type DeploymentsService interface {
	ListHosts(context.Context, string, *HostListOptions) (*Hosts, *atlas.Response, error)
	GetHost(context.Context, string, string) (*Host, *atlas.Response, error)
	GetHostByHostname(context.Context, string, string, int) (*Host, *atlas.Response, error)
	StartMonitoring(context.Context, string, *Host) (*Host, *atlas.Response, error)
	UpdateMonitoring(context.Context, string, string, *Host) (*Host, *atlas.Response, error)
	StopMonitoring(context.Context, string, string) (*atlas.Response, error)
	ListPartitions(context.Context, string, string, *atlas.ListOptions) (*atlas.ProcessDisksResponse, *atlas.Response, error)
	GetPartition(context.Context, string, string, string) (*atlas.ProcessDisk, *atlas.Response, error)
	ListDatabases(context.Context, string, string, *atlas.ListOptions) (*atlas.ProcessDatabasesResponse, *atlas.Response, error)
	GetDatabase(context.Context, string, string, string) (*atlas.ProcessDatabase, *atlas.Response, error)
}

// DeploymentsServiceOp provides an implementation of the DeploymentsService interface
type DeploymentsServiceOp service

var _ DeploymentsService = new(DeploymentsServiceOp)
