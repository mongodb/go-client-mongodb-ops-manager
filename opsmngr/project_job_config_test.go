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

func TestProjectJobConfigServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/admin/backup/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
			"daemonFilter" : [ {
			  "headRootDirectory" : "/data/backup",
			  "machine" : "http://backup.example.com"
			} ],
			"id" : "5628faffd4c606594adaa3b2",
			"kmipClientCertPassword" : "P4$$w0rD",
			"kmipClientCertPath" : "/certs/kmip",
			"labelFilter" : [ "l1", "l2" ],
			"links" : [ {
			  "href" : "https://{OPSMANAGER-HOST}:{PORT}/api/public/v1.0/admin/backup/groups/{PROJECT-ID}",
			  "rel" : "self"
			}, {
			  "href" : "https://{OPSMANAGER-HOST}:{PORT}/api/public/groups/{PROJECT-ID}",
			  "rel" : "http://mms.mongodb.com/group"
			} ],
			"oplogStoreFilter" : [ {
			  "id" : "5628faffd4c606594adaa3b2",
			  "type" : "oplogStore"
			} ],
			"snapshotStoreFilter" : [ {
			  "id" : "5628faffd4c606594adaa3b2",
			  "type" : "s3blockstore"
			} ],
			"syncStoreFilter" : [ "s1", "s2" ]
		  } ],
		  "totalCount" : 1
}`)
	})

	config, _, err := client.ProjectJobConfig.List(ctx, nil)
	if err != nil {
		t.Fatalf("ProjectJobConfig.List returned error: %v", err)
	}

	expected := &ProjectJobs{
		Results: []*ProjectJob{
			{
				AdminBackupConfig: AdminBackupConfig{
					ID: ID,
				},
				KMIPClientCertPassword: "P4$$w0rD",
				KMIPClientCertPath:     "/certs/kmip",
				LabelFilter:            []string{"l1", "l2"},
				SyncStoreFilter:        []string{"s1", "s2"},
				DaemonFilter: []*Machine{{
					Machine:           "http://backup.example.com",
					HeadRootDirectory: "/data/backup",
				}},
				OplogStoreFilter: []*StoreFilter{{
					ID:   "5628faffd4c606594adaa3b2",
					Type: "oplogStore",
				}},
				SnapshotStoreFilter: []*StoreFilter{{
					ID:   "5628faffd4c606594adaa3b2",
					Type: "s3blockstore",
				}},
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProjectJobConfigServiceOp_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/groups/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"daemonFilter" : [ {
			  "headRootDirectory" : "/data/backup",
			  "machine" : "http://backup.example.com"
			} ],
			"id" : "5628faffd4c606594adaa3b2",
			"kmipClientCertPassword" : "P4$$w0rD",
			"kmipClientCertPath" : "/certs/kmip",
			"labelFilter" : [ "l1", "l2" ],
			"links" : [ {
			  "href" : "https://{OPSMANAGER-HOST}:{PORT}/api/public/v1.0/admin/backup/groups/{PROJECT-ID}",
			  "rel" : "self"
			}, {
			  "href" : "https://{OPSMANAGER-HOST}:{PORT}/api/public/groups/{PROJECT-ID}",
			  "rel" : "http://mms.mongodb.com/group"
			} ],
			"oplogStoreFilter" : [ {
			  "id" : "5628faffd4c606594adaa3b2",
			  "type" : "oplogStore"
			} ],
			"snapshotStoreFilter" : [ {
			  "id" : "5628faffd4c606594adaa3b2",
			  "type" : "s3blockstore"
			} ],
			"syncStoreFilter" : [ "s1", "s2" ]
}`)
	})

	config, _, err := client.ProjectJobConfig.Get(ctx, ID)
	if err != nil {
		t.Fatalf("ProjectJobConfig.Get returned error: %v", err)
	}

	expected := &ProjectJob{

		AdminBackupConfig: AdminBackupConfig{
			ID: ID,
		},
		KMIPClientCertPassword: "P4$$w0rD",
		KMIPClientCertPath:     "/certs/kmip",
		LabelFilter:            []string{"l1", "l2"},
		SyncStoreFilter:        []string{"s1", "s2"},
		DaemonFilter: []*Machine{{
			Machine:           "http://backup.example.com",
			HeadRootDirectory: "/data/backup",
		}},
		OplogStoreFilter: []*StoreFilter{{
			ID:   "5628faffd4c606594adaa3b2",
			Type: "oplogStore",
		}},
		SnapshotStoreFilter: []*StoreFilter{{
			ID:   "5628faffd4c606594adaa3b2",
			Type: "s3blockstore",
		}},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBProjectJobConfigServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/groups/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			"daemonFilter" : [ {
			  "headRootDirectory" : "/data/backup",
			  "machine" : "http://backup.example.com"
			} ],
			"id" : "5628faffd4c606594adaa3b2",
			"kmipClientCertPassword" : "P4$$w0rD",
			"kmipClientCertPath" : "/certs/kmip",
			"labelFilter" : [ "l1", "l2" ],
			"links" : [ {
			  "href" : "https://{OPSMANAGER-HOST}:{PORT}/api/public/v1.0/admin/backup/groups/{PROJECT-ID}",
			  "rel" : "self"
			}, {
			  "href" : "https://{OPSMANAGER-HOST}:{PORT}/api/public/groups/{PROJECT-ID}",
			  "rel" : "http://mms.mongodb.com/group"
			} ],
			"oplogStoreFilter" : [ {
			  "id" : "5628faffd4c606594adaa3b2",
			  "type" : "oplogStore"
			} ],
			"snapshotStoreFilter" : [ {
			  "id" : "5628faffd4c606594adaa3b2",
			  "type" : "s3blockstore"
			} ],
			"syncStoreFilter" : [ "s1", "s2" ]
}`)
	})

	projectJob := &ProjectJob{

		AdminBackupConfig: AdminBackupConfig{
			ID: ID,
		},
		KMIPClientCertPassword: "P4$$w0rD",
		KMIPClientCertPath:     "/certs/kmip",
		LabelFilter:            []string{"l1", "l2"},
		SyncStoreFilter:        []string{"s1", "s2"},
		DaemonFilter: []*Machine{{
			Machine:           "http://backup.example.com",
			HeadRootDirectory: "/data/backup",
		}},
		OplogStoreFilter: []*StoreFilter{{
			ID:   "5628faffd4c606594adaa3b2",
			Type: "oplogStore",
		}},
		SnapshotStoreFilter: []*StoreFilter{{
			ID:   "5628faffd4c606594adaa3b2",
			Type: "s3blockstore",
		}},
	}

	config, _, err := client.ProjectJobConfig.Update(ctx, ID, projectJob)
	if err != nil {
		t.Fatalf("ProjectJobConfig.Update returned error: %v", err)
	}

	expected := &ProjectJob{

		AdminBackupConfig: AdminBackupConfig{
			ID: ID,
		},
		KMIPClientCertPassword: "P4$$w0rD",
		KMIPClientCertPath:     "/certs/kmip",
		LabelFilter:            []string{"l1", "l2"},
		SyncStoreFilter:        []string{"s1", "s2"},
		DaemonFilter: []*Machine{{
			Machine:           "http://backup.example.com",
			HeadRootDirectory: "/data/backup",
		}},
		OplogStoreFilter: []*StoreFilter{{
			ID:   "5628faffd4c606594adaa3b2",
			Type: "oplogStore",
		}},
		SnapshotStoreFilter: []*StoreFilter{{
			ID:   "5628faffd4c606594adaa3b2",
			Type: "s3blockstore",
		}},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}
