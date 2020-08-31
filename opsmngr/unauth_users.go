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

	"github.com/google/go-querystring/query"
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	unauthUsersBasePath = "unauth/users"
)

// UnauthUsersService is an interface for interfacing with unauthenticated APIs
type UnauthUsersService interface {
	CreateFirstUser(context.Context, *atlas.AtlasUser, *WhitelistOpts) (*CreateUserResponse, *atlas.Response, error)
}

// UnauthUsersServiceOp handles communication with the unauthenticated API
type UnauthUsersServiceOp service

// CreateFirstUser creates the first user for a new installation.
//
// See more: https://docs.opsmanager.mongodb.com/master/reference/api/user-create-first/
func (s *UnauthUsersServiceOp) CreateFirstUser(ctx context.Context, user *atlas.AtlasUser, opts *WhitelistOpts) (*CreateUserResponse, *atlas.Response, error) {
	// if a whitelist was not specified, do not pass the parameter
	basePath := unauthUsersBasePath

	if opts != nil {
		v, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		basePath = fmt.Sprintf("%s?%s", unauthUsersBasePath, v.Encode())
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, basePath, user)
	if err != nil {
		return nil, nil, err
	}

	root := new(CreateUserResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

type WhitelistOpts struct {
	Whitelist []string `url:"whitelist"`
}

// CreateUserResponse API response for the CreateFirstUser() call
type CreateUserResponse struct {
	APIKey             string           `json:"apiKey,omitempty"`
	ProgrammaticAPIKey atlas.APIKey     `json:"programmaticApiKey,omitempty"`
	User               *atlas.AtlasUser `json:"user"`
}
