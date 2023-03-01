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

const (
	alertConfigurationPath = "api/public/v1.0/groups/%s/alertConfigs"
	fieldNamesPath         = "api/public/v1.0/alertConfigs/matchers/fieldNames"
)

// AlertConfigurationsServiceOp handles communication with the AlertConfiguration related methods
// of the MongoDB Ops Manager API.
type AlertConfigurationsServiceOp service

var _ atlas.AlertConfigurationsService = &AlertConfigurationsServiceOp{}

// Create creates an alert configuration for the project associated to {GROUP-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/alert-configurations-create-config/
func (s *AlertConfigurationsServiceOp) Create(ctx context.Context, groupID string, createReq *atlas.AlertConfiguration) (*atlas.AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if createReq == nil {
		return nil, nil, atlas.NewArgError("createReq", "cannot be nil")
	}

	path := fmt.Sprintf(alertConfigurationPath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.AlertConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// EnableAnAlertConfig Enables/disables the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/alert-configurations-enable-disable-config/
func (s *AlertConfigurationsServiceOp) EnableAnAlertConfig(ctx context.Context, groupID, alertConfigID string, enabled *bool) (*atlas.AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, nil, atlas.NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, atlas.AlertConfiguration{Enabled: enabled})
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.AlertConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetAnAlertConfig gets the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/alert-configurations-get-config/
func (s *AlertConfigurationsServiceOp) GetAnAlertConfig(ctx context.Context, groupID, alertConfigID string) (*atlas.AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, nil, atlas.NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.AlertConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetOpenAlertsConfig gets all open alerts for the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/alert-configurations-get-open-alerts/
func (s *AlertConfigurationsServiceOp) GetOpenAlertsConfig(ctx context.Context, groupID, alertConfigID string) ([]atlas.AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, nil, atlas.NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s/alerts", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.AlertConfigurationsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}
	return root.Results, resp, err
}

// List gets all alert configurations for the project associated to {GROUP-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/alert-configurations-get-all-configs/
func (s *AlertConfigurationsServiceOp) List(ctx context.Context, groupID string, listOptions *ListOptions) ([]atlas.AlertConfiguration, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(alertConfigurationPath, groupID)

	// Add query params from listOptions
	path, err := setQueryParams(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.AlertConfigurationsResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if l := root.Links; l != nil {
		resp.Links = l
	}

	return root.Results, resp, nil
}

// Update the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/alert-configurations-update-config/
func (s *AlertConfigurationsServiceOp) Update(ctx context.Context, groupID, alertConfigID string, updateReq *atlas.AlertConfiguration) (*atlas.AlertConfiguration, *Response, error) {
	if updateReq == nil {
		return nil, nil, atlas.NewArgError("updateRequest", "cannot be nil")
	}
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, nil, atlas.NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, updateReq)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.AlertConfiguration)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Delete the alert configuration specified to {ALERT-CONFIG-ID} for the project associated to {GROUP-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/alert-configurations-delete-config/
func (s *AlertConfigurationsServiceOp) Delete(ctx context.Context, groupID, alertConfigID string) (*Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}
	if alertConfigID == "" {
		return nil, atlas.NewArgError("alertConfigID", "must be set")
	}

	basePath := fmt.Sprintf(alertConfigurationPath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, alertConfigID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// ListMatcherFields gets all field names that the matchers.fieldName parameter accepts when you create or update an Alert Configuration.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/alert-configurations-get-matchers-field-names/
func (s *AlertConfigurationsServiceOp) ListMatcherFields(ctx context.Context) ([]string, *Response, error) {
	req, err := s.Client.NewRequest(ctx, http.MethodGet, fieldNamesPath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := make([]string, 0)
	resp, err := s.Client.Do(ctx, req, &root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
