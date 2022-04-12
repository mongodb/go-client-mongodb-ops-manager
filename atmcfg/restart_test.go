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
)

func TestRestart(t *testing.T) {
	const clusterName = "restartTest"
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, false)
		Restart(config, clusterName)
		for i := range config.Processes {
			if config.Processes[i].LastRestart == "" {
				t.Errorf("TestRestart\n got=%#v", config.Processes[i].LastRestart)
			}
		}
	})
	t.Run("sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, false)
		Restart(config, clusterName)
		for i := range config.Processes {
			if config.Processes[i].LastRestart == "" {
				t.Errorf("TestRestart\n Got = %#v", config.Processes[i].LastRestart)
			}
		}
	})
}

func TestRestartProcessesByClusterName(t *testing.T) {
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, true)
		err := RestartProcessesByClusterName(config, clusterName, []string{"host0:27017"})
		if err != nil {
			t.Fatalf("RestartProcessesByClusterName() returned an unexpected error: %v", err)
		}
		if config.Processes[0].LastRestart == "" {
			t.Errorf("Got = %#v", config.Processes[0].LastRestart)
		}
	})
	t.Run("sharded cluster - one process", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := RestartProcessesByClusterName(config, clusterName, []string{"host2:27018"})
		if err != nil {
			t.Fatalf("RestartProcessesByClusterName() returned an unexpected error: %v", err)
		}
		if config.Processes[0].LastRestart != "" {
			t.Errorf("Got = %#v, want = %#v", config.Processes[0].LastRestart, "")
		}

		if config.Processes[1].LastRestart == "" {
			t.Errorf("Got = %#v", config.Processes[1].LastRestart)
		}
	})

	t.Run("sharded cluster - two processes", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := RestartProcessesByClusterName(config, clusterName, []string{"host0:27017", "host2:27018"})
		if err != nil {
			t.Fatalf("RestartProcessesByClusterName() returned an unexpected error: %v", err)
		}
		if config.Processes[0].LastRestart == "" {
			t.Errorf("Got = %#v", config.Processes[0].LastRestart)
		}

		if config.Processes[1].LastRestart == "" {
			t.Errorf("Got = %#v", config.Processes[1].LastRestart)
		}
	})
	t.Run("restart entire sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)

		err := RestartProcessesByClusterName(config, clusterName, nil)
		if err != nil {
			t.Fatalf("RestartProcessesByClusterName() returned an unexpected error: %v", err)
		}
		for i := range config.Processes {
			if config.Processes[i].LastRestart == "" {
				t.Errorf("Got = %#v", config.Processes[i].LastRestart)
			}
		}
	})
	t.Run("provide a process that does not exist", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)
		err := RestartProcessesByClusterName(config, clusterName, []string{"hostTest:21021"})

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
