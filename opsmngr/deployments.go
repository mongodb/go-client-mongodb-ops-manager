package opsmngr

import (
	"context"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
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
