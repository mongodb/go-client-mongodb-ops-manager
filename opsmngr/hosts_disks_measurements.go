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

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const processDiskMeasurementsPath = "/groups/%s/hosts/%s/disks/%s/measurements"

// HostDiskMeasurementsService is an interface for interfacing with the Process Disk Measurements
// endpoints of the MongoDB Atlas API.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/measures/get-disk-measurements/
type HostDiskMeasurementsService interface {
	List(context.Context, string, string, string, *atlas.ProcessMeasurementListOptions) (*atlas.ProcessDiskMeasurements, *atlas.Response, error)
}

// HostDiskMeasurementsServiceOp handles communication with the Process Disk Measurements related methods of the
// MongoDB API
type HostDiskMeasurementsServiceOp struct {
	Client atlas.RequestDoer
}

var _ HostDiskMeasurementsService = &HostDiskMeasurementsServiceOp{}

// List gets measurements for a specific host MongoDB disk.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/measures/get-disk-measurements/
func (s *HostDiskMeasurementsServiceOp) List(ctx context.Context, groupID, hostID string, diskName string, opts *atlas.ProcessMeasurementListOptions) (*atlas.ProcessDiskMeasurements, *atlas.Response, error) {

	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if hostID == "" {
		return nil, nil, atlas.NewArgError("hostID", "must be set")
	}

	if diskName == "" {
		return nil, nil, atlas.NewArgError("diskName", "must be set")
	}

	basePath := fmt.Sprintf(processDiskMeasurementsPath, groupID, hostID, diskName)

	//Add query params from listOptions
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.ProcessDiskMeasurements)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
