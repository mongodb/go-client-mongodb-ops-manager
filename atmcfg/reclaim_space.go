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
	"fmt"
	"time"

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

// ReclaimFreeSpace sets all process of a cluster to reclaim free space.
func ReclaimFreeSpace(out *opsmngr.AutomationConfig, clusterName string) {
	ReclaimFreeSpaceWithLastCompact(out, clusterName, "")
}

func ReclaimFreeSpaceWithLastCompact(out *opsmngr.AutomationConfig, clusterName, lastCompact string) {
	if lastCompact == "" {
		lastCompact = time.Now().Format(time.RFC3339)
	}
	newDeploymentAuthMechanisms(out)
	reclaimByReplicaSetName(out, clusterName, lastCompact)
	reclaimByShardName(out, clusterName, lastCompact)
}

// ReclaimFreeSpaceForProcessesByClusterName reclaims free space for a cluster. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func ReclaimFreeSpaceForProcessesByClusterName(out *opsmngr.AutomationConfig, clusterName, lastCompact string, processes []string) error {
	if len(processes) == 0 {
		ReclaimFreeSpaceWithLastCompact(out, clusterName, lastCompact)
		return nil
	}

	return reclaimFreeSpaceForProcesses(out, clusterName, lastCompact, processes)
}

func reclaimFreeSpaceForProcesses(out *opsmngr.AutomationConfig, clusterName, lastCompact string, processes []string) error {
	newDeploymentAuthMechanisms(out)
	processesMap := newProcessMap(processes)
	reclaimByReplicaSetNameAndProcesses(out, processesMap, clusterName, lastCompact)
	reclaimByShardNameAndProcesses(out, processesMap, clusterName, lastCompact)

	return newProcessNotFoundError(clusterName, processesMap)
}

func reclaimByReplicaSetName(out *opsmngr.AutomationConfig, clusterName, lastCompact string) {
	reclaimByReplicaSetNameAndProcesses(out, nil, clusterName, lastCompact)
}

func reclaimByReplicaSetNameAndProcesses(out *opsmngr.AutomationConfig, processesMap map[string]bool, clusterName, lastCompact string) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == clusterName
	})
	if found {
		rs := out.ReplicaSets[i]
		for _, m := range rs.Members {
			for k, p := range out.Processes {
				if p.Name == m.Host {
					setLastCompact(out.Processes[k], processesMap, lastCompact)
				}
			}
		}
	}
}

func setLastCompact(process *opsmngr.Process, processesMap map[string]bool, lastCompact string) {
	if len(processesMap) == 0 {
		process.LastCompact = lastCompact
		return
	}

	key := fmt.Sprintf("%s:%d", process.Hostname, process.Args26.NET.Port)
	if _, ok := processesMap[key]; ok {
		process.LastCompact = lastCompact
		processesMap[key] = true
	}
}

func reclaimByShardName(out *opsmngr.AutomationConfig, clusterName, lastCompact string) {
	reclaimByShardNameAndProcesses(out, nil, clusterName, lastCompact)
}

func reclaimByShardNameAndProcesses(out *opsmngr.AutomationConfig, processesMap map[string]bool, clusterName, lastCompact string) {
	i, found := search.ShardingConfig(out.Sharding, func(s *opsmngr.ShardingConfig) bool {
		return s.Name == clusterName
	})
	if found {
		s := out.Sharding[i]
		// compact shards
		for _, rs := range s.Shards {
			reclaimByReplicaSetNameAndProcesses(out, processesMap, rs.ID, lastCompact)
		}
		// compact config rs
		reclaimByReplicaSetNameAndProcesses(out, processesMap, s.ConfigServerReplica, lastCompact)
		// compact doesn't run on mongoses
	}
}
