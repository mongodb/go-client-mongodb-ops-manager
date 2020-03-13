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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	Version          = "0.1" // Version for client
	CloudURL         = "https://cloud.mongodb.com"
	DefaultBaseURL   = CloudURL + APIPublicV1Path                                                             // DefaultBaseURL API default base URL for cloud manager
	APIPublicV1Path  = "/api/public/v1.0/"                                                                    // DefaultAPIPath default root path for all API endpoints
	DefaultUserAgent = "go-client-ops-manager/" + Version + " (" + runtime.GOOS + "; " + runtime.GOARCH + ")" // DefaultUserAgent To be submitted by the client
	mediaType        = "application/json"
)

// Client manages communication with v1.0 API
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string

	Organizations         OrganizationsService
	Projects              ProjectsService
	AutomationConfig      AutomationConfigService
	AutomationStatus      AutomationStatusService
	UnauthUsers           UnauthUsersService
	AlertConfigurations   atlas.AlertConfigurationsService
	ContinuousSnapshots   atlas.ContinuousSnapshotsService
	ContinuousRestoreJobs atlas.ContinuousRestoreJobsService
	AllCusters            AllClustersService
	OpsManagerCheckpoints CheckpointsService

	onRequestCompleted atlas.RequestCompletionCallback
}

// NewClient returns a new Ops Manager API Client
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(DefaultBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: DefaultUserAgent,
	}

	c.Organizations = &OrganizationsServiceOp{Client: c}
	c.Projects = &ProjectsServiceOp{Client: c}
	c.AutomationConfig = &AutomationConfigServiceOp{Client: c}
	c.AutomationStatus = &AutomationStatusServiceOp{Client: c}
	c.AlertConfigurations = &atlas.AlertConfigurationsServiceOp{Client: c}
	c.UnauthUsers = &UnauthUsersServiceOp{Client: c}
	c.AllCusters = &AllClustersServiceOp{Client: c}
	c.ContinuousSnapshots = &atlas.ContinuousSnapshotsServiceOp{Client: c}
	c.ContinuousRestoreJobs = &atlas.ContinuousRestoreJobsServiceOp{Client: c}
	c.OpsManagerCheckpoints = &CheckpointsServiceOp{Client: c}

	return c
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new Ops Manager API client instance.
func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	c := NewClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// OnRequestCompleted sets the DO API request completion callback
func (c *Client) OnRequestCompleted(rc atlas.RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*atlas.Response, error) {
	resp, err := atlas.DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := &atlas.Response{Response: resp}

	err = atlas.CheckResponse(resp)
	if err != nil {
		return response, err
	}

	switch t := v.(type) {
	case nil:
	case io.Writer:
		_, err = io.Copy(t, resp.Body)
	default:
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return response, err
}
