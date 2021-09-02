// Copyright 2021 MongoDB Inc
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
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const versionPath = "api/private/unauth/version"

// ServiceVersionService is an interface for interfacing with the version private endpoint of the MongoDB Atlas API.
type ServiceVersionService interface {
	Get(context.Context) (*ServiceVersion, *Response, error)
}

type ServiceVersionServiceOp struct {
	Client atlas.PlainRequestDoer
}

// Get gets a compressed (.gz) log file that contains a range of log messages for a particular host.
func (s *ServiceVersionServiceOp) Get(ctx context.Context) (*ServiceVersion, *Response, error) {
	req, err := s.Client.NewPlainRequest(ctx, http.MethodGet, versionPath)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return nil, nil, err
	}

	version := resp.ServiceVersion()

	return version, resp, nil
}
