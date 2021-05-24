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

	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	performanceAdvisorPath                     = "groups/%s/hosts/%s/performanceAdvisor"
	performanceAdvisorNamespacesPath           = performanceAdvisorPath + "/namespaces"
	performanceAdvisorSlowQueryLogsPath        = performanceAdvisorPath + "/slowQueryLogs"
	performanceAdvisorSuggestedIndexesLogsPath = performanceAdvisorPath + "/suggestedIndexes"
)

// PerformanceAdvisorService is an interface of the Performance Advisor
// endpoints of the MongoDB Ops Manager API.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/performance-advisor/index.html
type PerformanceAdvisorService interface {
	GetNamespaces(context.Context, string, string, *atlas.NamespaceOptions) (*atlas.Namespaces, *atlas.Response, error)
	GetSlowQueries(context.Context, string, string, *atlas.SlowQueryOptions) (*atlas.SlowQueries, *atlas.Response, error)
	GetSuggestedIndexes(context.Context, string, string, *atlas.SuggestedIndexOptions) (*atlas.SuggestedIndexes, *atlas.Response, error)
}

// PerformanceAdvisorServiceOp handles communication with the Performance Advisor related methods of the MongoDB Atlas API.
type PerformanceAdvisorServiceOp service

var _ PerformanceAdvisorService = &PerformanceAdvisorServiceOp{}

// GetNamespaces retrieves the namespaces for collections experiencing slow queries for a specified host.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/performance-advisor/pa-namespaces-get-all/
func (s *PerformanceAdvisorServiceOp) GetNamespaces(ctx context.Context, groupID, processName string, opts *atlas.NamespaceOptions) (*atlas.Namespaces, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if processName == "" {
		return nil, nil, atlas.NewArgError("processName", "must be set")
	}

	path := fmt.Sprintf(performanceAdvisorNamespacesPath, groupID, processName)
	path, err := setQueryParams(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.Namespaces)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetSlowQueries gets log lines for slow queries as determined by the Performance Advisor.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/performance-advisor/get-slow-queries/
func (s *PerformanceAdvisorServiceOp) GetSlowQueries(ctx context.Context, groupID, processName string, opts *atlas.SlowQueryOptions) (*atlas.SlowQueries, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if processName == "" {
		return nil, nil, atlas.NewArgError("processName", "must be set")
	}

	path := fmt.Sprintf(performanceAdvisorSlowQueryLogsPath, groupID, processName)
	path, err := setQueryParams(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.SlowQueries)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// GetSuggestedIndexes gets suggested indexes as determined by the Performance Advisor.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/performance-advisor/get-suggested-indexes/
func (s *PerformanceAdvisorServiceOp) GetSuggestedIndexes(ctx context.Context, groupID, processName string, opts *atlas.SuggestedIndexOptions) (*atlas.SuggestedIndexes, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	if processName == "" {
		return nil, nil, atlas.NewArgError("processName", "must be set")
	}

	path := fmt.Sprintf(performanceAdvisorSuggestedIndexesLogsPath, groupID, processName)
	path, err := setQueryParams(path, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(atlas.SuggestedIndexes)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
