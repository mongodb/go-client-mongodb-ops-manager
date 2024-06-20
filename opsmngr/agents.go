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
	agentsBasePath = "api/public/v1.0/groups/%s/agents"
	componentsPath = "api/public/v1.0/softwareComponents/versions"
)

// AgentsService provides access to the agent related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents/
type AgentsService interface {
	ListAgentLinks(context.Context, string) (*Agents, *Response, error)
	ListAgentsByType(context.Context, string, string, *ListOptions) (*Agents, *Response, error)
	CreateAgentAPIKey(context.Context, string, *AgentAPIKeysRequest) (*AgentAPIKey, *Response, error)
	ListAgentAPIKeys(context.Context, string) ([]*AgentAPIKey, *Response, error)
	DeleteAgentAPIKey(context.Context, string, string) (*Response, error)
	GlobalVersions(context.Context) (*SoftwareVersions, *Response, error)
	ProjectVersions(context.Context, string) (*AgentVersions, *Response, error)
}

// AgentsServiceOp provides an implementation of the AgentsService interface.
type AgentsServiceOp service

var _ AgentsService = new(AgentsServiceOp)

// Agent represents an Ops Manager agent.
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

// Agents is a paginated collection of Agent.
type Agents struct {
	Links      []*Link  `json:"links"`
	Results    []*Agent `json:"results"`
	TotalCount int      `json:"totalCount"`
}

// SoftwareVersions is a set of software components and their expected current and minimum versions.
type SoftwareVersions struct {
	AutomationVersion         string  `json:"automationVersion"`
	AutomationMinimumVersion  string  `json:"automationMinimumVersion"`
	BiConnectorVersion        string  `json:"biConnectorVersion"`
	BiConnectorMinimumVersion string  `json:"biConnectorMinimumVersion"`
	Links                     []*Link `json:"links"`
	MongoDBToolsVersion       string  `json:"mongoDbToolsVersion"`
}

// AgentVersions is a set of available agents and agent versions for a project.
type AgentVersions struct {
	Count                       int             `json:"count"`
	Entries                     []*AgentVersion `json:"entries"`
	IsAnyAgentNotManaged        bool            `json:"isAnyAgentNotManaged"`
	IsAnyAgentVersionDeprecated bool            `json:"isAnyAgentVersionDeprecated"`
	IsAnyAgentVersionOld        bool            `json:"isAnyAgentVersionOld"`
	Links                       []*Link         `json:"links"`
	MinimumAgentVersionDetected string          `json:"minimumAgentVersionDetected"`
	MinimumVersion              string          `json:"minimumVersion"`
}

type AgentVersion struct {
	Address             string `json:"address"`
	ConfCount           int    `json:"confCount"`
	Hostname            string `json:"hostname"`
	HostnameShort       string `json:"hostnameShort"`
	IsManaged           bool   `json:"isManaged"`
	IsModule            bool   `json:"isModule"`
	IsPrimary           bool   `json:"isPrimary"`
	IsVersionOld        bool   `json:"isVersionOld"`
	IsVersionDeprecated bool   `json:"isVersionDeprecated"`
	LastConf            int64  `json:"lastConf"`
	LastPing            int64  `json:"lastPing"`
	MsSinceLastConf     int64  `json:"msSinceLastConf"`
	NumProcess          int    `json:"numProcess"`
	PingCount           int    `json:"pingCount"`
	PingState           string `json:"pingState"`
	TagName             string `json:"tagName"`
	Version             string `json:"version"`
}

// ListAgentLinks gets links to monitoring, backup, and automation agent resources for a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents-get-all/
func (s *AgentsServiceOp) ListAgentLinks(ctx context.Context, groupID string) (*Agents, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
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

// ListAgentsByType gets agents of a specified type (i.e. Monitoring, Backup, or Automation) for a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents-get-by-type/
func (s *AgentsServiceOp) ListAgentsByType(ctx context.Context, groupID, agentType string, listOptions *ListOptions) (*Agents, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	if agentType == "" {
		return nil, nil, NewArgError("agentType", "must be set")
	}
	basePath := fmt.Sprintf(agentsBasePath, groupID)
	path, err := setQueryParams(fmt.Sprintf("%s/%s", basePath, agentType), listOptions)
	if err != nil {
		return nil, nil, err
	}
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

// GlobalVersions returns the supported versions of different software components
// supported by the Ops Manager installation.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents/get-agent-versions-global/
func (s *AgentsServiceOp) GlobalVersions(ctx context.Context) (*SoftwareVersions, *Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodGet, componentsPath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SoftwareVersions)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// ProjectVersions returns a list of versions of all MongoDB Agents, in the provided project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents/get-agent-versions-per-project/
func (s *AgentsServiceOp) ProjectVersions(ctx context.Context, groupID string) (*AgentVersions, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	basePath := fmt.Sprintf(agentsBasePath, groupID)
	path := fmt.Sprintf("%s/versions", basePath)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AgentVersions)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
