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
	"io"
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	diagnosticsBasePath = "groups/%s/diagnostics"
)

// DiagnosticsService is an interface for interfacing with Diagnostic Archives in MongoDB Ops Manager APIs
//
// https://docs.opsmanager.mongodb.com/current/reference/api/diagnostic-archives/
type DiagnosticsService interface {
	Get(context.Context, string, *DiagnosticsListOpts, io.Writer) (*atlas.Response, error)
}

// DiagnosticsServiceOp provides an implementation of the DiagnosticsService interface.
type DiagnosticsServiceOp struct {
	Client atlas.GZipRequestDoer
}

// DiagnosticsListOpts query options for getting the archive.
type DiagnosticsListOpts struct {
	Limit   int64 `json:"limit,omitempty"`
	Minutes int64 `json:"minutes,omitempty"`
}

// Get retrieves the project’s diagnostics archive file.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/diagnostics/get-project-diagnostic-archive/
func (s *DiagnosticsServiceOp) Get(ctx context.Context, groupID string, opts *DiagnosticsListOpts, out io.Writer) (*atlas.Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(diagnosticsBasePath, groupID)
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, err
	}

	req, err := s.Client.NewGZipRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, out)

	return resp, err
}
