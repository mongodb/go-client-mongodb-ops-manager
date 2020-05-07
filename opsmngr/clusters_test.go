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
	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestClusters_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6b8cd3c380eef5349ef77gf7"

	path := fmt.Sprintf("/groups/%s/clusters", groupID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "totalCount": 1,
		  "results": [ {
			"id" : "533d7d4730040be257defe88",
			"typeName": "SHARDED_REPLICA_SET",
			"clusterName": "Animals",
			"groupId": "6b8cd3c380eef5349ef77gf7",
			"lastHeartbeat": "2014-04-03T15:26:58Z",
			"links": []
		  } ],
		  "links" : []
		}`)
	})

	clusters, _, err := client.Clusters.List(ctx, groupID, nil)
	if err != nil {
		t.Fatalf("Clusters.List returned error: %v", err)
	}

	expected := &Clusters{
		Results: []*Cluster{
			{
				ID:            "533d7d4730040be257defe88",
				GroupID:       "6b8cd3c380eef5349ef77gf7",
				TypeName:      "SHARDED_REPLICA_SET",
				ClusterName:   "Animals",
				LastHeartbeat: "2014-04-03T15:26:58Z",
				Links:         []*atlas.Link{},
			},
		},
		Links:      []*atlas.Link{},
		TotalCount: 1,
	}

	if diff := deep.Equal(clusters, expected); diff != nil {
		t.Error(diff)
	}
}

func TestClusters_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6b8cd3c380eef5349ef77gf7"
	clusterID := "533d7d4730040be257defe88"
	path := fmt.Sprintf("/groups/%s/clusters/%s", groupID, clusterID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"id": "533d7d4730040be257defe88",
			"typeName": "SHARDED_REPLICA_SET",
			"clusterName": "Animals",
			"groupId": "6b8cd3c380eef5349ef77gf7",
			"lastHeartbeat": "2014-04-03T15:26:58Z",
			"links" : []
		  }`)
	})

	cluster, _, err := client.Clusters.Get(ctx, groupID, clusterID)
	if err != nil {
		t.Fatalf("Clusters.Get returned error: %v", err)
	}

	expected := &Cluster{
		ID:            "533d7d4730040be257defe88",
		GroupID:       "6b8cd3c380eef5349ef77gf7",
		TypeName:      "SHARDED_REPLICA_SET",
		ClusterName:   "Animals",
		LastHeartbeat: "2014-04-03T15:26:58Z",
		Links:         []*atlas.Link{},
	}

	if diff := deep.Equal(cluster, expected); diff != nil {
		t.Error(diff)
	}
}

