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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const backupAdministratorDaemonBasePath = "admin/backup/daemon/configs"

// DaemonConfigService is an interface for using the Backup Daemon
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/backup-daemon-config/
type DaemonConfigService interface {
	List(context.Context, *atlas.ListOptions) (*Daemons, *Response, error)
	Get(context.Context, string) (*Daemon, *Response, error)
	Update(context.Context, string, *Daemon) (*Daemon, *Response, error)
	Delete(context.Context, string) (*Response, error)
}

// DaemonConfigServiceOp provides an implementation of the DaemonConfigService interface.
type DaemonConfigServiceOp service

var _ DaemonConfigService = &DaemonConfigServiceOp{}

// Daemon represents a Backup Daemon Configuration in the MongoDB Ops Manager API.
type Daemon struct {
	AdminBackupConfig
	BackupJobsEnabled           bool     `json:"backupJobsEnabled"`
	Configured                  bool     `json:"configured"`
	GarbageCollectionEnabled    bool     `json:"garbageCollectionEnabled"`
	ResourceUsageEnabled        bool     `json:"resourceUsageEnabled"`
	RestoreQueryableJobsEnabled bool     `json:"restoreQueryableJobsEnabled"`
	HeadDiskType                string   `json:"headDiskType,omitempty"`
	NumWorkers                  int64    `json:"numWorkers,omitempty"`
	Machine                     *Machine `json:"machine,omitempty"`
}

type Machine struct {
	Machine           string `json:"machine,omitempty"`
	HeadRootDirectory string `json:"headRootDirectory,omitempty"`
}

// Daemons represents a paginated collection of Daemon.
type Daemons struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Daemon     `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// Get retrieves a Daemon.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/daemonConfigs/get-one-backup-daemon-configuration-by-host/
func (s *DaemonConfigServiceOp) Get(ctx context.Context, daemonID string) (*Daemon, *Response, error) {
	if daemonID == "" {
		return nil, nil, atlas.NewArgError("daemonID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorDaemonBasePath, daemonID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Daemon)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List retrieves all the Daemons.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/daemonConfigs/get-all-backup-daemon-configurations/
func (s *DaemonConfigServiceOp) List(ctx context.Context, options *atlas.ListOptions) (*Daemons, *Response, error) {
	path, err := setQueryParams(backupAdministratorDaemonBasePath, options)
	if err != nil {
		return nil, nil, err
	}
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Daemons)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates a Daemon.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/daemonConfigs/update-one-backup-daemon-configuration/
func (s *DaemonConfigServiceOp) Update(ctx context.Context, daemonID string, daemon *Daemon) (*Daemon, *Response, error) {
	if daemonID == "" {
		return nil, nil, atlas.NewArgError("daemonID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorDaemonBasePath, daemonID)
	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, daemon)
	if err != nil {
		return nil, nil, err
	}

	root := new(Daemon)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete removes a Daemon.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/admin/backup/daemonConfigs/delete-one-backup-daemon-configuration/
func (s *DaemonConfigServiceOp) Delete(ctx context.Context, daemonID string) (*Response, error) {
	if daemonID == "" {
		return nil, atlas.NewArgError("daemonID", "must be set")
	}

	path := fmt.Sprintf("%s/%s", backupAdministratorDaemonBasePath, daemonID)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
