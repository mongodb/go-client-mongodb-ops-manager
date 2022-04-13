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

func TestSuspend(t *testing.T) {
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, false)
		Suspend(config, clusterName)
		if !config.Processes[0].ManualMode {
			t.Errorf("Got = %#v, want = %#v", config.Processes[0].Disabled, true)
		}
	})
	t.Run("sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, false)
		Suspend(config, clusterName)
		for i := range config.Processes {
			if !config.Processes[i].ManualMode {
				t.Errorf("Got = %#v, want = %#v", config.Processes[i].Disabled, true)
			}
		}
	})
}

func TestSuspendProcessesByClusterName(t *testing.T) {
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, false)
		err := SuspendProcessesByClusterName(config, clusterName, []string{"host0:27017"})
		if err != nil {
			t.Fatalf("ShutdownProcessesByClusterName() returned an unexpected error: %v", err)
		}
		if !config.Processes[0].ManualMode {
			t.Errorf("Got = %#v, want = %#v", config.Processes[0].Disabled, true)
		}
	})
	t.Run("sharded cluster - one process", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, false)
		err := SuspendProcessesByClusterName(config, clusterName, []string{"host2:27018"})
		if err != nil {
			t.Fatalf("ShutdownProcessesByClusterName() returned an unexpected error: %v", err)
		}

		if config.Processes[0].ManualMode {
			t.Errorf("Got = %#v, want = %#v", config.Processes[0].Disabled, false)
		}

		if !config.Processes[1].ManualMode {
			t.Errorf("Got = %#v, want = %#v", config.Processes[1].Disabled, true)
		}
	})
	t.Run("sharded cluster - two processes", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, false)
		err := SuspendProcessesByClusterName(config, clusterName, []string{"host2:27018", "host0:27017"})
		if err != nil {
			t.Fatalf("ShutdownProcessesByClusterName() returned an unexpected error: %v", err)
		}

		if !config.Processes[0].ManualMode {
			t.Errorf("Got = %#v, want = %#v", config.Processes[0].Disabled, true)
		}

		if !config.Processes[1].ManualMode {
			t.Errorf("Got = %#v, want = %#v", config.Processes[1].Disabled, true)
		}
	})
	t.Run("shutdown entire sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)

		err := SuspendProcessesByClusterName(config, clusterName, nil)
		if err != nil {
			t.Fatalf("ShutdownProcessesByClusterName() returned an unexpected error: %v", err)
		}

		for i := range config.Processes {
			if !config.Processes[i].ManualMode {
				t.Errorf("Got = %#v, want = %#v", config.Processes[i].Disabled, true)
			}
		}
	})
	t.Run("provide a process that does not exist", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)

		err := SuspendProcessesByClusterName(config, clusterName, []string{"hostTest:21021"})
		if !errors.Is(err, ErrProcessNotFound) {
			t.Fatalf("Got = %#v, want = %#v", err, ErrProcessNotFound)
		}

		for i := range config.Processes {
			if !config.Processes[i].ManualMode {
				t.Errorf("Got = %#v, want = %#v", config.Processes[i].Disabled, false)
			}
		}
	})
}
