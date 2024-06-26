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
	hostMeasurementsBasePath = "api/public/v1.0/groups/%s/hosts/%s/measurements"
)

// Host measurements provide data on the state of the MongoDB process.
// The Monitoring collects host measurements through the MongoDB serverStatus and dbStats commands.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/measures/get-host-process-system-measurements/
func (s *MeasurementsServiceOp) Host(ctx context.Context, projectID, hostID string, listOptions *ProcessMeasurementListOptions) (*ProcessMeasurements, *Response, error) {
	if projectID == "" {
		return nil, nil, NewArgError("projectID", "must be set")
	}
	if hostID == "" {
		return nil, nil, NewArgError("hostID", "must be set")
	}

	basePath := fmt.Sprintf(hostMeasurementsBasePath, projectID, hostID)
	path, err := setQueryParams(basePath, listOptions)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessMeasurements)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}
