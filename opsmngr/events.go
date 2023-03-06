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
	eventsPathProjects     = "api/public/v1.0/groups/%s/events"
	eventsPathOrganization = "api/public/v1.0/orgs/%s/events"
)

// EventsServiceOp handles communication with the Event related methods
// of the MongoDB Ops Manager API.
type EventsServiceOp service

var _ atlas.EventsService = &EventsServiceOp{}

type (
	EventListOptions = atlas.EventListOptions
	EventResponse    = atlas.EventResponse
	Event            = atlas.Event
)

// ListOrganizationEvents lists all events in the organization associated to {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/events/get-all-events-for-org/
func (s *EventsServiceOp) ListOrganizationEvents(ctx context.Context, orgID string, listOptions *EventListOptions) (*EventResponse, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	path := fmt.Sprintf(eventsPathOrganization, orgID)

	// Add query params from listOptions
	path, err := setQueryParams(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(EventResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// GetOrganizationEvent gets the event specified to {EVENT-ID} from the organization associated to {ORG-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/events/get-one-event-for-org/
func (s *EventsServiceOp) GetOrganizationEvent(ctx context.Context, orgID, eventID string) (*Event, *Response, error) {
	if orgID == "" {
		return nil, nil, atlas.NewArgError("orgID", "must be set")
	}
	if eventID == "" {
		return nil, nil, atlas.NewArgError("eventID", "must be set")
	}
	basePath := fmt.Sprintf(eventsPathOrganization, orgID)
	path := fmt.Sprintf("%s/%s", basePath, eventID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Event)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// ListProjectEvents lists all events in the project associated to {PROJECT-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/events/get-all-events-for-project/
func (s *EventsServiceOp) ListProjectEvents(ctx context.Context, groupID string, listOptions *atlas.EventListOptions) (*EventResponse, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(eventsPathProjects, groupID)

	// Add query params from listOptions
	path, err := setQueryParams(path, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(EventResponse)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, nil
}

// GetProjectEvent gets the alert specified to {EVENT-ID} from the project associated to {PROJECT-ID}.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/events/get-one-event-for-project/
func (s *EventsServiceOp) GetProjectEvent(ctx context.Context, groupID, eventID string) (*Event, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	if eventID == "" {
		return nil, nil, atlas.NewArgError("eventID", "must be set")
	}
	basePath := fmt.Sprintf(eventsPathProjects, groupID)
	path := fmt.Sprintf("%s/%s", basePath, eventID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Event)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
