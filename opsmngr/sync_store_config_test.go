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

func TestSyncStoreConfigServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/sync/mongoConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
					  "assignmentEnabled" : true,
					  "encryptedCredentials" : false,
					  "id" : "5628faffd4c606594adaa3b2",
					  "labels" : [ "l1", "l2" ],
					  "maxCapacityGB" : 8,
					  "uri" : "mongodb://localhost:27017",
					  "ssl" : true,
					  "usedSize" : 222,
					  "writeConcern" : "W2"
		  } ],
		  "totalCount" : 1
}`)
	})

	config, _, err := client.SyncStoreConfig.List(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.List returned error: %v", err)
	}

	assignmentEnabled := true
	encryptedCredentials := false
	ssl := true
	maxCapacityGB := int64(8)

	expected := &BackupStores{
		Results: []*BackupStore{
			{
				MaxCapacityGB: &maxCapacityGB,
				AdminBackupConfig: AdminBackupConfig{
					ID:                   ID,
					AssignmentEnabled:    &assignmentEnabled,
					EncryptedCredentials: &encryptedCredentials,
					URI:                  "mongodb://localhost:27017",
					Labels:               []string{"l1", "l2"},
					SSL:                  &ssl,
					WriteConcern:         "W2",
					UsedSize:             222,
				},
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestSyncStoreConfigServiceOp_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/sync/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
					  "assignmentEnabled" : true,
					  "encryptedCredentials" : false,
					  "id" : "5628faffd4c606594adaa3b2",
					  "labels" : [ "l1", "l2" ],
					  "maxCapacityGB" : 8,
					  "uri" : "mongodb://localhost:27017",
					  "ssl" : true,
					  "usedSize" : 222,
					  "writeConcern" : "W2"
}`)
	})

	config, _, err := client.SyncStoreConfig.Get(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.Get returned error: %v", err)
	}

	assignmentEnabled := true
	encryptedCredentials := false
	ssl := true
	maxCapacityGB := int64(8)

	expected := &BackupStore{
		MaxCapacityGB: &maxCapacityGB,
		AdminBackupConfig: AdminBackupConfig{
			ID:                   ID,
			AssignmentEnabled:    &assignmentEnabled,
			EncryptedCredentials: &encryptedCredentials,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  &ssl,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestSyncStoreConfigServiceOp_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/sync/mongoConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
					  "assignmentEnabled" : true,
					  "encryptedCredentials" : false,
					  "id" : "5628faffd4c606594adaa3b2",
					  "labels" : [ "l1", "l2" ],
					  "maxCapacityGB" : 8,
					  "uri" : "mongodb://localhost:27017",
					  "ssl" : true,
					  "usedSize" : 222,
					  "writeConcern" : "W2"
}`)
	})

	assignmentEnabled := true
	encryptedCredentials := false
	ssl := true
	maxCapacityGB := int64(8)

	sync := &BackupStore{
		MaxCapacityGB: &maxCapacityGB,
		AdminBackupConfig: AdminBackupConfig{
			ID:                   ID,
			AssignmentEnabled:    &assignmentEnabled,
			EncryptedCredentials: &encryptedCredentials,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  &ssl,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}

	config, _, err := client.SyncStoreConfig.Create(ctx, sync)
	if err != nil {
		t.Fatalf("BackupAdministrator.Create returned error: %v", err)
	}

	expected := &BackupStore{
		MaxCapacityGB: &maxCapacityGB,
		AdminBackupConfig: AdminBackupConfig{
			ID:                   ID,
			AssignmentEnabled:    &assignmentEnabled,
			EncryptedCredentials: &encryptedCredentials,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  &ssl,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestSyncStoreConfigServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/sync/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
					  "assignmentEnabled" : true,
					  "encryptedCredentials" : false,
					  "id" : "5628faffd4c606594adaa3b2",
					  "labels" : [ "l1", "l2" ],
					  "maxCapacityGB" : 8,
					  "uri" : "mongodb://localhost:27017",
					  "ssl" : true,
					  "usedSize" : 222,
					  "writeConcern" : "W2"
}`)
	})

	assignmentEnabled := true
	encryptedCredentials := false
	ssl := true
	maxCapacityGB := int64(8)

	sync := &BackupStore{
		MaxCapacityGB: &maxCapacityGB,
		AdminBackupConfig: AdminBackupConfig{
			ID:                   ID,
			AssignmentEnabled:    &assignmentEnabled,
			EncryptedCredentials: &encryptedCredentials,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  &ssl,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}

	config, _, err := client.SyncStoreConfig.Update(ctx, ID, sync)
	if err != nil {
		t.Fatalf("BackupAdministrator.Update returned error: %v", err)
	}

	expected := &BackupStore{
		MaxCapacityGB: &maxCapacityGB,
		AdminBackupConfig: AdminBackupConfig{
			ID:                   ID,
			AssignmentEnabled:    &assignmentEnabled,
			EncryptedCredentials: &encryptedCredentials,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  &ssl,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestSyncStoreConfigServiceOp_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/sync/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.SyncStoreConfig.Delete(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.Delete returned error: %v", err)
	}
}
