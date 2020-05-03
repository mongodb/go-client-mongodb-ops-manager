// Copyright 2019 MongoDB Inc
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
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

const jsonBlob = `{
  "auth" : {
    "authoritativeSet" : false,
    "autoAuthMechanism" : "MONGODB-CR",
    "disabled" : true
  },
  "processes" : [ {
    "args2_6" : {
      "net" : {
        "port" : 27000
      },
      "replication" : {
        "replSetName" : "myReplicaSet"
      },
      "storage" : {
        "dbPath" : "/data/rs1",
        "wiredTiger" : {
          "collectionConfig" : { },
          "engineConfig" : {
            "cacheSizeGB" : 0.5
          },
          "indexConfig" : { }
        }
      },
      "systemLog" : {
        "destination" : "file",
        "path" : "/data/rs1/mongodb.log"
      }
    },
    "authSchemaVersion" : 5,
    "disabled" : false,
    "featureCompatibilityVersion" : "4.2",
    "hostname" : "host0",
    "logRotate" : {
      "sizeThresholdMB" : 1000.0,
      "timeThresholdHrs" : 24
    },
    "manualMode" : false,
    "name" : "myReplicaSet_1",
    "processType" : "mongod",
    "version" : "4.2.2"
  }, {
    "args2_6" : {
      "net" : {
        "port" : 27010
      },
      "replication" : {
        "replSetName" : "myReplicaSet"
      },
      "storage" : {
        "dbPath" : "/data/rs2",
        "wiredTiger" : {
          "collectionConfig" : { },
          "engineConfig" : {
            "cacheSizeGB" : 0.5
          },
          "indexConfig" : { }
        }
      },
      "systemLog" : {
        "destination" : "file",
        "path" : "/data/rs2/mongodb.log"
      }
    },
    "authSchemaVersion" : 5,
    "disabled" : false,
    "featureCompatibilityVersion" : "4.2",
    "hostname" : "host1",
    "logRotate" : {
      "sizeThresholdMB" : 1000.0,
      "timeThresholdHrs" : 24
    },
    "manualMode" : false,
    "name" : "myReplicaSet_2",
    "processType" : "mongod",
    "version" : "4.2.2"
  }, {
    "args2_6" : {
      "net" : {
        "port" : 27020
      },
      "replication" : {
        "replSetName" : "myReplicaSet"
      },
      "storage" : {
        "dbPath" : "/data/rs3",
        "wiredTiger" : {
          "collectionConfig" : { },
          "engineConfig" : {
            "cacheSizeGB" : 0.5
          },
          "indexConfig" : { }
        }
      },
      "systemLog" : {
        "destination" : "file",
        "path" : "/data/rs3/mongodb.log"
      }
    },
    "authSchemaVersion" : 5,
    "disabled" : false,
    "featureCompatibilityVersion" : "4.2",
    "hostname" : "host0",
    "logRotate" : {
      "sizeThresholdMB" : 1000.0,
      "timeThresholdHrs" : 24
    },
    "manualMode" : false,
    "name" : "myReplicaSet_3",
    "processType" : "mongod",
    "version" : "4.2.2"
  } ],
  "replicaSets" : [ {
    "_id" : "myReplicaSet",
    "members" : [ {
      "_id" : 0,
      "arbiterOnly" : false,
      "buildIndexes" : true,
      "hidden" : false,
      "host" : "myReplicaSet_1",
      "priority" : 1.0,
      "slaveDelay" : 0,
      "votes" : 1
    }, {
      "_id" : 1,
      "arbiterOnly" : false,
      "buildIndexes" : true,
      "hidden" : false,
      "host" : "myReplicaSet_2",
      "priority" : 1.0,
      "slaveDelay" : 0,
      "votes" : 1
    }, {
      "_id" : 2,
      "arbiterOnly" : false,
      "buildIndexes" : true,
      "hidden" : false,
      "host" : "myReplicaSet_3",
      "priority" : 1.0,
      "slaveDelay" : 0,
      "votes" : 1
    } ],
    "protocolVersion" : "1",
    "settings" : { }
  } ],
  "uiBaseUrl" : null,
  "version" : 1
}`

