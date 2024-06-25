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

func TestDaemonConfigServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/admin/backup/daemon/configs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
			 "assignmentEnabled" : true,
			 "backupJobsEnabled" : false,
			 "configured" : true,
			 "garbageCollectionEnabled" : true,
			 "headDiskType" : "SSD",
			 "id" : "5628faffd4c606594adaa3b2",
			 "labels" : [ "l1", "l2" ],
			 "machine" : {
			   "headRootDirectory" : "/data/backup/",
			   "machine" : "localhost"
			 },
			 "numWorkers" : 50,
			 "resourceUsageEnabled" : true,
			 "restoreJobsEnabled" : false,
			 "restoreQueryableJobsEnabled" : true
		  } ],
		  "totalCount" : 1
}`)
	})

	config, _, err := client.DaemonConfig.List(ctx, nil)
	if err != nil {
		t.Fatalf("DaemonConfig.List returned error: %v", err)
	}

	assignmentEnabled := true

	expected := &Daemons{
		Results: []*Daemon{
			{
				AdminBackupConfig: AdminBackupConfig{
					ID:                ID,
					AssignmentEnabled: &assignmentEnabled,
					Labels:            []string{"l1", "l2"},
				},
				BackupJobsEnabled:           false,
				Configured:                  true,
				GarbageCollectionEnabled:    true,
				ResourceUsageEnabled:        true,
				RestoreQueryableJobsEnabled: true,
				HeadDiskType:                "SSD",
				NumWorkers:                  50,
				Machine: &Machine{
					Machine:           "localhost",
					HeadRootDirectory: "/data/backup/",
				},
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDaemonConfigServiceOp_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/daemon/configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			 "assignmentEnabled" : true,
			 "backupJobsEnabled" : false,
			 "configured" : true,
			 "garbageCollectionEnabled" : true,
			 "headDiskType" : "SSD",
			 "id" : "5628faffd4c606594adaa3b2",
			 "labels" : [ "l1", "l2" ],
			 "machine" : {
			   "headRootDirectory" : "/data/backup/",
			   "machine" : "localhost"
			 },
			 "numWorkers" : 50,
			 "resourceUsageEnabled" : true,
			 "restoreJobsEnabled" : false,
			 "restoreQueryableJobsEnabled" : true
}`)
	})

	config, _, err := client.DaemonConfig.Get(ctx, ID)
	if err != nil {
		t.Fatalf("DaemonConfig.Get returned error: %v", err)
	}
	assignmentEnabled := true

	expected := &Daemon{
		AdminBackupConfig: AdminBackupConfig{
			ID:                ID,
			AssignmentEnabled: &assignmentEnabled,
			Labels:            []string{"l1", "l2"},
		},
		BackupJobsEnabled:           false,
		Configured:                  true,
		GarbageCollectionEnabled:    true,
		ResourceUsageEnabled:        true,
		RestoreQueryableJobsEnabled: true,
		HeadDiskType:                "SSD",
		NumWorkers:                  50,
		Machine: &Machine{
			Machine:           "localhost",
			HeadRootDirectory: "/data/backup/",
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDaemonConfigServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/daemon/configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			 "assignmentEnabled" : true,
			 "backupJobsEnabled" : false,
			 "configured" : true,
			 "garbageCollectionEnabled" : true,
			 "headDiskType" : "SSD",
			 "id" : "5628faffd4c606594adaa3b2",
			 "labels" : [ "l1", "l2" ],
			 "machine" : {
			   "headRootDirectory" : "/data/backup/",
			   "machine" : "localhost"
			 },
			 "numWorkers" : 50,
			 "resourceUsageEnabled" : true,
			 "restoreJobsEnabled" : false,
			 "restoreQueryableJobsEnabled" : true
}`)
	})

	assignmentEnabled := true

	deamon := &Daemon{
		AdminBackupConfig: AdminBackupConfig{
			ID:                ID,
			AssignmentEnabled: &assignmentEnabled,
			Labels:            []string{"l1", "l2"},
		},
		BackupJobsEnabled:           false,
		Configured:                  true,
		GarbageCollectionEnabled:    true,
		ResourceUsageEnabled:        true,
		RestoreQueryableJobsEnabled: true,
		HeadDiskType:                "SSD",
		NumWorkers:                  50,
		Machine: &Machine{
			Machine:           "localhost",
			HeadRootDirectory: "/data/backup/",
		},
	}

	config, _, err := client.DaemonConfig.Update(ctx, ID, deamon)
	if err != nil {
		t.Fatalf("DaemonConfig.Update returned error: %v", err)
	}

	expected := &Daemon{
		AdminBackupConfig: AdminBackupConfig{
			ID:                ID,
			AssignmentEnabled: &assignmentEnabled,
			Labels:            []string{"l1", "l2"},
		},
		BackupJobsEnabled:           false,
		Configured:                  true,
		GarbageCollectionEnabled:    true,
		ResourceUsageEnabled:        true,
		RestoreQueryableJobsEnabled: true,
		HeadDiskType:                "SSD",
		NumWorkers:                  50,
		Machine: &Machine{
			Machine:           "localhost",
			HeadRootDirectory: "/data/backup/",
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDaemonConfigServiceOp_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/daemon/configs/%s", ID), func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.DaemonConfig.Delete(ctx, ID)
	if err != nil {
		t.Fatalf("DaemonConfig.Delete returned error: %v", err)
	}
}
