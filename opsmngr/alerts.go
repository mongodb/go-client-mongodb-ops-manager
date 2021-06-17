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
	"fmt"
	"net/http"

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const alertPath = "api/public/v1.0/groups/%s/alerts"

// AlertsServiceOp provides an implementation of AlertsService.
type AlertsServiceOp service

var _ atlas.AlertsService = &AlertsServiceOp{}

// Get gets the alert specified to {ALERT-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alerts-get-alert/
func (s *AlertsServiceOp) Get(ctx context.Context, groupID, alertID string) (*atlas.Alert, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if alertID == "" {
		return nil, nil, atlas.NewArgError("alertID", "must be set")
	}

	basePath := fmt.Sprintf(alertPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Alert)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// List gets all alert for the project associated to {GROUP-ID}.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alerts-get-all-alerts/
func (s *AlertsServiceOp) List(ctx context.Context, groupID string, listOptions *atlas.AlertsListOptions) (*atlas.AlertsResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(alertPath, groupID)

	// Add query params from listOptions
	path, err := setQueryParams(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.AlertsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root, resp, nil
}

// Acknowledge allows to acknowledge an alert.
//
// See more: https://docs.atlas.mongodb.com/reference/api/alerts-acknowledge-alert/
func (s *AlertsServiceOp) Acknowledge(ctx context.Context, groupID, alertID string, params *atlas.AcknowledgeRequest) (*atlas.Alert, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if alertID == "" {
		return nil, nil, atlas.NewArgError("alertID", "must be set")
	}

	if params == nil {
		return nil, nil, atlas.NewArgError("params", "must be set")
	}

	basePath := fmt.Sprintf(alertPath, groupID)

	path := fmt.Sprintf("%s/%s", basePath, alertID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, params)

	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Alert)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
