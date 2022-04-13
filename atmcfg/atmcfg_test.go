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

package atmcfg

import (
	"testing"

	"go.mongodb.org/ops-manager/opsmngr"
)

const clusterName = "cluster_1"

func TestIsGoalState(t *testing.T) {
	tests := []struct {
		name string
		s    *opsmngr.AutomationStatus
		want bool
	}{
		{
			name: "reached goal",
			s: &opsmngr.AutomationStatus{
				Processes: []opsmngr.ProcessStatus{
					{
						LastGoalVersionAchieved: 0,
						Name:                    "test",
						Hostname:                "test",
					},
				},
				GoalVersion: 0,
			},
			want: true,
		},
		{
			name: "pending goal",
			s: &opsmngr.AutomationStatus{
				Processes: []opsmngr.ProcessStatus{
					{
						LastGoalVersionAchieved: 0,
						Name:                    "test",
						Hostname:                "test",
					},
				},
				GoalVersion: 1,
			},
			want: false,
		},
		{
			name: "empty",
			s:    &opsmngr.AutomationStatus{},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGoalState(tt.s); got != tt.want {
				t.Errorf("IsGoalState() = %v, want %v", got, tt.want)
			}
		})
	}
}
