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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

type (
	ProcessMeasurements           = atlas.ProcessMeasurements
	ProcessDiskMeasurements       = atlas.ProcessDiskMeasurements
	ProcessDatabaseMeasurements   = atlas.ProcessDatabaseMeasurements
	ProcessMeasurementListOptions = atlas.ProcessMeasurementListOptions
)

// MeasurementsService provides access to the measurement related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/measurements/
type MeasurementsService interface {
	Host(context.Context, string, string, *ProcessMeasurementListOptions) (*ProcessMeasurements, *Response, error)
	Disk(context.Context, string, string, string, *ProcessMeasurementListOptions) (*ProcessDiskMeasurements, *Response, error)
	Database(context.Context, string, string, string, *ProcessMeasurementListOptions) (*ProcessDatabaseMeasurements, *Response, error)
}

// MeasurementsServiceOp provides an implementation of the MeasurementsService interface.
type MeasurementsServiceOp service

var _ MeasurementsService = &MeasurementsServiceOp{}
