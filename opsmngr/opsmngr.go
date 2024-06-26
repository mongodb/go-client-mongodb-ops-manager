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
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	DefaultBaseURL = "https://cloud.mongodb.com/"
	userAgent      = "go-ops-manager"
	jsonMediaType  = "application/json"
	gzipMediaType  = "application/gzip"
	plainMediaType = "text/plain"
	// ClientVersion of the current API client. Should be set to the next version planned to be released.
	ClientVersion = "0.56.0"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Doer basic interface of a client to be able to do a request.
type Doer interface {
	Do(context.Context, *http.Request, interface{}) (*Response, error)
}

// Completer interface for clients with callback.
type Completer interface {
	OnRequestCompleted(RequestCompletionCallback)
}

// RequestDoer minimum interface for any service of the client.
type RequestDoer interface {
	Doer
	Completer
	NewRequest(context.Context, string, string, interface{}) (*http.Request, error)
}

// GZipRequestDoer minimum interface for any service of the client that should handle gzip downloads.
type GZipRequestDoer interface {
	Doer
	Completer
	NewGZipRequest(context.Context, string, string) (*http.Request, error)
}

// PlainRequestDoer minimum interface for any service of the client that should handle plain text.
type PlainRequestDoer interface {
	Doer
	Completer
	NewPlainRequest(context.Context, string, string) (*http.Request, error)
}

// Client manages communication with Ops Manager API.
type Client struct {
	client    HTTPClient
	BaseURL   *url.URL
	UserAgent string

	// copy raw server response to the Response struct
	withRaw bool

	Organizations          OrganizationsService
	Projects               ProjectsService
	Users                  UsersService
	Teams                  TeamsService
	Automation             AutomationService
	UnauthUsers            UnauthUsersService
	AlertConfigurations    AlertConfigurationsService
	Alerts                 AlertsService
	ContinuousSnapshots    ContinuousSnapshotsService
	ContinuousRestoreJobs  ContinuousRestoreJobsService
	Events                 EventsService
	OrganizationAPIKeys    APIKeysService
	ProjectAPIKeys         ProjectAPIKeysService
	AccessListAPIKeys      AccessListAPIKeysService
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
	LiveMigration          LiveDataMigrationService
	ServiceVersion         ServiceVersionService

	onRequestCompleted  RequestCompletionCallback
	onResponseProcessed ResponseProcessedCallback
}

// RequestCompletionCallback defines the type of the request callback function.
type RequestCompletionCallback func(*http.Request, *http.Response)

// ResponseProcessedCallback defines the type of the after request completion callback function.
type ResponseProcessedCallback func(*Response)

type service struct {
	Client RequestDoer
}

// Response is a MongoDB Ops Manager response. This wraps the standard http.Response returned from MongoDB Ops Manager  API.
type Response struct {
	*http.Response

	// Links that were returned with the response.
	Links []*Link `json:"links"`

	// Raw data from server response
	Raw []byte `json:"-"`
}

// ListOptions specifies the optional parameters to List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	PageNum int `url:"pageNum,omitempty"`

	// For paginated result sets, the number of results to include per page.
	ItemsPerPage int `url:"itemsPerPage,omitempty"`

	// Flag that indicates whether Ops Manager returns the totalCount parameter in the response body.
	IncludeCount bool `url:"includeCount,omitempty"`
}

func (resp *Response) getLinkByRef(ref string) *Link {
	for i := range resp.Links {
		if resp.Links[i].Rel == ref {
			return resp.Links[i]
		}
	}
	return nil
}

func (resp *Response) getCurrentPageLink() (*Link, error) {
	if link := resp.getLinkByRef("self"); link != nil {
		return link, nil
	}
	return nil, errors.New("no self link found")
}

// CurrentPage gets the current page for list pagination request.
func (resp *Response) CurrentPage() (int, error) {
	link, err := resp.getCurrentPageLink()
	if err != nil {
		return 0, err
	}

	pageNumStr, err := link.getHrefQueryParam("pageNum")
	if err != nil {
		return 0, err
	}

	pageNum, err := strconv.Atoi(pageNumStr)
	if err != nil {
		return 0, fmt.Errorf("error getting current page: %w", err)
	}

	return pageNum, nil
}

