// Copyright 2019 MongoDB Inc
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

package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	backupAgentConfigBasePath = "groups/%s/automationConfig/backupAgentConfig"
)

// GetBackupAgentConfig retrieves the current backup agent configuration for a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/#get-backup-attributes
func (s *AutomationServiceOp) GetBackupAgentConfig(ctx context.Context, groupID string) (*AgentConfig, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	basePath := fmt.Sprintf(backupAgentConfigBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, basePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AgentConfig)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

type AgentConfig struct {
	LogPath                 string                  `json:"logPath,omitempty"`
	LogPathWindows          string                  `json:"logPathWindows,omitempty"`
	LogRotate               *LogRotate              `json:"logRotate,omitempty"`
	ConfigOverrides         *map[string]string      `json:"configOverrides,omitempty"`
	HasPassword             *bool                   `json:"hasPassword,omitempty"`
	Username                string                  `json:"username,omitempty"`
	Password                string                  `json:"password,omitempty"`
	SSLPEMKeyFile           string                  `json:"sslPEMKeyFile,omitempty"`
	SSLPEMKeyFileWindows    string                  `json:"sslPEMKeyFileWindows,omitempty"`
	SSLPEMKeyPwd            string                  `json:"sslPEMKeyPwd,omitempty"`
	KerberosPrincipal       string                  `json:"kerberosPrincipal,omitempty"`
	KerberosKeytab          string                  `json:"kerberosKeytab,omitempty"`
	KerberosWindowsUsername string                  `json:"kerberosWindowsUsername,omitempty"`
	KerberosWindowsPassword string                  `json:"kerberosWindowsPassword,omitempty"`
	HasSslPEMKeyPwd         *bool                   `json:"hasSslPEMKeyPwd,omitempty"`
	LDAPGroupDN             string                  `json:"ldapGroupDN,omitempty"`
	URLs                    *map[string]interface{} `json:"urls,omitempty"`
}
