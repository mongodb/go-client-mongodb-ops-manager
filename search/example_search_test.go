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

package search_test

import (
	"fmt"

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

const rsName = "myReplicaSet_1"

var slaveDelay = float64(0)

var fixture = &opsmngr.AutomationConfig{
	Auth: opsmngr.Auth{
		AutoAuthMechanism: "MONGODB-CR",
		Disabled:          true,
		AuthoritativeSet:  false,
		UsersWanted: []*opsmngr.MongoDBUser{
			{
				Mechanisms: &[]string{"SCRAM-SHA-1"},
				Roles: []*opsmngr.Role{
					{
						Role:     "test",
						Database: "test",
					},
				},
				Username: "test",
				Database: "test",
			},
		},
	},
	Processes: []*opsmngr.Process{
		{
			Name:                        rsName,
			ProcessType:                 "mongod",
			Version:                     "4.2.2",
			AuthSchemaVersion:           5,
			FeatureCompatibilityVersion: "4.2",
			Disabled:                    false,
			ManualMode:                  false,
			Hostname:                    "host0",
			Args26: opsmngr.Args26{
				NET: opsmngr.Net{
					Port: 27000,
				},
				Storage: &opsmngr.Storage{
					DBPath: "/data/rs1",
				},
				SystemLog: opsmngr.SystemLog{
					Destination: "file",
					Path:        "/data/rs1/mongodb.log",
				},
				Replication: &opsmngr.Replication{
					ReplSetName: "myReplicaSet",
				},
			},
			LogRotate: &opsmngr.LogRotate{
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
			Args26: opsmngr.Args26{
				NET: opsmngr.Net{
					Port: 27010,
				},
				Storage: &opsmngr.Storage{
					DBPath: "/data/rs2",
				},
				SystemLog: opsmngr.SystemLog{
					Destination: "file",
					Path:        "/data/rs2/mongodb.log",
				},
				Replication: &opsmngr.Replication{
					ReplSetName: "myReplicaSet",
				},
			},
			LogRotate: &opsmngr.LogRotate{
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
			Args26: opsmngr.Args26{
				NET: opsmngr.Net{
					Port: 27020,
				},
				Storage: &opsmngr.Storage{
					DBPath: "/data/rs3",
				},
				SystemLog: opsmngr.SystemLog{
					Destination: "file",
					Path:        "/data/rs3/mongodb.log",
				},
				Replication: &opsmngr.Replication{
					ReplSetName: "myReplicaSet",
				},
			},
			LogRotate: &opsmngr.LogRotate{
				SizeThresholdMB:  1000.0,
				TimeThresholdHrs: 24,
			},
			LastGoalVersionAchieved: 0,
			Cluster:                 "",
		},
	},
	ReplicaSets: []*opsmngr.ReplicaSet{
		{
			ID:              "myReplicaSet",
			ProtocolVersion: "1",
			Members: []opsmngr.Member{
				{
					ID:           0,
					ArbiterOnly:  false,
					BuildIndexes: true,
					Hidden:       false,
					Host:         rsName,
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
	IndexConfigs: []*opsmngr.IndexConfig{
		{
			DBName:         "test",
			CollectionName: "test",
			RSName:         rsName,
			Key: [][]string{
				{"test", "1"},
			},
			Options:   nil,
			Collation: nil,
		},
	},
	Version: 1,
}

// This example demonstrates searching a list of processes by name.
func ExampleProcesses() {
	a := fixture.Processes
	const x = "myReplicaSet_2"
	i, found := search.Processes(a, func(p *opsmngr.Process) bool { return p.Name == x })
	if i < len(a) && found {
		fmt.Printf("found %v at index %d\n", x, i)
	} else {
		fmt.Printf("%s not found\n", x)
	}
	// Output:
	// found myReplicaSet_2 at index 1
}

// This example demonstrates searching a list of replica sets by ID.
func ExampleReplicaSets() {
	a := fixture.ReplicaSets
	const x = "myReplicaSet"
	i, found := search.ReplicaSets(a, func(r *opsmngr.ReplicaSet) bool { return r.ID == x })
	if i < len(a) && found {
		fmt.Printf("found %v at index %d\n", x, i)
	} else {
		fmt.Printf("%s not found\n", x)
	}
	// Output:
	// found myReplicaSet at index 0
}

// This example demonstrates searching a list of replica sets by ID.
func ExampleShardingConfig() {
	a := []*opsmngr.ShardingConfig{{Name: "myCluster"}}
	const x = "myCluster"
	i, found := search.ShardingConfig(a, func(r *opsmngr.ShardingConfig) bool { return r.Name == x })
	if i < len(a) && found {
		fmt.Printf("found %v at index %d\n", x, i)
	} else {
		fmt.Printf("%s not found\n", x)
	}
	// Output:
	// found myCluster at index 0
}

// This example demonstrates searching a list of members by host.
func ExampleMembers() {
	a := fixture.ReplicaSets[0].Members
	const x = "myReplicaSet_2"
	i, found := search.Members(a, func(m opsmngr.Member) bool { return m.Host == x })
	if i < len(a) && found {
		fmt.Printf("found %v at index %d\n", x, i)
	} else {
		fmt.Printf("%s not found\n", x)
	}
	// Output:
	// found myReplicaSet_2 at index 1
}

// This example demonstrates searching a list of db users by username.
func ExampleMongoDBUsers() {
	a := fixture.Auth.UsersWanted
	const x = "test"
	i, found := search.MongoDBUsers(a, func(m *opsmngr.MongoDBUser) bool { return m.Username == x })
	if i < len(a) && found {
		fmt.Printf("found %v at index %d\n", x, i)
	} else {
		fmt.Printf("%s not found\n", x)
	}
	// Output:
	// found test at index 0
}

// This example demonstrates searching an index by RSName.
func ExampleMongoDBIndexes() {
	a := fixture
	const x = "myReplicaSet_1"
	i, found := search.MongoDBIndexes(a.IndexConfigs, func(r *opsmngr.IndexConfig) bool {
		return r.RSName == x
	})

	if i < len(a.IndexConfigs) && found {
		fmt.Printf("found %v at index %d\n", x, i)
	} else {
		fmt.Printf("%s not found\n", x)
	}
	// Output:
	// found myReplicaSet_1 at index 0
}
