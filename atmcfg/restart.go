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

package atmcfg // Package atmcfg import "go.mongodb.org/ops-manager/atmcfg"

import (
	"fmt"
	"time"

	"go.mongodb.org/ops-manager/search"

	"go.mongodb.org/ops-manager/opsmngr"
)

// Restart sets all process of a cluster to restart.
func Restart(out *opsmngr.AutomationConfig, name string) {
	newDeploymentAuthMechanisms(out)
	lastRestart := time.Now().Format(time.RFC3339)
	restartByReplicaSetName(out, name, lastRestart)
	restartByShardName(out, name, lastRestart)
}

// RestartProcessesByClusterName restart the entire cluster or its processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func RestartProcessesByClusterName(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	if len(processes) == 0 {
		Restart(out, clusterName)
		return nil
	}
	return restartProcesses(out, clusterName, processes)
}

// restartProcesses restarts processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func restartProcesses(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	processesMap := newProcessMap(processes)
	restartByNameAndProcesses(out, clusterName, processesMap)

	return newProcessNotFoundError(clusterName, processesMap)
}

func restartByNameAndProcesses(out *opsmngr.AutomationConfig, clusterName string, processesMap map[string]bool) {
	lastRestart := time.Now().Format(time.RFC3339)
	newDeploymentAuthMechanisms(out)
	restartByReplicaSetNameAndProcesses(out, processesMap, clusterName, lastRestart)
	restartByShardNameAndProcesses(out, processesMap, clusterName, lastRestart)
}

func restartByReplicaSetName(out *opsmngr.AutomationConfig, clusterName, lastRestart string) {
	restartByReplicaSetNameAndProcesses(out, nil, clusterName, lastRestart)
}

func restartByReplicaSetNameAndProcesses(out *opsmngr.AutomationConfig, processesMap map[string]bool, clusterName, lastRestart string) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == clusterName
	})
	if found {
		rs := out.ReplicaSets[i]
		for _, m := range rs.Members {
			for k, p := range out.Processes {
				if p.Name == m.Host {
					setLastRestart(out.Processes[k], processesMap, lastRestart)
				}
			}
		}
	}
}

func setLastRestart(process *opsmngr.Process, processesMap map[string]bool, lastRestart string) {
	if len(processesMap) == 0 {
		process.LastRestart = lastRestart
		return
	}

	key := fmt.Sprintf("%s:%d", process.Hostname, process.Args26.NET.Port)
	if _, ok := processesMap[key]; ok {
		process.LastRestart = lastRestart
		processesMap[key] = true
	}
}

func restartByShardName(out *opsmngr.AutomationConfig, clusterName, lastRestart string) {
	restartByShardNameAndProcesses(out, nil, clusterName, lastRestart)
}

func restartByShardNameAndProcesses(out *opsmngr.AutomationConfig, processesMap map[string]bool, clusterName, lastRestart string) {
	i, found := search.ShardingConfig(out.Sharding, func(s *opsmngr.ShardingConfig) bool {
		return s.Name == clusterName
	})
	if found {
		s := out.Sharding[i]
		// restart shards
		for _, rs := range s.Shards {
			restartByReplicaSetNameAndProcesses(out, processesMap, rs.ID, lastRestart)
		}
		// restart config rs
		restartByReplicaSetNameAndProcesses(out, processesMap, s.ConfigServerReplica, lastRestart)
		// restart mongos
		for i := range out.Processes {
			if out.Processes[i].Cluster == clusterName {
				setLastRestart(out.Processes[i], processesMap, lastRestart)
			}
		}
	}
}
