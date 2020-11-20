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

package opsmngr // import "go.mongodb.org/ops-manager/opsmngr"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const (
	cloudURL       = "https://cloud.mongodb.com/"
	defaultBaseURL = cloudURL + APIPublicV1Path
	userAgent      = "go-ops-manager"
	jsonMediaType  = "application/json"
	gzipMediaType  = "application/gzip"
	// APIPublicV1Path specifies the v1 api path
	APIPublicV1Path = "api/public/v1.0/"
)

// Client manages communication with Ops Manager API
type Client struct {
	client    *http.Client
	BaseURL   *url.URL
	UserAgent string

	// copy raw atlas server response to the Response struct
	withRaw bool

	Organizations          OrganizationsService
	Projects               ProjectsService
	Users                  UsersService
	Teams                  TeamsService
	Automation             AutomationService
	UnauthUsers            UnauthUsersService
	AlertConfigurations    atlas.AlertConfigurationsService
	Alerts                 atlas.AlertsService
	ContinuousSnapshots    atlas.ContinuousSnapshotsService
	ContinuousRestoreJobs  atlas.ContinuousRestoreJobsService
	Events                 atlas.EventsService
	OrganizationAPIKeys    atlas.APIKeysService
	ProjectAPIKeys         atlas.ProjectAPIKeysService
	WhitelistAPIKeys       atlas.WhitelistAPIKeysService
	Agents                 AgentsService
	Checkpoints            CheckpointsService
	GlobalAlerts           GlobalAlertsService
	Deployments            DeploymentsService
	Measurements           MeasurementsService
	Clusters               ClustersService
	Logs                   LogsService
	LogCollections         LogCollectionService
	Diagnostics            DiagnosticsService
	GlobalAPIKeys          GlobalAPIKeysService
	GlobalAPIKeysWhitelist GlobalAPIKeyWhitelistsService
	MaintenanceWindows     MaintenanceWindowsService
	PerformanceAdvisor     PerformanceAdvisorService
	VersionManifest        VersionManifestService
	BackupConfigs          BackupConfigsService
	ProjectJobConfig       ProjectJobConfigService
	BlockstoreConfig       BlockstoreConfigService
	FileSystemStoreConfig  FileSystemStoreConfigService
	S3BlockstoreConfig     S3BlockstoreConfigService
	OplogStoreConfig       OplogStoreConfigService
	SyncStoreConfig        SyncStoreConfigService
	DaemonConfig           DaemonConfigService
	SnapshotSchedule       SnapshotScheduleService
	FeatureControlPolicies FeatureControlPoliciesService
	ServerUsage            ServerUsageService
	ServerUsageReport      ServerUsageReportService

	onRequestCompleted atlas.RequestCompletionCallback
}

type service struct {
	Client atlas.RequestDoer
}

