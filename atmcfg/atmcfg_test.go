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

package atmcfg

import (
	"testing"

	"github.com/mongodb/go-client-mongodb-ops-manager/opsmngr"
)

func automationConfigWithOneReplicaSet(name string, disabled bool) *opsmngr.AutomationConfig {
	return &opsmngr.AutomationConfig{
		Processes: []*opsmngr.Process{
			{
				Args26: opsmngr.Args26{
					NET: opsmngr.Net{
						Port: 27017,
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
				AuthSchemaVersion:           5,
				Name:                        name + "_0",
				Disabled:                    disabled,
				FeatureCompatibilityVersion: "4.2",
				Hostname:                    "host0",
				LogRotate: &opsmngr.LogRotate{
					SizeThresholdMB:  1000,
					TimeThresholdHrs: 24,
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
						SlaveDelay:   0,
						Votes:        1,
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
			Users:             make([]*opsmngr.MongoDBUser, 0),
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
			Users: []*opsmngr.MongoDBUser{
				mongoDBUsers(),
			},
		},
	}
}

func mongoDBUsers() *opsmngr.MongoDBUser {
	return &opsmngr.MongoDBUser{
		Mechanisms: []string{"SCRAM-SHA-1"},
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

func TestShutdown(t *testing.T) {
	name := "cluster_1"
	config := automationConfigWithOneReplicaSet(name, false)

	Shutdown(config, name)
	if !config.Processes[0].Disabled {
		t.Errorf("TestShutdown\n got=%#v\nwant=%#v\n", config.Processes[0].Disabled, true)
	}
}

func TestStartup(t *testing.T) {
	name := "cluster_1"
	cloud := automationConfigWithOneReplicaSet(name, true)

	Startup(cloud, name)
	if cloud.Processes[0].Disabled {
		t.Errorf("TestStartup\n got=%#v\nwant=%#v\n", cloud.Processes[0].Disabled, false)
	}
}

func TestAddUser(t *testing.T) {
	config := automationConfigWithoutMongoDBUsers()
	u := mongoDBUsers()
	AddUser(config, u)
	if len(config.Auth.Users) != 1 {
		t.Error("User not added\n")
	}
}

func TestAddIndexConfig(t *testing.T) {
	newIndex := &opsmngr.IndexConfig{
		DBName:         "test",
		CollectionName: "test",
		RSName:         "myReplicaSet",
		Key: [][]string{
			{
				"test1", "test",
			},
		},
		Options:   nil,
		Collation: nil,
	}
	t.Run("AutomationConfig not initialized", func(t *testing.T) {
		err := AddIndexConfig(nil, newIndex)
		if err == nil {
			t.Error("AddIndexConfig should return an exception")
		}
	})

	t.Run("empty IndexConfig", func(t *testing.T) {
		a := &opsmngr.AutomationConfig{}
		err := AddIndexConfig(a, newIndex)
		if err != nil {
			t.Fatalf("AddIndexConfig unexpected error: %v", err)
		}
		if len(a.IndexConfigs) != 1 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("add an index with different keys", func(t *testing.T) {
		config := automationConfigWithIndexConfig()
		err := AddIndexConfig(config, newIndex)
		if err != nil {
			t.Fatalf("AutomationConfig() returned an unexpected error: %v", err)
		}
		if len(config.IndexConfigs) != 2 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("add an index with different rsName", func(t *testing.T) {
		newIndex := &opsmngr.IndexConfig{
			DBName:         "test",
			CollectionName: "test",
			RSName:         "myReplicaSet_1",
			Key: [][]string{
				{
					"test1", "test",
				},
			},
			Options:   nil,
			Collation: nil,
		}
		config := automationConfigWithIndexConfig()
		err := AddIndexConfig(config, newIndex)
		if err != nil {
			t.Fatalf("AutomationConfig() returned an unexpected error: %v", err)
		}
		if len(config.IndexConfigs) != 2 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("trying to add an index that is already in the AutomationConfig", func(t *testing.T) {
		config := automationConfigWithIndexConfig()
		index := &opsmngr.IndexConfig{
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
		}
		err := AddIndexConfig(config, index)

		if err != nil {
			t.Fatalf("AutomationConfig() returned an unexpected error: %v", err)
		}

		if len(config.IndexConfigs) != 1 {
			t.Error("the same indexConfig has been added to the AutomationConfig")
		}
	})

}

func TestRemoveUser(t *testing.T) {
	config := automationConfigWithMongoDBUsers()
	t.Run("user exists", func(t *testing.T) {
		u := mongoDBUsers()
		err := RemoveUser(config, u.Username, u.Database)
		if err != nil {
			t.Fatalf("RemoveUser unexpecter err: %#v\n", err)
		}
		if len(config.Auth.Users) != 0 {
			t.Error("User not removed\n")
		}
	})
	t.Run("user does not exists", func(t *testing.T) {
		err := RemoveUser(config, "random", "random")
		if err == nil {
			t.Fatal("RemoveUser should return an error\n")
		}
	})
}
