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

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

// Suspend suspends automation for all processes.
func Suspend(out *opsmngr.AutomationConfig, clusterName string) {
	setManualModeByClusterName(out, clusterName, true)
}

// SuspendProcessesByClusterName suspends automation for the entire cluster or its processes.
// Processes are provided in the format {"hostname:port","hostname2:port2"}.
func SuspendProcessesByClusterName(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	if len(processes) == 0 {
		Suspend(out, clusterName)
		return nil
	}
	return suspendProcesses(out, clusterName, processes)
}

// suspendProcesses suspends automation for the processes.
// Processes are provided in the format {"hostname:port","hostname2:port2"}.
func suspendProcesses(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	processesMap := newProcessMap(processes)
	setManualModeByNameAndProcesses(out, clusterName, processesMap, true)

	return newProcessNotFoundError(clusterName, processesMap)
}

func setManualModeByClusterName(out *opsmngr.AutomationConfig, name string, manualMode bool) {
	newDeploymentAuthMechanisms(out)
	setManualModeByReplicaSetName(out, name, manualMode)
	setManualModeByShardName(out, name, manualMode)
}

func setManualModeByReplicaSetName(out *opsmngr.AutomationConfig, name string, manualMode bool) {
	setManualModeByReplicaSetNameAndProcesses(out, name, nil, manualMode)
}

func setManualModeByShardName(out *opsmngr.AutomationConfig, name string, manualMode bool) {
	setManualModeByShardNameAndProcesses(out, name, nil, manualMode)
}

func setManualModeByReplicaSetNameAndProcesses(out *opsmngr.AutomationConfig, name string, processesMap map[string]bool, manualMode bool) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == name
	})
	if !found {
		return
	}
	rs := out.ReplicaSets[i]
	for _, m := range rs.Members {
		for k, p := range out.Processes {
			if p.Name == m.Host {
				setManualMode(out.Processes[k], processesMap, manualMode)
			}
		}
	}

}

func setManualModeByShardNameAndProcesses(out *opsmngr.AutomationConfig, name string, processesMap map[string]bool, manualMode bool) {
	i, found := search.ShardingConfig(out.Sharding, func(s *opsmngr.ShardingConfig) bool {
		return s.Name == name
	})
	if !found {
		return
	}

	s := out.Sharding[i]
	// manual mode per shards
	for _, rs := range s.Shards {
		setManualModeByReplicaSetNameAndProcesses(out, rs.ID, processesMap, manualMode)
	}
	// manual mode config rs
	setManualModeByReplicaSetNameAndProcesses(out, s.ConfigServerReplica, processesMap, manualMode)
	// manual mode mongos
	for j := range out.Processes {
		if out.Processes[j].Cluster == name {
			setManualMode(out.Processes[j], processesMap, manualMode)
		}
	}
}

func setManualMode(process *opsmngr.Process, processesMap map[string]bool, manualMode bool) {
	if len(processesMap) == 0 {
		process.ManualMode = manualMode
		return
	}

	key := fmt.Sprintf("%s:%d", process.Hostname, process.Args26.NET.Port)
	if _, ok := processesMap[key]; ok {
		process.ManualMode = manualMode
		processesMap[key] = true
	}
}

func setManualModeByNameAndProcesses(out *opsmngr.AutomationConfig, clusterName string, processesMap map[string]bool, manualMode bool) {
	newDeploymentAuthMechanisms(out)
	setManualModeByReplicaSetNameAndProcesses(out, clusterName, processesMap, manualMode)
	setManualModeByShardNameAndProcesses(out, clusterName, processesMap, manualMode)
}
