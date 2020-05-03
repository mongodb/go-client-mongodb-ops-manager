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

package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	automationStatusBasePath = "groups/%s/automationStatus"
)

// GetStatus lets you see whether each MongoDB process is up-to-date with the current automation configuration.
// The endpoint returns the goalVersion field to report the current version of the automation configuration
// and the lastGoalVersionAchieved fields to report the versions of the configuration running on each server.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/automation-status/
func (s *AutomationServiceOp) GetStatus(ctx context.Context, groupID string) (*AutomationStatus, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	basePath := fmt.Sprintf(automationStatusBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, basePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AutomationStatus)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

type AutomationStatus struct {
	Processes   []ProcessStatus `json:"processes"`
	GoalVersion int             `json:"goalVersion"`
}

type ProcessStatus struct {
	Plan                    []string `json:"plan"`
	LastGoalVersionAchieved int      `json:"lastGoalVersionAchieved"`
	Name                    string   `json:"name"`
	Hostname                string   `json:"hostname"`
}