// NewClient returns a new Ops Manager API client. If a nil httpClient is
// provided, a http.DefaultClient will be used. To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the https://github.com/mongodb-forks/digest).
func NewClient(httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(DefaultBaseURL)

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
	c.AlertConfigurations = &AlertConfigurationsServiceOp{Client: c}
	c.UnauthUsers = &UnauthUsersServiceOp{Client: c}
	c.ContinuousSnapshots = &ContinuousSnapshotsServiceOp{Client: c}
	c.ContinuousRestoreJobs = &ContinuousRestoreJobsServiceOp{Client: c}
	c.Agents = &AgentsServiceOp{Client: c}
	c.Checkpoints = &CheckpointsServiceOp{Client: c}
	c.Alerts = &AlertsServiceOp{Client: c}
	c.GlobalAlerts = &GlobalAlertsServiceOp{Client: c}
	c.Events = &EventsServiceOp{Client: c}
	c.Deployments = &DeploymentsServiceOp{Client: c}
	c.Measurements = &MeasurementsServiceOp{Client: c}
	c.Clusters = &ClustersServiceOp{Client: c}
	c.Logs = &LogsServiceOp{Client: c}
	c.LogCollections = &LogCollectionServiceOp{Client: c}
	c.Diagnostics = &DiagnosticsServiceOp{Client: c}
	c.OrganizationAPIKeys = &APIKeysServiceOp{Client: c}
	c.ProjectAPIKeys = &ProjectAPIKeysOp{Client: c}
	c.AccessListAPIKeys = &AccessListAPIKeysServiceOp{Client: c}
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
	c.LiveMigration = &LiveDataMigrationServiceOp{Client: c}
	c.ServiceVersion = &ServiceVersionServiceOp{Client: c}

	return c
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// Options turns a list of ClientOpt instances into a ClientOpt.
func Options(opts ...ClientOpt) ClientOpt {
	return func(c *Client) error {
		for _, opt := range opts {
			if err := opt(c); err != nil {
				return err
			}
		}
		return nil
	}
}

// New returns a new Ops Manager API client instance.
func New(httpClient HTTPClient, opts ...ClientOpt) (*Client, error) {
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

// SetWithRaw is a client option for getting raw Ops Manager server response within Response structure.
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
	req, err := c.newRequest(ctx, urlStr, method, body)
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
	req, err := c.newRequest(ctx, urlStr, method, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", gzipMediaType)

	return req, nil
}

// NewPlainRequest creates an API request that accepts plain text.
// A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash.
func (c *Client) NewPlainRequest(ctx context.Context, method, urlStr string) (*http.Request, error) {
	req, err := c.newRequest(ctx, urlStr, method, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", plainMediaType)

	return req, nil
}

func (c *Client) newRequest(ctx context.Context, urlStr, method string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		c.BaseURL.Path += "/"
	}
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.Reader
	if body != nil {
		if buf, err = newEncodedBody(body); err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// newEncodedBody returns an ReadWriter object containing the body of the http request.
func newEncodedBody(body interface{}) (io.Reader, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(body)
	return buf, err
}

// OnRequestCompleted sets the DO API request completion callback.
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// OnResponseProcessed sets the DO API request completion callback.
func (c *Client) OnResponseProcessed(rc ResponseProcessedCallback) {
	c.onResponseProcessed = rc
}

// ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	// Response that caused this error
	Response *http.Response
	// ErrorCode is the error code
	ErrorCode string `json:"errorCode"`
	// HTTPCode status code.
	HTTPCode int `json:"error"` //nolint:tagliatelle // used as in the API
	// Reason is short description of the error, which is simply the HTTP status phrase.
	Reason string `json:"reason"`
	// Detail is more detailed description of the error.
	Detail string `json:"detail,omitempty"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d (request %q) %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.ErrorCode, r.Detail)
}

func (r *ErrorResponse) Is(target error) bool {
	var v *ErrorResponse

	return errors.As(target, &v) &&
		r.ErrorCode == v.ErrorCode &&
		r.HTTPCode == v.HTTPCode &&
		r.Reason == v.Reason &&
		r.Detail == v.Detail
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func (resp *Response) CheckResponse(body io.ReadCloser) error {
	if c := resp.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: resp.Response}
	data, err := io.ReadAll(body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			log.Printf("[DEBUG] unmarshal error response: %s", err)
			errorResponse.Reason = string(data)
		}
	}

	return errorResponse
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
// The provided ctx must be non-nil, if it is nil an error is returned. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
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

	defer resp.Body.Close()

	response := &Response{Response: resp}

	defer func() {
		if c.onResponseProcessed != nil {
			c.onResponseProcessed(response)
		}
	}()

	body := resp.Body

	if c.withRaw {
		raw := new(bytes.Buffer)
		_, err = io.Copy(raw, body)
		if err != nil {
			return response, err
		}

		response.Raw = raw.Bytes()
		body = io.NopCloser(raw)
	}

	err = response.CheckResponse(body)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, _ = io.Copy(w, body)
		} else {
			decErr := json.NewDecoder(body).Decode(v)
			if errors.Is(decErr, io.EOF) {
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

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(
	ctx context.Context,
	client *http.Client,
	req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

// ServiceVersion represents version information.
type ServiceVersion struct {
	GitHash string
	Version string
}

// String serializes VersionInfo into string.
func (v *ServiceVersion) String() string {
	return fmt.Sprintf("gitHash=%s; versionString=%s", v.GitHash, v.Version)
}

func parseVersionInfo(s string) *ServiceVersion {
	if s == "" {
		return nil
	}

	var result ServiceVersion
	pairs := strings.Split(s, ";")
	for _, pair := range pairs {
		keyvalue := strings.Split(strings.TrimSpace(pair), "=")
		switch keyvalue[0] {
		case "gitHash":
			result.GitHash = keyvalue[1]
		case "versionString":
			result.Version = keyvalue[1]
		}
	}
	return &result
}

// ServiceVersion parses version information returned in the response.
func (resp *Response) ServiceVersion() *ServiceVersion {
	return parseVersionInfo(resp.Header.Get("X-MongoDB-Service-Version"))
}

func pointer[T any](x T) *T {
	return &x
}
