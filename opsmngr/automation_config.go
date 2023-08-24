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
func (s *AutomationServiceOp) GetConfig(ctx context.Context, groupID string) (*AutomationConfig, *Response, error) {
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
func (s *AutomationServiceOp) UpdateConfig(ctx context.Context, groupID string, updateRequest *AutomationConfig) (*Response, error) {
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
// See more: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/automation-config-parameters/
type AutomationConfig struct {
	AgentVersion              *map[string]interface{}   `json:"agentVersion,omitempty"`
	AtlasProxies              *[]interface{}            `json:"atlasProxies,omitempty"`
	AtlasUISes                []*map[string]interface{} `json:"atlasUISes"`
	Filebeat                  *map[string]interface{}   `json:"filebeat,omitempty"`
	Auth                      Auth                      `json:"auth"`
	BackupVersions            []*ConfigVersion          `json:"backupVersions"`
	Balancer                  *map[string]interface{}   `json:"balancer"`
	ClusterWideConfigurations *map[string]interface{}   `json:"clusterWideConfigurations,omitempty"`
	CPSModules                []*map[string]interface{} `json:"cpsModules"`
	DBCheckModules            []*map[string]interface{} `json:"dbCheckModules"`
	IndexConfigs              []*IndexConfig            `json:"indexConfigs"`
	Kerberos                  *map[string]interface{}   `json:"kerberos,omitempty"`
	LDAP                      *map[string]interface{}   `json:"ldap,omitempty"`
	MaintainedEnvoys          []*map[string]interface{} `json:"maintainedEnvoys"`
	MongoDBToolsVersion       *map[string]interface{}   `json:"mongoDbToolsVersion,omitempty"`
	MongoDBVersions           []*map[string]interface{} `json:"mongoDbVersions,omitempty"`
	MongoSQLDs                []*map[string]interface{} `json:"mongosqlds"` //nolint:tagliatelle // correct from API
	MonitoringVersions        []*ConfigVersion          `json:"monitoringVersions,omitempty"`
	OnlineArchiveModules      []*map[string]interface{} `json:"onlineArchiveModules"`
	Mongots                   []*map[string]interface{} `json:"mongots"`
	Options                   *map[string]interface{}   `json:"options"`
	Processes                 []*Process                `json:"processes"`
	ReplicaSets               []*ReplicaSet             `json:"replicaSets"`
	Roles                     []*map[string]interface{} `json:"roles"`
	Sharding                  []*ShardingConfig         `json:"sharding"`
	SSL                       *SSL                      `json:"ssl,omitempty"` // Deprecated: prefer TLS
	TLS                       *SSL                      `json:"tls,omitempty"`
	UIBaseURL                 *string                   `json:"uiBaseUrl,omitempty"`
	Version                   int                       `json:"version,omitempty"`
}

type ConfigVersion struct {
	Name     string `json:"name"`
	Hostname string `json:"hostname"`
}

// ShardingConfig sharded clusters configuration.
type ShardingConfig struct {
	Collections         []*map[string]interface{} `json:"collections"`
	ConfigServerReplica string                    `json:"configServerReplica"`
	Draining            []string                  `json:"draining"`
	ManagedSharding     bool                      `json:"managedSharding"`
	Name                string                    `json:"name"`
	Shards              []*Shard                  `json:"shards"`
	Tags                []*map[string]interface{} `json:"tags"`
}

// Shard details.
type Shard struct {
	ID   string   `json:"_id"` //nolint:tagliatelle // correct from API
	RS   string   `json:"rs"`
	Tags []string `json:"tags"`
}

// IndexConfig represents a new index requests for a given database and collection.
//
// See: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/automation-config-parameters/#indexes
type IndexConfig struct {
	DBName         string                  `json:"dbName"`              // DBName of the database that is indexed
	CollectionName string                  `json:"collectionName"`      // CollectionName that is indexed
	RSName         string                  `json:"rsName"`              // RSName that the index is built on
	Key            [][]string              `json:"key"`                 // Key array of keys to index and their type, sorting of keys is important for an index
	Options        *atlas.IndexOptions     `json:"options,omitempty"`   // Options MongoDB index options
	Collation      *atlas.CollationOptions `json:"collation,omitempty"` // Collation Mongo collation index options
}

// SSL config properties.
//
// See: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/automation-config-parameters/#tls
type SSL struct {
	AutoPEMKeyFilePath    string `json:"autoPEMKeyFilePath,omitempty"` //nolint:tagliatelle // correct from API
	AutoPEMKeyFilePwd     string `json:"autoPEMKeyFilePwd,omitempty"`  //nolint:tagliatelle // correct from API
	CAFilePath            string `json:"CAFilePath,omitempty"`         //nolint:tagliatelle // correct from API
	ClientCertificateMode string `json:"clientCertificateMode,omitempty"`
}

// Auth authentication config.
//
// See: https://docs.opsmanager.mongodb.com/current/reference/api/automation-config/automation-config-parameters/#authentication
type Auth struct {
	AuthoritativeSet         bool           `json:"authoritativeSet"`             // AuthoritativeSet indicates if the MongoDBUsers should be synced with the current list of UsersWanted
	AutoAuthMechanism        string         `json:"autoAuthMechanism"`            // AutoAuthMechanism is the currently active agent authentication mechanism. This is a read only field
	AutoAuthMechanisms       []string       `json:"autoAuthMechanisms,omitempty"` // AutoAuthMechanisms is a list of auth mechanisms the Automation Agent is able to use
	AutoAuthRestrictions     []interface{}  `json:"autoAuthRestrictions"`
	AutoKerberosKeytabPath   string         `json:"autoKerberosKeytabPath,omitempty"`
	AutoLdapGroupDN          string         `json:"autoLdapGroupDN,omitempty"`          //nolint:tagliatelle // AutoLdapGroupDN follows go convention while tag is correct from API
	AutoPwd                  string         `json:"autoPwd,omitempty"`                  // AutoPwd is a required field when going from `Disabled=false` to `Disabled=true`
	AutoUser                 string         `json:"autoUser,omitempty"`                 // AutoUser is the MongoDB Automation Agent user, when x509 is enabled, it should be set to the subject of the AA's certificate
	DeploymentAuthMechanisms []string       `json:"deploymentAuthMechanisms,omitempty"` // DeploymentAuthMechanisms is a list of possible auth mechanisms that can be used within deployments
	Disabled                 bool           `json:"disabled"`                           // Disabled indicates if auth is disabled
	Key                      string         `json:"key,omitempty"`                      // Key is the contents of the Keyfile, the automation agent will ensure this a Keyfile with these contents exists at the `Keyfile` path
	Keyfile                  string         `json:"keyfile,omitempty"`                  // Keyfile is the path to a keyfile with read & write permissions. It is a required field if `Disabled=false`
	KeyfileWindows           string         `json:"keyfileWindows,omitempty"`           // KeyfileWindows is required if `Disabled=false` even if the value is not used.
	NewAutoPwd               string         `json:"newAutoPwd,omitempty"`               // NewAutoPwd is a new password that the Automation uses when connecting to an instance.
	UsersDeleted             []*MongoDBUser `json:"usersDeleted"`                       // UsersDeleted are objects that define the authenticated users to be deleted from specified databases or from all databases
	UsersWanted              []*MongoDBUser `json:"usersWanted"`                        // UsersWanted is a list which contains the desired users at the project level.
}

// Args26 part of the internal Process struct.
type Args26 struct {
	AuditLog           *AuditLog               `json:"auditLog,omitempty"` // AuditLog configuration for audit logs
	BasisTech          *map[string]interface{} `json:"basisTech,omitempty"`
	NET                Net                     `json:"net"` // NET configuration for db connection (ports)
	OperationProfiling *map[string]interface{} `json:"operationProfiling,omitempty"`
	ProcessManagement  *map[string]interface{} `json:"processManagement,omitempty"`
	Replication        *Replication            `json:"replication,omitempty"` // Replication configuration for ReplicaSets, omit this field if setting Sharding
	SetParameter       *map[string]interface{} `json:"setParameter,omitempty"`
	Security           *map[string]interface{} `json:"security,omitempty"`
	Sharding           *Sharding               `json:"sharding,omitempty"` // Replication configuration for sharded clusters, omit this field if setting Replication
	Storage            *Storage                `json:"storage,omitempty"`  // Storage configuration for dbpath, config servers don't define this
	SNMP               *map[string]interface{} `json:"snmp,omitempty"`
	SystemLog          SystemLog               `json:"systemLog"` // SystemLog configuration for the dblog
}

// MongoDBUser database user.
type MongoDBUser struct {
	AuthenticationRestrictions []AuthenticationRestriction `json:"authenticationRestrictions"`
	CustomData                 interface{}                 `json:"customData,omitempty"`
	Database                   string                      `json:"db"`                //nolint:tagliatelle // Database is a better name
	InitPassword               string                      `json:"initPwd,omitempty"` // The cleartext password to be assigned to the user
	Mechanisms                 *[]string                   `json:"mechanisms,omitempty"`
	Password                   string                      `json:"pwd,omitempty"` //nolint:tagliatelle // Password is a better name than just pwd
	Roles                      []*Role                     `json:"roles"`
	ScramSha256Creds           *ScramShaCreds              `json:"scramSha256Creds,omitempty"`
	ScramSha1Creds             *ScramShaCreds              `json:"scramSha1Creds,omitempty"`
	Username                   string                      `json:"user"` //nolint:tagliatelle // Username is a better name than just user
}

// AuthenticationRestriction of a database user.
type AuthenticationRestriction struct {
	ClientSource  []string `json:"clientSource"`
	ServerAddress []string `json:"serverAddress"`
}

// Role of a database user.
type Role struct {
	Role     string `json:"role"`
	Database string `json:"db"` //nolint:tagliatelle // Database is a better name than just db
}

// ScramShaCreds configuration.
type ScramShaCreds struct {
	IterationCount int    `json:"iterationCount"`
	Salt           string `json:"salt"`
	ServerKey      string `json:"serverKey"`
	StoredKey      string `json:"storedKey"`
}

// Member configs.
type Member struct {
	ID                 int                `json:"_id"` //nolint:tagliatelle // correct from API
	ArbiterOnly        bool               `json:"arbiterOnly"`
	BuildIndexes       bool               `json:"buildIndexes"`
	Hidden             bool               `json:"hidden"`
	Horizons           *map[string]string `json:"horizons,omitempty"` // Horizons are managed by Kubernetes Operator and should not be manually edited
	Host               string             `json:"host"`
	Priority           float64            `json:"priority"`
	SlaveDelay         *float64           `json:"slaveDelay,omitempty"`         // Deprecated: since 5.0+ use SecondaryDelaySecs instead
	SecondaryDelaySecs *float64           `json:"secondaryDelaySecs,omitempty"` // SecondaryDelaySecs replaces SlaveDelay since 5.0+
	Tags               *map[string]string `json:"tags,omitempty"`
	Votes              float64            `json:"votes"`
}

// ReplicaSet configuration.
type ReplicaSet struct {
	ID                                 string                  `json:"_id"` //nolint:tagliatelle // correct from API
	ProtocolVersion                    string                  `json:"protocolVersion,omitempty"`
	Members                            []Member                `json:"members"`
	Settings                           *map[string]interface{} `json:"settings,omitempty"`
	WriteConcernMajorityJournalDefault string                  `json:"writeConcernMajorityJournalDefault,omitempty"`
	Force                              *Force                  `json:"force,omitempty"`
}

type Force struct {
	CurrentVersion int `json:"currentVersion"`
}

// TLS defines TLS parameters for Net.
type TLS struct {
	CAFile                     string `json:"CAFile,omitempty"` //nolint:tagliatelle // correct from API
	CertificateKeyFile         string `json:"certificateKeyFile,omitempty"`
	CertificateKeyFilePassword string `json:"certificateKeyFilePassword,omitempty"`
	CertificateSelector        string `json:"certificateSelector,omitempty"`
	ClusterCertificateSelector string `json:"clusterCertificateSelector,omitempty"`
	ClusterFile                string `json:"clusterFile,omitempty"`
	ClusterPassword            string `json:"clusterPassword,omitempty"`
	CRLFile                    string `json:"CRLFile,omitempty"` //nolint:tagliatelle // correct from API
	DisabledProtocols          string `json:"disabledProtocols,omitempty"`
	FIPSMode                   *bool  `json:"FIPSMode,omitempty"` //nolint:tagliatelle // correct from API
	Mode                       string `json:"mode,omitempty"`
	PEMKeyFile                 string `json:"PEMKeyFile,omitempty"` //nolint:tagliatelle // correct from API
}

// Net part of the internal Process struct.
type Net struct {
	BindIP                 *string                 `json:"bindIp,omitempty"`
	BindIPAll              *bool                   `json:"bindIpAll,omitempty"`
	Compression            *map[string]interface{} `json:"compression,omitempty"`
	HTTP                   *map[string]interface{} `json:"http,omitempty"` // Deprecated: deprecated since 3.2 and removed in 3.6
	IPV6                   *bool                   `json:"ipv6,omitempty"`
	ListenBacklog          string                  `json:"listenBacklog,omitempty"`
	MaxIncomingConnections *int                    `json:"maxIncomingConnections,omitempty"`
	Port                   int                     `json:"port,omitempty"`
	ServiceExecutor        string                  `json:"serviceExecutor,omitempty"`
	SSL                    *TLS                    `json:"ssl,omitempty"` // Deprecated: deprecated since 4.4 use TLS instead
	TLS                    *TLS                    `json:"tls,omitempty"`
	TransportLayer         string                  `json:"transportLayer,omitempty"`
	UnixDomainSocket       *map[string]interface{} `json:"unixDomainSocket,omitempty"`
}

// Storage part of the internal Process struct.
type Storage struct {
	DBPath                 string                  `json:"dbPath,omitempty"`
	DirectoryPerDB         *bool                   `json:"directoryPerDB,omitempty"` //nolint:tagliatelle // DirectoryPerDB follows go convention while directoryPerDB is correct from API
	Engine                 string                  `json:"engine,omitempty"`
	IndexBuildRetry        *bool                   `json:"indexBuildRetry,omitempty"`
	InMemory               *map[string]interface{} `json:"inMemory,omitempty"`
	Journal                *map[string]interface{} `json:"journal,omitempty"`
	NSSize                 *int                    `json:"nsSize,omitempty"`
	OplogMinRetentionHours *float64                `json:"oplogMinRetentionHours,omitempty"`
	PreallocDataFiles      *bool                   `json:"preallocDataFiles,omitempty"`
	Quota                  *map[string]interface{} `json:"quota,omitempty"`
	RepairPath             string                  `json:"repairPath,omitempty"`
	SmallFiles             *bool                   `json:"smallFiles,omitempty"`
	SyncPeriodSecs         *float64                `json:"syncPeriodSecs,omitempty"`
	WiredTiger             *map[string]interface{} `json:"wiredTiger,omitempty"`
}

// Replication is part of the internal Process struct.
type Replication struct {
	EnableMajorityReadConcern *bool  `json:"enableMajorityReadConcern,omitempty"`
	OplogSizeMB               *int   `json:"oplogSizeMB,omitempty"` //nolint:tagliatelle // Bytes vs bits
	ReplSetName               string `json:"replSetName,omitempty"`
}

// Sharding is part of the internal Process struct.
type Sharding struct {
	ArchiveMovedChunks *bool  `json:"archiveMovedChunks,omitempty"`
	AutoSplit          *bool  `json:"autoSplit,omitempty"`
	ChunkSize          *int   `json:"chunkSize,omitempty"`
	ClusterRole        string `json:"clusterRole,omitempty"`
}

// SystemLog part of the internal Process struct.
type SystemLog struct {
	Destination     string `json:"destination,omitempty"`
	Path            string `json:"path,omitempty"`
	LogAppend       bool   `json:"logAppend,omitempty"`
	Verbosity       int    `json:"verbosity,omitempty"`
	Quiet           bool   `json:"quiet,omitempty"`
	SyslogFacility  string `json:"syslogFacility,omitempty"`
	LogRotate       string `json:"logRotate,omitempty"`
	TimeStampFormat string `json:"timeStampFormat,omitempty"`
}

// AuditLog part of the internal Process struct.
type AuditLog struct {
	Destination string `json:"destination,omitempty"`
	Path        string `json:"path,omitempty"`
	Format      string `json:"format,omitempty"`
	Filter      string `json:"filter,omitempty"`
}

// LogRotate part of the internal Process struct.
type LogRotate struct {
	MaxUncompressed    *int     `json:"maxUncompressed,omitempty"`
	NumTotal           *int     `json:"numTotal,omitempty"`
	NumUncompressed    *int     `json:"numUncompressed,omitempty"`
	PercentOfDiskspace *float64 `json:"percentOfDiskspace,omitempty"`
	SizeThresholdMB    float64  `json:"sizeThresholdMB,omitempty"` //nolint:tagliatelle // Bytes vs bits
	TimeThresholdHrs   int      `json:"timeThresholdHrs,omitempty"`
}

type DefaultReadConcern struct {
	Level string `json:"level"`
}

type DefaultWriteConcern struct {
	W        interface{} `json:"w,omitempty"` // W can be string or number
	J        *bool       `json:"j,omitempty"`
	Wtimeout int         `json:"wtimeout"`
}

type DefaultRWConcern struct {
	DefaultReadConcern  *DefaultReadConcern  `json:"defaultReadConcern,omitempty"`
	DefaultWriteConcern *DefaultWriteConcern `json:"defaultWriteConcern,omitempty"`
}

// Process represents a single process in a deployment.
type Process struct {
	Alias                                      string             `json:"alias,omitempty"`
	Args26                                     Args26             `json:"args2_6"` //nolint:tagliatelle // correct from API
	AuditLogRotate                             *LogRotate         `json:"auditLogRotate,omitempty"`
	AuthSchemaVersion                          int                `json:"authSchemaVersion,omitempty"`
	BackupPITRestoreType                       string             `json:"backupPitRestoreType,omitempty"`
	BackupRestoreCertificateValidationHostname string             `json:"backupRestoreCertificateValidationHostname,omitempty"`
	BackupRestoreCheckpointTimestamp           interface{}        `json:"backupRestoreCheckpointTimestamp,omitempty"`
	BackupRestoreDesiredTime                   interface{}        `json:"backupRestoreDesiredTime,omitempty"`
	BackupRestoreFilterList                    interface{}        `json:"backupRestoreFilterList,omitempty"`
	BackupRestoreJobID                         string             `json:"backupRestoreJobId,omitempty"`
	BackupRestoreURL                           string             `json:"backupRestoreUrl,omitempty"`
	BackupRestoreOplogBaseURL                  string             `json:"backupRestoreOplogBaseUrl,omitempty"`
	BackupRestoreOplog                         interface{}        `json:"backupRestoreOplog,omitempty"`
	BackupRestoreRsVersion                     *int               `json:"backupRestoreRsVersion,omitempty"`
	BackupRestoreSourceGroupID                 string             `json:"backupRestoreSourceGroupId,omitempty"`
	BackupRestoreSourceRsID                    string             `json:"backupRestoreSourceRsId,omitempty"`
	BackupRestoreSystemRolesUUID               string             `json:"backupRestoreSystemRolesUUID,omitempty"` //nolint:tagliatelle // correct from API
	BackupRestoreSystemUsersUUID               string             `json:"backupRestoreSystemUsersUUID,omitempty"` //nolint:tagliatelle // correct from API
	BackupRestoreVerificationKey               string             `json:"backupRestoreVerificationKey,omitempty"`
	CPUAffinity                                []int              `json:"cpuAffinity,omitempty"`
	Cluster                                    string             `json:"cluster,omitempty"`
	CredentialsVersion                         *int               `json:"credentialsVersion,omitempty"`
	DefaultRWConcern                           *DefaultRWConcern  `json:"defaultRWConcern,omitempty"` //nolint:tagliatelle // correct from API
	Disabled                                   bool               `json:"disabled"`
	FeatureCompatibilityVersion                string             `json:"featureCompatibilityVersion,omitempty"`
	FullVersion                                interface{}        `json:"fullVersion,omitempty"`
	Horizons                                   *map[string]string `json:"horizons,omitempty"`
	Hostname                                   string             `json:"hostname,omitempty"`
	KMIPProxyPort                              *int               `json:"kmipProxyPort,omitempty"`
	LastCompact                                string             `json:"lastCompact,omitempty"`
	LastResync                                 string             `json:"lastResync,omitempty"`
	LastRestart                                string             `json:"lastRestart,omitempty"`
	LastGoalVersionAchieved                    int                `json:"lastGoalVersionAchieved,omitempty"`
	LastKMIPMasterKeyRotation                  string             `json:"lastKmipMasterKeyRotation,omitempty"`
	LogLevel                                   *int               `json:"logLevel,omitempty"`
	LogRotate                                  *LogRotate         `json:"logRotate,omitempty"`
	ManualMode                                 bool               `json:"manualMode"`
	Name                                       string             `json:"name,omitempty"`
	NumCores                                   int                `json:"numCores,omitempty"`
	Plan                                       []string           `json:"plan,omitempty"`
	ProcessType                                string             `json:"processType,omitempty"`
	Version                                    string             `json:"version,omitempty"`
}
