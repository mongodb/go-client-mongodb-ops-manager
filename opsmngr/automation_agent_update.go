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

// UpdateAgentVersion updates the MongoDB Agent and tools to the latest versions available at the time of the request.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/#update-agend-versions-example
func (s *AutomationServiceOp) UpdateAgentVersion(ctx context.Context, groupID string) (*AutomationConfigAgent, *Response, error) {
	if groupID == "" {
		return nil, nil, NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(automationConfigBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "updateAgentVersions")

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	agent := new(AutomationConfigAgent)
	resp, err := s.Client.Do(ctx, req, agent)

	return agent, resp, err
}

// AutomationConfigAgent components versions.
type AutomationConfigAgent struct {
	AutomationAgentVersion string `json:"automationAgentVersion"`
	BiConnectorVersion     string `json:"biConnectorVersion"`
}
