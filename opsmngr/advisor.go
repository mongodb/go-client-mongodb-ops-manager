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
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	upgradeCheckPath = "api/private/v1.0/migration/check"
)

// AdvisorService ...
type AdvisorService interface {
	CheckUpgrade(context.Context, string) (*[]UpgradeCheckStep, *Response, error)
}

type AdvisorServiceOp service

var _ AdvisorService = &AdvisorServiceOp{}

// UpgradeCheckResult ...
type UpgradeCheckResult struct {
	Result json.RawMessage `json:"result,omitempty"`
	Error  string          `json:"error,omitempty"`
}

type RawRequirement struct {
	StepType       string   `json:"stepType"`
	BaseVersion    *string  `json:"baseVersion"`
	CurrentVersion string   `json:"currentVersion"`
	TargetVersion  string   `json:"targetVersion"`
	Hostnames      []string `json:"hostnames"`
	Note           string   `json:"note"`
}

type RawStep struct {
	OMCurrentVersion             string           `json:"omCurrentVersion"`
	OMTargetVersion              string           `json:"omTargetVersion"`
	OMHosts                      []string         `json:"omHosts"`
	OSUpgradeRequirements        []RawRequirement `json:"osUpgradeRequirements"`
	MongoUpgradeRequirements     []RawRequirement `json:"mongoUpgradeRequirements"`
	AgentPostUpgradeRequirements []RawRequirement `json:"agentPostUpgradeRequirements"`
}

type UpgradeCheckStep struct {
	OpsManager      OpsManagerStep        `json:"ops_manager,omitempty"`
	OperatingSystem []OperatingSystemStep `json:"operating_system,omitempty"`
	MongoDB         []MongoDBStep         `json:"mongodb,omitempty"`
	Agent           []AgentStep           `json:"agent,omitempty"`
}

type OpsManagerStep struct {
	CurrentVersion string   `json:"current_version,omitempty"`
	TargetVersion  string   `json:"target_version,omitempty"`
	Hosts          []string `json:"hosts,omitempty"`
}
type OperatingSystemStep struct {
	BaseVersion    string   `json:"base_version,omitempty"`
	CurrentVersion string   `json:"current_version,omitempty"`
	TargetVersion  string   `json:"target_version,omitempty"`
	Hosts          []string `json:"hosts,omitempty"`
}
type MongoDBStep struct {
	CurrentVersion string   `json:"current_version,omitempty"`
	TargetVersion  string   `json:"target_version,omitempty"`
	Hosts          []string `json:"hosts,omitempty"`
}
type AgentStep struct {
	CurrentVersion string   `json:"current_version,omitempty"`
	TargetVersion  string   `json:"target_version,omitempty"`
	Hosts          []string `json:"hosts,omitempty"`
}

func (s *AdvisorServiceOp) CheckUpgrade(ctx context.Context, version string) (*[]UpgradeCheckStep, *Response, error) {
	if version == "" {
		return nil, nil, NewArgError("targetVersion", "must be set")
	}

	path := fmt.Sprintf("%s/%s", upgradeCheckPath, version)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(UpgradeCheckResult)
	resp, err := s.Client.Do(ctx, req, root)

	if err != nil {
		return nil, resp, err
	}

	steps, err := parseUpgradeSteps([]byte(root.Result))
	return &steps, resp, err
}

func fold(reqs []RawRequirement) (cur, tgt string, hosts []string) {
	if len(reqs) == 0 {
		return
	}
	cur, tgt = reqs[0].CurrentVersion, reqs[0].TargetVersion
	for _, r := range reqs {
		hosts = append(hosts, r.Hostnames...)
	}
	return
}

func groupOS(reqs []RawRequirement) []OperatingSystemStep {
	type triple struct{ base, cur, tgt string }
	tmp := map[triple][]string{}
	for _, r := range reqs {
		base := ""
		if r.BaseVersion != nil {
			base = *r.BaseVersion
		}
		t := triple{base, r.CurrentVersion, r.TargetVersion}
		tmp[t] = append(tmp[t], r.Hostnames...)
	}
	out := make([]OperatingSystemStep, 0, len(tmp))
	for k, hosts := range tmp {
		out = append(out, OperatingSystemStep{
			BaseVersion:    k.base,
			CurrentVersion: k.cur,
			TargetVersion:  k.tgt,
			Hosts:          hosts,
		})
	}
	return out
}

func groupMongo(reqs []RawRequirement) []MongoDBStep {
	type pair struct{ cur, tgt string }
	m := map[pair][]string{}
	for _, r := range reqs {
		p := pair{r.CurrentVersion, r.TargetVersion}
		m[p] = append(m[p], r.Hostnames...)
	}
	out := make([]MongoDBStep, 0, len(m))
	for p, hosts := range m {
		out = append(out, MongoDBStep{p.cur, p.tgt, hosts})
	}
	return out
}

func groupAgent(reqs []RawRequirement) []AgentStep {
	type pair struct{ cur, tgt string }
	m := map[pair][]string{}
	for _, r := range reqs {
		p := pair{r.CurrentVersion, r.TargetVersion}
		m[p] = append(m[p], r.Hostnames...)
	}
	out := make([]AgentStep, 0, len(m))
	for p, hosts := range m {
		out = append(out, AgentStep{p.cur, p.tgt, hosts})
	}
	return out
}

func parseUpgradeSteps(data []byte) ([]UpgradeCheckStep, error) {
	var raw []RawStep
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	steps := make([]UpgradeCheckStep, 0, len(raw))
	for _, r := range raw {
		steps = append(steps, UpgradeCheckStep{
			OpsManager: OpsManagerStep{
				CurrentVersion: r.OMCurrentVersion,
				TargetVersion:  r.OMTargetVersion,
				Hosts:          r.OMHosts,
			},
			OperatingSystem: groupOS(r.OSUpgradeRequirements),
			MongoDB:         groupMongo(r.MongoUpgradeRequirements),
			Agent:           groupAgent(r.AgentPostUpgradeRequirements),
		})
	}
	return steps, nil
}
