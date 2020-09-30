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

package opsmngr

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

func TestBackupConfigsServiceOp_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/groups/%s/backupConfigs/%s", projectID, clusterID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			  "clusterId" : "1",
			  "encryptionEnabled" : false,
			  "groupId" : "1",
			  "SSLEnabled" : false,
			  "statusName" : "STARTED",
			  "storageEngineName" : "WIRED_TIGER"
}`)
	})

	config, _, err := client.BackupConfigs.Get(ctx, projectID, clusterID)
	if err != nil {
		t.Fatalf("BackupConfigs.Get returned error: %v", err)
	}

	encryptionEnabled := false
	sslEnabled := false

	expected := &BackupConfig{
		GroupID:           "1",
		ClusterID:         "1",
		StatusName:        "STARTED",
		StorageEngineName: "WIRED_TIGER",
		EncryptionEnabled: &encryptionEnabled,
		SSLEnabled:        &sslEnabled,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupConfigsServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/groups/%s/backupConfigs", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
				"totalCount" : 1,
				  "results" : [ {
					"groupId" : "1",
					"clusterId" : "1",
					"statusName" : "STARTED",
					"storageEngineName" : "WIRED_TIGER",
					"SSLEnabled" : false
				  } ]
}`)
	})

	config, _, err := client.BackupConfigs.List(ctx, projectID, nil)
	if err != nil {
		t.Fatalf("Get.List returned error: %v", err)
	}

	sslEnabled := false

	expected := &BackupConfigs{
		Results: []*BackupConfig{
			{
				GroupID:           "1",
				ClusterID:         "1",
				StatusName:        "STARTED",
				StorageEngineName: "WIRED_TIGER",
				SSLEnabled:        &sslEnabled,
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupConfigsServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/groups/%s/backupConfigs/%s", projectID, clusterID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprint(w, `{
			  "clusterId" : "1",
			  "encryptionEnabled" : false,
			  "groupId" : "1",
			  "SSLEnabled" : false,
			  "statusName" : "STARTED",
			  "storageEngineName" : "WIRED_TIGER"
}`)
	})

	encryptionEnabled := false
	sslEnabled := false

	backupConfig := &BackupConfig{
		GroupID:           "1",
		ClusterID:         "1",
		StatusName:        "STARTED",
		StorageEngineName: "WIRED_TIGER",
		EncryptionEnabled: &encryptionEnabled,
		SSLEnabled:        &sslEnabled,
	}

	config, _, err := client.BackupConfigs.Update(ctx, projectID, clusterID, backupConfig)
	if err != nil {
		t.Fatalf("BackupConfigs.Update returned error: %v", err)
	}

	expected := &BackupConfig{
		GroupID:           "1",
		ClusterID:         "1",
		StatusName:        "STARTED",
		StorageEngineName: "WIRED_TIGER",
		EncryptionEnabled: &encryptionEnabled,
		SSLEnabled:        &sslEnabled,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}
