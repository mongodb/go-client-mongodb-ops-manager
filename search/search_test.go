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

	"github.com/mongodb/go-client-mongodb-ops-manager/opsmngr"
	"github.com/mongodb/go-client-mongodb-ops-manager/search"
)

func TestProcesses(t *testing.T) {
	processes := fixture.Processes
	t.Run("value exists", func(t *testing.T) {
		_, e := search.Processes(processes, func(p *opsmngr.Process) bool {
			return p.Name == "myReplicaSet_1"
		})
		if !e {
			t.Error("Processes() should find the value")
		}
	})

	t.Run("value does not exists", func(t *testing.T) {
		i, e := search.Processes(processes, func(p *opsmngr.Process) bool {
			return p.Name == "myReplicaSet_4"
		})
		if e {
			t.Errorf("Processes() found at: %d", i)
		}
	})
}

func TestMembers(t *testing.T) {
	members := fixture.ReplicaSets[0].Members
	t.Run("value exists", func(t *testing.T) {
		_, e := search.Members(members, func(p opsmngr.Member) bool {
			return p.Host == "myReplicaSet_1"
		})
		if !e {
			t.Error("Members() should find the value")
		}
	})

	t.Run("value does not exists", func(t *testing.T) {
		i, e := search.Members(members, func(p opsmngr.Member) bool {
			return p.Host == "myReplicaSet_4"
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
		i, e := search.ReplicaSets(rs, func(p *opsmngr.ReplicaSet) bool {
			return p.ID == "other"
		})
		if e {
			t.Errorf("ReplicaSets() found at: %d", i)
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
			return p.Username == "other"
		})
		if e {
			t.Errorf("MongoDBUsers() found at: %d", i)
		}
	})
}
