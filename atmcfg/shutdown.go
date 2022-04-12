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

import "go.mongodb.org/ops-manager/opsmngr"

// Shutdown disables all processes of the given cluster name.
func Shutdown(out *opsmngr.AutomationConfig, clusterName string) {
	setDisabledByClusterName(out, clusterName, true)
}

// ShutdownProcessesByClusterName disables the entire cluster or its processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func ShutdownProcessesByClusterName(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	if len(processes) == 0 {
		Shutdown(out, clusterName)
		return nil
	}
	return shutdownProcesses(out, clusterName, processes)
}

// shutdownProcesses disables processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func shutdownProcesses(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	processesMap := newProcessMap(processes)
	setDisableByNameAndProcesses(out, clusterName, processesMap, true)

	return newProcessNotFoundError(clusterName, processesMap)
}