func TestAAutomation_GetConfig(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, jsonBlob)
	})

	config, _, err := client.Automation.GetConfig(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.GetConfig returned error: %v", err)
	}

	expected := &AutomationConfig{
		Auth: Auth{
			AutoAuthMechanism: "MONGODB-CR",
			Disabled:          true,
			AuthoritativeSet:  false,
		},
		Processes: []*Process{
			{
				Name:                        "myReplicaSet_1",
				ProcessType:                 "mongod",
				Version:                     "4.2.2",
				AuthSchemaVersion:           5,
				FeatureCompatibilityVersion: "4.2",
				Disabled:                    false,
				ManualMode:                  false,
				Hostname:                    "host0",
				Args26: Args26{
					NET: Net{
						Port: 27000,
					},
					Storage: &Storage{
						DBPath: "/data/rs1",
					},
					SystemLog: SystemLog{
						Destination: "file",
						Path:        "/data/rs1/mongodb.log",
					},
					Replication: &Replication{
						ReplSetName: "myReplicaSet",
					},
				},
				LogRotate: &LogRotate{
					SizeThresholdMB:  1000.0,
					TimeThresholdHrs: 24,
				},
				LastGoalVersionAchieved: 0,
				Cluster:                 "",
			},
			{
				Name:                        "myReplicaSet_2",
				ProcessType:                 "mongod",
				Version:                     "4.2.2",
				AuthSchemaVersion:           5,
				FeatureCompatibilityVersion: "4.2",
				Disabled:                    false,
				ManualMode:                  false,
				Hostname:                    "host1",
				Args26: Args26{
					NET: Net{
						Port: 27010,
					},
					Storage: &Storage{
						DBPath: "/data/rs2",
					},
					SystemLog: SystemLog{
						Destination: "file",
						Path:        "/data/rs2/mongodb.log",
					},
					Replication: &Replication{
						ReplSetName: "myReplicaSet",
					},
				},
				LogRotate: &LogRotate{
					SizeThresholdMB:  1000.0,
					TimeThresholdHrs: 24,
				},
				LastGoalVersionAchieved: 0,
				Cluster:                 "",
			},
			{
				Name:                        "myReplicaSet_3",
				ProcessType:                 "mongod",
				Version:                     "4.2.2",
				AuthSchemaVersion:           5,
				FeatureCompatibilityVersion: "4.2",
				Disabled:                    false,
				ManualMode:                  false,
				Hostname:                    "host0",
				Args26: Args26{
					NET: Net{
						Port: 27020,
					},
					Storage: &Storage{
						DBPath: "/data/rs3",
					},
					SystemLog: SystemLog{
						Destination: "file",
						Path:        "/data/rs3/mongodb.log",
					},
					Replication: &Replication{
						ReplSetName: "myReplicaSet",
					},
				},
				LogRotate: &LogRotate{
					SizeThresholdMB:  1000.0,
					TimeThresholdHrs: 24,
				},
				LastGoalVersionAchieved: 0,
				Cluster:                 "",
			},
		},
		ReplicaSets: []*ReplicaSet{
			{
				ID:              "myReplicaSet",
				ProtocolVersion: "1",
				Members: []Member{
					{
						ID:           0,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_1",
						Priority:     1,
						SlaveDelay:   0,
						Votes:        1,
					},
					{
						ID:           1,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_2",
						Priority:     1,
						SlaveDelay:   0,
						Votes:        1,
					},
					{
						ID:           2,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_3",
						Priority:     1,
						SlaveDelay:   0,
						Votes:        1,
					},
				},
			},
		},
		Version: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAutomation_UpdateConfig(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"
	clusterName := "myReplicaSet"
	updateRequest := &AutomationConfig{
		Auth: Auth{
			AutoAuthMechanism: "MONGODB-CR",
			Disabled:          true,
			AuthoritativeSet:  false,
		},
		Processes: []*Process{
			{
				Name:                        "myReplicaSet_1",
				ProcessType:                 "mongod",
				Version:                     "4.2.2",
				AuthSchemaVersion:           5,
				FeatureCompatibilityVersion: "4.2",
				Disabled:                    false,
				ManualMode:                  false,
				Hostname:                    "host0",
				Args26: Args26{
					NET: Net{
						Port: 27000,
					},
					Storage: &Storage{
						DBPath: "/data/rs1",
					},
					SystemLog: SystemLog{
						Destination: "file",
						Path:        "/data/rs1/mongodb.log",
					},
					Replication: &Replication{
						ReplSetName: "myReplicaSet",
					},
				},
				LogRotate: &LogRotate{
					SizeThresholdMB:  1000.0,
					TimeThresholdHrs: 24,
				},
				LastGoalVersionAchieved: 0,
				Cluster:                 "",
			},
			{
				Name:                        "myReplicaSet_2",
				ProcessType:                 "mongod",
				Version:                     "4.2.2",
				AuthSchemaVersion:           5,
				FeatureCompatibilityVersion: "4.2",
				Disabled:                    false,
				ManualMode:                  false,
				Hostname:                    "host1",
				Args26: Args26{
					NET: Net{
						Port: 27010,
					},
					Storage: &Storage{
						DBPath: "/data/rs2",
					},
					SystemLog: SystemLog{
						Destination: "file",
						Path:        "/data/rs2/mongodb.log",
					},
					Replication: &Replication{
						ReplSetName: clusterName,
					},
				},
				LogRotate: &LogRotate{
					SizeThresholdMB:  1000.0,
					TimeThresholdHrs: 24,
				},
				LastGoalVersionAchieved: 0,
				Cluster:                 "",
			},
			{
				Name:                        "myReplicaSet_3",
				ProcessType:                 "mongod",
				Version:                     "4.2.2",
				AuthSchemaVersion:           5,
				FeatureCompatibilityVersion: "4.2",
				Disabled:                    false,
				ManualMode:                  false,
				Hostname:                    "host0",
				Args26: Args26{
					NET: Net{
						Port: 27020,
					},
					Storage: &Storage{
						DBPath: "/data/rs3",
					},
					SystemLog: SystemLog{
						Destination: "file",
						Path:        "/data/rs3/mongodb.log",
					},
					Replication: &Replication{
						ReplSetName: clusterName,
					},
				},
				LogRotate: &LogRotate{
					SizeThresholdMB:  1000.0,
					TimeThresholdHrs: 24,
				},
				LastGoalVersionAchieved: 0,
			},
		},
		ReplicaSets: []*ReplicaSet{
			{
				ID:              clusterName,
				ProtocolVersion: "1",
				Members: []Member{
					{
						ID:           0,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_1",
						Priority:     1,
						SlaveDelay:   0,
						Votes:        1,
					},
					{
						ID:           1,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_2",
						Priority:     1,
						SlaveDelay:   0,
						Votes:        1,
					},
					{
						ID:           2,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_3",
						Priority:     1,
						SlaveDelay:   0,
						Votes:        1,
					},
				},
			},
		},
		Version: 1,
	}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		_, _ = fmt.Fprint(w, `{}`)
	})

	_, err := client.Automation.UpdateConfig(ctx, projectID, updateRequest)
	if err != nil {
		t.Fatalf("Automation.UpdateConfig returned error: %v", err)
	}
}

func TestAutomation_UpdateAgentVersion(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/automationConfig/updateAgentVersions", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w,
			`{
				  "automationAgentVersion": "10.2.7.5898",
                  "biConnectorVersion": "2.6.1"
				}`,
		)
	})

	agent, _, err := client.Automation.UpdateAgentVersion(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.UpdateAgentVersion returned error: %v", err)
	}

	expected := &AutomationConfigAgent{
		AutomationAgentVersion: "10.2.7.5898",
		BiConnectorVersion:     "2.6.1",
	}

	if diff := deep.Equal(agent, expected); diff != nil {
		t.Error(diff)
	}

}
