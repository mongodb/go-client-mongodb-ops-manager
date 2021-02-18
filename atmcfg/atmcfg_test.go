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

package atmcfg

import (
	"testing"

	"go.mongodb.org/ops-manager/opsmngr"
)

const clusterName = "cluster_1"

func TestShutdown(t *testing.T) {
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, false)
		Shutdown(config, clusterName)
		if !config.Processes[0].Disabled {
			t.Errorf("TestShutdown\n got=%#v\nwant=%#v\n", config.Processes[0].Disabled, true)
		}
	})
	t.Run("sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, false)
		Shutdown(config, clusterName)
		for i := range config.Processes {
			if !config.Processes[i].Disabled {
				t.Errorf("TestShutdown\n got=%#v\nwant=%#v\n", config.Processes[i].Disabled, true)
			}
		}
	})
}

func TestStartup(t *testing.T) {
	t.Run("replica set", func(t *testing.T) {
		cloud := automationConfigWithOneReplicaSet(clusterName, true)

		Startup(cloud, clusterName)
		if cloud.Processes[0].Disabled {
			t.Errorf("TestStartup\n got=%#v\nwant=%#v\n", cloud.Processes[0].Disabled, false)
		}
	})
	t.Run("sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, true)

		Startup(config, clusterName)
		for i := range config.Processes {
			if config.Processes[i].Disabled {
				t.Errorf("TestStartup\n got=%#v\nwant=%#v\n", config.Processes[i].Disabled, false)
			}
		}
	})
}

func TestRemoveByClusterName(t *testing.T) {
	t.Run("replica set", func(t *testing.T) {
		config := automationConfigWithOneReplicaSet(clusterName, false)

		RemoveByClusterName(config, clusterName)
		if len(config.Processes) != 0 {
			t.Errorf("RemoveByClusterName\n got=%#v\nwant=0\n", len(config.Processes))
		}
	})
	t.Run("sharded cluster", func(t *testing.T) {
		config := automationConfigWithOneShardedCluster(clusterName, false)

		RemoveByClusterName(config, clusterName)
		if len(config.Processes) != 0 {
			t.Errorf("RemoveByClusterName\n got=%#v\nwant=0\n", len(config.Processes))
		}
	})
}

