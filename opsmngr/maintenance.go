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
	maintenanceWindowsBasePath = "api/public/v1.0/groups/%s/maintenanceWindows"
)

// MaintenanceWindowsService is an interface for interfacing with the Maintenance Windows
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/maintenance-windows/
type MaintenanceWindowsService interface {
	Get(context.Context, string, string) (*MaintenanceWindow, *Response, error)
	List(context.Context, string) (*MaintenanceWindows, *Response, error)
	Create(context.Context, string, *MaintenanceWindow) (*MaintenanceWindow, *Response, error)
	Update(context.Context, string, string, *MaintenanceWindow) (*MaintenanceWindow, *Response, error)
	Delete(context.Context, string, string) (*Response, error)
}

// MaintenanceWindow represents MongoDB Maintenance Windows.
type MaintenanceWindow struct {
	ID             string   `json:"id,omitempty"`
	GroupID        string   `json:"groupId,omitempty"`
	Created        string   `json:"created,omitempty"`
	StartDate      string   `json:"startDate,omitempty"`
	EndDate        string   `json:"endDate,omitempty"`
	Updated        string   `json:"updated,omitempty"`
	AlertTypeNames []string `json:"alertTypeNames,omitempty"`
	Description    string   `json:"description,omitempty"`
}

// MaintenanceWindows is the response from the MaintenanceWindowsService.List.
type MaintenanceWindows struct {
	Links      []*atlas.Link        `json:"links,omitempty"`
	Results    []*MaintenanceWindow `json:"results,omitempty"`
	TotalCount int                  `json:"totalCount,omitempty"`
}

// MaintenanceWindowsServiceOp handles communication with the MaintenanceWindows related methods
// of the OpsManager Atlas API.
type MaintenanceWindowsServiceOp service

var _ MaintenanceWindowsService = &MaintenanceWindowsServiceOp{}

// List gets all maintenance windows.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/maintenance-windows-get-all/
func (s *MaintenanceWindowsServiceOp) List(ctx context.Context, groupID string) (*MaintenanceWindows, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(maintenanceWindowsBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(MaintenanceWindows)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Get gets a maintenance window.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/maintenance-windows-get-one/
func (s *MaintenanceWindowsServiceOp) Get(ctx context.Context, groupID, maintenanceWindowID string) (*MaintenanceWindow, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if maintenanceWindowID == "" {
		return nil, nil, atlas.NewArgError("maintenanceWindowID", "must be set")
	}

	basePath := fmt.Sprintf(maintenanceWindowsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, maintenanceWindowID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(MaintenanceWindow)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Create creates one maintenance window.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/maintenance-windows-create-one/
func (s *MaintenanceWindowsServiceOp) Create(ctx context.Context, groupID string, maintenanceWindow *MaintenanceWindow) (*MaintenanceWindow, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	path := fmt.Sprintf(maintenanceWindowsBasePath, groupID)
	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, maintenanceWindow)
	if err != nil {
		return nil, nil, err
	}

	root := new(MaintenanceWindow)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Update updates one maintenance window with an end date in the future.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/maintenance-windows-update-one/
func (s *MaintenanceWindowsServiceOp) Update(ctx context.Context, groupID, maintenanceWindowID string, maintenanceWindow *MaintenanceWindow) (*MaintenanceWindow, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if maintenanceWindowID == "" {
		return nil, nil, atlas.NewArgError("maintenanceWindowID", "must be set")
	}

	basePath := fmt.Sprintf(maintenanceWindowsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, maintenanceWindowID)

	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, maintenanceWindow)
	if err != nil {
		return nil, nil, err
	}

	root := new(MaintenanceWindow)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Delete removes one maintenance window with an end date in the future.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/maintenance-windows-delete-one/
func (s *MaintenanceWindowsServiceOp) Delete(ctx context.Context, groupID, maintenanceWindowID string) (*Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}

	if maintenanceWindowID == "" {
		return nil, atlas.NewArgError("maintenanceWindowID", "must be set")
	}

	basePath := fmt.Sprintf(maintenanceWindowsBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, maintenanceWindowID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
