// Copyright 2023 MongoDB Inc
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
	"fmt"
	"time"

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

// StartInitialSync the MongoDB Agent checks whether the specified timestamp is later than the time of the last resync,
// and if confirmed, starts init sync on the secondary nodes in a rolling fashion.
// The MongoDB Agent waits until you ask the primary node to become the secondary with the rs.stepDown() method, and then starts init sync on this node.
//
// Warning: Use this method with caution. During initial sync, Automation removes the entire contents of the node’s dbPath directory.
//
// See also: https://www.mongodb.com/docs/manual/core/replica-set-sync/#replica-set-initial-sync
func StartInitialSync(out *opsmngr.AutomationConfig, clusterName string) {
	StartInitialSyncAt(out, clusterName, "")
}

// StartInitialSyncAt specify the type to start the initial sync at,
//
// Warning: Use this method with caution. During initial sync, Automation removes the entire contents of the node’s dbPath directory.
//
// See also: https://www.mongodb.com/docs/manual/core/replica-set-sync/#replica-set-initial-sync
func StartInitialSyncAt(out *opsmngr.AutomationConfig, clusterName, lastResync string) {
	if lastResync == "" {
		lastResync = time.Now().Format(time.RFC3339)
	}
	newDeploymentAuthMechanisms(out)
	syncByReplicaSetName(out, clusterName, lastResync)
	syncByShardName(out, clusterName, lastResync)
}

// StartInitialSyncAtForProcessesByClusterName trigger initial sync for a cluster. Processes are provided in the format {"hostname:port","hostname2:port2"}.
//
// See also: https://www.mongodb.com/docs/manual/core/replica-set-sync/#replica-set-initial-sync
func StartInitialSyncAtForProcessesByClusterName(out *opsmngr.AutomationConfig, clusterName, lastResync string, processes []string) error {
	if len(processes) == 0 {
		StartInitialSyncAt(out, clusterName, lastResync)
		return nil
	}

	return initialSyncForProcesses(out, clusterName, lastResync, processes)
}

func initialSyncForProcesses(out *opsmngr.AutomationConfig, clusterName, lastResync string, processes []string) error {
	newDeploymentAuthMechanisms(out)
	processesMap := newProcessMap(processes)
	syncByReplicaSetNameAndProcesses(out, processesMap, clusterName, lastResync)
	syncByShardNameAndProcesses(out, processesMap, clusterName, lastResync)

	return newProcessNotFoundError(clusterName, processesMap)
}

func syncByReplicaSetName(out *opsmngr.AutomationConfig, clusterName, lastResync string) {
	syncByReplicaSetNameAndProcesses(out, nil, clusterName, lastResync)
}

func syncByReplicaSetNameAndProcesses(out *opsmngr.AutomationConfig, processesMap map[string]bool, clusterName, lastResync string) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == clusterName
	})
	if found {
		rs := out.ReplicaSets[i]
		for _, m := range rs.Members {
			for k, p := range out.Processes {
				if p.Name == m.Host {
					seLastResync(out.Processes[k], processesMap, lastResync)
				}
			}
		}
	}
}

func seLastResync(process *opsmngr.Process, processesMap map[string]bool, lastResync string) {
	if len(processesMap) == 0 {
		process.LastResync = lastResync
		return
	}

	key := fmt.Sprintf("%s:%d", process.Hostname, process.Args26.NET.Port)
	if _, ok := processesMap[key]; ok {
		process.LastResync = lastResync
		processesMap[key] = true
	}
}

func syncByShardName(out *opsmngr.AutomationConfig, clusterName, lastResync string) {
	syncByShardNameAndProcesses(out, nil, clusterName, lastResync)
}

func syncByShardNameAndProcesses(out *opsmngr.AutomationConfig, processesMap map[string]bool, clusterName, lastResync string) {
	i, found := search.ShardingConfig(out.Sharding, func(s *opsmngr.ShardingConfig) bool {
		return s.Name == clusterName
	})
	if found {
		s := out.Sharding[i]
		// sync shards
		for _, rs := range s.Shards {
			syncByReplicaSetNameAndProcesses(out, processesMap, rs.ID, lastResync)
		}
		// sync config rs
		syncByReplicaSetNameAndProcesses(out, processesMap, s.ConfigServerReplica, lastResync)
		// sync doesn't run on mongoses
	}
}
