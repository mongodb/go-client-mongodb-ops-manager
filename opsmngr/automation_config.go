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

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

const (
	automationConfigBasePath = "groups/%s/automationConfig"
)

// AutomationConfigService is an interface for interfacing with the Automation Config
// endpoints of the MongoDB CLoud API.
// See more: https://docs.cloudmanager.mongodb.com/reference/api/automation-config/
type AutomationConfigService interface {
	Get(context.Context, string) (*AutomationConfig, *atlas.Response, error)
	Update(context.Context, string, *AutomationConfig) (*atlas.Response, error)
	UpdateAgent(context.Context, string) (*AutomationConfigAgent, *atlas.Response, error)
}

// AutomationConfigServiceOp handles communication with the Automation config related methods of the MongoDB Cloud API
type AutomationConfigServiceOp service

// See more: https://docs.cloudmanager.mongodb.com/reference/api/automation-config/#get-the-automation-configuration
func (s *AutomationConfigServiceOp) Get(ctx context.Context, groupID string) (*AutomationConfig, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}
	basePath := fmt.Sprintf(automationConfigBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodGet, basePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AutomationConfig)
	resp, err := s.Client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// See more: https://docs.cloudmanager.mongodb.com/reference/api/automation-config/#update-the-automation-configuration
func (s *AutomationConfigServiceOp) Update(ctx context.Context, groupID string, updateRequest *AutomationConfig) (*atlas.Response, error) {
	if groupID == "" {
		return nil, atlas.NewArgError("groupID", "must be set")
	}
	basePath := fmt.Sprintf(automationConfigBasePath, groupID)

	req, err := s.Client.NewRequest(ctx, http.MethodPut, basePath, updateRequest)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(ctx, req, nil)

	return resp, err
}

// See more: https://docs.cloudmanager.mongodb.com/reference/api/automation-config/#update-the-automation-configuration
func (s *AutomationConfigServiceOp) UpdateAgent(ctx context.Context, groupID string) (*AutomationConfigAgent, *atlas.Response, error) {
	if groupID == "" {
		return nil, nil, atlas.NewArgError("groupID", "must be set")
	}

	basePath := fmt.Sprintf(automationConfigBasePath, groupID)
	path := fmt.Sprintf("%s/%s", basePath, "updateAgentVersions")

	req, err := s.Client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	agent := new(AutomationConfigAgent)
	resp, err := s.Client.Do(ctx, req, agent)

	return agent, resp, err
}

var _ AutomationConfigService = new(AutomationConfigServiceOp)

type AutomationConfigAgent struct {
	AutomationAgentVersion string `json:"automationAgentVersion"`
	BiConnectorVersion     string `json:"biConnectorVersion"`
}

type AutomationConfig struct {
	AgentVersion       *map[string]interface{}   `json:"agentVersion,omitempty"`
	Auth               Auth                      `json:"auth"`
	BackupVersions     []*map[string]interface{} `json:"backupVersions,omitempty"`
	Balancer           *map[string]interface{}   `json:"balancer,omitempty"`
	CPSModules         []*map[string]interface{} `json:"cpsModules,omitempty"`
	IndexConfigs       []*IndexConfig            `json:"indexConfigs,omitempty"`
	Kerberos           *map[string]interface{}   `json:"kerberos,omitempty"`
	LDAP               *map[string]interface{}   `json:"ldap,omitempty"`
	MongoDBVersions    []*map[string]interface{} `json:"mongoDbVersions,omitempty"`
	MongoSQLDs         []*map[string]interface{} `json:"mongosqlds,omitempty"`
	MonitoringVersions []*map[string]interface{} `json:"monitoringVersions,omitempty"`
	MongoTs            []*map[string]interface{} `json:"mongots,omitempty"`
	Options            *Options                  `json:"options"`
	Processes          []*Process                `json:"processes,omitempty"`
	ReplicaSets        []*ReplicaSet             `json:"replicaSets,omitempty"`
	Roles              []*map[string]interface{} `json:"roles,omitempty"`
	Sharding           []*map[string]interface{} `json:"sharding,omitempty"`
	SSL                *SSL                      `json:"ssl,omitempty"`
	UIBaseURL          string                    `json:"uiBaseUrl,omitempty"`
	Version            int                       `json:"version,omitempty"`
}

// IndexConfig represents a new index requests for a given database and collection.
type IndexConfig struct {
	DBName         string                  `json:"dbName"`              // Database that is indexed
	CollectionName string                  `json:"collectionName"`      // Collection that is indexed
	RSName         string                  `json:"rsName"`              // The replica set that the index is built on
	Key            [][]string              `json:"key"`                 // Keys array of keys to index and their type, sorting of keys is important for an index
	Options        *atlas.IndexOptions     `json:"options,omitempty"`   // Options MongoDB index options
	Collation      *atlas.CollationOptions `json:"collation,omitempty"` // Collation Mongo collation index options
}

// SSL config properties
type SSL struct {
	AutoPEMKeyFilePath    string `json:"autoPEMKeyFilePath,omitempty"`
	CAFilePath            string `json:"CAFilePath,omitempty"`
	ClientCertificateMode string `json:"clientCertificateMode,omitempty"`
}

