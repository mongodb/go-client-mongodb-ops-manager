// Copyright 2019 MongoDB Inc
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
	"io"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	logsBasePath = "groups/%s/logCollectionJobs"
)

// LogsService is an interface for interfacing with the Log Collection Jobs
// endpoints of the MongoDB Ops Manager API.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collection/
type LogsService interface {
	List(context.Context, string, *LogListOptions) (*Logs, *atlas.Response, error)
	Get(context.Context, string, string, *LogListOptions) (*Log, *atlas.Response, error)
	Retry(context.Context, string, string) (*atlas.Response, error)
	Download(context.Context, string, string, io.Writer) (*atlas.Response, error)
	Create(context.Context, string, *Log) (*Log, *atlas.Response, error)
	Extend(context.Context, string, string, *Log) (*atlas.Response, error)
	Delete(context.Context, string, string) (*atlas.Response, error)
}

// LogsServiceOp handles communication with the Log Collection Jobs related methods of the
// MongoDB Ops Manager API
type LogsServiceOp struct {
	Client RequestDoer
}

var _ LogsService = &LogsServiceOp{}

// Log represents a Log Collection Job in the MongoDB Ops Manager API.
type Log struct {
	ID                         string      `json:"id,omitempty"`
	GroupID                    string      `json:"groupId,omitempty"`
	UserID                     string      `json:"userId,omitempty"`
	CreationDate               string      `json:"creationDate,omitempty"`
	ExpirationDate             string      `json:"expirationDate,omitempty"`
	Status                     string      `json:"status,omitempty"`
	ResourceType               string      `json:"resourceType,omitempty"`
	ResourceName               string      `json:"resourceName,omitempty"`
	RootResourceName           string      `json:"rootResourceName,omitempty"`
	RootResourceType           string      `json:"rootResourceType,omitempty"`
	URL                        string      `json:"downloadUrl,omitempty"`
	Redacted                   *bool        `json:"redacted,omitempty"`
	LogTypes                   []string    `json:"logTypes,omitempty"`
	SizeRequestedPerFileBytes  int64       `json:"sizeRequestedPerFileBytes,omitempty"`
	UncompressedDiskSpaceBytes int64       `json:"uncompressedSizeTotalBytes,omitempty"`
	ChildJob                   []*ChildJob `json:"childJobs,omitempty"` //included if verbose is true
}

// ChildJob represents a ChildJob in the MongoDB Ops Manager API.
type ChildJob struct {
	AutomationAgentID          string `json:"automationAgentId,omitempty"`
	ErrorMessage               string `json:"errorMessage,omitempty"`
	FinishDate                 string `json:"finishDate"`
	HostName                   string `json:"hostName"`
	LogCollectionType          string `json:"logCollectionType"`
	Path                       string `json:"path"`
	StartDate                  string `json:"startDate"`
	Status                     string `json:"status"`
	UncompressedDiskSpaceBytes int64  `json:"uncompressedDiskSpaceBytes"`
}

// Logs represents a array of Logs
type Logs struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Log        `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// LogListOptions specifies the optional parameters to List methods that
// support pagination.
type LogListOptions struct {
	atlas.ListOptions
	Verbose bool `url:"verbose,omitempty"`
}

// List gets all logs.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-get-all/
func (s *LogsServiceOp) List(ctx context.Context, groupID string, opts *LogListOptions) (*Logs, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(logsBasePath, groupID)
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Logs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets a log.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-get-one/
func (s *LogsServiceOp) Get(ctx context.Context, groupID, jobID string, opts *LogListOptions) (*Log, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if jobID == "" {
		return nil, nil, atlas.NewArgError("jobID", "must be set")
	}

	basePath := fmt.Sprintf(logsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, jobID)

	path, err := setQueryParams(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Log)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create creates a log.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-submit/
func (s *LogsServiceOp) Create(ctx context.Context, groupID string, log *Log) (*Log, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if log == nil {
		return nil, nil, atlas.NewArgError("log", "must be set")
	}

	path := fmt.Sprintf(logsBasePath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, log)
	if err != nil {
		return nil, nil, err
	}

	root := new(Log)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Extend extends the expiration data of a log.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-update-one/
func (s *LogsServiceOp) Extend(ctx context.Context, groupID, jobID string, log *Log) (*atlas.Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}

	if jobID == "" {
		return nil, atlas.NewArgError("jobID", "must be set")
	}

	if log == nil {
		return nil, atlas.NewArgError("log", "must be set")
	}

	basePath := fmt.Sprintf(logsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, jobID)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, log)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// Retry retries a single failed log collection job.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-retry/
func (s *LogsServiceOp) Retry(ctx context.Context, groupID, jobID string) (*atlas.Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}

	if jobID == "" {
		return nil, atlas.NewArgError("jobID", "must be set")
	}

	basePath := fmt.Sprintf(logsBasePath, groupID)
	path := fmt.Sprintf("%s/%s/retry", basePath, jobID)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// Delete removes a log.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-delete-one/
func (s *LogsServiceOp) Delete(ctx context.Context, groupID, jobID string) (*atlas.Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}

	if jobID == "" {
		return nil, atlas.NewArgError("jobID", "must be set")
	}

	basePath := fmt.Sprintf(logsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, jobID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// Download allows to download a log.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-download-job/
func (s *LogsServiceOp) Download(ctx context.Context, groupID, jobID string, out io.Writer) (*atlas.Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}

	if jobID == "" {
		return nil, atlas.NewArgError("jobID", "must be set")
	}

	basePath := fmt.Sprintf(logsBasePath, groupID)
	path := fmt.Sprintf("%s/%s/download", basePath, jobID)

	req, err := s.Client.NewGZipRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, out)

	return resp, err
}
