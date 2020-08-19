package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const whitelistAPIKeysPath = "admin/whitelist"

// GlobalAPIKeyWhitelistsService provides access to the global alerts related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/global-api-key-whitelists/
type GlobalAPIKeyWhitelistsService interface {
	List(context.Context, *atlas.ListOptions) (*GlobalWhitelistAPIKeys, *atlas.Response, error)
	Get(context.Context, string) (*GlobalWhitelistAPIKey, *atlas.Response, error)
	Create(context.Context, []*WhitelistAPIKeysReq) (*GlobalWhitelistAPIKeys, *atlas.Response, error)
	Delete(context.Context, string) (*atlas.Response, error)
}

// GlobalAPIKeyWhitelistsServiceOp provides an implementation of the GlobalAPIKeyWhitelistsService interface
type GlobalAPIKeyWhitelistsServiceOp service

var _ GlobalAPIKeyWhitelistsService = &GlobalAPIKeyWhitelistsServiceOp{}

// GlobalWhitelistAPIKey represents a Whitelist API key.
type GlobalWhitelistAPIKey struct {
	ID          string `json:"id,omitempty"`
	CidrBlock   string `json:"cidrBlock,omitempty"`
	Created     string `json:"created,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Updated     string `json:"updated,omitempty"`
}

// GlobalWhitelistAPIKeys represents all Whitelist API keys.
type GlobalWhitelistAPIKeys struct {
	Results    []*GlobalWhitelistAPIKey `json:"results,omitempty"`    // Includes one GlobalWhitelistAPIKey object for each item detailed in the results array section.
	Links      []*atlas.Link            `json:"links,omitempty"`      // One or more links to sub-resources and/or related resources.
	TotalCount int                      `json:"totalCount,omitempty"` // Count of the total number of items in the result set. It may be greater than the number of objects in the results array if the entire result set is paginated.
}

type WhitelistAPIKeysReq struct {
	CidrBlock   string `json:"cidrBlock"`
	Description string `json:"description"`
}

// List all global whitelist entries.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/get-all-global-whitelist/
func (s *GlobalAPIKeyWhitelistsServiceOp) List(ctx context.Context, listOptions *atlas.ListOptions) (*GlobalWhitelistAPIKeys, *atlas.Response, error) {
	path, err := setQueryParams(whitelistAPIKeysPath, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalWhitelistAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Get one global whitelist entry.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/get-one-global-whitelist/
func (s *GlobalAPIKeyWhitelistsServiceOp) Get(ctx context.Context, id string) (*GlobalWhitelistAPIKey, *atlas.Response, error) {
	if id == "" {
		return nil, nil, atlas.NewArgError("id", "must be set")
	}

	path := fmt.Sprintf("%s/%s", whitelistAPIKeysPath, id)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalWhitelistAPIKey)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Create one global whitelist entry.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/create-one-global-whitelist/
func (s *GlobalAPIKeyWhitelistsServiceOp) Create(ctx context.Context, createRequest []*WhitelistAPIKeysReq) (*GlobalWhitelistAPIKeys, *atlas.Response, error) {
	if createRequest == nil {
		return nil, nil, atlas.NewArgError("createRequest", "cannot be nil")
	}

	req, err := s.Client.NewRequest(ctx, http.MethodPost, whitelistAPIKeysPath, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalWhitelistAPIKeys)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the global whitelist entry specified by id.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/api-keys/global/delete-one-global-whitelist/
func (s *GlobalAPIKeyWhitelistsServiceOp) Delete(ctx context.Context, id string) (*atlas.Response, error) {
	if id == "" {
		return nil, atlas.NewArgError("id", "must be set")
	}

	path := fmt.Sprintf("%s/%s", whitelistAPIKeysPath, id)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
