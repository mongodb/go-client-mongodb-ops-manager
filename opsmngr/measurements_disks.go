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

const hostDiskMeasurementsPath = "api/public/v1.0/groups/%s/hosts/%s/disks/%s/measurements"

// Disk measurements provide data on IOPS, disk use, and disk latency on the disk partitions for hosts running MongoDB that the Automations collect.
// You must run Ops Manager Automation to retrieve disk measurements.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/measures/get-disk-measurements/
func (s *MeasurementsServiceOp) Disk(ctx context.Context, groupID, hostID, diskName string, opts *ProcessMeasurementListOptions) (*ProcessDiskMeasurements, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	if hostID == "" {
		return nil, nil, NewArgError("hostID", "must be set")
	}

	if diskName == "" {
		return nil, nil, NewArgError("diskName", "must be set")
	}

	basePath := fmt.Sprintf(hostDiskMeasurementsPath, groupID, hostID, diskName)

	// Add query params from listOptions
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(ProcessDiskMeasurements)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
