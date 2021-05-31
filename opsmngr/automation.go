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
)

const (
	automationConfigBasePath = "groups/%s/automationConfig"
)

// AutomationService provides access to the automation related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/nav/automation/
type AutomationService interface {
	GetConfig(context.Context, string) (*AutomationConfig, *Response, error)
	UpdateConfig(context.Context, string, *AutomationConfig) (*Response, error)
	UpdateAgentVersion(context.Context, string) (*AutomationConfigAgent, *Response, error)
	GetBackupAgentConfig(context.Context, string) (*AgentConfig, *Response, error)
	GetMonitoringAgentConfig(context.Context, string) (*AgentConfig, *Response, error)
	GetStatus(context.Context, string) (*AutomationStatus, *Response, error)
}

// AutomationServiceOp provides an implementation of the AutomationService interface.
type AutomationServiceOp service

var _ AutomationService = new(AutomationServiceOp)
