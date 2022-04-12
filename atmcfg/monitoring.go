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
	"errors"
	"fmt"

	"go.mongodb.org/ops-manager/opsmngr"
)

var ErrMonitoringEnabled = errors.New("monitoring already enabled")

const monitoringVersion = "7.2.0.488-1" // monitoringVersion last monitoring version released

// EnableMonitoring enables monitoring for the given hostname.
func EnableMonitoring(out *opsmngr.AutomationConfig, hostname string) error {
	for _, v := range out.MonitoringVersions {
		if v.Hostname == hostname {
			return fmt.Errorf("%w for '%s'", ErrMonitoringEnabled, hostname)
		}
	}
	out.MonitoringVersions = append(out.MonitoringVersions, &opsmngr.ConfigVersion{
		Name:     monitoringVersion,
		Hostname: hostname,
	})
	return nil
}

// DisableMonitoring disables monitoring for the given hostname.
func DisableMonitoring(out *opsmngr.AutomationConfig, hostname string) error {
	for i, v := range out.MonitoringVersions {
		if v.Hostname == hostname {
			out.MonitoringVersions = append(out.MonitoringVersions[:i], out.MonitoringVersions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no monitoring for '%s'", hostname)
}
