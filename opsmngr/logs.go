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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	logsBasePath = "api/public/v1.0/groups/%s/logCollectionJobs"
)

// LogCollectionService is an interface for interfacing with the Log Collection Jobs
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collection/
type LogCollectionService interface {
	List(context.Context, string, *LogListOptions) (*LogCollectionJobs, *Response, error)
	Get(context.Context, string, string, *LogListOptions) (*LogCollectionJob, *Response, error)
	Retry(context.Context, string, string) (*Response, error)
	Create(context.Context, string, *LogCollectionJob) (*LogCollectionJob, *Response, error)
	Extend(context.Context, string, string, *LogCollectionJob) (*Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// LogCollectionServiceOp provides an implementation of the DiagnosticsService interface.
type LogCollectionServiceOp service

var _ LogCollectionService = &LogCollectionServiceOp{}

// LogsService is an interface for interfacing with the Log Collection Jobs
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collection/
type LogsService interface {
	Download(context.Context, string, string, io.Writer) (*Response, error)
}

// LogsServiceOp handles communication with the Log Collection Jobs download method of the
// MongoDB Ops Manager API.
type LogsServiceOp struct {
	Client atlas.GZipRequestDoer
}

var _ LogsService = &LogsServiceOp{}

// LogCollectionJob represents a Log Collection Job in the MongoDB Ops Manager API.
type LogCollectionJob struct {
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
	DownloadURL                string      `json:"downloadUrl,omitempty"`
	Redacted                   *bool       `json:"redacted,omitempty"`
	LogTypes                   []string    `json:"logTypes,omitempty"`
	SizeRequestedPerFileBytes  int64       `json:"sizeRequestedPerFileBytes,omitempty"`
	UncompressedSizeTotalBytes int64       `json:"uncompressedSizeTotalBytes,omitempty"`
	ChildJobs                  []*ChildJob `json:"childJobs,omitempty"` // ChildJobs included if verbose is true
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

// LogCollectionJobs represents a array of LogCollectionJobs.
type LogCollectionJobs struct {
	Links      []*atlas.Link       `json:"links"`
	Results    []*LogCollectionJob `json:"results"`
	TotalCount int                 `json:"totalCount"`
}

// LogListOptions specifies the optional parameters to List methods that
// support pagination.
type LogListOptions struct {
	atlas.ListOptions
	Verbose bool `url:"verbose,omitempty"`
}

// List gets all collection log jobs.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-get-all/
func (s *LogCollectionServiceOp) List(ctx context.Context, groupID string, opts *LogListOptions) (*LogCollectionJobs, *Response, error) {
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

	root := new(LogCollectionJobs)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets a log collection job.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-get-one/
func (s *LogCollectionServiceOp) Get(ctx context.Context, groupID, jobID string, opts *LogListOptions) (*LogCollectionJob, *Response, error) {
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

	root := new(LogCollectionJob)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create creates a log collection job.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-submit/
func (s *LogCollectionServiceOp) Create(ctx context.Context, groupID string, log *LogCollectionJob) (*LogCollectionJob, *Response, error) {
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

	root := new(LogCollectionJob)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Extend extends the expiration data of a log collection job.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-update-one/
func (s *LogCollectionServiceOp) Extend(ctx context.Context, groupID, jobID string, log *LogCollectionJob) (*Response, error) {
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
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-retry/
func (s *LogCollectionServiceOp) Retry(ctx context.Context, groupID, jobID string) (*Response, error) {
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

// Delete removes a log collection job.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-delete-one/
func (s *LogCollectionServiceOp) Delete(ctx context.Context, groupID, jobID string) (*Response, error) {
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

// Download allows to download a log from a collection job.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/log-collections/log-collections-download-job/
func (s *LogsServiceOp) Download(ctx context.Context, groupID, jobID string, out io.Writer) (*Response, error) {
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
