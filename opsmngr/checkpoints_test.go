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

func TestCheckpoints_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6b8cd3c380eef5349ef77gf7"
	clusterID := "6b8cd61180eef547110159d9"

	path := fmt.Sprintf("/groups/%s/clusters/%s/checkpoints", groupID, clusterID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
		  "links":[
			{
			  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/6b8cd3c380eef5349ef77gf7/clusters/Cluster0/checkpoints?pageNum=1&itemsPerPage=100",
			  "rel":"self"
			}
		  ],
		  "results":[
			{
			  "clusterId":"6b8cd61180eef547110159d9",
			  "completed":"2018-02-08T23:20:25Z",
			  "groupId":"6b8cd3c380eef5349ef77gf7",
			  "id":"5a7cdb3980eef53de5bffdcf",
			  "links":[
				{
				  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/6b8cd3c380eef5349ef77gf7/clusters/Cluster0/checkpoints",
				  "rel":"self"
				}
			  ],
			  "parts":[
				{
				  "replicaSetName":"Cluster0-shard-1",
				  "shardName":"Cluster0-shard-1",
				  "tokenDiscovered":true,
				  "tokenTimestamp":{
					"date":"2018-02-08T23:20:25Z",
					"increment":1
				  },
				  "typeName":"REPLICA_SET"
				},
				{
				  "replicaSetName":"Cluster0-shard-0",
				  "shardName":"Cluster0-shard-0",
				  "tokenDiscovered":true,
				  "tokenTimestamp":{
					"date":"2018-02-08T23:20:25Z",
					"increment":1
				  },
				  "typeName":"REPLICA_SET"
				},
				{
				  "replicaSetName":"Cluster0-config-0",
				  "tokenDiscovered":true,
				  "tokenTimestamp":{
					"date":"2018-02-08T23:20:25Z",
					"increment":2
				  },
				  "typeName":"CONFIG_SERVER_REPLICA_SET"
				}
			  ],
			  "restorable":true,
			  "started":"2018-02-08T23:20:25Z",
			  "timestamp":"2018-02-08T23:19:37Z"
			},
			{
			  "clusterId":"6b8cd61180eef547110159d9",
			  "completed":"2018-02-09T14:50:33Z",
			  "groupId":"6b8cd3c380eef5349ef77gf7",
			  "id":"5a7db53987d9d64fe298ff46",
			  "links":[
				{
				  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/6b8cd3c380eef5349ef77gf7/clusters/Cluster0/checkpoints?pretty=true",
				  "rel":"self"
				}
			  ],
			  "parts":[
				{
				  "replicaSetName":"Cluster0-shard-1",
				  "shardName":"Cluster0-shard-1",
				  "tokenDiscovered":true,
				  "tokenTimestamp":{
					"date":"2018-02-09T14:50:33Z",
					"increment":1
				  },
				  "typeName":"REPLICA_SET"
				},
				{
				  "replicaSetName":"Cluster0-shard-0",
				  "shardName":"Cluster0-shard-0",
				  "tokenDiscovered":true,
				  "tokenTimestamp":{
					"date":"2018-02-09T14:50:33Z",
					"increment":2
				  },
				  "typeName":"REPLICA_SET"
				},
				{
				  "replicaSetName":"Cluster0-config-0",
				  "tokenDiscovered":true,
				  "tokenTimestamp":{
					"date":"2018-02-09T14:50:33Z",
					"increment":4
				  },
				  "typeName":"CONFIG_SERVER_REPLICA_SET"
				}
			  ],
			  "restorable":true,
			  "started":"2018-02-09T14:50:33Z",
			  "timestamp":"2018-02-09T14:50:18Z"
			}
		  ],
		  "totalCount":2
		}`,
		)
	})

	snapshots, _, err := client.Checkpoints.List(ctx, groupID, clusterID, nil)
	if err != nil {
		t.Fatalf("Checkpoints.List returned error: %v", err)
	}

	expected := &atlas.Checkpoints{
		Results: []*atlas.Checkpoint{
			{
				ClusterID: "6b8cd61180eef547110159d9",
				Completed: "2018-02-08T23:20:25Z",
				GroupID:   "6b8cd3c380eef5349ef77gf7",
				ID:        "5a7cdb3980eef53de5bffdcf",
				Links: []*atlas.Link{
					{
						Rel:  "self",
						Href: "https://cloud.mongodb.com/api/public/v1.0/groups/6b8cd3c380eef5349ef77gf7/clusters/Cluster0/checkpoints",
					},
				},
				Parts: []*atlas.Part{
					{
						ReplicaSetName: "Cluster0-shard-1",
						TypeName:       "REPLICA_SET",
						CheckpointPart: atlas.CheckpointPart{
							ShardName:       "Cluster0-shard-1",
							TokenDiscovered: true,
							TokenTimestamp: atlas.SnapshotTimestamp{
								Date:      "2018-02-08T23:20:25Z",
								Increment: 1,
							}},
					},
					{
						ReplicaSetName: "Cluster0-shard-0",
						TypeName:       "REPLICA_SET",
						CheckpointPart: atlas.CheckpointPart{
							ShardName:       "Cluster0-shard-0",
							TokenDiscovered: true,
							TokenTimestamp: atlas.SnapshotTimestamp{
								Date:      "2018-02-08T23:20:25Z",
								Increment: 1,
							}},
					},
					{
						ReplicaSetName: "Cluster0-config-0",
						TypeName:       "CONFIG_SERVER_REPLICA_SET",
						CheckpointPart: atlas.CheckpointPart{
							TokenDiscovered: true,
							TokenTimestamp: atlas.SnapshotTimestamp{
								Date:      "2018-02-08T23:20:25Z",
								Increment: 2,
							}},
					},
				},
				Restorable: true,
				Started:    "2018-02-08T23:20:25Z",
				Timestamp:  "2018-02-08T23:19:37Z",
			},
			{
				ClusterID: "6b8cd61180eef547110159d9",
				Completed: "2018-02-09T14:50:33Z",
				GroupID:   "6b8cd3c380eef5349ef77gf7",
				ID:        "5a7db53987d9d64fe298ff46",
				Links: []*atlas.Link{
					{
						Rel:  "self",
						Href: "https://cloud.mongodb.com/api/public/v1.0/groups/6b8cd3c380eef5349ef77gf7/clusters/Cluster0/checkpoints?pretty=true",
					},
				},
				Parts: []*atlas.Part{
					{
						ReplicaSetName: "Cluster0-shard-1",
						TypeName:       "REPLICA_SET",
						CheckpointPart: atlas.CheckpointPart{
							ShardName:       "Cluster0-shard-1",
							TokenDiscovered: true,
							TokenTimestamp: atlas.SnapshotTimestamp{
								Date:      "2018-02-09T14:50:33Z",
								Increment: 1,
							}},
					},
					{
						ReplicaSetName: "Cluster0-shard-0",
						TypeName:       "REPLICA_SET",
						CheckpointPart: atlas.CheckpointPart{
							ShardName:       "Cluster0-shard-0",
							TokenDiscovered: true,
							TokenTimestamp: atlas.SnapshotTimestamp{
								Date:      "2018-02-09T14:50:33Z",
								Increment: 2,
							}},
					},
					{
						ReplicaSetName: "Cluster0-config-0",
						TypeName:       "CONFIG_SERVER_REPLICA_SET",
						CheckpointPart: atlas.CheckpointPart{
							TokenDiscovered: true,
							TokenTimestamp: atlas.SnapshotTimestamp{
								Date:      "2018-02-09T14:50:33Z",
								Increment: 4,
							}},
					},
				},
				Restorable: true,
				Started:    "2018-02-09T14:50:33Z",
				Timestamp:  "2018-02-09T14:50:18Z",
			},
		},
		Links: []*atlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/6b8cd3c380eef5349ef77gf7/clusters/Cluster0/checkpoints?pageNum=1&itemsPerPage=100",
				Rel:  "self",
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(snapshots, expected); diff != nil {
		t.Error(diff)
	}
}

func TestCheckpoints_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groupID := "6b8cd3c380eef5349ef77gf7"
	clusterID := "6b8cd61180eef547110159d9"
	checkpointID := "6b8cd61180eef547110159d9"
	path := fmt.Sprintf("/groups/%s/clusters/%s/checkpoints/%s", groupID, clusterID, checkpointID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
		  "clusterId":"6b8cd61180eef547110159d9",
		  "completed":"2018-02-08T23:20:25Z",
		  "groupId":"6b8cd3c380eef5349ef77gf7",
		  "id":"5a7cdb3980eef53de5bffdcf",
		  "links":[
			{
			  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/6b8cd3c380eef5349ef77gf7/clusters/Cluster0/checkpoints",
			  "rel":"self"
			}
		  ],
		  "parts":[
			{
			  "replicaSetName":"Cluster0-shard-1",
			  "shardName":"Cluster0-shard-1",
			  "tokenDiscovered":true,
			  "tokenTimestamp":{
				"date":"2018-02-08T23:20:25Z",
				"increment":1
			  },
			  "typeName":"REPLICA_SET"
			},
			{
			  "replicaSetName":"Cluster0-shard-0",
			  "shardName":"Cluster0-shard-0",
			  "tokenDiscovered":true,
			  "tokenTimestamp":{
				"date":"2018-02-08T23:20:25Z",
				"increment":1
			  },
			  "typeName":"REPLICA_SET"
			},
			{
			  "replicaSetName":"Cluster0-config-0",
			  "tokenDiscovered":true,
			  "tokenTimestamp":{
				"date":"2018-02-08T23:20:25Z",
				"increment":2
			  },
			  "typeName":"CONFIG_SERVER_REPLICA_SET"
			}
		  ],
		  "restorable":true,
		  "started":"2018-02-08T23:20:25Z",
		  "timestamp":"2018-02-08T23:19:37Z"
		}`)
	})

	checkpoint, _, err := client.Checkpoints.Get(ctx, groupID, clusterID, checkpointID)
	if err != nil {
		t.Fatalf("Checkpoints.Get returned error: %v", err)
	}

	expected := &atlas.Checkpoint{
		ClusterID: "6b8cd61180eef547110159d9",
		Completed: "2018-02-08T23:20:25Z",
		GroupID:   "6b8cd3c380eef5349ef77gf7",
		ID:        "5a7cdb3980eef53de5bffdcf",
		Links: []*atlas.Link{
			{
				Rel:  "self",
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/6b8cd3c380eef5349ef77gf7/clusters/Cluster0/checkpoints",
			},
		},
		Parts: []*atlas.Part{

			{
				ReplicaSetName: "Cluster0-shard-1",
				TypeName:       "REPLICA_SET",
				CheckpointPart: atlas.CheckpointPart{
					ShardName:       "Cluster0-shard-1",
					TokenDiscovered: true,
					TokenTimestamp: atlas.SnapshotTimestamp{
						Date:      "2018-02-08T23:20:25Z",
						Increment: 1,
					},
				},
			},
			{
				ReplicaSetName: "Cluster0-shard-0",
				TypeName:       "REPLICA_SET",
				CheckpointPart: atlas.CheckpointPart{
					ShardName:       "Cluster0-shard-0",
					TokenDiscovered: true,
					TokenTimestamp: atlas.SnapshotTimestamp{
						Date:      "2018-02-08T23:20:25Z",
						Increment: 1,
					}},
			},
			{
				ReplicaSetName: "Cluster0-config-0",
				TypeName:       "CONFIG_SERVER_REPLICA_SET",
				CheckpointPart: atlas.CheckpointPart{
					TokenDiscovered: true,
					TokenTimestamp: atlas.SnapshotTimestamp{
						Date:      "2018-02-08T23:20:25Z",
						Increment: 2,
					}},
			},
		},
		Restorable: true,
		Started:    "2018-02-08T23:20:25Z",
		Timestamp:  "2018-02-08T23:19:37Z",
	}

	if diff := deep.Equal(checkpoint, expected); diff != nil {
		t.Error(diff)
	}
}
