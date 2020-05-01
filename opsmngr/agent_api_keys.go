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
	agentAPIKeysBasePath = "groups/%s/agentapikeys"
)

type AgentAPIKeysService interface {
	Create(context.Context, string, *AgentAPIKeysRequest) (*AgentAPIKey, *atlas.Response, error)
	List(context.Context, string) ([]*AgentAPIKey, *atlas.Response, error)
	Delete(context.Context, string, string) (*atlas.Response, error)
}

type AgentAPIKeysServiceOp service

type AgentAPIKey struct {
	ID            string  `json:"_id"`
	Key           string  `json:"key"`
	Desc          string  `json:"desc"`
	CreatedTime   int64   `json:"createdTime"`
	CreatedUserID *string `json:"createdUserId"`
	CreatedIPAddr *string `json:"createdIpAddr"`
	CreatedBy     string  `json:"createdBy"`
}

type AgentAPIKeysRequest struct {
	Desc string `json:"desc"`
}

// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agentapikeys/create-one-agent-api-key/
func (s *AgentAPIKeysServiceOp) Create(ctx context.Context, projectID string, agent *AgentAPIKeysRequest) (*AgentAPIKey, *atlas.Response, error) {
	path := fmt.Sprintf(agentAPIKeysBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, agent)
	if err != nil {
		return nil, nil, err
	}

	root := new(AgentAPIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agentapikeys/get-all-agent-api-keys-for-project/
func (s *AgentAPIKeysServiceOp) List(ctx context.Context, groupID string) ([]*AgentAPIKey, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(agentAPIKeysBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := make([]*AgentAPIKey, 0)
	resp, err := s.Client.Do(ctx, req, &root)

	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// See more: hhttps://docs.opsmanager.mongodb.com/current/reference/api/agentapikeys/delete-one-agent-api-key/
func (s *AgentAPIKeysServiceOp) Delete(ctx context.Context, groupID, agentAPIKey string) (*atlas.Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}
	if agentAPIKey == "" {
		return nil, atlas.NewArgError("agentAPIKey", "must be set")
	}
	basePath := fmt.Sprintf(agentAPIKeysBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, agentAPIKey)
	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}
