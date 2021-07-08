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

package atmcfg // import "go.mongodb.org/ops-manager/atmcfg"

import (
	"crypto/sha1" //nolint:gosec // mongodb scram-sha-1 supports this tho is not recommended
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/ops-manager/opsmngr"
	"go.mongodb.org/ops-manager/search"
)

func setDisabledByClusterName(out *opsmngr.AutomationConfig, name string, disabled bool) {
	newDeploymentAuthMechanisms(out)
	setDisabledByReplicaSetName(out, name, disabled)
	setDisabledByShardName(out, name, disabled)
}

func newDeploymentAuthMechanisms(out *opsmngr.AutomationConfig) {
	// This value may not be present and is mandatory
	if out.Auth.DeploymentAuthMechanisms == nil {
		out.Auth.DeploymentAuthMechanisms = make([]string, 0)
	}
}

func setDisabledByReplicaSetName(out *opsmngr.AutomationConfig, name string, disabled bool) {
	setDisabledByReplicaSetNameAndProcesses(out, name, nil, disabled)
}

func setDisabledByShardName(out *opsmngr.AutomationConfig, name string, disabled bool) {
	setDisabledByShardNameAndProcesses(out, name, nil, disabled)
}

func setDisabledByReplicaSetNameAndProcesses(out *opsmngr.AutomationConfig, name string, processesMap map[string]bool, disabled bool) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == name
	})
	if found {
		rs := out.ReplicaSets[i]
		for _, m := range rs.Members {
			for k, p := range out.Processes {
				if p.Name == m.Host {
					setDisable(out.Processes[k], processesMap, disabled)
				}
			}
		}
	}
}

func setDisabledByShardNameAndProcesses(out *opsmngr.AutomationConfig, name string, processesMap map[string]bool, disabled bool) {
	i, found := search.ShardingConfig(out.Sharding, func(s *opsmngr.ShardingConfig) bool {
		return s.Name == name
	})
	if found {
		s := out.Sharding[i]
		// disable shards
		for _, rs := range s.Shards {
			setDisabledByReplicaSetNameAndProcesses(out, rs.ID, processesMap, disabled)
		}
		// disable config rs
		setDisabledByReplicaSetNameAndProcesses(out, s.ConfigServerReplica, processesMap, disabled)
		// disable mongos
		for i := range out.Processes {
			if out.Processes[i].Cluster == name {
				setDisable(out.Processes[i], processesMap, disabled)
			}
		}
	}
}

func setDisable(process *opsmngr.Process, processesMap map[string]bool, disabled bool) {
	if len(processesMap) > 0 {
		key := fmt.Sprintf("%s:%d", process.Hostname, process.Args26.NET.Port)
		if _, ok := processesMap[key]; ok {
			process.Disabled = disabled
			processesMap[key] = true
		}
	} else {
		process.Disabled = disabled
	}
}

func setDisableProcessByProcessNameAndPort(out *opsmngr.AutomationConfig, clusterName string, processesMap map[string]bool, disabled bool) {
	newDeploymentAuthMechanisms(out)
	for k, p := range out.Processes {
		key := fmt.Sprintf("%s:%d", p.Hostname, p.Args26.NET.Port)
		if _, ok := processesMap[key]; ok && strings.Contains(out.Processes[k].Name, clusterName) {
			out.Processes[k].Disabled = disabled
			processesMap[key] = true
		}
	}
}

// Shutdown disables all processes of the given cluster name.
func Shutdown(out *opsmngr.AutomationConfig, clusterName string) {
	setDisabledByClusterName(out, clusterName, true)
}

// ShutdownProcessesByClusterName disables the entire cluster or its processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func ShutdownProcessesByClusterName(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	if len(processes) == 0 {
		Shutdown(out, clusterName)
		return nil
	}
	return shutdownProcesses(out, clusterName, processes)
}

