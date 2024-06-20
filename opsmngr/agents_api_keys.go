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
)

const (
	agentAPIKeysBasePath = "api/public/v1.0/groups/%s/agentapikeys" //nolint:gosec // This is a path
)

// AgentAPIKey defines the structure for an Agent API key.
type AgentAPIKey struct {
	ID            string  `json:"_id"` //nolint:tagliatelle // correct from API
	Key           string  `json:"key"`
	Desc          string  `json:"desc"`
	CreatedTime   int64   `json:"createdTime"`
	CreatedUserID *string `json:"createdUserId"`
	CreatedIPAddr *string `json:"createdIpAddr"`
	CreatedBy     string  `json:"createdBy"`
}

// AgentAPIKeysRequest a creation request for Agent API keys.
type AgentAPIKeysRequest struct {
	Desc string `json:"desc"`
}

// CreateAgentAPIKey creates a new agent API key.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agentapikeys/create-one-agent-api-key/
func (s *AgentsServiceOp) CreateAgentAPIKey(ctx context.Context, projectID string, agent *AgentAPIKeysRequest) (*AgentAPIKey, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
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

// ListAgentAPIKeys lists agent API keys.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agentapikeys/get-all-agent-api-keys-for-project/
func (s *AgentsServiceOp) ListAgentAPIKeys(ctx context.Context, projectID string) ([]*AgentAPIKey, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(agentAPIKeysBasePath, projectID)

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

// DeleteAgentAPIKey removes an agent API key from a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agentapikeys/delete-one-agent-api-key/
func (s *AgentsServiceOp) DeleteAgentAPIKey(ctx context.Context, projectID, agentAPIKey string) (*Response, error) {
	if projectID == "" {
		return nil, NewArgError("projectID", "must be set")
	}
	if agentAPIKey == "" {
		return nil, NewArgError("agentAPIKey", "must be set")
	}
	basePath := fmt.Sprintf(agentAPIKeysBasePath, projectID)
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
