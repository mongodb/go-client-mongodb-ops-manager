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

	mux.HandleFunc(fmt.Sprintf("/groups/%s/backupConfigs/%s/snapshotSchedule", groupID, clusterID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			  "clusterId" : "6b8cd61180eef547110159d9",
			  "dailySnapshotRetentionDays" : 7,
			  "groupId" : "5c8100bcf2a30b12ff88258f",
			  "monthlySnapshotRetentionMonths" : 13,
			  "pointInTimeWindowHours": 24,
			  "snapshotIntervalHours" : 6,
			  "snapshotRetentionDays" : 2,
			  "weeklySnapshotRetentionWeeks" : 4
}`)
	})

	snapshot, _, err := client.SnapshotSchedule.Get(ctx, groupID, clusterID)
	if err != nil {
		t.Fatalf("SnapshotSchedule.Get returned error: %v", err)
	}

	expected := &SnapshotSchedule{
		ClusterID:                      clusterID,
		GroupID:                        groupID,
		MonthlySnapshotRetentionMonths: 13,
		PointInTimeWindowHours:         24,
		SnapshotIntervalHours:          6,
		SnapshotRetentionDays:          2,
		WeeklySnapshotRetentionWeeks:   4,
		DailySnapshotRetentionDays:     7,
	}
	if diff := deep.Equal(snapshot, expected); diff != nil {
		t.Error(diff)
	}
}

func TestSnapshotScheduleServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/groups/%s/backupConfigs/%s/snapshotSchedule", groupID, clusterID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprint(w, `{
			  "clusterId" : "6b8cd61180eef547110159d9",
			  "dailySnapshotRetentionDays" : 7,
			  "groupId" : "5c8100bcf2a30b12ff88258f",
			  "monthlySnapshotRetentionMonths" : 13,
			  "pointInTimeWindowHours": 24,
			  "snapshotIntervalHours" : 6,
			  "snapshotRetentionDays" : 2,
			  "weeklySnapshotRetentionWeeks" : 4
}`)
	})

	snapshotSchedule := &SnapshotSchedule{
		ClusterID:                      clusterID,
		GroupID:                        groupID,
		MonthlySnapshotRetentionMonths: 13,
		PointInTimeWindowHours:         24,
		SnapshotIntervalHours:          6,
		SnapshotRetentionDays:          2,
		WeeklySnapshotRetentionWeeks:   4,
		DailySnapshotRetentionDays:     7,
	}

	snapshot, _, err := client.SnapshotSchedule.Update(ctx, snapshotSchedule.GroupID, snapshotSchedule.ClusterID, snapshotSchedule)
	if err != nil {
		t.Fatalf("SnapshotSchedule.Update returned error: %v", err)
	}

	expected := &SnapshotSchedule{
		ClusterID:                      clusterID,
		GroupID:                        groupID,
		MonthlySnapshotRetentionMonths: 13,
		PointInTimeWindowHours:         24,
		SnapshotIntervalHours:          6,
		SnapshotRetentionDays:          2,
		WeeklySnapshotRetentionWeeks:   4,
		DailySnapshotRetentionDays:     7,
	}
	if diff := deep.Equal(snapshot, expected); diff != nil {
		t.Error(diff)
	}
}
