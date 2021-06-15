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

const statusBlob = `{
  "processes": [
    {
      "plan": [],
      "lastGoalVersionAchieved": 2,
      "name": "shardedCluster_myShard_0_0",
      "hostname": "testDeploy-0"
    },
    {
      "plan": [],
      "lastGoalVersionAchieved": 2,
      "name": "shardedCluster_myShard_0_1",
      "hostname": "testDeploy-1"
    },
    {
      "plan": ["Download", "Start", "WaitRsInit"],
      "lastGoalVersionAchieved": 2,
      "name": "shardedCluster_myShard_0_2",
      "hostname": "testDeploy-2"
    }
  ],
  "goalVersion": 2
}`

func TestAutomation_GetStatus(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationStatus", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, statusBlob)
	})

	config, _, err := client.Automation.GetStatus(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.GetStatus returned error: %v", err)
	}

	expected := &AutomationStatus{
		GoalVersion: 2,
		Processes: []ProcessStatus{
			{
				Name:                    "shardedCluster_myShard_0_0",
				Hostname:                "testDeploy-0",
				Plan:                    []string{},
				LastGoalVersionAchieved: 2,
			},
			{
				Name:                    "shardedCluster_myShard_0_1",
				Hostname:                "testDeploy-1",
				Plan:                    []string{},
				LastGoalVersionAchieved: 2,
			},
			{
				Name:                    "shardedCluster_myShard_0_2",
				Plan:                    []string{"Download", "Start", "WaitRsInit"},
				Hostname:                "testDeploy-2",
				LastGoalVersionAchieved: 2,
			},
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}
