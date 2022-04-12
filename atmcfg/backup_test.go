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
	"testing"

	"go.mongodb.org/ops-manager/opsmngr"
)

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
