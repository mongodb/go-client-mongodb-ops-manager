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

// AgentService is an interface for interfacing with the Agents
// endpoints of the MongoDB Cloud API.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/agents/
type AgentService interface {
	List(context.Context, string, string) (*Agents, *atlas.Response, error)
	ListLinks(context.Context, string) (*Agents, *atlas.Response, error)
}

// AutomationConfigServiceOp handles communication with the Automation config related methods of the MongoDB Cloud API
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
func (s *AgentsServiceOp) ListLinks(ctx context.Context, projectID string) (*Agents, *atlas.Response, error) {
	path := fmt.Sprintf(agentsBasePath, projectID)

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
func (s *AgentsServiceOp) List(ctx context.Context, projectID, agentType string) (*Agents, *atlas.Response, error) {
	basePath := fmt.Sprintf(agentsBasePath, projectID)
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
