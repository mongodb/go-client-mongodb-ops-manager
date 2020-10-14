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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	availablePoliciesBasePath  = "groups/availablePolicies"
	controlledFeaturesBasePath = "groups/%s/controlledFeature"
)

// FeatureControlPoliciesService provides access to the Feature Control Policies related functions in the Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/feature-control-policies/
type FeatureControlPoliciesService interface {
	List(context.Context, string, *atlas.ListOptions) (*FeaturePolicy, *atlas.Response, error)
	Update(context.Context, string, *FeaturePolicy) (*FeaturePolicy, *atlas.Response, error)
	ListSupportedPolicies(context.Context, *atlas.ListOptions) (*FeaturePolicy, *atlas.Response, error)
}

// AgentsServiceOp provides an implementation of the AgentsService interface
type FeatureControlPoliciesServiceOp service

var _ FeatureControlPoliciesService = new(FeatureControlPoliciesServiceOp)

// Agent represents an Ops Manager agent
type FeaturePolicy struct {
	Created                  string                    `json:"created,omitempty"`
	Updated                  string                    `json:"updated,omitempty"`
	ExternalManagementSystem *ExternalManagementSystem `json:"externalManagementSystem,omitempty"`
	Policies                 []*Policy                 `json:"policies,omitempty"`
}

// ExternalManagementSystem contains parameters for the external system that manages this Ops Manager Project
type ExternalManagementSystem struct {
	Name     string `json:"name"`
	SystemID string `json:"systemId,omitempty"`
	Version  string `json:"version,omitempty"`
}

// Policy contains policies that the external system applies to this Ops Manager Project
type Policy struct {
	Policy         string   `json:"policy,omitempty"`
	DisabledParams []string `json:"disabledParams,omitempty"`
}

// List retrieves the policies that have been set on a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/controlled-features/get-controlled-features-for-one-project/
func (s *FeatureControlPoliciesServiceOp) List(ctx context.Context, groupID string, opts *atlas.ListOptions) (*FeaturePolicy, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(controlledFeaturesBasePath, groupID)

	path, err := setQueryParams(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FeaturePolicy)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Update updates the Feature Control Policies for one Project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/controlled-features/update-controlled-features-for-one-project/
func (s *FeatureControlPoliciesServiceOp) Update(ctx context.Context, groupID string, policy *FeaturePolicy) (*FeaturePolicy, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	path := fmt.Sprintf(controlledFeaturesBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, path, policy)
	if err != nil {
		return nil, nil, err
	}

	root := new(FeaturePolicy)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// ListSupportedPolicies retrieves all supported policies by Ops Manager.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/controlled-features/get-all-feature-control-policies/
func (s *FeatureControlPoliciesServiceOp) ListSupportedPolicies(ctx context.Context, opts *atlas.ListOptions) (*FeaturePolicy, *atlas.Response, error) {
	path, err := setQueryParams(availablePoliciesBasePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(FeaturePolicy)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}