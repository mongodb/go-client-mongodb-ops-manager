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

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	diagnosticsBasePath  = "groups/%s/diagnostics"
	diagnosticsAdminPath = "admin/diagnostics"
)

// DiagnosticsService is an interface for interfacing with Diagnostic Archives in MongoDB Ops Manager APIs
// https://docs.opsmanager.mongodb.com/current/reference/api/diagnostic-archives/
type DiagnosticsService interface {
	Get(context.Context, string, io.Writer) (*atlas.Response, error)
	List(context.Context, io.Writer) (*atlas.Response, error)
}

type DiagnosticsServiceOp struct {
	Client atlas.GZipRequestDoer
}

// Get retrieves the project’s diagnostics archive file.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/diagnostics/get-project-diagnostic-archive/
func (s *DiagnosticsServiceOp) Get(ctx context.Context, projectID string, out io.Writer) (*atlas.Response, error) {
	path := fmt.Sprintf(diagnosticsBasePath, projectID)
	req, err := s.Client.NewGZipRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, out)

	return resp, err
}

// List retrieves all projects’ diagnostics archive file.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/diagnostics/get-global-diagnostic-archive/
func (s *DiagnosticsServiceOp) List(ctx context.Context, out io.Writer) (*atlas.Response, error) {
	req, err := s.Client.NewGZipRequest(ctx, http.MethodGet, diagnosticsAdminPath)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, out)

	return resp, err
}
