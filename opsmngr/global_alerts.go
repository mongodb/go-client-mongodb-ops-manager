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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	globalAlertsBasePath = "globalAlerts"
)

// GlobalAlertsService provides access to the global alerts related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/global-alerts/
type GlobalAlertsService interface {
	Get(context.Context, string) (*GlobalAlert, *Response, error)
	List(context.Context, *atlas.AlertsListOptions) (*GlobalAlerts, *Response, error)
	Acknowledge(context.Context, string, *atlas.AcknowledgeRequest) (*GlobalAlert, *Response, error)
}

// GlobalAlertsServiceOp provides an implementation of the GlobalAlertsService interface.
type GlobalAlertsServiceOp service

// GlobalAlert configuration struct.
type GlobalAlert struct {
	atlas.Alert
	SourceTypeName string        `json:"sourceTypeName,omitempty"`
	Tags           []string      `json:"tags,omitempty"`
	Links          []*atlas.Link `json:"links,omitempty"`
	HostID         string        `json:"hostId,omitempty"`
	ClusterID      string        `json:"clusterId,omitempty"`
}

// GlobalAlerts collection of configurations.
type GlobalAlerts struct {
	Links      []*atlas.Link  `json:"links"`
	Results    []*GlobalAlert `json:"results"`
	TotalCount int            `json:"totalCount"`
}

// Get gets a global alert.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/global-alerts/
func (s *GlobalAlertsServiceOp) Get(ctx context.Context, alertID string) (*GlobalAlert, *Response, error) {
	if alertID == "" {
		return nil, nil, atlas.NewArgError("alertID", "must be set")
	}
	path := fmt.Sprintf("%s/%s", globalAlertsBasePath, alertID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalAlert)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List gets all global alerts.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/global-alerts/
func (s *GlobalAlertsServiceOp) List(ctx context.Context, opts *atlas.AlertsListOptions) (*GlobalAlerts, *Response, error) {
	path, err := setQueryParams(globalAlertsBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalAlerts)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Acknowledge acknowledges a global alert.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/global-alerts/
func (s *GlobalAlertsServiceOp) Acknowledge(ctx context.Context, alertID string, body *atlas.AcknowledgeRequest) (*GlobalAlert, *Response, error) {
	if alertID == "" {
		return nil, nil, atlas.NewArgError("alertID", "must be set")
	}
	path := fmt.Sprintf("%s/%s", globalAlertsBasePath, alertID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, nil, err
	}

	root := new(GlobalAlert)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
