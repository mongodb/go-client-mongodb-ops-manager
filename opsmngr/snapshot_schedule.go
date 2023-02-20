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

const snapshotScheduleBasePath = "api/public/v1.0/groups/%s/backupConfigs/%s/snapshotSchedule"

// SnapshotScheduleService is an interface for using the Snapshot schedule
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/snapshot-schedule/
type SnapshotScheduleService interface {
	Get(context.Context, string, string) (*SnapshotSchedule, *Response, error)
	Update(context.Context, string, string, *SnapshotSchedule) (*SnapshotSchedule, *Response, error)
}

// SnapshotScheduleServiceOp provides an implementation of the SnapshotScheduleService interface.
type SnapshotScheduleServiceOp service

var _ SnapshotScheduleService = &SnapshotScheduleServiceOp{}

type SnapshotSchedule struct {
	ClusterID                      string        `json:"clusterId"`
	GroupID                        string        `json:"groupId"`
	ReferenceTimeZoneOffset        string        `json:"referenceTimeZoneOffset,omitempty"`
	DailySnapshotRetentionDays     *int          `json:"dailySnapshotRetentionDays,omitempty"`
	ClusterCheckpointIntervalMin   int           `json:"clusterCheckpointIntervalMin,omitempty"`
	Links                          []*atlas.Link `json:"links,omitempty"`
	MonthlySnapshotRetentionMonths *int          `json:"monthlySnapshotRetentionMonths,omitempty"`
	PointInTimeWindowHours         *int          `json:"pointInTimeWindowHours,omitempty"`
	ReferenceHourOfDay             *int          `json:"referenceHourOfDay,omitempty"`
	ReferenceMinuteOfHour          *int          `json:"referenceMinuteOfHour,omitempty"`
	SnapshotIntervalHours          int           `json:"snapshotIntervalHours,omitempty"`
	SnapshotRetentionDays          int           `json:"snapshotRetentionDays,omitempty"`
	WeeklySnapshotRetentionWeeks   *int          `json:"weeklySnapshotRetentionWeeks,omitempty"`
}

// Get gets the snapshot schedule for an instance.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/backup/get-snapshot-schedule/
func (s *SnapshotScheduleServiceOp) Get(ctx context.Context, groupID, clusterID string) (*SnapshotSchedule, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if clusterID == "" {
		return nil, nil, atlas.NewArgError("clusterID", "must be set")
	}

	path := fmt.Sprintf(snapshotScheduleBasePath, groupID, clusterID)
	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(SnapshotSchedule)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates the parameters of snapshot creation and retention.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/backup/update-one-snapshot-schedule-by-cluster-id/
func (s *SnapshotScheduleServiceOp) Update(ctx context.Context, groupID, clusterID string, snapshot *SnapshotSchedule) (*SnapshotSchedule, *Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if clusterID == "" {
		return nil, nil, atlas.NewArgError("clusterID", "must be set")
	}

	path := fmt.Sprintf(snapshotScheduleBasePath, groupID, clusterID)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, snapshot)
	if err != nil {
		return nil, nil, err
	}

	root := new(SnapshotSchedule)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
