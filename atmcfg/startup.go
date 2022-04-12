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

import "go.mongodb.org/ops-manager/opsmngr"

// StartupProcessesByClusterName enables the entire cluster or its processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func StartupProcessesByClusterName(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	if len(processes) == 0 {
		Startup(out, clusterName)
		return nil
	}
	return startupProcess(out, clusterName, processes)
}

// StartupProcess enables processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func startupProcess(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	processesMap := newProcessMap(processes)
	setDisableByNameAndProcesses(out, clusterName, processesMap, false)

	return newProcessNotFoundError(clusterName, processesMap)
}

// Startup enables all processes of the given cluster name.
func Startup(out *opsmngr.AutomationConfig, name string) {
	setDisabledByClusterName(out, name, false)
}
