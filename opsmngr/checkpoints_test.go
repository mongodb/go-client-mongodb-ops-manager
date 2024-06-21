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

const clusterID = "6b8cd61180eef547110159d9" //nolint:gosec // not a credential

func TestCheckpoints_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/groups/%s/clusters/%s/checkpoints", projectID, clusterID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{
		  "links":[
			{
			  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/clusters/Cluster0/checkpoints?pageNum=1&itemsPerPage=100",
			  "rel":"self"
			}
		  ],
		  "results":[
			{
			  "clusterId":"6b8cd61180eef547110159d9",
			  "completed":"2018-02-08T23:20:25Z",
			  "groupId":"%[1]s",
			  "id":"5a7cdb3980eef53de5bffdcf",
			  "links":[
				{
				  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/clusters/Cluster0/checkpoints",
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
			  "groupId":"%[1]s",
			  "id":"5a7db53987d9d64fe298ff46",
			  "links":[
				{
				  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/clusters/Cluster0/checkpoints?pretty=true",
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
			projectID,
		)
	})

	snapshots, _, err := client.Checkpoints.List(ctx, projectID, clusterID, nil)
	if err != nil {
		t.Fatalf("Checkpoints.List returned error: %v", err)
	}

	expected := &Checkpoints{
		Results: []*Checkpoint{
			{
				ClusterID: clusterID,
				Completed: "2018-02-08T23:20:25Z",
				GroupID:   projectID,
				ID:        "5a7cdb3980eef53de5bffdcf",
				Links: []*Link{
					{
						Rel:  "self",
						Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/clusters/Cluster0/checkpoints",
					},
				},
				Parts: []*Part{
					{
						ReplicaSetName: "Cluster0-shard-1",
						TypeName:       "REPLICA_SET",
						CheckpointPart: CheckpointPart{
							ShardName:       "Cluster0-shard-1",
							TokenDiscovered: true,
							TokenTimestamp: SnapshotTimestamp{
								Date:      "2018-02-08T23:20:25Z",
								Increment: 1,
							}},
					},
					{
						ReplicaSetName: "Cluster0-shard-0",
						TypeName:       "REPLICA_SET",
						CheckpointPart: CheckpointPart{
							ShardName:       "Cluster0-shard-0",
							TokenDiscovered: true,
							TokenTimestamp: SnapshotTimestamp{
								Date:      "2018-02-08T23:20:25Z",
								Increment: 1,
							}},
					},
					{
						ReplicaSetName: "Cluster0-config-0",
						TypeName:       "CONFIG_SERVER_REPLICA_SET",
						CheckpointPart: CheckpointPart{
							TokenDiscovered: true,
							TokenTimestamp: SnapshotTimestamp{
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
				ClusterID: clusterID,
				Completed: "2018-02-09T14:50:33Z",
				GroupID:   projectID,
				ID:        "5a7db53987d9d64fe298ff46",
				Links: []*Link{
					{
						Rel:  "self",
						Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/clusters/Cluster0/checkpoints?pretty=true",
					},
				},
				Parts: []*Part{
					{
						ReplicaSetName: "Cluster0-shard-1",
						TypeName:       "REPLICA_SET",
						CheckpointPart: CheckpointPart{
							ShardName:       "Cluster0-shard-1",
							TokenDiscovered: true,
							TokenTimestamp: SnapshotTimestamp{
								Date:      "2018-02-09T14:50:33Z",
								Increment: 1,
							}},
					},
					{
						ReplicaSetName: "Cluster0-shard-0",
						TypeName:       "REPLICA_SET",
						CheckpointPart: CheckpointPart{
							ShardName:       "Cluster0-shard-0",
							TokenDiscovered: true,
							TokenTimestamp: SnapshotTimestamp{
								Date:      "2018-02-09T14:50:33Z",
								Increment: 2,
							}},
					},
					{
						ReplicaSetName: "Cluster0-config-0",
						TypeName:       "CONFIG_SERVER_REPLICA_SET",
						CheckpointPart: CheckpointPart{
							TokenDiscovered: true,
							TokenTimestamp: SnapshotTimestamp{
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
		Links: []*Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/clusters/Cluster0/checkpoints?pageNum=1&itemsPerPage=100",
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

	const checkpointID = "6b8cd61180eef547110159d9" //nolint:gosec // not a credential
	path := fmt.Sprintf("/api/public/v1.0/groups/%s/clusters/%s/checkpoints/%s", projectID, clusterID, checkpointID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{
		  "clusterId":"6b8cd61180eef547110159d9",
		  "completed":"2018-02-08T23:20:25Z",
		  "groupId":"%[1]s",
		  "id":"5a7cdb3980eef53de5bffdcf",
		  "links":[
			{
			  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/clusters/Cluster0/checkpoints",
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
		}`, projectID)
	})

	checkpoint, _, err := client.Checkpoints.Get(ctx, projectID, clusterID, checkpointID)
	if err != nil {
		t.Fatalf("Checkpoints.Get returned error: %v", err)
	}

	expected := &Checkpoint{
		ClusterID: clusterID,
		Completed: "2018-02-08T23:20:25Z",
		GroupID:   projectID,
		ID:        "5a7cdb3980eef53de5bffdcf",
		Links: []*Link{
			{
				Rel:  "self",
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/clusters/Cluster0/checkpoints",
			},
		},
		Parts: []*Part{

			{
				ReplicaSetName: "Cluster0-shard-1",
				TypeName:       "REPLICA_SET",
				CheckpointPart: CheckpointPart{
					ShardName:       "Cluster0-shard-1",
					TokenDiscovered: true,
					TokenTimestamp: SnapshotTimestamp{
						Date:      "2018-02-08T23:20:25Z",
						Increment: 1,
					},
				},
			},
			{
				ReplicaSetName: "Cluster0-shard-0",
				TypeName:       "REPLICA_SET",
				CheckpointPart: CheckpointPart{
					ShardName:       "Cluster0-shard-0",
					TokenDiscovered: true,
					TokenTimestamp: SnapshotTimestamp{
						Date:      "2018-02-08T23:20:25Z",
						Increment: 1,
					}},
			},
			{
				ReplicaSetName: "Cluster0-config-0",
				TypeName:       "CONFIG_SERVER_REPLICA_SET",
				CheckpointPart: CheckpointPart{
					TokenDiscovered: true,
					TokenTimestamp: SnapshotTimestamp{
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
