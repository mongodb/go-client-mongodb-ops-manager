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
)

const backupVersion = "7.8.1.1109-1" // backupVersion last backup version released

// EnableBackup enables backup for the given hostname.
func EnableBackup(out *opsmngr.AutomationConfig, hostname string) error {
	for _, v := range out.BackupVersions {
		if v.Hostname == hostname {
			return fmt.Errorf("backup already enabled for '%s'", hostname)
		}
	}
	out.BackupVersions = append(out.BackupVersions, &opsmngr.ConfigVersion{
		Name:     backupVersion,
		Hostname: hostname,
	})
	return nil
}

// DisableBackup disables backup for the given hostname.
func DisableBackup(out *opsmngr.AutomationConfig, hostname string) error {
	for i, v := range out.BackupVersions {
		if v.Hostname == hostname {
			out.BackupVersions = append(out.BackupVersions[:i], out.BackupVersions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no backup for '%s'", hostname)
}
