package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	hostsDisksBasePath = "groups/%s/hosts/%s/disks"
)

// HostDisksService is an interface for interfacing with Hosts in MongoDB Ops Manager APIs
// https://docs.opsmanager.mongodb.com/current/reference/api/disks/
type HostDisksService interface {
	Get(context.Context, string, string, string) (*atlas.ProcessDisk, *atlas.Response, error)
	List(context.Context, string, string, *atlas.ListOptions) (*atlas.ProcessDisksResponse, *atlas.Response, error)
}

type HostDisksServiceOp struct {
	Client atlas.RequestDoer
}

// Get gets the MongoDB disks with the specified host ID and partition name.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/disk-get-one/
func (s *HostDisksServiceOp) Get(ctx context.Context, projectID, hostID, partitionName string) (*atlas.ProcessDisk, *atlas.Response, error) {
	basePath := fmt.Sprintf(hostsDisksBasePath, projectID, hostID)
	path := fmt.Sprintf("%s/%s", basePath, partitionName)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.ProcessDisk)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List lists all MongoDB partitions in a host.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/disks-get-all/
func (s *HostDisksServiceOp) List(ctx context.Context, projectID, hostID string, opts *atlas.ListOptions) (*atlas.ProcessDisksResponse, *atlas.Response, error) {
	basePath := fmt.Sprintf(hostsDisksBasePath, projectID, hostID)
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.ProcessDisksResponse)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
