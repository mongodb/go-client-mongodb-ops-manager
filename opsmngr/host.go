package opsmngr

import (
	"context"
	"fmt"
	"net/http"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	hostBasePath = "groups/%s/hosts"
)

// HostService is an interface for interfacing with Hosts in MongoDB Ops Manager APIs
type HostService interface {
	Get(context.Context, string, string) (*Host, *atlas.Response, error)
	GetByHostname(context.Context, string, string, int) (*Host, *atlas.Response, error)
	List(context.Context, string, *HostListOptions) (*Hosts, *atlas.Response, error)
	Monitoring(context.Context, string, *Host) (*Host, *atlas.Response, error)
	UpdateMonitoring(context.Context, string, string, *Host) (*Host, *atlas.Response, error)
	StopMonitoring(context.Context, string, string) (*atlas.Response, error)
}

type HostServiceOp struct {
	Client atlas.RequestDoer
}

type Host struct {
	Aliases            []string      `json:"aliases,omitempty"`
	AuthMechanismName  string        `json:"authMechanismName,omitempty"`
	ClusterID          string        `json:"clusterId,omitempty"`
	Created            string        `json:"created,omitempty"`
	GroupID            string        `json:"groupId,omitempty"`
	Hostname           string        `json:"hostname"`
	ID                 string        `json:"id,omitempty"`
	IPAddress          string        `json:"ipAddress,omitempty"`
	LastPing           string        `json:"lastPing,omitempty"`
	LastRestart        string        `json:"lastRestart,omitempty"`
	ReplicaSetName     string        `json:"replicaSetName,omitempty"`
	ReplicaStateName   string        `json:"replicaStateName,omitempty"`
	ShardName          string        `json:"shardName,omitempty"`
	TypeName           string        `json:"typeName,omitempty"`
	Version            string        `json:"version,omitempty"`
	Username           string        `json:"username,omitempty"`
	Password           string        `json:"password,omitempty"`
	Deactivated        bool          `json:"deactivated,omitempty"`
	HasStartupWarnings bool          `json:"hasStartupWarnings,omitempty"`
	Hidden             bool          `json:"hidden,omitempty"`
	HiddenSecondary    bool          `json:"hiddenSecondary,omitempty"`
	HostEnabled        bool          `json:"hostEnabled,omitempty"`
	JournalingEnabled  bool          `json:"journalingEnabled,omitempty"`
	LowUlimit          bool          `json:"lowUlimit,omitempty"`
	MuninEnabled       bool          `json:"muninEnabled,omitempty"`
	LogsEnabled        *bool         `json:"logsEnabled,omitempty"`
	AlertsEnabled      *bool         `json:"alertsEnabled,omitempty"`
	ProfilerEnabled    *bool         `json:"profilerEnabled,omitempty"`
	SSLEnabled         *bool         `json:"sslEnabled,omitempty"`
	LastDataSizeBytes  float64       `json:"lastDataSizeBytes,omitempty"`
	LastIndexSizeBytes float64       `json:"lastIndexSizeBytes,omitempty"`
	MuninPort          *int32        `json:"muninPort,omitempty"`
	Port               int32         `json:"port"`
	SlaveDelaySec      int64         `json:"slaveDelaySec,omitempty"`
	UptimeMsec         int64         `json:"uptimeMsec,omitempty"`
	Links              []*atlas.Link `json:"links,omitempty"`
}

type Hosts struct {
	Links      []*atlas.Link `json:"links"`
	Results    []*Host       `json:"results"`
	TotalCount int           `json:"totalCount"`
}

type HostListOptions struct {
	atlas.ListOptions
	ClusterID string `url:"clusterId,omitempty"`
}

// Get gets the MongoDB process with the specified host ID.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/hosts/get-one-host-by-id/
func (s *HostServiceOp) Get(ctx context.Context, projectID, hostID string) (*Host, *atlas.Response, error) {

	basePath := fmt.Sprintf(hostBasePath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, hostID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Host)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// GetByHostname gets a single MongoDB process by its hostname and port combination.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/hosts/get-one-host-by-id/
func (s *HostServiceOp) GetByHostname(ctx context.Context, projectID, hostname string, port int) (*Host, *atlas.Response, error) {

	basePath := fmt.Sprintf(hostBasePath, projectID)
	path := fmt.Sprintf("%s/byName/%s:%d", basePath, hostname, port)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Host)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// List lists all MongoDB hosts in a project.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/hosts/get-all-hosts-in-group/
func (s *HostServiceOp) List(ctx context.Context, projectID string, opts *HostListOptions) (*Hosts, *atlas.Response, error) {

	basePath := fmt.Sprintf(hostBasePath, projectID)
	path, err := setQueryParams(basePath, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(Hosts)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// Monitoring starts monitoring a new MongoDB process.
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/hosts/create-one-host/
func (s *HostServiceOp) Monitoring(ctx context.Context, projectID string, host *Host) (*Host, *atlas.Response, error) {
	path := fmt.Sprintf(hostBasePath, projectID)

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, host)
	if err != nil {
		return nil, nil, err
	}

	root := new(Host)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// UpdateMonitoring updates the configuration of a monitored MongoDB process..
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/hosts/update-one-host/
func (s *HostServiceOp) UpdateMonitoring(ctx context.Context, projectID, hostID string, host *Host) (*Host, *atlas.Response, error) {
	basePath := fmt.Sprintf(hostBasePath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, hostID)
	req, err := s.Client.NewRequest(ctx, http.MethodPatch, path, host)
	if err != nil {
		return nil, nil, err
	}

	root := new(Host)
	resp, err := s.Client.Do(ctx, req, root)

	return root, resp, err
}

// StopMonitoring stops the Monitoring from monitoring the MongoDB process on the hostname and port you specify..
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/hosts/delete-one-host/
func (s *HostServiceOp) StopMonitoring(ctx context.Context, projectID, hostID string) (*atlas.Response, error) {
	basePath := fmt.Sprintf(hostBasePath, projectID)
	path := fmt.Sprintf("%s/%s", basePath, hostID)

	req, err := s.Client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}
