// Copyright 2022 MongoDB Inc
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
	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

// RemoveByClusterName removes a cluster and its associated processes from the config.
// This won't shutdown any running process.
func RemoveByClusterName(out *opsmngr.AutomationConfig, name string) {
	// This value may not be present and is mandatory
	if out.Auth.DeploymentAuthMechanisms == nil {
		out.Auth.DeploymentAuthMechanisms = make([]string, 0)
	}
	removeByReplicaSetName(out, name)
	removeByShardName(out, name)
}

func removeByReplicaSetName(out *opsmngr.AutomationConfig, name string) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == name
	})
	if found {
		rs := out.ReplicaSets[i]
		out.ReplicaSets = append(out.ReplicaSets[:i], out.ReplicaSets[i+1:]...)
		for _, m := range rs.Members {
			processes := []*opsmngr.Process{}
			for _, p := range out.Processes {
				if p.Name != m.Host {
					processes = append(processes, p)
				}
			}
			out.Processes = processes
		}
	}
}

func removeByShardName(out *opsmngr.AutomationConfig, name string) {
	i, found := search.ShardingConfig(out.Sharding, func(rs *opsmngr.ShardingConfig) bool {
		return rs.Name == name
	})
	if found {
		s := out.Sharding[i]
		out.Sharding = append(out.Sharding[:i], out.Sharding[i+1:]...)
		// remove shards
		for _, rs := range s.Shards {
			removeByReplicaSetName(out, rs.ID)
		}
		// remove config rs
		removeByReplicaSetName(out, s.ConfigServerReplica)
		// remove mongos
		processes := []*opsmngr.Process{}
		for _, p := range out.Processes {
			if p.Cluster != name {
				processes = append(processes, p)
			}
		}
		out.Processes = processes
	}
}
