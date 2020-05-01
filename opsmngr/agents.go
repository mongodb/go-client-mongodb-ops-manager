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

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	agentsBasePath = "groups/%s/agents"
)

// AgentsService is an interface for interfacing with the Agents
// endpoints of the MongoDB Cloud API.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents/
type AgentsService interface {
	List(context.Context, string, string) (*Agents, *atlas.Response, error)
	ListLinks(context.Context, string) (*Agents, *atlas.Response, error)
}

// AgentsServiceOp handles communication with the Agent related methods of the MongoDB Cloud API
type AgentsServiceOp struct {
	Client atlas.RequestDoer
}

type Agent struct {
	TypeName  string  `json:"typeName"`
	Hostname  string  `json:"hostname"`
	ConfCount int64   `json:"confCount"`
	LastConf  string  `json:"lastConf"`
	StateName string  `json:"stateName"`
	PingCount int64   `json:"pingCount"`
	IsManaged bool    `json:"isManaged"`
	LastPing  string  `json:"lastPing"`
	Tag       *string `json:"tag"`
}

type Agents struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Agent      `json:"results"`
	TotalCount int           `json:"totalCount"`
}

// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents-get-all/
func (s *AgentsServiceOp) ListLinks(ctx context.Context, groupID string) (*Agents, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(agentsBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Agents)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents-get-by-type/
func (s *AgentsServiceOp) List(ctx context.Context, groupID, agentType string) (*Agents, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if agentType == "" {
		return nil, nil, atlas.NewArgError("agentType", "must be set")
	}
	basePath := fmt.Sprintf(agentsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, agentType)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Agents)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
