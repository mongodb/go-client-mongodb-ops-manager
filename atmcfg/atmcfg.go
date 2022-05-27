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

package atmcfg // Package atmcfg import "go.mongodb.org/ops-manager/atmcfg"

import (
	"fmt"

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

func newDeploymentAuthMechanisms(out *opsmngr.AutomationConfig) {
	// This value may not be present and is mandatory
	if out.Auth.DeploymentAuthMechanisms == nil {
		out.Auth.DeploymentAuthMechanisms = make([]string, 0)
	}
}

func setDisabledByClusterName(out *opsmngr.AutomationConfig, name string, disabled bool) {
	newDeploymentAuthMechanisms(out)
	setDisabledByReplicaSetName(out, name, disabled)
	setDisabledByShardName(out, name, disabled)
}

func setDisabledByReplicaSetName(out *opsmngr.AutomationConfig, name string, disabled bool) {
	setDisabledByReplicaSetNameAndProcesses(out, name, nil, disabled)
}

func setDisabledByShardName(out *opsmngr.AutomationConfig, name string, disabled bool) {
	setDisabledByShardNameAndProcesses(out, name, nil, disabled)
}

func setDisabledByReplicaSetNameAndProcesses(out *opsmngr.AutomationConfig, name string, processesMap map[string]bool, disabled bool) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == name
	})
	if found {
		rs := out.ReplicaSets[i]
		for _, m := range rs.Members {
			for k, p := range out.Processes {
				if p.Name == m.Host {
					setDisable(out.Processes[k], processesMap, disabled)
				}
			}
		}
	}
}

func setDisabledByShardNameAndProcesses(out *opsmngr.AutomationConfig, name string, processesMap map[string]bool, disabled bool) {
	i, found := search.ShardingConfig(out.Sharding, func(s *opsmngr.ShardingConfig) bool {
		return s.Name == name
	})
	if found {
		s := out.Sharding[i]
		// disable shards
		for _, rs := range s.Shards {
			setDisabledByReplicaSetNameAndProcesses(out, rs.ID, processesMap, disabled)
		}
		// disable config rs
		setDisabledByReplicaSetNameAndProcesses(out, s.ConfigServerReplica, processesMap, disabled)
		// disable mongos
		for i := range out.Processes {
			if out.Processes[i].Cluster == name {
				setDisable(out.Processes[i], processesMap, disabled)
			}
		}
	}
}

func setDisable(process *opsmngr.Process, processesMap map[string]bool, disabled bool) {
	if len(processesMap) == 0 {
		process.Disabled = disabled
		return
	}

	key := fmt.Sprintf("%s:%d", process.Hostname, process.Args26.NET.Port)
	if _, ok := processesMap[key]; ok {
		process.Disabled = disabled
		processesMap[key] = true
	}
}

func setDisableByNameAndProcesses(out *opsmngr.AutomationConfig, clusterName string, processesMap map[string]bool, disabled bool) {
	newDeploymentAuthMechanisms(out)
	setDisabledByReplicaSetNameAndProcesses(out, clusterName, processesMap, disabled)
	setDisabledByShardNameAndProcesses(out, clusterName, processesMap, disabled)
}

func stringInSlice(a []string, x string) bool {
	for _, b := range a {
		if b == x {
			return true
		}
	}
	return false
}

func newProcessMap(processes []string) map[string]bool {
	processesMap := map[string]bool{}
	for _, hostnameAndPort := range processes {
		processesMap[hostnameAndPort] = false
	}

	return processesMap
}

func IsGoalState(s *opsmngr.AutomationStatus) bool {
	for _, p := range s.Processes {
		if p.LastGoalVersionAchieved != s.GoalVersion {
			return false
		}
	}
	return true
}
