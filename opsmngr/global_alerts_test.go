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

func TestGlobalAlerts_List(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	mux.HandleFunc("/api/public/v1.0/globalAlerts", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{
			  "links" : [],
			  "results" : [ {
						"alertConfigId" : "573b7d12e4b0979a262467c1",
						"created" : "2016-10-18T08:08:08Z",
						"currentValue" : {
						  "number" : 143.4739833843463,
						  "units" : "RAW"
						},
						"eventTypeName" : "OUTSIDE_METRIC_THRESHOLD",
						"groupId" : "1",
						"hostId" : "63f42376fb735471fe40ec54a7",
						"hostnameAndPort" : "replicaset-shard-00-02:27017",
						"id" : "573b7d2de4b02fd2c93423a6",
						"links" : [],
						"metricName" : "OPCOUNTER_CMD",
						"lastNotified" : "2016-10-18T19:29:54Z",
						"replicaSetName" : "replicaSet-shard-0",
						"resolved" : "2016-10-18T21:30:04Z",
						"status" : "CLOSED",
						"tags" : [ ],
						"updated" : "2016-10-18T21:30:04Z"
					  }],
			  "totalCount" : 1
		}`)
	})

	opts := AlertsListOptions{
		Status: "CLOSED",
	}

	alerts, _, err := client.GlobalAlerts.List(ctx, &opts)
	if err != nil {
		t.Fatalf("client.GlobalAlerts returned error: %v", err)
	}

	currentValueNumber := 143.4739833843463
	expected := &GlobalAlerts{
		Links: []*Link{},
		Results: []*GlobalAlert{
			{
				Alert: Alert{
					ID:              "573b7d2de4b02fd2c93423a6",
					GroupID:         "1",
					AlertConfigID:   "573b7d12e4b0979a262467c1",
					EventTypeName:   "OUTSIDE_METRIC_THRESHOLD",
					Created:         "2016-10-18T08:08:08Z",
					Updated:         "2016-10-18T21:30:04Z",
					Resolved:        "2016-10-18T21:30:04Z",
					Status:          "CLOSED",
					LastNotified:    "2016-10-18T19:29:54Z",
					HostnameAndPort: "replicaset-shard-00-02:27017",
					MetricName:      "OPCOUNTER_CMD",
					CurrentValue: &CurrentValue{
						Number: &currentValueNumber,
						Units:  "RAW",
					},
					ReplicaSetName: "replicaSet-shard-0",
				},
				Tags:   []string{},
				Links:  []*Link{},
				HostID: "63f42376fb735471fe40ec54a7",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}

func TestGlobalAlerts_Get(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/globalAlerts/%s", alertID)

	mux.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{
			    	  "alertConfigId" : "5730f5e1e4b030a9634a3f69",
					  "clusterId" : "572a00f2e4b051814b144e90",
					  "clusterName" : "shardedCluster",
					  "created" : "2016-10-09T06:16:36Z",
					  "eventTypeName" : "OPLOG_BEHIND",
					  "groupId" : "1",
					  "id" : "57b76ddc96e8215c017ceafb",
					  "links" : [],
					  "lastNotified" : "2016-10-10T20:42:32Z",
					  "replicaSetName" : "shardedCluster-shard-0",
					  "sourceTypeName" : "REPLICA_SET",
					  "status" : "OPEN",
					  "tags" : [ ],
					  "updated" : "2016-10-10T20:42:32Z"
		}`)
	})

	alerts, _, err := client.GlobalAlerts.Get(ctx, alertID)
	if err != nil {
		t.Fatalf("client.GlobalAlerts returned error: %v", err)
	}

	expected := &GlobalAlert{
		Alert: Alert{
			ID:              alertID,
			GroupID:         "1",
			AlertConfigID:   "5730f5e1e4b030a9634a3f69",
			EventTypeName:   "OPLOG_BEHIND",
			Created:         "2016-10-09T06:16:36Z",
			Updated:         "2016-10-10T20:42:32Z",
			Status:          "OPEN",
			LastNotified:    "2016-10-10T20:42:32Z",
			ReplicaSetName:  "shardedCluster-shard-0",
			ClusterName:     "shardedCluster",
			Matchers:        nil,
			MetricThreshold: nil,
			Notifications:   nil,
		},
		Tags:           []string{},
		Links:          []*Link{},
		SourceTypeName: "REPLICA_SET",
		ClusterID:      "572a00f2e4b051814b144e90",
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}
func TestGlobalAlerts_Acknowledge(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/globalAlerts/%s", alertID)

	mux.HandleFunc(path, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprint(w, `{
					  "alertConfigId" : "5730f5e1e4b030a9634a3f69",
					  "clusterId" : "572a00f2e4b051814b144e90",
					  "clusterName" : "shardedCluster",
					  "created" : "2016-10-09T06:16:36Z",
					  "eventTypeName" : "OPLOG_BEHIND",
					  "groupId" : "1",
					  "id" : "57b76ddc96e8215c017ceafb",
					  "links" : [],
					  "lastNotified" : "2016-10-10T20:42:32Z",
					  "replicaSetName" : "shardedCluster-shard-0",
					  "sourceTypeName" : "REPLICA_SET",
					  "status" : "OPEN",
					  "acknowledgedUntil" : "2016-11-01T00:00:00Z",
					  "acknowledgingUsername" : "admin@example.com",
					  "tags" : [ ],
					  "updated" : "2016-10-10T22:03:11Z"
		}`)
	})

	until := "2016-11-01T00:00:00-0400"
	body := &AcknowledgeRequest{
		AcknowledgedUntil: &until,
	}

	alerts, _, err := client.GlobalAlerts.Acknowledge(ctx, alertID, body)
	if err != nil {
		t.Fatalf("client.GlobalAlerts returned error: %v", err)
	}

	expected := &GlobalAlert{
		Alert: Alert{
			ID:                    alertID,
			GroupID:               "1",
			AlertConfigID:         "5730f5e1e4b030a9634a3f69",
			EventTypeName:         "OPLOG_BEHIND",
			Created:               "2016-10-09T06:16:36Z",
			Updated:               "2016-10-10T22:03:11Z",
			Status:                "OPEN",
			LastNotified:          "2016-10-10T20:42:32Z",
			ReplicaSetName:        "shardedCluster-shard-0",
			ClusterName:           "shardedCluster",
			AcknowledgedUntil:     "2016-11-01T00:00:00Z",
			AcknowledgingUsername: "admin@example.com",
		},
		Tags:           []string{},
		Links:          []*Link{},
		SourceTypeName: "REPLICA_SET",
		ClusterID:      "572a00f2e4b051814b144e90",
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}