// shutdownProcesses disables processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func shutdownProcesses(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	processesMap := map[string]bool{}
	for _, hostnameAndPort := range processes {
		processesMap[hostnameAndPort] = false
	}

	newDeploymentAuthMechanisms(out)
	setDisabledByReplicaSetNameAndProcesses(out, clusterName, processesMap, true)
	setDisabledByShardNameAndProcesses(out, clusterName, processesMap, true)

	return newProcessNotFoundError(clusterName, processesMap)
}

// newProcessNotFoundError returns an error if a process was not found.
func newProcessNotFoundError(clusterName string, processesMap map[string]bool) error {
	var processesNotFound []string
	for k, v := range processesMap {
		if !v {
			processesNotFound = append(processesNotFound, k)
		}
	}

	if len(processesNotFound) > 0 {
		return fmt.Errorf(`process not found in "%v": %v`, clusterName, processesNotFound)
	}

	return nil
}

// StartupProcessesByClusterName enables the entire cluster or its processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func StartupProcessesByClusterName(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	if len(processes) == 0 {
		Startup(out, clusterName)
		return nil
	}
	return startupProcess(out, clusterName, processes)
}

// StartupProcess enables processes. Processes are provided in the format {"hostname:port","hostname2:port2"}.
func startupProcess(out *opsmngr.AutomationConfig, clusterName string, processes []string) error {
	processesMap := map[string]bool{}
	for _, hostnameAndPort := range processes {
		processesMap[hostnameAndPort] = false
	}

	setDisableProcessByProcessNameAndPort(out, clusterName, processesMap, false)

	return newProcessNotFoundError(clusterName, processesMap)
}

// Startup enables all processes of the given cluster name.
func Startup(out *opsmngr.AutomationConfig, name string) {
	setDisabledByClusterName(out, name, false)
}

const monitoringVersion = "7.2.0.488-1" // Last monitoring version released

// EnableMonitoring enables monitoring for the given hostname.
func EnableMonitoring(out *opsmngr.AutomationConfig, hostname string) error {
	for _, v := range out.MonitoringVersions {
		if v.Hostname == hostname {
			return fmt.Errorf("monitoring already enabled for '%s'", hostname)
		}
	}
	out.MonitoringVersions = append(out.MonitoringVersions, &opsmngr.ConfigVersion{
		Name:     monitoringVersion,
		Hostname: hostname,
	})
	return nil
}

