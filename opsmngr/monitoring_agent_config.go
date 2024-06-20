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
)

const (
	monitoringAgentConfigBasePath = "api/public/v1.0/groups/%s/automationConfig/monitoringAgentConfig"
)

// GetMonitoringAgentConfig retrieves the current monitoring agent configuration for a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/#get-monitoring-attributes
func (s *AutomationServiceOp) GetMonitoringAgentConfig(ctx context.Context, groupID string) (*AgentConfig, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}
	basePath := fmt.Sprintf(monitoringAgentConfigBasePath, groupID)

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
