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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime"

	"github.com/google/go-querystring/query"
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
	Alerts                atlas.AlertsService
	ContinuousSnapshots   atlas.ContinuousSnapshotsService
	ContinuousRestoreJobs atlas.ContinuousRestoreJobsService
	Events                atlas.EventsService
	AllClusters           AllClustersService
	Agents                AgentsService
	AgentAPIKeys          AgentAPIKeysService
	Checkpoints           CheckpointsService
	GlobalAlerts          GlobalAlertsService
	SystemMeasurements    SystemMeasurementsService

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
	c.AllClusters = &AllClustersServiceOp{Client: c}
	c.ContinuousSnapshots = &atlas.ContinuousSnapshotsServiceOp{Client: c}
	c.ContinuousRestoreJobs = &atlas.ContinuousRestoreJobsServiceOp{Client: c}
	c.Agents = &AgentsServiceOp{Client: c}
	c.AgentAPIKeys = &AgentAPIKeysServiceOp{Client: c}
	c.Checkpoints = &CheckpointsServiceOp{Client: c}
	c.Alerts = &atlas.AlertsServiceOp{Client: c}
	c.GlobalAlerts = &GlobalAlertsServiceOp{Client: c}
	c.Events = &atlas.EventsServiceOp{Client: c}
	c.SystemMeasurements = &ServiceMeasurementsServiceOp{Client: c}

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

// OptionSkipVerify will set the Insecure Skip which means that TLS certs will not be
// verified for validity.
func OptionSkipVerify(c *Client) error {
	if c.client.Transport == nil {
		c.client.Transport = http.DefaultTransport
	}
	transport, ok := c.client.Transport.(*http.Transport)
	if !ok {
		return errors.New("client.Transport not an *http.Transport")
	}
	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	c.client.Transport = transport

	return nil
}

// OptionCAValidate will use the CA certificate, passed as a string, to validate the
// certificates presented by Ops Manager.
func OptionCAValidate(ca string) ClientOpt {
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM([]byte(ca))
	TLSClientConfig := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            caCertPool,
	}

	return func(c *Client) error {
		if c.client.Transport == nil {
			c.client.Transport = http.DefaultTransport
		}
		transport, ok := c.client.Transport.(*http.Transport)
		if !ok {
			return errors.New("client.Transport not an *http.Transport")
		}
		transport.TLSClientConfig = TLSClientConfig
		c.client.Transport = transport

		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", mediaType)
	}
	req.Header.Add("Accept", mediaType)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
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
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err)
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer resp.Body.Close()

	response := &atlas.Response{Response: resp}

	err = atlas.CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, _ = io.Copy(w, resp.Body)
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}
	return response, err
}

func setQueryParams(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}
