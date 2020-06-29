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

const hostDatabaseMeasurementsPath = "groups/%s/hosts/%s/databases/%s/measurements"

// Database measurements provide statistics on database performance and storage.
// The Monitoring collects database measurements through the dbStats command.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/measures/get-database-measurements/
func (s *MeasurementsServiceOp) Database(ctx context.Context, groupID, hostID, databaseName string, opts *atlas.ProcessMeasurementListOptions) (*atlas.ProcessDatabaseMeasurements, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if hostID == "" {
		return nil, nil, atlas.NewArgError("hostID", "must be set")
	}

	if databaseName == "" {
		return nil, nil, atlas.NewArgError("diskName", "must be set")
	}

	basePath := fmt.Sprintf(hostDatabaseMeasurementsPath, groupID, hostID, databaseName)

	// Add query params from listOptions
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.ProcessDatabaseMeasurements)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
