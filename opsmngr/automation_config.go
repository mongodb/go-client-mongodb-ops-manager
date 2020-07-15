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

// GetConfig retrieves the current automation configuration for a project.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/#get-the-automation-configuration
func (s *AutomationServiceOp) GetConfig(ctx context.Context, groupID string) (*AutomationConfig, *atlas.Response, error) {
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

// UpdateConfig updates a project’s automation configuration.
// When you submit updates, Ops Manager makes internal modifications to the data
// and then saves your new configuration version.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/#update-the-automation-configuration
func (s *AutomationServiceOp) UpdateConfig(ctx context.Context, groupID string, updateRequest *AutomationConfig) (*atlas.Response, error) {
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

// AutomationConfig represents an Ops Manager project automation config.
//
// See more: https://docs.opsmanager.mongodb.com/current/reference/cluster-configuration/
type AutomationConfig struct {
	AgentVersion         *map[string]interface{}   `json:"agentVersion,omitempty"`
	Auth                 Auth                      `json:"auth"`
	BackupVersions       []*map[string]interface{} `json:"backupVersions"`
	Balancer             *map[string]interface{}   `json:"balancer"`
	CPSModules           []*map[string]interface{} `json:"cpsModules"`
	IndexConfigs         []*IndexConfig            `json:"indexConfigs"`
	Kerberos             *map[string]interface{}   `json:"kerberos,omitempty"`
	LDAP                 *map[string]interface{}   `json:"ldap,omitempty"`
	MongoDBVersions      []*map[string]interface{} `json:"mongoDbVersions,omitempty"`
	MongoSQLDs           []*map[string]interface{} `json:"mongosqlds"`
	MonitoringVersions   []*map[string]interface{} `json:"monitoringVersions"`
	OnlineArchiveModules []*map[string]interface{} `json:"onlineArchiveModules"`
	MongoTS              []*map[string]interface{} `json:"mongots"`
	Options              *map[string]interface{}   `json:"options"`
	Processes            []*Process                `json:"processes"`
	ReplicaSets          []*ReplicaSet             `json:"replicaSets"`
	Roles                []*map[string]interface{} `json:"roles"`
	Sharding             []*ShardingConfig         `json:"sharding"`
	SSL                  *SSL                      `json:"ssl,omitempty"`
	TLS                  *SSL                      `json:"tls,omitempty"`
	UIBaseURL            *string                   `json:"uiBaseUrl"`
	Version              int                       `json:"version,omitempty"`
}

type ShardingConfig struct {
	Collections         []*map[string]interface{} `json:"collections"`
	ConfigServerReplica string                    `json:"configServerReplica"`
	Draining            []string                  `json:"draining"`
	ManagedSharding     bool                      `json:"managedSharding"`
	Name                string                    `json:"name"`
	Shards              []*Shard                  `json:"shards"`
	Tags                []string                  `json:"tags"`
}

type Shard struct {
	ID   string   `json:"_id"`
	RS   string   `json:"rs"`
	Tags []string `json:"tags"`
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
	AuthoritativeSet         bool           `json:"authoritativeSet"`             // AuthoritativeSet indicates if the MongoDBUsers should be synced with the current list of Users
	AutoAuthMechanism        string         `json:"autoAuthMechanism"`            // AutoAuthMechanism is the currently active agent authentication mechanism. This is a read only field
	AutoAuthMechanisms       []string       `json:"autoAuthMechanisms,omitempty"` // AutoAuthMechanisms is a list of auth mechanisms the Automation Agent is able to use
	AutoAuthRestrictions     []interface{}  `json:"autoAuthRestrictions"`
	AutoKerberosKeytabPath   string         `json:"autoKerberosKeytabPath,omitempty"`
	AutoLdapGroupDN          string         `json:"autoLdapGroupDN,omitempty"`
	AutoPwd                  string         `json:"autoPwd,omitempty"`                  // AutoPwd is a required field when going from `Disabled=false` to `Disabled=true`
	AutoUser                 string         `json:"autoUser,omitempty"`                 // AutoUser is the MongoDB Automation Agent user, when x509 is enabled, it should be set to the subject of the AA's certificate
	DeploymentAuthMechanisms []string       `json:"deploymentAuthMechanisms,omitempty"` // DeploymentAuthMechanisms is a list of possible auth mechanisms that can be used within deployments
	Disabled                 bool           `json:"disabled"`
	Key                      string         `json:"key,omitempty"`            // Key is the contents of the KeyFile, the automation agent will ensure this a KeyFile with these contents exists at the `KeyFile` path
	KeyFile                  string         `json:"keyfile,omitempty"`        // KeyFile is the path to a keyfile with read & write permissions. It is a required field if `Disabled=false`
	KeyFileWindows           string         `json:"keyfileWindows,omitempty"` // KeyFileWindows is required if `Disabled=false` even if the value is not used
	Users                    []*MongoDBUser `json:"usersWanted,omitempty"`    // Users is a list which contains the desired users at the project level.
	UsersDelete              []*MongoDBUser `json:"usersDeleted,omitempty"`
}

// Args26 part of the internal Process struct
type Args26 struct {
	AuditLog          *AuditLog               `json:"auditLog,omitempty"` // AuditLog configuration for audit logs
	NET               Net                     `json:"net"`                // NET configuration for db connection (ports)
	ProcessManagement *map[string]interface{} `json:"processManagement,omitempty"`
	Replication       *Replication            `json:"replication,omitempty"` // Replication configuration for ReplicaSets, omit this field if setting Sharding
	Sharding          *Sharding               `json:"sharding,omitempty"`    // Replication configuration for sharded clusters, omit this field if setting Replication
	Storage           *Storage                `json:"storage,omitempty"`     // Storage configuration for dbpath, config servers don't define this
	SystemLog         SystemLog               `json:"systemLog"`             // SystemLog configuration for the dblog
}

type MongoDBUser struct {
	AuthenticationRestrictions []string       `json:"authenticationRestrictions"`
	Database                   string         `json:"db"`
	InitPassword               string         `json:"initPwd,omitempty"` // The cleartext password to be assigned to the user
	Mechanisms                 []string       `json:"mechanisms"`
	Password                   string         `json:"pwd,omitempty"`
	Roles                      []*Role        `json:"roles"`
	Username                   string         `json:"user"`
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
	ID           int                     `json:"_id"`
	ArbiterOnly  bool                    `json:"arbiterOnly"`
	BuildIndexes bool                    `json:"buildIndexes"`
	Hidden       bool                    `json:"hidden"`
	Host         string                  `json:"host"`
	Priority     float64                 `json:"priority"`
	SlaveDelay   float64                 `json:"slaveDelay"`
	Tags         *map[string]interface{} `json:"tags"`
	Votes        float64                 `json:"votes"`
}

// ReplicaSet configs
type ReplicaSet struct {
	ID                                 string                  `json:"_id"`
	ProtocolVersion                    string                  `json:"protocolVersion,omitempty"`
	Members                            []Member                `json:"members"`
	Settings                           *map[string]interface{} `json:"settings,omitempty"`
	WriteConcernMajorityJournalDefault string                  `json:"writeConcernMajorityJournalDefault,omitempty"`
}

// TLS defines TLS parameters for Net
type TLS struct {
	Mode               string `json:"mode,omitempty"`
	PEMKeyFile         string `json:"PEMKeyFile,omitempty"`
	CAFile             string `json:"CAFile,omitempty"`
	CertificateKeyFile string `json:"certificateKeyFile,omitempty"`
}

// Net part of the internal Process struct
type Net struct {
	Port                   int     `json:"port,omitempty"`
	BindIP                 *string `json:"bindIp,omitempty"`
	BindIPAll              *bool   `json:"bindIpAll,omitempty"`
	IPV6                   *bool   `json:"ipv6,omitempty"`
	MaxIncomingConnections *int    `json:"maxIncomingConnections,omitempty"`
	SSL                    *TLS    `json:"ssl,omitempty"`
	TLS                    *TLS    `json:"tls,omitempty"`
}

// Storage part of the internal Process struct
type Storage struct {
	DBPath  string                  `json:"dbPath,omitempty"`
	Engine  string                  `json:"engine,omitempty"`
	Journal *map[string]interface{} `json:"journal,omitempty"`
}

// Replication is part of the internal Process struct
type Replication struct {
	ReplSetName string `json:"replSetName,omitempty"`
}

// Sharding is part of the internal Process struct
type Sharding struct {
	ClusterRole string `json:"clusterRole,omitempty"`
}

// SystemLog part of the internal Process struct
type SystemLog struct {
	Destination string `json:"destination,omitempty"`
	Path        string `json:"path,omitempty"`
	LogAppend   bool   `json:"logAppend,omitempty"`
}

// AuditLog part of the internal Process struct
type AuditLog struct {
	Destination string `json:"destination,omitempty"`
	Path        string `json:"path,omitempty"`
	Format      string `json:"format,omitempty"`
}

// LogRotate part of the internal Process struct
type LogRotate struct {
	SizeThresholdMB  float64 `json:"sizeThresholdMB,omitempty"`
	TimeThresholdHrs int     `json:"timeThresholdHrs,omitempty"`
}

// Process represents a single process in a deployment
type Process struct {
	Args26                      Args26             `json:"args2_6"`
	AuthSchemaVersion           int                `json:"authSchemaVersion,omitempty"`
	LastGoalVersionAchieved     int                `json:"lastGoalVersionAchieved,omitempty"`
	Name                        string             `json:"name,omitempty"`
	Cluster                     string             `json:"cluster,omitempty"`
	FeatureCompatibilityVersion string             `json:"featureCompatibilityVersion,omitempty"`
	Hostname                    string             `json:"hostname,omitempty"`
	LogRotate                   *LogRotate         `json:"logRotate,omitempty"`
	Plan                        []string           `json:"plan,omitempty"`
	ProcessType                 string             `json:"processType,omitempty"`
	Version                     string             `json:"version,omitempty"`
	Disabled                    bool               `json:"disabled"`
	ManualMode                  bool               `json:"manualMode"`
	NumCores                    int                `json:"numCores"`
	Horizons                    *map[string]string `json:"horizons,omitempty"`
}