func TestClusters_ListAll(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	mux.HandleFunc("/clusters", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
				 "results": [
				{
				  "groupName": "AtlasGroup1",
				  "orgName": "TestAtlasOrg1",
				  "planType": "Atlas",
				  "groupId": "5e5fbc29e76c9a4be2ed3d39",
				  "clusters": [
					{
					  "backupEnabled": true,
					  "authEnabled": true,
					  "alertCount": 0,
					  "versions": [
						"3.4.2"
					  ],
					  "sslEnabled": true,
					  "name": "AtlasCluster1",
					  "dataSizeBytes": 1000000,
					  "nodeCount": 7,
					  "availability": "unavailable",
					  "clusterId": "5e5fbc29e76c9a4be2ed3d4d",
					  "type": "sharded cluster"
					},
					{
					  "backupEnabled": false,
					  "authEnabled": true,
					  "alertCount": 0,
					  "versions": [
						"3.4.1"
					  ],
					  "sslEnabled": true,
					  "name": "AtlasReplSet1",
					  "dataSizeBytes": 1300000,
					  "nodeCount": 2,
					  "availability": "dead",
					  "clusterId": "5e5fbc29e76c9a4be2ed3d4f",
					  "type": "replica set"
					}
				  ],
				  "orgId": "5e5fbc29e76c9a4be2ed3d36"
				},
				{
				  "groupName": "CloudGroup1",
				  "orgName": "TestCloudOrg1",
				  "planType": "Cloud Manager",
				  "groupId": "5e5fbc29e76c9a4be2ed3d38",
				  "clusters": [
					{
					  "backupEnabled": true,
					  "authEnabled": false,
					  "alertCount": 0,
					  "versions": [
						"3.4.1",
						"2.4.3"
					  ],
					  "sslEnabled": false,
					  "name": "cluster1",
					  "dataSizeBytes": 1000000,
					  "nodeCount": 6,
					  "availability": "warning",
					  "clusterId": "5e5fbc29e76c9a4be2ed3d42",
					  "type": "sharded cluster"
					},
					{
					  "backupEnabled": true,
					  "authEnabled": true,
					  "alertCount": 0,
					  "versions": [
						"3.4.1"
					  ],
					  "sslEnabled": true,
					  "name": "replica_set",
					  "dataSizeBytes": 500000,
					  "nodeCount": 2,
					  "availability": "available",
					  "clusterId": "5e5fbc29e76c9a4be2ed3d3c",
					  "type": "replica set"
					},
					{
					  "backupEnabled": false,
					  "authEnabled": false,
					  "alertCount": 0,
					  "versions": [
						"2.4.3"
					  ],
					  "sslEnabled": true,
					  "name": "standalone:27017",
					  "dataSizeBytes": 2000000,
					  "nodeCount": 1,
					  "availability": "unavailable",
					  "clusterId": "da303f3fec69b2100bacf10dd9e6d5e0",
					  "type": "standalone"
					}
				  ],
				  "orgId": "5e5fbc29e76c9a4be2ed3d34",
				  "tags": [
					"some tag 1",
					"some tag 2"
				  ]
				}
			  ],
			  "links": [],
			  "totalCount": 2
		}`)
	})

	clusters, _, err := client.Clusters.ListAll(ctx)
	if err != nil {
		t.Fatalf("Clusters.ListAll returned error: %v", err)
	}

	expected := &AllClustersProjects{
		Links: []*atlas.Link{},
		Results: []*AllClustersProject{
			{
				GroupName: "AtlasGroup1",
				OrgName:   "TestAtlasOrg1",
				PlanType:  "Atlas",
				GroupID:   "5e5fbc29e76c9a4be2ed3d39",
				OrgID:     "5e5fbc29e76c9a4be2ed3d36",
				Clusters: []AllClustersCluster{
					{
						ClusterID:     "5e5fbc29e76c9a4be2ed3d4d",
						Name:          "AtlasCluster1",
						Type:          "sharded cluster",
						Availability:  "unavailable",
						Versions:      []string{"3.4.2"},
						BackupEnabled: true,
						AuthEnabled:   true,
						SSLEnabled:    true,
						AlertCount:    0,
						DataSizeBytes: 1000000,
						NodeCount:     7,
					},
					{
						ClusterID:     "5e5fbc29e76c9a4be2ed3d4f",
						Name:          "AtlasReplSet1",
						Type:          "replica set",
						Availability:  "dead",
						Versions:      []string{"3.4.1"},
						BackupEnabled: false,
						AuthEnabled:   true,
						SSLEnabled:    true,
						AlertCount:    0,
						DataSizeBytes: 1300000,
						NodeCount:     2,
					},
				},
			},
			{
				GroupName: "CloudGroup1",
				OrgName:   "TestCloudOrg1",
				PlanType:  "Cloud Manager",
				GroupID:   "5e5fbc29e76c9a4be2ed3d38",
				OrgID:     "5e5fbc29e76c9a4be2ed3d34",
				Tags:      []string{"some tag 1", "some tag 2"},
				Clusters: []AllClustersCluster{
					{
						ClusterID:     "5e5fbc29e76c9a4be2ed3d42",
						Name:          "cluster1",
						Type:          "sharded cluster",
						Availability:  "warning",
						Versions:      []string{"3.4.1", "2.4.3"},
						BackupEnabled: true,
						AuthEnabled:   false,
						SSLEnabled:    false,
						AlertCount:    0,
						DataSizeBytes: 1000000,
						NodeCount:     6,
					},
					{
						ClusterID:     "5e5fbc29e76c9a4be2ed3d3c",
						Name:          "replica_set",
						Type:          "replica set",
						Availability:  "available",
						Versions:      []string{"3.4.1"},
						BackupEnabled: true,
						AuthEnabled:   true,
						SSLEnabled:    true,
						AlertCount:    0,
						DataSizeBytes: 500000,
						NodeCount:     2,
					},
					{
						ClusterID:     "da303f3fec69b2100bacf10dd9e6d5e0",
						Name:          "standalone:27017",
						Type:          "standalone",
						Availability:  "unavailable",
						Versions:      []string{"2.4.3"},
						BackupEnabled: false,
						AuthEnabled:   false,
						SSLEnabled:    true,
						AlertCount:    0,
						DataSizeBytes: 2000000,
						NodeCount:     1,
					},
				},
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(clusters, expected); diff != nil {
		t.Error(diff)
	}
}