// DisableMonitoring disables monitoring for the given hostname.
func DisableMonitoring(out *opsmngr.AutomationConfig, hostname string) error {
	for i, v := range out.MonitoringVersions {
		if v.Hostname == hostname {
			out.MonitoringVersions = append(out.MonitoringVersions[:i], out.MonitoringVersions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no monitoring for '%s'", hostname)
}

const backupVersion = "7.8.1.1109-1" // Last backup version released

// EnableBackup enables backup for the given hostname.
func EnableBackup(out *opsmngr.AutomationConfig, hostname string) error {
	for _, v := range out.BackupVersions {
		if v.Hostname == hostname {
			return fmt.Errorf("backup already enabled for '%s'", hostname)
		}
	}
	out.BackupVersions = append(out.BackupVersions, &opsmngr.ConfigVersion{
		Name:     backupVersion,
		Hostname: hostname,
	})
	return nil
}

// DisableBackup disables backup for the given hostname.
func DisableBackup(out *opsmngr.AutomationConfig, hostname string) error {
	for i, v := range out.BackupVersions {
		if v.Hostname == hostname {
			out.BackupVersions = append(out.BackupVersions[:i], out.BackupVersions[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("no backup for '%s'", hostname)
}

// RemoveByClusterName removes a cluster and its associated processes from the config.
// This won't shutdown any running process.
func RemoveByClusterName(out *opsmngr.AutomationConfig, name string) {
	// This value may not be present and is mandatory
	if out.Auth.DeploymentAuthMechanisms == nil {
		out.Auth.DeploymentAuthMechanisms = make([]string, 0)
	}
	removeByReplicaSetName(out, name)
	removeByShardName(out, name)
}

func removeByReplicaSetName(out *opsmngr.AutomationConfig, name string) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == name
	})
	if found {
		rs := out.ReplicaSets[i]
		out.ReplicaSets = append(out.ReplicaSets[:i], out.ReplicaSets[i+1:]...)
		for _, m := range rs.Members {
			for k, p := range out.Processes {
				if p.Name == m.Host {
					out.Processes = append(out.Processes[:k], out.Processes[k+1:]...)
				}
			}
		}
	}
}

func removeByShardName(out *opsmngr.AutomationConfig, name string) {
	i, found := search.ShardingConfig(out.Sharding, func(rs *opsmngr.ShardingConfig) bool {
		return rs.Name == name
	})
	if found {
		s := out.Sharding[i]
		out.Sharding = append(out.Sharding[:i], out.Sharding[i+1:]...)
		// remove shards
		for _, rs := range s.Shards {
			removeByReplicaSetName(out, rs.ID)
		}
		// remove config rs
		removeByReplicaSetName(out, s.ConfigServerReplica)
		// remove mongos
		for j := range out.Processes {
			if out.Processes[j].Cluster == name {
				out.Processes = append(out.Processes[:j], out.Processes[j+1:]...)
			}
		}
	}
}

// AddUser adds a opsmngr.MongoDBUser to the opsmngr.AutomationConfi.
func AddUser(out *opsmngr.AutomationConfig, u *opsmngr.MongoDBUser) {
	out.Auth.UsersWanted = append(out.Auth.UsersWanted, u)
}

// ConfigureScramCredentials creates both SCRAM-SHA-1 and SCRAM-SHA-256 credentials.
// Use this method to guarantee that password can be updated later.
func ConfigureScramCredentials(user *opsmngr.MongoDBUser, password string) error {
	scram256Creds, err := newScramSha256Creds(user, password)
	if err != nil {
		return err
	}

	scram1Creds, err := newScramSha1Creds(user, password)
	if err != nil {
		return err
	}
	user.ScramSha256Creds = scram256Creds
	user.ScramSha1Creds = scram1Creds
	return nil
}

func newScramSha1Creds(user *opsmngr.MongoDBUser, password string) (*opsmngr.ScramShaCreds, error) {
	scram1Salt, err := generateSalt(sha1.New)
	if err != nil {
		return nil, fmt.Errorf("error generating scramSha1 salt: %w", err)
	}
	scram1Creds, err := newScramShaCreds(scram1Salt, user.Username, password, mongoCR)
	if err != nil {
		return nil, fmt.Errorf("error generating scramSha1Creds: %w", err)
	}
	return scram1Creds, nil
}

func newScramSha256Creds(user *opsmngr.MongoDBUser, password string) (*opsmngr.ScramShaCreds, error) {
	scram256Salt, err := generateSalt(sha256.New)
	if err != nil {
		return nil, fmt.Errorf("error generating scramSha256 salt: %w", err)
	}
	scram256Creds, err := newScramShaCreds(scram256Salt, user.Username, password, scramSha256)
	if err != nil {
		return nil, fmt.Errorf("error generating scramSha256 creds: %w", err)
	}
	return scram256Creds, nil
}

// ErrUnsupportedMechanism means the provided mechanism wasn't valid.
var ErrUnsupportedMechanism = errors.New("unrecognized SCRAM-SHA format")

// newScramShaCreds takes a plain text password and a specified mechanism name and generates
// the ScramShaCreds which will be embedded into a MongoDBUser.
func newScramShaCreds(salt []byte, username, password, mechanism string) (*opsmngr.ScramShaCreds, error) {
	if mechanism != scramSha256 && mechanism != mongoCR {
		return nil, fmt.Errorf("%w %s", ErrUnsupportedMechanism, mechanism)
	}
	var hashConstructor hashingFunc
	iterations := 0
	if mechanism == scramSha256 {
		hashConstructor = sha256.New
		iterations = scramSha256Iterations
	} else if mechanism == mongoCR {
		hashConstructor = sha1.New
		iterations = scramSha1Iterations

		// MONGODB-CR/SCRAM-SHA-1 requires the hash of the password being passed computeScramCredentials
		// instead of the plain text password.
		var err error
		password, err = md5Hex(username + ":mongo:" + password)
		if err != nil {
			return nil, err
		}
	}
	base64EncodedSalt := base64.StdEncoding.EncodeToString(salt)
	return computeScramCredentials(hashConstructor, iterations, base64EncodedSalt, password)
}

// AddIndexConfig adds an opsmngr.IndexConfig to the opsmngr.AutomationConfig.
func AddIndexConfig(out *opsmngr.AutomationConfig, newIndex *opsmngr.IndexConfig) error {
	if out == nil {
		return errors.New("the Automation Config has not been initialized")
	}
	_, exists := search.MongoDBIndexes(out.IndexConfigs, compareIndexConfig(newIndex))

	if exists {
		return errors.New("index already exists")
	}
	out.IndexConfigs = append(out.IndexConfigs, newIndex)

	return nil
}

// compareIndexConfig returns a function that compares two indexConfig struts.
func compareIndexConfig(newIndex *opsmngr.IndexConfig) func(index *opsmngr.IndexConfig) bool {
	return func(index *opsmngr.IndexConfig) bool {
		if newIndex.RSName == index.RSName && newIndex.CollectionName == index.CollectionName && newIndex.DBName == index.DBName && len(newIndex.Key) == len(index.Key) {
			// if keys are equal the two indexes are considered to be the same
			for i := 0; i < len(newIndex.Key); i++ {
				if newIndex.Key[i][0] != index.Key[i][0] || newIndex.Key[i][1] != index.Key[i][1] {
					return false
				}
			}
			return true
		}
		return false
	}
}

// RemoveUser removes a MongoDBUser from the authentication config.
func RemoveUser(out *opsmngr.AutomationConfig, username, database string) error {
	pos, found := search.MongoDBUsers(out.Auth.UsersWanted, func(p *opsmngr.MongoDBUser) bool {
		return p.Username == username && p.Database == database
	})
	if !found {
		return fmt.Errorf("user '%s' not found for '%s'", username, database)
	}
	out.Auth.UsersWanted = append(out.Auth.UsersWanted[:pos], out.Auth.UsersWanted[pos+1:]...)
	return nil
}

const (
	automationAgentName            = "mms-automation"
	keyLength                      = 500
	mongoCR                        = "MONGODB-CR"
	scramSha256                    = "SCRAM-SHA-256"
	atmAgentWindowsKeyFilePath     = "%SystemDrive%\\MMSAutomation\\versions\\keyfile"
	atmAgentKeyFilePathInContainer = "/var/lib/mongodb-mms-automation/keyfile"
)

// EnableMechanism allows you to enable a given set of authentication mechanisms to an opsmngr.AutomationConfig.
// This method currently only supports MONGODB-CR, and SCRAM-SHA-256.
func EnableMechanism(out *opsmngr.AutomationConfig, m []string) error {
	out.Auth.Disabled = false
	for _, v := range m {
		if v != mongoCR && v != scramSha256 {
			return fmt.Errorf("unsupported mechanism %s", v)
		}
		if v == scramSha256 && out.Auth.AutoAuthMechanism == "" {
			out.Auth.AutoAuthMechanism = v
		}
		if !stringInSlice(out.Auth.DeploymentAuthMechanisms, v) {
			out.Auth.DeploymentAuthMechanisms = append(out.Auth.DeploymentAuthMechanisms, v)
		}
		if !stringInSlice(out.Auth.AutoAuthMechanisms, v) {
			out.Auth.AutoAuthMechanisms = append(out.Auth.AutoAuthMechanisms, v)
		}
	}

	if out.Auth.AutoUser == "" && out.Auth.AutoPwd == "" {
		if err := setAutoUser(out); err != nil {
			return err
		}
	}

	if out.Auth.Key == "" {
		var err error
		if out.Auth.Key, err = generateRandomBase64String(keyLength); err != nil {
			return err
		}
	}
	if out.Auth.Keyfile == "" {
		out.Auth.Keyfile = atmAgentKeyFilePathInContainer
	}
	if out.Auth.KeyfileWindows == "" {
		out.Auth.KeyfileWindows = atmAgentWindowsKeyFilePath
	}

	return nil
}

func setAutoUser(out *opsmngr.AutomationConfig) error {
	var err error
	out.Auth.AutoUser = automationAgentName
	if out.Auth.AutoPwd, err = generateRandomASCIIString(keyLength); err != nil {
		return err
	}

	return nil
}

func stringInSlice(a []string, x string) bool {
	for _, b := range a {
		if b == x {
			return true
		}
	}
	return false
}

func restartByReplicaSetName(out *opsmngr.AutomationConfig, name, lastRestart string) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == name
	})
	if found {
		rs := out.ReplicaSets[i]
		for _, m := range rs.Members {
			for k, p := range out.Processes {
				if p.Name == m.Host {
					out.Processes[k].LastRestart = lastRestart
				}
			}
		}
	}
}

func restartByShardName(out *opsmngr.AutomationConfig, name, lastRestart string) {
	i, found := search.ShardingConfig(out.Sharding, func(s *opsmngr.ShardingConfig) bool {
		return s.Name == name
	})
	if found {
		s := out.Sharding[i]
		// restart shards
		for _, rs := range s.Shards {
			restartByReplicaSetName(out, rs.ID, lastRestart)
		}
		// restart config rs
		restartByReplicaSetName(out, s.ConfigServerReplica, lastRestart)
		// restart mongos
		for i := range out.Processes {
			if out.Processes[i].Cluster == name {
				out.Processes[i].LastRestart = lastRestart
			}
		}
	}
}

// Restart sets all process of a cluster to restart.
func Restart(out *opsmngr.AutomationConfig, name string) {
	// This value may not be present and is mandatory
	if out.Auth.DeploymentAuthMechanisms == nil {
		out.Auth.DeploymentAuthMechanisms = make([]string, 0)
	}
	lastRestart := time.Now().Format(time.RFC3339)
	restartByReplicaSetName(out, name, lastRestart)
	restartByShardName(out, name, lastRestart)
}

// ReclaimFreeSpace sets all process of a cluster to reclaim free space.
func ReclaimFreeSpace(out *opsmngr.AutomationConfig, name string) {
	// This value may not be present and is mandatory
	if out.Auth.DeploymentAuthMechanisms == nil {
		out.Auth.DeploymentAuthMechanisms = make([]string, 0)
	}
	lastCompact := time.Now().Format(time.RFC3339)
	reclaimByReplicaSetName(out, name, lastCompact)
	reclaimByShardName(out, name, lastCompact)
}

func reclaimByReplicaSetName(out *opsmngr.AutomationConfig, name, lastCompact string) {
	i, found := search.ReplicaSets(out.ReplicaSets, func(rs *opsmngr.ReplicaSet) bool {
		return rs.ID == name
	})
	if found {
		rs := out.ReplicaSets[i]
		for _, m := range rs.Members {
			for k, p := range out.Processes {
				if p.Name == m.Host {
					out.Processes[k].LastCompact = lastCompact
				}
			}
		}
	}
}

func reclaimByShardName(out *opsmngr.AutomationConfig, name, lastCompact string) {
	i, found := search.ShardingConfig(out.Sharding, func(s *opsmngr.ShardingConfig) bool {
		return s.Name == name
	})
	if found {
		s := out.Sharding[i]
		// compact shards
		for _, rs := range s.Shards {
			reclaimByReplicaSetName(out, rs.ID, lastCompact)
		}
		// compact config rs
		reclaimByReplicaSetName(out, s.ConfigServerReplica, lastCompact)
		// compact doesn't run on mongoses
	}
}
