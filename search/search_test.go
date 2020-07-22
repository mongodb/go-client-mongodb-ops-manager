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
	"testing"

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

func TestProcesses(t *testing.T) {
	processes := fixture.Processes
	t.Run("value exists", func(t *testing.T) {
		want := rsName
		_, e := search.Processes(processes, func(p *opsmngr.Process) bool {
			return p.Name == want
		})
		if !e {
			t.Error("Processes() should find the value")
		}
	})

	t.Run("value does not exists", func(t *testing.T) {
		want := "other_rs"
		i, e := search.Processes(processes, func(p *opsmngr.Process) bool {
			return p.Name == want
		})
		if e {
			t.Errorf("Processes() found at: %d", i)
		}
	})
}

func TestMembers(t *testing.T) {
	members := fixture.ReplicaSets[0].Members
	t.Run("value exists", func(t *testing.T) {
		want := rsName
		_, e := search.Members(members, func(p opsmngr.Member) bool {
			return p.Host == want
		})
		if !e {
			t.Error("Members() should find the value")
		}
	})

	t.Run("value does not exists", func(t *testing.T) {
		want := "other_rs_!"
		i, e := search.Members(members, func(p opsmngr.Member) bool {
			return p.Host == want
		})
		if e {
			t.Errorf("Members() found at: %d", i)
		}
	})
}

func TestReplicaSets(t *testing.T) {
	rs := fixture.ReplicaSets
	t.Run("value exists", func(t *testing.T) {
		_, e := search.ReplicaSets(rs, func(p *opsmngr.ReplicaSet) bool {
			return p.ID == "myReplicaSet"
		})
		if !e {
			t.Error("ReplicaSets() should find the value")
		}
	})

	t.Run("value does not exists", func(t *testing.T) {
		want := "other_rs"
		i, e := search.ReplicaSets(rs, func(p *opsmngr.ReplicaSet) bool {
			return p.ID == want
		})
		if e {
			t.Errorf("ReplicaSets() found at: %d", i)
		}
	})
}

func TestShardingConfig(t *testing.T) {
	rs := []*opsmngr.ShardingConfig{{Name: "myCluster"}}
	t.Run("value exists", func(t *testing.T) {
		_, e := search.ShardingConfig(rs, func(p *opsmngr.ShardingConfig) bool {
			return p.Name == "myCluster"
		})
		if !e {
			t.Error("ShardingConfig() should find the value")
		}
	})

	t.Run("value does not exists", func(t *testing.T) {
		i, e := search.ShardingConfig(rs, func(p *opsmngr.ShardingConfig) bool {
			return p.Name == "other_shard"
		})
		if e {
			t.Errorf("ShardingConfig() found at: %d", i)
		}
	})
}

func TestMongoDBUsers(t *testing.T) {
	users := fixture.Auth.Users
	t.Run("value exists", func(t *testing.T) {
		_, e := search.MongoDBUsers(users, func(p *opsmngr.MongoDBUser) bool {
			return p.Username == "test"
		})
		if !e {
			t.Error("MongoDBUsers() should find the value")
		}
	})

	t.Run("value does not exists", func(t *testing.T) {
		i, e := search.MongoDBUsers(users, func(p *opsmngr.MongoDBUser) bool {
			return p.Username == "other_user"
		})
		if e {
			t.Errorf("MongoDBUsers() found at: %d", i)
		}
	})
}

func TestMongoDBIndexes(t *testing.T) {
	indexConfigs := fixture.IndexConfigs
	t.Run("value exists", func(t *testing.T) {
		_, e := search.MongoDBIndexes(indexConfigs, func(p *opsmngr.IndexConfig) bool {
			return p.RSName == "myReplicaSet_1"
		})
		if !e {
			t.Error("MongoDBIndexes() should find the value")
		}
	})

	t.Run("value does not exists", func(t *testing.T) {
		i, e := search.MongoDBIndexes(indexConfigs, func(p *opsmngr.IndexConfig) bool {
			return p.RSName == "myReplicaSet_4"
		})
		if e {
			t.Errorf("MongoDBIndexes() found at: %d", i)
		}
	})
}