// Auth authentication config
type Auth struct {
	Users                    []*MongoDBUser `json:"usersWanted,omitempty"` // Users is a list which contains the desired users at the project level.
	Disabled                 bool           `json:"disabled"`
	AuthoritativeSet         bool           `json:"authoritativeSet"`                   // AuthoritativeSet indicates if the MongoDBUsers should be synced with the current list of Users
	AutoAuthMechanisms       []string       `json:"autoAuthMechanisms,omitempty"`       // AutoAuthMechanisms is a list of auth mechanisms the Automation Agent is able to use
	AutoAuthMechanism        string         `json:"autoAuthMechanism"`                  // AutoAuthMechanism is the currently active agent authentication mechanism. This is a read only field
	DeploymentAuthMechanisms []string       `json:"deploymentAuthMechanisms,omitempty"` // DeploymentAuthMechanisms is a list of possible auth mechanisms that can be used within deployments
	AutoUser                 string         `json:"autoUser,omitempty"`                 // AutoUser is the MongoDB Automation Agent user, when x509 is enabled, it should be set to the subject of the AA's certificate
	Key                      string         `json:"key,omitempty"`                      // Key is the contents of the KeyFile, the automation agent will ensure this a KeyFile with these contents exists at the `KeyFile` path
	KeyFile                  string         `json:"keyfile,omitempty"`                  // KeyFile is the path to a keyfile with read & write permissions. It is a required field if `Disabled=false`
	KeyFileWindows           string         `json:"keyfileWindows,omitempty"`           // KeyFileWindows is required if `Disabled=false` even if the value is not used
	AutoPwd                  string         `json:"autoPwd,omitempty"`                  // AutoPwd is a required field when going from `Disabled=false` to `Disabled=true`
}

type MongoDBUser struct {
	Mechanisms                 []string       `json:"mechanisms"`
	Roles                      []*Role        `json:"roles"`
	Username                   string         `json:"user"`
	Database                   string         `json:"db"`
	AuthenticationRestrictions []string       `json:"authenticationRestrictions,omitempty"`
	InitPassword               string         `json:"initPwd,omitempty"` // The cleartext password to be assigned to the user
	ScramSha256Creds           *ScramShaCreds `json:"scramSha256Creds,omitempty"`
	ScramSha1Creds             *ScramShaCreds `json:"scramSha1Creds,omitempty"`
}

type Role struct {
	Role     string `json:"role"`
	Database string `json:"db"`
}

type ScramShaCreds struct {
	IterationCount int    `json:"iterationCount"`
	Salt           string `json:"salt"`
	ServerKey      string `json:"serverKey"`
	StoredKey      string `json:"storedKey"`
}

// Member configs
type Member struct {
	ID           int     `json:"_id"`
	ArbiterOnly  bool    `json:"arbiterOnly"`
	BuildIndexes bool    `json:"buildIndexes"`
	Hidden       bool    `json:"hidden"`
	Host         string  `json:"host"`
	Priority     float64 `json:"priority"`
	SlaveDelay   float64 `json:"slaveDelay"`
	Votes        float64 `json:"votes"`
}

// ReplicaSet configs
type ReplicaSet struct {
	ID              string   `json:"_id"`
	ProtocolVersion string   `json:"protocolVersion,omitempty"`
	Members         []Member `json:"members"`
}

// Options configs
type Options struct {
	DownloadBase string `json:"downloadBase"`
}

// NetSSL defines SSL parameters for Net
type NetSSL struct {
	Mode       string `json:"mode"`
	PEMKeyFile string `json:"PEMKeyFile"`
}

// Net part of the internal Process struct
type Net struct {
	Port int     `json:"port,omitempty"`
	SSL  *NetSSL `json:"ssl,omitempty"`
}

// Storage part of the internal Process struct
type Storage struct {
	DBPath string `json:"dbPath,omitempty"`
}

// Replication is part of the internal Process struct
type Replication struct {
	ReplSetName string `json:"replSetName,omitempty"`
}

// Sharding is part of the internal Process struct
type Sharding struct {
	ClusterRole string `json:"clusterRole"`
}

// SystemLog part of the internal Process struct
type SystemLog struct {
	Destination string `json:"destination,omitempty"`
	Path        string `json:"path,omitempty"`
}

// Args26 part of the internal Process struct
type Args26 struct {
	NET         Net          `json:"net"`                   // NET configuration for db connection (ports)
	Replication *Replication `json:"replication,omitempty"` // Replication configuration for ReplicaSets, omit this field if setting Sharding
	Sharding    *Sharding    `json:"sharding,omitempty"`    // Replication configuration for sharded clusters, omit this field if setting Replication
	Storage     *Storage     `json:"storage,omitempty"`     // Storage configuration for dbpath, config servers don't define this
	SystemLog   SystemLog    `json:"systemLog"`             // SystemLog configuration for the dblog
}

// LogRotate part of the internal Process struct
type LogRotate struct {
	SizeThresholdMB  float64 `json:"sizeThresholdMB,omitempty"`
	TimeThresholdHrs int     `json:"timeThresholdHrs,omitempty"`
}

// Process represents a single process in a deployment
type Process struct {
	Args26                      Args26     `json:"args2_6"`
	AuthSchemaVersion           int        `json:"authSchemaVersion,omitempty"`
	LastGoalVersionAchieved     int        `json:"lastGoalVersionAchieved,omitempty"`
	Name                        string     `json:"name,omitempty"`
	Cluster                     string     `json:"cluster,omitempty"`
	FeatureCompatibilityVersion string     `json:"featureCompatibilityVersion,omitempty"`
	Hostname                    string     `json:"hostname,omitempty"`
	LogRotate                   *LogRotate `json:"logRotate,omitempty"`
	Plan                        []string   `json:"plan,omitempty"`
	ProcessType                 string     `json:"processType,omitempty"`
	Version                     string     `json:"version,omitempty"`
	Disabled                    bool       `json:"disabled,omitempty"`
	ManualMode                  bool       `json:"manualMode,omitempty"`
}
