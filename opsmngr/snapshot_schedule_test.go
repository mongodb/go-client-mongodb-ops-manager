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
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

func TestSnapshotScheduleServiceOp_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/backupConfigs/%s/snapshotSchedule", projectID, clusterID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{
			  "clusterId" : "6b8cd61180eef547110159d9",
			  "dailySnapshotRetentionDays" : 7,
			  "groupId" : "%[1]s",
			  "monthlySnapshotRetentionMonths" : 13,
			  "pointInTimeWindowHours": 24,
			  "snapshotIntervalHours" : 6,
			  "snapshotRetentionDays" : 2,
			  "weeklySnapshotRetentionWeeks" : 4
}`, projectID)
	})

	snapshot, _, err := client.SnapshotSchedule.Get(ctx, projectID, clusterID)
	if err != nil {
		t.Fatalf("SnapshotSchedule.Get returned error: %v", err)
	}

	pointInTimeWindowHours := 24

	expected := &SnapshotSchedule{
		ClusterID:                      clusterID,
		GroupID:                        projectID,
		MonthlySnapshotRetentionMonths: pointer(13),
		PointInTimeWindowHours:         &pointInTimeWindowHours,
		SnapshotIntervalHours:          6,
		SnapshotRetentionDays:          2,
		WeeklySnapshotRetentionWeeks:   pointer(4),
		DailySnapshotRetentionDays:     pointer(7),
	}
	if diff := deep.Equal(snapshot, expected); diff != nil {
		t.Error(diff)
	}
}

func TestSnapshotScheduleServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/backupConfigs/%s/snapshotSchedule", projectID, clusterID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprintf(w, `{
			  "clusterId" : "6b8cd61180eef547110159d9",
			  "dailySnapshotRetentionDays" : 7,
			  "groupId" : "%[1]s",
			  "monthlySnapshotRetentionMonths" : 13,
			  "pointInTimeWindowHours": 24,
			  "snapshotIntervalHours" : 6,
			  "snapshotRetentionDays" : 2,
			  "weeklySnapshotRetentionWeeks" : 4
}`, projectID)
	})

	pointInTimeWindowHours := 24

	snapshotSchedule := &SnapshotSchedule{
		ClusterID:                      clusterID,
		GroupID:                        projectID,
		MonthlySnapshotRetentionMonths: pointer(13),
		PointInTimeWindowHours:         &pointInTimeWindowHours,
		SnapshotIntervalHours:          6,
		SnapshotRetentionDays:          2,
		WeeklySnapshotRetentionWeeks:   pointer(4),
		DailySnapshotRetentionDays:     pointer(7),
	}

	snapshot, _, err := client.SnapshotSchedule.Update(ctx, snapshotSchedule.GroupID, snapshotSchedule.ClusterID, snapshotSchedule)
	if err != nil {
		t.Fatalf("SnapshotSchedule.Update returned error: %v", err)
	}

	expected := &SnapshotSchedule{
		ClusterID:                      clusterID,
		GroupID:                        projectID,
		MonthlySnapshotRetentionMonths: pointer(13),
		PointInTimeWindowHours:         &pointInTimeWindowHours,
		SnapshotIntervalHours:          6,
		SnapshotRetentionDays:          2,
		WeeklySnapshotRetentionWeeks:   pointer(4),
		DailySnapshotRetentionDays:     pointer(7),
	}
	if diff := deep.Equal(snapshot, expected); diff != nil {
		t.Error(diff)
	}
}
