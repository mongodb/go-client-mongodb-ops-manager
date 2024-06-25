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

func TestFileSystemStoreConfigServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	assignmentEnabled := true
	loadFactor := int64(50)

	mux.HandleFunc("/api/public/v1.0/admin/backup/snapshot/fileSystemConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
			"assignmentEnabled" : true,
			"id" : "5628faffd4c606594adaa3b2",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 50,
			"mmapv1CompressionSetting" : "NONE",
			"storePath" : "/data/backup",
			"wtCompressionSetting" : "ZLIB"
		  } ],
		  "totalCount" : 1
}`)
	})

	config, _, err := client.FileSystemStoreConfig.List(ctx, nil)
	if err != nil {
		t.Fatalf("FileSystemStoreConfig.List returned error: %v", err)
	}

	expected := &FileSystemStoreConfigurations{
		Results: []*FileSystemStoreConfiguration{
			{
				BackupStore: BackupStore{
					AdminBackupConfig: AdminBackupConfig{
						ID:                ID,
						Labels:            []string{"l1", "l2"},
						AssignmentEnabled: &assignmentEnabled,
					},
					LoadFactor: &loadFactor,
				},
				MMAPV1CompressionSetting: "NONE",
				StorePath:                "/data/backup",
				WTCompressionSetting:     "ZLIB",
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestFileSystemStoreConfigServiceOp_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	assignmentEnabled := true
	loadFactor := int64(50)

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/snapshot/fileSystemConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"id" : "5628faffd4c606594adaa3b2",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 50,
			"mmapv1CompressionSetting" : "NONE",
			"storePath" : "/data/backup",
			"wtCompressionSetting" : "ZLIB"
}`)
	})

	config, _, err := client.FileSystemStoreConfig.Get(ctx, ID)
	if err != nil {
		t.Fatalf("FileSystemStoreConfig.Get returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		BackupStore: BackupStore{
			AdminBackupConfig: AdminBackupConfig{
				ID:                ID,
				Labels:            []string{"l1", "l2"},
				AssignmentEnabled: &assignmentEnabled,
			},
			LoadFactor: &loadFactor,
		},
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestFileSystemStoreConfigServiceOp_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	assignmentEnabled := true
	loadFactor := int64(50)

	mux.HandleFunc("/api/public/v1.0/admin/backup/snapshot/fileSystemConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"id" : "5628faffd4c606594adaa3b2",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 50,
			"mmapv1CompressionSetting" : "NONE",
			"storePath" : "/data/backup",
			"wtCompressionSetting" : "ZLIB"
}`)
	})

	fileSystem := &FileSystemStoreConfiguration{
		BackupStore: BackupStore{
			AdminBackupConfig: AdminBackupConfig{
				ID:                ID,
				Labels:            []string{"l1", "l2"},
				AssignmentEnabled: &assignmentEnabled,
			},
			LoadFactor: &loadFactor,
		},
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
	}

	config, _, err := client.FileSystemStoreConfig.Create(ctx, fileSystem)
	if err != nil {
		t.Fatalf("FileSystemStoreConfig.Create returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		BackupStore: BackupStore{
			AdminBackupConfig: AdminBackupConfig{
				ID:                ID,
				Labels:            []string{"l1", "l2"},
				AssignmentEnabled: &assignmentEnabled,
			},
			LoadFactor: &loadFactor,
		},
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestFileSystemStoreConfigServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	assignmentEnabled := true
	loadFactor := int64(50)

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/snapshot/fileSystemConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"id" : "5628faffd4c606594adaa3b2",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 50,
			"mmapv1CompressionSetting" : "NONE",
			"storePath" : "/data/backup",
			"wtCompressionSetting" : "ZLIB"
}`)
	})

	fileSystem := &FileSystemStoreConfiguration{
		BackupStore: BackupStore{
			AdminBackupConfig: AdminBackupConfig{
				ID:                ID,
				Labels:            []string{"l1", "l2"},
				AssignmentEnabled: &assignmentEnabled,
			},
			LoadFactor: &loadFactor,
		},
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
	}

	config, _, err := client.FileSystemStoreConfig.Update(ctx, ID, fileSystem)
	if err != nil {
		t.Fatalf("FileSystemStoreConfig.Update returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		BackupStore: BackupStore{
			AdminBackupConfig: AdminBackupConfig{
				ID:                ID,
				Labels:            []string{"l1", "l2"},
				AssignmentEnabled: &assignmentEnabled,
			},
			LoadFactor:    &loadFactor,
			MaxCapacityGB: nil,
			Provisioned:   nil,
			SyncSource:    "",
			Username:      "",
		},
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestFileSystemStoreConfigServiceOp_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/snapshot/fileSystemConfigs/%s", ID), func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.FileSystemStoreConfig.Delete(ctx, ID)
	if err != nil {
		t.Fatalf("FileSystemStoreConfig.Delete returned error: %v", err)
	}
}
