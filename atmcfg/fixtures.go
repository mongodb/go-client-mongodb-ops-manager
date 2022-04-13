// Copyright 2021 MongoDB Inc
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

import "go.mongodb.org/ops-manager/opsmngr"

const (
	defaultMongoPort        = 27017
	defaultSizeThresholdMB  = 1000
	defaultTimeThresholdHrs = 24
	authSchemaVersion       = 5
)

func automationConfigWithOneReplicaSet(name string, disabled bool) *opsmngr.AutomationConfig {
	slaveDelay := float64(0)
	return &opsmngr.AutomationConfig{
		Processes: []*opsmngr.Process{
			{
				Args26: opsmngr.Args26{
					NET: opsmngr.Net{
						Port: defaultMongoPort,
					},
					Replication: &opsmngr.Replication{
						ReplSetName: name,
					},
					Sharding: nil,
					Storage: &opsmngr.Storage{
						DBPath: "/data/db/",
					},
					SystemLog: opsmngr.SystemLog{
						Destination: "file",
						Path:        "/data/db/mongodb.log",
					},
				},
				AuthSchemaVersion:           authSchemaVersion,
				Name:                        name + "_0",
				Disabled:                    disabled,
				ManualMode:                  disabled,
				FeatureCompatibilityVersion: "4.2",
				Hostname:                    "host0",
				LogRotate: &opsmngr.LogRotate{
					SizeThresholdMB:  defaultSizeThresholdMB,
					TimeThresholdHrs: defaultTimeThresholdHrs,
				},
				ProcessType: "mongod",
				Version:     "4.2.2",
			},
		},
		ReplicaSets: []*opsmngr.ReplicaSet{
			{
				ID:              name,
				ProtocolVersion: "1",
				Members: []opsmngr.Member{
					{
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         name + "_0",
						Priority:     1,
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
				},
			},
		},
	}
}

func automationConfigWithOneShardedCluster(name string, disabled bool) *opsmngr.AutomationConfig {
	configRSPort := defaultMongoPort + 1
	mongosPort := configRSPort + 1
	slaveDelay := float64(0)
	return &opsmngr.AutomationConfig{
		Processes: []*opsmngr.Process{
			{
				Args26: opsmngr.Args26{
					NET: opsmngr.Net{
						Port: defaultMongoPort,
					},
					Replication: &opsmngr.Replication{
						ReplSetName: name,
					},
					Sharding: nil,
					Storage: &opsmngr.Storage{
						DBPath: "/data/db/",
					},
					SystemLog: opsmngr.SystemLog{
						Destination: "file",
						Path:        "/data/db/mongodb.log",
					},
				},
				AuthSchemaVersion:           authSchemaVersion,
				Name:                        name + "_shard_0_0",
				Disabled:                    disabled,
				ManualMode:                  disabled,
				FeatureCompatibilityVersion: "4.2",
				Hostname:                    "host0",
				LogRotate: &opsmngr.LogRotate{
					SizeThresholdMB:  defaultSizeThresholdMB,
					TimeThresholdHrs: defaultTimeThresholdHrs,
				},
				ProcessType: "mongod",
				Version:     "4.2.2",
			},
			{
				Args26: opsmngr.Args26{
					NET: opsmngr.Net{
						Port: configRSPort,
					},
					Replication: &opsmngr.Replication{
						ReplSetName: name + "_configRS",
					},
					Sharding: &opsmngr.Sharding{
						ClusterRole: "configsvr",
					},
					Storage: &opsmngr.Storage{
						DBPath: "/data/db/",
					},
					SystemLog: opsmngr.SystemLog{
						Destination: "file",
						Path:        "/data/db/mongodb.log",
					},
				},
				AuthSchemaVersion:           authSchemaVersion,
				Name:                        name + "_configRS_0",
				Disabled:                    disabled,
				ManualMode:                  disabled,
				FeatureCompatibilityVersion: "4.2",
				Hostname:                    "host2",
				LogRotate: &opsmngr.LogRotate{
					SizeThresholdMB:  defaultSizeThresholdMB,
					TimeThresholdHrs: defaultTimeThresholdHrs,
				},
				ProcessType: "mongod",
				Version:     "4.2.2",
			},
			{
				Args26: opsmngr.Args26{
					NET: opsmngr.Net{
						Port: mongosPort,
					},
					Replication: nil,
					Sharding:    nil,
					Storage:     nil,
					SystemLog: opsmngr.SystemLog{
						Destination: "file",
						Path:        "/data/db/mongos.log",
					},
				},
				AuthSchemaVersion:           authSchemaVersion,
				Cluster:                     name,
				Name:                        name + "_mongos_0",
				Disabled:                    disabled,
				ManualMode:                  disabled,
				FeatureCompatibilityVersion: "4.2",
				Hostname:                    "host1",
				LogRotate: &opsmngr.LogRotate{
					SizeThresholdMB:  defaultSizeThresholdMB,
					TimeThresholdHrs: defaultTimeThresholdHrs,
				},
				ProcessType: "mongos",
				Version:     "4.2.2",
			},
		},
		ReplicaSets: []*opsmngr.ReplicaSet{
			{
				ID:              name + "_shard_0",
				ProtocolVersion: "1",
				Members: []opsmngr.Member{
					{
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         name + "_shard_0_0",
						Priority:     1,
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
				},
			},
			{
				ID:              name + "_configRS",
				ProtocolVersion: "1",
				Members: []opsmngr.Member{
					{
						ArbiterOnly:  false,
						BuildIndexes: true,
						Hidden:       false,
						Host:         name + "_configRS_0",
						Priority:     1,
						SlaveDelay:   &slaveDelay,
						Votes:        1,
					},
				},
			},
		},
		Sharding: []*opsmngr.ShardingConfig{
			{
				Name:                name,
				ConfigServerReplica: name + "_configRS",
				Shards: []*opsmngr.Shard{
					{
						ID: name + "_shard_0",
						RS: name + "_shard_0",
					},
				},
			},
		},
	}
}

func automationConfigWithoutMongoDBUsers() *opsmngr.AutomationConfig {
	return &opsmngr.AutomationConfig{
		Auth: opsmngr.Auth{
			AutoAuthMechanism: "MONGODB-CR",
			Disabled:          true,
			AuthoritativeSet:  false,
			UsersWanted:       make([]*opsmngr.MongoDBUser, 0),
		},
	}
}

func automationConfigWithIndexConfig() *opsmngr.AutomationConfig {
	return &opsmngr.AutomationConfig{
		IndexConfigs: []*opsmngr.IndexConfig{
			{
				DBName:         "test",
				CollectionName: "test",
				RSName:         "myReplicaSet",
				Key: [][]string{
					{
						"test", "test",
					},
				},
				Options:   nil,
				Collation: nil,
			}},
	}
}

func automationConfigWithMongoDBUsers() *opsmngr.AutomationConfig {
	return &opsmngr.AutomationConfig{
		Auth: opsmngr.Auth{
			AutoAuthMechanism: "MONGODB-CR",
			Disabled:          true,
			AuthoritativeSet:  false,
			UsersWanted: []*opsmngr.MongoDBUser{
				mongoDBUsers(),
			},
		},
	}
}

func mongoDBUsers() *opsmngr.MongoDBUser {
	return &opsmngr.MongoDBUser{
		Mechanisms: &[]string{"SCRAM-SHA-1"},
		Roles: []*opsmngr.Role{
			{
				Role:     "test",
				Database: "test",
			},
		},
		Username: "test",
		Database: "test",
	}
}
