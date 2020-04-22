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
	client, mux, _, teardown := setup()
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
	client, mux, _, teardown := setup()
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