// NewClient returns a new Ops Manager API client. If a nil httpClient is
// provided, a http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the https://github.com/mongodb-forks/digest).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}

	c.Organizations = &OrganizationsServiceOp{Client: c}
	c.Projects = &ProjectsServiceOp{Client: c}
	c.Users = &UsersServiceOp{Client: c}
	c.Teams = &TeamsServiceOp{Client: c}
	c.Automation = &AutomationServiceOp{Client: c}
	c.AlertConfigurations = &atlas.AlertConfigurationsServiceOp{Client: c}
	c.UnauthUsers = &UnauthUsersServiceOp{Client: c}
	c.ContinuousSnapshots = &atlas.ContinuousSnapshotsServiceOp{Client: c}
	c.ContinuousRestoreJobs = &atlas.ContinuousRestoreJobsServiceOp{Client: c}
	c.Agents = &AgentsServiceOp{Client: c}
	c.Checkpoints = &CheckpointsServiceOp{Client: c}
	c.Alerts = &atlas.AlertsServiceOp{Client: c}
	c.GlobalAlerts = &GlobalAlertsServiceOp{Client: c}
	c.Events = &atlas.EventsServiceOp{Client: c}
	c.Deployments = &DeploymentsServiceOp{Client: c}
	c.Measurements = &MeasurementsServiceOp{Client: c}
	c.Clusters = &ClustersServiceOp{Client: c}
	c.Logs = &LogsServiceOp{Client: c}
	c.LogCollections = &LogCollectionServiceOp{Client: c}
	c.Diagnostics = &DiagnosticsServiceOp{Client: c}
	c.OrganizationAPIKeys = &atlas.APIKeysServiceOp{Client: c}
	c.ProjectAPIKeys = &atlas.ProjectAPIKeysOp{Client: c}
	c.WhitelistAPIKeys = &atlas.WhitelistAPIKeysServiceOp{Client: c}
	c.GlobalAPIKeys = &GlobalAPIKeysServiceOp{Client: c}
	c.GlobalAPIKeysWhitelist = &GlobalAPIKeyWhitelistsServiceOp{Client: c}
	c.MaintenanceWindows = &MaintenanceWindowsServiceOp{Client: c}
	c.PerformanceAdvisor = &PerformanceAdvisorServiceOp{Client: c}
	c.VersionManifest = &VersionManifestServiceOp{Client: c}
	c.BackupConfigs = &BackupConfigsServiceOp{Client: c}
	c.ProjectJobConfig = &ProjectJobConfigServiceOp{Client: c}
	c.BlockstoreConfig = &BlockstoreConfigServiceOp{Client: c}
	c.FileSystemStoreConfig = &FileSystemStoreConfigServiceOp{Client: c}
	c.S3BlockstoreConfig = &S3BlockstoreConfigServiceOp{Client: c}
	c.OplogStoreConfig = &OplogStoreConfigServiceOp{Client: c}
	c.SyncStoreConfig = &SyncStoreConfigServiceOp{Client: c}
	c.DaemonConfig = &DaemonConfigServiceOp{Client: c}
	c.SnapshotSchedule = &SnapshotScheduleServiceOp{Client: c}
	c.FeatureControlPolicies = &FeatureControlPoliciesServiceOp{Client: c}
	c.ServerUsage = &ServerUsageServiceOp{Client: c}
	c.ServerUsageReport = &ServerUsageReportServiceOp{Client: c}

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

// SetWithRaw is a client option for getting raw atlas server response within Response structure.
func SetWithRaw() ClientOpt {
	return func(c *Client) error {
		c.withRaw = true
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

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("base URL must have a trailing slash, but %q does not", c.BaseURL)
	}
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err = enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", jsonMediaType)
	}
	req.Header.Add("Accept", jsonMediaType)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// NewGZipRequest creates an API gzip request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewGZipRequest(ctx context.Context, method, urlStr string) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf("base URL must have a trailing slash, but %q does not", c.BaseURL)
	}
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", gzipMediaType)
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
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*atlas.Response, error) {
	if ctx == nil {
		return nil, errors.New("context must be non-nil")
	}

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
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

	defer func() {
		// Ensure the response body is fully read and closed
		// before we reconnect, so that we reuse the same TCP connection.
		// Close the previous response's body. But read at least some of
		// the body so if it's small the underlying TCP connection will be
		// re-used. No need to check for errors: if it fails, the Transport
		// won't reuse it anyway.
		const maxBodySlurpSize = 2 << 10
		if resp.ContentLength == -1 || resp.ContentLength <= maxBodySlurpSize {
			_, _ = io.CopyN(ioutil.Discard, resp.Body, maxBodySlurpSize)
		}

		resp.Body.Close()
	}()

	response := &atlas.Response{Response: resp}

	err = atlas.CheckResponse(resp)
	if err != nil {
		return response, err
	}

	body := resp.Body

	if c.withRaw {
		raw := new(bytes.Buffer)
		_, err = io.Copy(raw, body)
		if err != nil {
			return response, err
		}

		response.Raw = raw.Bytes()
		body = ioutil.NopCloser(raw)
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, _ = io.Copy(w, body)
		} else {
			decErr := json.NewDecoder(body).Decode(v)
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
