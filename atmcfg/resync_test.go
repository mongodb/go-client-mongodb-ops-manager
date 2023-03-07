// Copyright 2023 MongoDB Inc
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

package atmcfg

import (
	"errors"
	"testing"
	"time"

	"go.mongodb.org/ops-manager/opsmngr"
)

func TestStartInitialSync(t *testing.T) {
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, false)
		StartInitialSync(config, clusterName)
		for i := range config.Processes {
			if config.Processes[i].LastResync == "" {
				t.Errorf("StartInitialSync\n got=%#v", config.Processes[i].LastResync)
			}
		}
	})
	t.Run("sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, false)
		StartInitialSync(config, clusterName)
		for i := range config.Processes {
			isLastResyncEmpty(t, config.Processes[i], "", -1)
			if config.Processes[i].ProcessType == mongos && config.Processes[i].LastResync != "" {
				t.Errorf("StartInitialSync\n got=%#v", config.Processes[i].LastResync)
			}
		}
	})
}

func TestStartInitialSyncAtForProcessesByClusterName(t *testing.T) {
	lastResync := time.Now().Format(time.RFC3339)
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, true)
		err := StartInitialSyncAtForProcessesByClusterName(config, clusterName, lastResync, []string{"host0:27017"})
		if err != nil {
			t.Fatalf("StartInitialSyncAtForProcessesByClusterName() returned an unexpected error: %v", err)
		}

		if config.Processes[0].LastResync == "" {
			t.Errorf("Got = %#v", config.Processes[0].LastResync)
		}
	})
	t.Run("sharded cluster - two processes", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := StartInitialSyncAtForProcessesByClusterName(config, clusterName, lastResync, []string{"host0:27017", "host2:27018"})
		if err != nil {
			t.Fatalf("StartInitialSyncAtForProcessesByClusterName() returned an unexpected error: %v", err)
		}
		for i := range config.Processes {
			isLastResyncEmpty(t, config.Processes[i], "host2", 27017)
			isLastResyncEmpty(t, config.Processes[i], "host0", 27018)
		}
	})
	t.Run("resync for entire sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := StartInitialSyncAtForProcessesByClusterName(config, clusterName, lastResync, nil)
		if err != nil {
			t.Fatalf("StartInitialSyncAtForProcessesByClusterName() returned an unexpected error: %v", err)
		}
		for i := range config.Processes {
			isLastResyncEmpty(t, config.Processes[i], "", -1)
		}
	})
	t.Run("provide a process that does not exist", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := StartInitialSyncAtForProcessesByClusterName(config, clusterName, lastResync, []string{"hostTest:21021"})
		if !errors.Is(err, ErrProcessNotFound) {
			t.Fatalf("Got = %#v, want = %#v", err, ErrProcessNotFound)
		}

		for i := range config.Processes {
			if config.Processes[i].LastRestart != "" {
				t.Errorf("Got = %#v, want = %#v", config.Processes[i].LastResync, "")
			}
		}
	})
}

func isLastResyncEmpty(t *testing.T, process *opsmngr.Process, hostname string, portN int) {
	t.Helper()
	if hostname != "" && process.Args26.NET.Port != portN && process.Hostname != hostname {
		return
	}

	if process.ProcessType == mongod && process.LastResync == "" {
		t.Errorf("got=%#v", process.LastResync)
	}
}
