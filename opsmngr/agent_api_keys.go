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
	Create(context.Context, string) (*AgentAPIKey, *atlas.Response, error)
	List(context.Context, string) (*AgentAPIKeys, *atlas.Response, error)
	Delete(context.Context, string, string) (*atlas.Response, error)
}

type AgentAPIKeysServiceOp struct {
	Client atlas.RequestDoer
}

type AgentAPIKey struct {
	ID            string  `json:"id"`
	Key           string  `json:"key"`
	Desc          string  `json:"desc"`
	CreatedTime   string  `json:"createdTime"`
	CreatedUserID *string `json:"createdUserId"`
	CreatedIPAddr *string `json:"createdIpAddr"`
	CreatedBy     string  `json:"createdBy"`
}

type AgentAPIKeys struct {
	Links      []*atlas.Link  `json:"links"`
	Results    []*AgentAPIKey `json:"results"`
	TotalCount int            `json:"totalCount"`
}

func (s *AgentAPIKeysServiceOp) Create(ctx context.Context, projectID string) (*AgentAPIKey, *atlas.Response, error) {
	path := fmt.Sprintf(agentAPIKeysBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
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

func (s *AgentAPIKeysServiceOp) List(ctx context.Context, projectID string) (*AgentAPIKeys, *atlas.Response, error) {
	path := fmt.Sprintf(agentAPIKeysBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AgentAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

func (s *AgentAPIKeysServiceOp) Delete(ctx context.Context, projectID, agentAPIKey string) (*atlas.Response, error) {
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
