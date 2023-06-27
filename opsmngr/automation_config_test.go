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
      "setParameter": {
        "opensslCipherConfig": "HIGH:EXPORT:!aNULL@STRENGTH;!DES:!3DES"
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
        "port" : 27010,
		"tls" : {
			"FIPSMode": true
		}
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
    "numCores": 2,
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
    "settings" : { },
	"force": { 
		"currentVersion": -1
	}
  } ],
  "version" : 1
}`

func TestAutomation_GetConfig(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, jsonBlob)
	})

	config, _, err := client.Automation.GetConfig(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.GetConfig returned error: %v", err)
	}

	slaveDelay := float64(0)
	fipsMode := true
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
					SetParameter: &map[string]interface{}{
						"opensslCipherConfig": "HIGH:EXPORT:!aNULL@STRENGTH;!DES:!3DES",
					},
					Storage: &Storage{
						DBPath: "/data/rs1",
						WiredTiger: &map[string]interface{}{
							"collectionConfig": map[string]interface{}{},
							"engineConfig": map[string]interface{}{
								"cacheSizeGB": 0.5,
							},
							"indexConfig": map[string]interface{}{},
						},
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
				NumCores:                0,
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
						TLS:  &TLS{FIPSMode: &fipsMode},
					},
					Storage: &Storage{
						DBPath: "/data/rs2",
						WiredTiger: &map[string]interface{}{
							"collectionConfig": map[string]interface{}{},
							"engineConfig": map[string]interface{}{
								"cacheSizeGB": 0.5,
							},
							"indexConfig": map[string]interface{}{},
						},
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
				NumCores:                0,
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
						WiredTiger: &map[string]interface{}{
							"collectionConfig": map[string]interface{}{},
							"engineConfig": map[string]interface{}{
								"cacheSizeGB": 0.5,
							},
							"indexConfig": map[string]interface{}{},
						},
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
				NumCores:                2,
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
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
					{
						ID:           1,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_2",
						Priority:     1,
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
					{
						ID:           2,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_3",
						Priority:     1,
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
				},
				Settings: &map[string]interface{}{},
				Force:    &Force{CurrentVersion: -1},
			},
		},
		Version: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAutomation_UpdateConfig(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	clusterName := "myReplicaSet"
	slaveDelay := float64(0)
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
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
					{
						ID:           1,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_2",
						Priority:     1,
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
					{
						ID:           2,
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         "myReplicaSet_3",
						Priority:     1,
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
				},
			},
		},
		Version: 1,
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
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

func TestAutomation_MongoDBUserNoMechanism(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
  "auth" : {
    "authoritativeSet" : false,
    "autoAuthMechanism" : "MONGODB-CR",
    "disabled" : true,
	"usersWanted": [
		{
			"db": "admin"
		}
	]
  }}`)
	})

	config, _, err := client.Automation.GetConfig(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.GetConfig returned error: %v", err)
	}

	expected := &AutomationConfig{
		Auth: Auth{
			AuthoritativeSet:  false,
			AutoAuthMechanism: "MONGODB-CR",
			Disabled:          true,
			UsersWanted: []*MongoDBUser{
				{
					Database: "admin",
				},
			},
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAutomation_MongoDBUserEmptyMechanism(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
  "auth" : {
    "authoritativeSet" : false,
    "autoAuthMechanism" : "MONGODB-CR",
    "disabled" : true,
	"usersWanted": [
		{
			"db": "admin",
			"mechanisms": []	
		}
	]
  }}`)
	})

	config, _, err := client.Automation.GetConfig(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.GetConfig returned error: %v", err)
	}

	expected := &AutomationConfig{
		Auth: Auth{
			AuthoritativeSet:  false,
			AutoAuthMechanism: "MONGODB-CR",
			Disabled:          true,
			UsersWanted: []*MongoDBUser{
				{
					Database:   "admin",
					Mechanisms: &[]string{},
				},
			},
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAutomation_UpdateMongoDBUserEmptyMechanism(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	updateRequest := &AutomationConfig{
		Auth: Auth{
			AuthoritativeSet:  false,
			AutoAuthMechanism: "MONGODB-CR",
			Disabled:          true,
			UsersWanted: []*MongoDBUser{
				{
					Database:   "admin",
					Mechanisms: &[]string{},
				},
				{
					Database: "admin",
				},
			},
		},
	}
	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"auth": map[string]interface{}{
				"authoritativeSet":     false,
				"autoAuthMechanism":    "MONGODB-CR",
				"autoAuthRestrictions": interface{}(nil),
				"disabled":             true,
				"usersDeleted":         interface{}(nil),
				"usersWanted": []interface{}{
					map[string]interface{}{
						"authenticationRestrictions": interface{}(nil),
						"db":                         "admin",
						"roles":                      interface{}(nil),
						"user":                       "",
						"mechanisms":                 []interface{}{},
					},
					map[string]interface{}{
						"authenticationRestrictions": interface{}(nil),
						"db":                         "admin",
						"roles":                      interface{}(nil),
						"user":                       "",
					},
				},
			},
			"backupVersions":       interface{}(nil),
			"balancer":             interface{}(nil),
			"cpsModules":           interface{}(nil),
			"indexConfigs":         interface{}(nil),
			"mongosqlds":           interface{}(nil),
			"mongots":              interface{}(nil),
			"onlineArchiveModules": interface{}(nil),
			"options":              interface{}(nil),
			"processes":            interface{}(nil),
			"replicaSets":          interface{}(nil),
			"roles":                interface{}(nil),
			"sharding":             interface{}(nil),
		}
		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}
		t.Logf("v=%#v\n", v)
		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}
		_, _ = fmt.Fprint(w, `{}`)
	})

	_, err := client.Automation.UpdateConfig(ctx, projectID, updateRequest)
	if err != nil {
		t.Fatalf("Automation.UpdateConfig returned error: %v", err)
	}
}

func TestAutomation_GetTlsConfig(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
  "tls" : {
    	"CAFilePath" : "/etc/ssl/rootCA.pem",
		"autoPEMKeyFilePath" : "/etc/ssl/client.pem",
		"autoPEMKeyFilePwd" : "password",
		"clientCertificateMode" : "OPTIONAL"
  }}`)
	})

	config, _, err := client.Automation.GetConfig(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.GetConfig returned error: %v", err)
	}

	expected := &AutomationConfig{
		TLS: &SSL{
			CAFilePath:            "/etc/ssl/rootCA.pem",
			AutoPEMKeyFilePath:    "/etc/ssl/client.pem",
			AutoPEMKeyFilePwd:     "password",
			ClientCertificateMode: "OPTIONAL",
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAutomation_Sharding(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
  "sharding": [
    {
      "collections": [
        {
          "_id": "EncounterDB.ipsum",
          "dropped": false,
          "key": [
            [
              "subjectReference.resourceId",
              1
            ]
          ],
          "presplitHashedZones": false,
          "unique": false
        }
      ],
      "configServerReplica": "lorem",
      "draining": [],
      "managedSharding": true,
      "name": "lorem",
      "shards": [
        {
          "_id": "lorem01_0",
          "rs": "lorem01_0",
          "tags": [
            "RECENT"
          ]
        }
      ],
      "tags": [
        {
          "max": [
            {
              "field": "encounterPeriod.start",
              "fieldType": "date",
              "value": "1593561600000"
            }
          ],
          "min": [
            {
              "field": "encounterPeriod.start",
              "fieldType": "minKey",
              "value": "1"
            }
          ],
          "ns": "EncounterDB.ipsum2",
          "tag": "PREVIOUS"
        }
      ]
    }
  ]}`)
	})

	config, _, err := client.Automation.GetConfig(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.GetConfig returned error: %v", err)
	}

	expected := &AutomationConfig{
		Sharding: []*ShardingConfig{
			{
				Collections: []*map[string]interface{}{
					{
						"_id":     "EncounterDB.ipsum",
						"dropped": false,
						"key": []interface{}{
							[]interface{}{
								"subjectReference.resourceId",
								float64(1),
							},
						},
						"presplitHashedZones": false,
						"unique":              false,
					},
				},
				ConfigServerReplica: "lorem",
				Draining:            []string{},
				ManagedSharding:     pointer(true),
				Name:                "lorem",
				Shards: []*Shard{
					{
						ID:   "lorem01_0",
						RS:   "lorem01_0",
						Tags: []string{"RECENT"},
					},
				},
				Tags: []*map[string]interface{}{
					{
						"max": []interface{}{
							map[string]interface{}{
								"field":     "encounterPeriod.start",
								"fieldType": "date",
								"value":     "1593561600000",
							},
						},
						"min": []interface{}{
							map[string]interface{}{
								"field":     "encounterPeriod.start",
								"fieldType": "minKey",
								"value":     "1",
							},
						},
						"ns":  "EncounterDB.ipsum2",
						"tag": "PREVIOUS",
					},
				},
			},
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}
