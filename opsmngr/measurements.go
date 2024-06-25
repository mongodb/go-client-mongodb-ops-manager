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
)

// ProcessMeasurements represents a MongoDB Process Measurements.
type ProcessMeasurements struct {
	End          string          `json:"end"`
	Granularity  string          `json:"granularity"`
	GroupID      string          `json:"groupId"`
	HostID       string          `json:"hostId"`
	Links        []*Link         `json:"links,omitempty"`
	Measurements []*Measurements `json:"measurements"`
	ProcessID    string          `json:"processId"`
	Start        string          `json:"start"`
}

// ProcessDiskMeasurements represents a MongoDB Process Disk Measurements.
type ProcessDiskMeasurements struct {
	*ProcessMeasurements
	PartitionName string `json:"partitionName"`
}

// ProcessDatabaseMeasurements represents a MongoDB process database measurements.
type ProcessDatabaseMeasurements struct {
	*ProcessMeasurements
	DatabaseName string `json:"databaseName"`
}

// ProcessMeasurementListOptions contains the list of options for Process Measurements.
type ProcessMeasurementListOptions struct {
	*ListOptions
	Granularity string   `url:"granularity"`
	Period      string   `url:"period,omitempty"`
	Start       string   `url:"start,omitempty"`
	End         string   `url:"end,omitempty"`
	M           []string `url:"m,omitempty"`
}

// Measurements represents a MongoDB Measurement.
type Measurements struct {
	DataPoints []*DataPoints `json:"dataPoints,omitempty"`
	Name       string        `json:"name"`
	Units      string        `json:"units"`
}

// DataPoints represents a MongoDB DataPoints.
type DataPoints struct {
	Timestamp string   `json:"timestamp"`
	Value     *float32 `json:"value"`
}

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
