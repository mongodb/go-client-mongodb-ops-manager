// Copyright 2022 MongoDB Inc
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

const (
	mongod = "mongod"
	mongos = "mongos"
)

func TestReclaimFreeSpace(t *testing.T) {
	const clusterName = "reclaimTest"
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, false)
		ReclaimFreeSpace(config, clusterName)
		for i := range config.Processes {
			if config.Processes[i].LastCompact == "" {
				t.Errorf("ReclaimFreeSpace\n got=%#v", config.Processes[i].LastRestart)
			}
		}
	})
	t.Run("sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, false)
		ReclaimFreeSpace(config, clusterName)
		for i := range config.Processes {
			isLastCompactEmpty(t, config.Processes[i], "", -1)
			if config.Processes[i].ProcessType == mongos && config.Processes[i].LastCompact != "" {
				t.Errorf("ReclaimFreeSpace\n got=%#v", config.Processes[i].LastRestart)
			}
		}
	})
}

func TestReclaimFreeSpaceForProcessesByClusterName(t *testing.T) {
	lastCompact := time.Now().Format(time.RFC3339)
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, true)
		err := ReclaimFreeSpaceForProcessesByClusterName(config, clusterName, lastCompact, []string{"host0:27017"})
		if err != nil {
			t.Fatalf("ReclaimFreeSpaceForProcessesByClusterName() returned an unexpected error: %v", err)
		}

		if config.Processes[0].LastCompact == "" {
			t.Errorf("Got = %#v", config.Processes[0].LastRestart)
		}
	})
	t.Run("sharded cluster - two processes", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := ReclaimFreeSpaceForProcessesByClusterName(config, clusterName, lastCompact, []string{"host0:27017", "host2:27018"})
		if err != nil {
			t.Fatalf("ReclaimFreeSpaceForProcessesByClusterName() returned an unexpected error: %v", err)
		}
		for i := range config.Processes {
			isLastCompactEmpty(t, config.Processes[i], "host2", 27017)
			isLastCompactEmpty(t, config.Processes[i], "host0", 27018)
		}
	})
	t.Run("reclaim free space for entire sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := ReclaimFreeSpaceForProcessesByClusterName(config, clusterName, lastCompact, nil)
		if err != nil {
			t.Fatalf("ReclaimFreeSpaceForProcessesByClusterName() returned an unexpected error: %v", err)
		}
		for i := range config.Processes {
			isLastCompactEmpty(t, config.Processes[i], "", -1)
		}
	})
	t.Run("provide a process that does not exist", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := ReclaimFreeSpaceForProcessesByClusterName(config, clusterName, lastCompact, []string{"hostTest:21021"})
		if !errors.Is(err, ErrProcessNotFound) {
			t.Fatalf("Got = %#v, want = %#v", err, ErrProcessNotFound)
		}

		for i := range config.Processes {
			if config.Processes[i].LastRestart != "" {
				t.Errorf("Got = %#v, want = %#v", config.Processes[i].LastRestart, "")
			}
		}
	})
}

func isLastCompactEmpty(t *testing.T, process *opsmngr.Process, hostname string, portN int) {
	t.Helper()
	if hostname != "" && process.Args26.NET.Port != portN && process.Hostname != hostname {
		return
	}

	if process.ProcessType == mongod && process.LastCompact == "" {
		t.Errorf("ReclaimFreeSpace\n got=%#v", process.LastCompact)
	}
}