func TestAddIndexConfig(t *testing.T) {
	newIndex := &opsmngr.IndexConfig{
		DBName:         "test",
		CollectionName: "test",
		RSName:         "myReplicaSet",
		Key: [][]string{
			{
				"test1", "test",
			},
		},
		Options:   nil,
		Collation: nil,
	}
	t.Run("AutomationConfig not initialized", func(t *testing.T) {
		err := AddIndexConfig(nil, newIndex)
		if err == nil {
			t.Error("AddIndexConfig should return an error")
		}
	})

	t.Run("empty IndexConfig", func(t *testing.T) {
		a := &opsmngr.AutomationConfig{}
		err := AddIndexConfig(a, newIndex)
		if err != nil {
			t.Fatalf("AddIndexConfig unexpected error: %v", err)
		}
		if len(a.IndexConfigs) != 1 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("add an index with different keys", func(t *testing.T) {
		config := automationConfigWithIndexConfig()
		err := AddIndexConfig(config, newIndex)
		if err != nil {
			t.Fatalf("AutomationConfig() returned an unexpected error: %v", err)
		}
		if len(config.IndexConfigs) != 2 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("add an index with different rsName", func(t *testing.T) {
		newIndex := &opsmngr.IndexConfig{
			DBName:         "test",
			CollectionName: "test",
			RSName:         "myReplicaSet_1",
			Key: [][]string{
				{
					"test1", "test",
				},
			},
			Options:   nil,
			Collation: nil,
		}
		config := automationConfigWithIndexConfig()
		err := AddIndexConfig(config, newIndex)
		if err != nil {
			t.Fatalf("AutomationConfig() returned an unexpected error: %v", err)
		}
		if len(config.IndexConfigs) != 2 {
			t.Error("indexConfig has not been added to the AutomationConfig")
		}
	})

	t.Run("trying to add an index that is already in the AutomationConfig", func(t *testing.T) {
		config := automationConfigWithIndexConfig()
		index := &opsmngr.IndexConfig{
			DBName:         "test",
			CollectionName: "test",
			RSName:         "myReplicaSet",
			Key: [][]string{
				{
					"test", "test",
				},
			},
			Options:   nil,
			Collation: nil,
		}
		err := AddIndexConfig(config, index)

		if err == nil {
			t.Fatalf("AddIndexConfig should return an error")
		}
	})
}

func TestAddUser(t *testing.T) {
	config := automationConfigWithoutMongoDBUsers()
	u := mongoDBUsers()
	AddUser(config, u)
	if len(config.Auth.Users) != 1 {
		t.Error("User not added\n")
	}
}

func TestRemoveUser(t *testing.T) {
	config := automationConfigWithMongoDBUsers()
	t.Run("user exists", func(t *testing.T) {
		u := mongoDBUsers()
		err := RemoveUser(config, u.Username, u.Database)
		if err != nil {
			t.Fatalf("RemoveUser unexpecter err: %#v\n", err)
		}
		if len(config.Auth.Users) != 0 {
			t.Error("User not removed\n")
		}
	})
	t.Run("user does not exists", func(t *testing.T) {
		err := RemoveUser(config, "random", "random")
		if err == nil {
			t.Fatal("RemoveUser should return an error\n")
		}
	})
}

func TestEnableMechanism(t *testing.T) {
	config := automationConfigWithoutMongoDBUsers()
	t.Run("enable invalid", func(t *testing.T) {
		if e := EnableMechanism(config, []string{"invalid"}); e == nil {
			t.Fatalf("EnableMechanism() expected an error but got none\n")
		}
	})
	t.Run("enable SCRAM-SHA-256", func(t *testing.T) {
		if e := EnableMechanism(config, []string{"SCRAM-SHA-256"}); e != nil {
			t.Fatalf("EnableMechanism() unexpected error: %v\n", e)
		}

		if config.Auth.Disabled {
			t.Error("config.Auth.Disabled is true\n")
		}

		if config.Auth.AutoAuthMechanisms[0] != "SCRAM-SHA-256" {
			t.Error("AutoAuthMechanisms not set\n")
		}

		if config.Auth.AutoUser == "" || config.Auth.AutoPwd == "" {
			t.Error("config.Auth.Auto* not set\n")
		}

		if config.Auth.Key == "" || config.Auth.KeyFileWindows == "" || config.Auth.KeyFile == "" {
			t.Error("config.Auth.Key* not set\n")
		}

		if len(config.Auth.Users) != 0 {
			t.Errorf("expected 0 user got: %d\n", len(config.Auth.Users))
		}
	})
}

func TestConfigureScramCredentials(t *testing.T) {
	u := &opsmngr.MongoDBUser{
		Username: "test",
	}
	if err := ConfigureScramCredentials(u, "password"); err != nil {
		t.Fatalf("ConfigureScramCredentials() unexpected error: %v\n", err)
	}
	if u.ScramSha1Creds == nil {
		t.Fatalf("ConfigureScramCredentials() unexpected error: %v\n", u.ScramSha1Creds)
	}
	if u.ScramSha256Creds == nil {
		t.Fatalf("ConfigureScramCredentials() unexpected error: %v\n", u.ScramSha256Creds)
	}
}

func TestEnableMonitoring(t *testing.T) {
	type args struct {
		out      *opsmngr.AutomationConfig
		hostname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty config",
			args: args{
				out:      &opsmngr.AutomationConfig{},
				hostname: "test",
			},
			wantErr: false,
		},
		{
			name: "empty config",
			args: args{
				out: &opsmngr.AutomationConfig{
					MonitoringVersions: []*opsmngr.ConfigVersion{
						{
							Name:     "1",
							Hostname: "test",
						},
					},
				},
				hostname: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		out := tt.args.out
		hostname := tt.args.hostname
		wantErr := tt.wantErr
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := EnableMonitoring(out, hostname); (err != nil) != wantErr {
				t.Errorf("EnableMonitoring() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

func TestDisableMonitoring(t *testing.T) {
	type args struct {
		out      *opsmngr.AutomationConfig
		hostname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty config",
			args: args{
				out:      &opsmngr.AutomationConfig{},
				hostname: "test",
			},
			wantErr: true,
		},
		{
			name: "empty config",
			args: args{
				out: &opsmngr.AutomationConfig{
					MonitoringVersions: []*opsmngr.ConfigVersion{
						{
							Name:     "1",
							Hostname: "test",
						},
					},
				},
				hostname: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		out := tt.args.out
		hostname := tt.args.hostname
		wantErr := tt.wantErr
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := DisableMonitoring(out, hostname); (err != nil) != wantErr {
				t.Errorf("EnableMonitoring() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

func TestEnableBackup(t *testing.T) {
	type args struct {
		out      *opsmngr.AutomationConfig
		hostname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty config",
			args: args{
				out:      &opsmngr.AutomationConfig{},
				hostname: "test",
			},
			wantErr: false,
		},
		{
			name: "empty config",
			args: args{
				out: &opsmngr.AutomationConfig{
					BackupVersions: []*opsmngr.ConfigVersion{
						{
							Name:     "1",
							Hostname: "test",
						},
					},
				},
				hostname: "test",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		out := tt.args.out
		hostname := tt.args.hostname
		wantErr := tt.wantErr
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := EnableBackup(out, hostname); (err != nil) != wantErr {
				t.Errorf("EnableBackup() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

func TestDisableBackup(t *testing.T) {
	type args struct {
		out      *opsmngr.AutomationConfig
		hostname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty config",
			args: args{
				out:      &opsmngr.AutomationConfig{},
				hostname: "test",
			},
			wantErr: true,
		},
		{
			name: "empty config",
			args: args{
				out: &opsmngr.AutomationConfig{
					BackupVersions: []*opsmngr.ConfigVersion{
						{
							Name:     "1",
							Hostname: "test",
						},
					},
				},
				hostname: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		out := tt.args.out
		hostname := tt.args.hostname
		wantErr := tt.wantErr
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := DisableBackup(out, hostname); (err != nil) != wantErr {
				t.Errorf("EnableBackup() error = %v, wantErr %v", err, wantErr)
			}
		})
	}
}

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
				t.Errorf("TestRestart\n got=%#v", config.Processes[i].LastRestart)
			}
		}
	})
}

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
			if config.Processes[i].ProcessType == "mongod" && config.Processes[i].LastCompact == "" {
				t.Errorf("ReclaimFreeSpace\n got=%#v", config.Processes[i].LastRestart)
			}
			if config.Processes[i].ProcessType == "mongos" && config.Processes[i].LastCompact != "" {
				t.Errorf("ReclaimFreeSpace\n got=%#v", config.Processes[i].LastRestart)
			}
		}
	})
}
