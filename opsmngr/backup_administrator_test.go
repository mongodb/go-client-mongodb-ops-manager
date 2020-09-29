// Copyright 2019 MongoDB Inc
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

func TestBackupAdministratorServiceOp_ListBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/mongoConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
			"assignmentEnabled" : true,
			"encryptedCredentials" : false,
			"id" : "5628faffd4c606594adaa3b2",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 2,
			"maxCapacityGB" : 8,
			"uri" : "mongodb://localhost:27017",
			"ssl" : true,
			"usedSize" : 222,
			"writeConcern" : "W2"
		  } ],
		  "totalCount" : 1
}`)
	})

	config, _, err := client.BackupAdministrator.ListBlockstores(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.ListBlockstores returned error: %v", err)
	}

	expected := &Blockstores{
		Results: []*Blockstore{
			{
				LoadFactor:    2,
				MaxCapacityGB: 8,
				AdminConfig: AdminConfig{
					ID:                   ID,
					AssignmentEnabled:    true,
					EncryptedCredentials: false,
					URI:                  "mongodb://localhost:27017",
					Labels:               []string{"l1", "l2"},
					SSL:                  true,
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

func TestBackupAdministratorServiceOp_GetBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"encryptedCredentials" : false,
			"id" : "5628faffd4c606594adaa3b2",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 2,
			"maxCapacityGB" : 8,
			"uri" : "mongodb://localhost:27017",
			"ssl" : true,
			"usedSize" : 222,
			"writeConcern" : "W2"
}`)
	})

	config, _, err := client.BackupAdministrator.GetBlockstore(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetBlockstore returned error: %v", err)
	}

	expected := &Blockstore{
		LoadFactor:    2,
		MaxCapacityGB: 8,
		AdminConfig: AdminConfig{
			ID:                   ID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  true,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_CreateBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/mongoConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"encryptedCredentials" : false,
			"id" : "5628faffd4c606594adaa3b2",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 2,
			"maxCapacityGB" : 8,
			"uri" : "mongodb://localhost:27017",
			"ssl" : true,
			"usedSize" : 222,
			"writeConcern" : "W2"
}`)
	})

	blockstore := &Blockstore{
		LoadFactor:    2,
		MaxCapacityGB: 8,
		AdminConfig: AdminConfig{
			ID:                   ID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  true,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}

	config, _, err := client.BackupAdministrator.CreateBlockstore(ctx, blockstore)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetBlockstore returned error: %v", err)
	}

	expected := &Blockstore{
		LoadFactor:    2,
		MaxCapacityGB: 8,
		AdminConfig: AdminConfig{
			ID:                   ID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  true,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_UpdateBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"encryptedCredentials" : false,
			"id" : "5628faffd4c606594adaa3b2",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 2,
			"maxCapacityGB" : 8,
			"uri" : "mongodb://localhost:27017",
			"ssl" : true,
			"usedSize" : 222,
			"writeConcern" : "W2"
}`)
	})

	blockstore := &Blockstore{
		LoadFactor:    2,
		MaxCapacityGB: 8,
		AdminConfig: AdminConfig{
			ID:                   ID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  true,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}

	config, _, err := client.BackupAdministrator.UpdateBlockstore(ctx, ID, blockstore)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateBlockstore returned error: %v", err)
	}

	expected := &Blockstore{
		LoadFactor:    2,
		MaxCapacityGB: 8,
		AdminConfig: AdminConfig{
			ID:                   ID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			URI:                  "mongodb://localhost:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  true,
			WriteConcern:         "W2",
			UsedSize:             222,
		},
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_DeleteBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteBlockstore(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteBlockstore returned error: %v", err)
	}
}

func TestBackupAdministratorServiceOp_ListFileSystemStoreConfigurations(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/fileSystemConfigs", func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.BackupAdministrator.ListFileSystemStores(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.ListFileSystemStores returned error: %v", err)
	}

	expected := &FileSystemStoreConfigurations{
		Results: []*FileSystemStoreConfiguration{
			{
				AdminConfig: AdminConfig{
					ID:     ID,
					Labels: []string{"l1", "l2"},
				},
				LoadFactor:               50,
				MMAPV1CompressionSetting: "NONE",
				StorePath:                "/data/backup",
				WTCompressionSetting:     "ZLIB",
				AssignmentEnabled:        true,
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_GetFileSystemStoreConfiguration(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/fileSystemConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.BackupAdministrator.GetFileSystemStore(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetFileSystemStoreConfigurations returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		AdminConfig: AdminConfig{
			ID:     ID,
			Labels: []string{"l1", "l2"},
		},
		LoadFactor:               50,
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
		AssignmentEnabled:        true,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_CreateFileSystemStoreConfiguration(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/fileSystemConfigs", func(w http.ResponseWriter, r *http.Request) {
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
		AdminConfig: AdminConfig{
			ID:     ID,
			Labels: []string{"l1", "l2"},
		},
		LoadFactor:               50,
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
		AssignmentEnabled:        true,
	}

	config, _, err := client.BackupAdministrator.CreateFileSystemStore(ctx, fileSystem)
	if err != nil {
		t.Fatalf("BackupAdministrator.CreateFileSystemStoreConfigurations returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		AdminConfig: AdminConfig{
			ID:     ID,
			Labels: []string{"l1", "l2"},
		},
		LoadFactor:               50,
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
		AssignmentEnabled:        true,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_UpdateFileSystemStoreConfiguration(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/fileSystemConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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
		AdminConfig: AdminConfig{
			ID:     ID,
			Labels: []string{"l1", "l2"},
		},
		LoadFactor:               50,
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
		AssignmentEnabled:        true,
	}

	config, _, err := client.BackupAdministrator.UpdateFileSystemStore(ctx, ID, fileSystem)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateFileSystemStore returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		AdminConfig: AdminConfig{
			ID:     ID,
			Labels: []string{"l1", "l2"},
		},
		LoadFactor:               50,
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
		AssignmentEnabled:        true,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_DeleteFileSystemStoreConfiguration(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/fileSystemConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteFileSystemStore(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteFileSystemStore returned error: %v", err)
	}
}

func TestBackupAdministratorServiceOp_ListS3Blockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/s3Configs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
			 "assignmentEnabled" : true,
			  "acceptedTos": true,
			  "assignmentEnabled": true,
			  "awsAccessKey": "5628faffd4c606594adaa3b2",
			  "awsSecretKey": "5628faffd4c606594adaa3b2",
			  "encryptedCredentials": true,
			  "id": "5628faffd4c606594adaa3b2",
			  "labels": [
				"l1",
				"l2"
			  ],
			  "loadFactor": 50,
			  "pathStyleAccessEnabled": false,
			  "s3AuthMethod": "KEYS",
			  "s3BucketEndpoint": "http://example.com/backupbucket",
			  "s3BucketName": "bucketname",
			  "s3MaxConnections": 50,
			  "sseEnabled": true,
			  "ssl": false,
			  "uri": "mongodb://127.0.0.1:27017",
			  "writeConcern": "W2"
		  } ],
		  "totalCount" : 1
}`)
	})

	config, _, err := client.BackupAdministrator.ListS3Blockstores(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.ListS3Blockstores returned error: %v", err)
	}

	expected := &S3Blockstores{
		Results: []*S3Blockstore{
			{
				Blockstore: Blockstore{
					LoadFactor: 50,
					AdminConfig: AdminConfig{
						ID:                   ID,
						AssignmentEnabled:    true,
						EncryptedCredentials: true,
						URI:                  "mongodb://127.0.0.1:27017",
						Labels:               []string{"l1", "l2"},
						SSL:                  false,
						WriteConcern:         "W2",
						UsedSize:             0,
					},
				},
				AWSAccessKey:           "5628faffd4c606594adaa3b2",
				AWSSecretKey:           "5628faffd4c606594adaa3b2",
				PathStyleAccessEnabled: false,
				S3AuthMethod:           "KEYS",
				S3BucketEndpoint:       "http://example.com/backupbucket",
				S3BucketName:           "bucketname",
				S3MaxConnections:       50,
				AcceptedTos:            true,
				SSEEnabled:             true,
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_GetS3Blockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/s3Configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			 "assignmentEnabled" : true,
			  "acceptedTos": true,
			  "assignmentEnabled": true,
			  "awsAccessKey": "5628faffd4c606594adaa3b2",
			  "awsSecretKey": "5628faffd4c606594adaa3b2",
			  "encryptedCredentials": false,
			  "id": "5628faffd4c606594adaa3b2",
			  "labels": [
				"l1",
				"l2"
			  ],
			  "loadFactor": 50,
			  "pathStyleAccessEnabled": false,
			  "s3AuthMethod": "KEYS",
			  "s3BucketEndpoint": "http://example.com/backupbucket",
			  "s3BucketName": "bucketname",
			  "s3MaxConnections": 50,
			  "sseEnabled": true,
			  "ssl": false,
			  "uri": "mongodb://127.0.0.1:27017",
			  "writeConcern": "W2"
}`)
	})

	config, _, err := client.BackupAdministrator.GetS3Blockstore(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetS3Blockstore returned error: %v", err)
	}

	expected := &S3Blockstore{
		Blockstore: Blockstore{
			LoadFactor: 50,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  false,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: false,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            true,
		SSEEnabled:             true,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_CreateS3BlockstoreBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/s3Configs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
			 "assignmentEnabled" : true,
			  "acceptedTos": true,
			  "assignmentEnabled": true,
			  "awsAccessKey": "5628faffd4c606594adaa3b2",
			  "awsSecretKey": "5628faffd4c606594adaa3b2",
			  "encryptedCredentials": false,
			  "id": "5628faffd4c606594adaa3b2",
			  "labels": [
				"l1",
				"l2"
			  ],
			  "loadFactor": 50,
			  "pathStyleAccessEnabled": false,
			  "s3AuthMethod": "KEYS",
			  "s3BucketEndpoint": "http://example.com/backupbucket",
			  "s3BucketName": "bucketname",
			  "s3MaxConnections": 50,
			  "sseEnabled": true,
			  "ssl": false,
			  "uri": "mongodb://127.0.0.1:27017",
			  "writeConcern": "W2"
}`)
	})

	blockstore := &S3Blockstore{
		Blockstore: Blockstore{
			LoadFactor: 50,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  false,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: false,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            true,
		SSEEnabled:             true,
	}

	config, _, err := client.BackupAdministrator.CreateS3Blockstore(ctx, blockstore)
	if err != nil {
		t.Fatalf("BackupAdministrator.CreateS3Blockstore returned error: %v", err)
	}

	expected := &S3Blockstore{
		Blockstore: Blockstore{
			LoadFactor: 50,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  false,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: false,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            true,
		SSEEnabled:             true,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_UpdateS3Blockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/s3Configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			 "assignmentEnabled" : true,
			  "acceptedTos": true,
			  "assignmentEnabled": true,
			  "awsAccessKey": "5628faffd4c606594adaa3b2",
			  "awsSecretKey": "5628faffd4c606594adaa3b2",
			  "encryptedCredentials": false,
			  "id": "5628faffd4c606594adaa3b2",
			  "labels": [
				"l1",
				"l2"
			  ],
			  "loadFactor": 50,
			  "pathStyleAccessEnabled": false,
			  "s3AuthMethod": "KEYS",
			  "s3BucketEndpoint": "http://example.com/backupbucket",
			  "s3BucketName": "bucketname",
			  "s3MaxConnections": 50,
			  "sseEnabled": true,
			  "ssl": false,
			  "uri": "mongodb://127.0.0.1:27017",
			  "writeConcern": "W2"
}`)
	})

	blockstore := &S3Blockstore{
		Blockstore: Blockstore{
			LoadFactor: 50,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  false,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: false,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            true,
		SSEEnabled:             true,
	}

	config, _, err := client.BackupAdministrator.UpdateS3Blockstore(ctx, ID, blockstore)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateS3Blockstore returned error: %v", err)
	}

	expected := &S3Blockstore{
		Blockstore: Blockstore{
			LoadFactor: 50,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  false,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: false,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            true,
		SSEEnabled:             true,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_DeleteS3Blockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/s3Configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteS3Blockstore(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteS3Blockstore returned error: %v", err)
	}
}

func TestBackupAdministratorServiceOp_ListOplogs(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/oplog/mongoConfigs", func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.BackupAdministrator.ListOplog(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.ListOplog returned error: %v", err)
	}

	expected := &Oplogs{
		Results: []*Oplog{
			{
				Blockstore: Blockstore{
					LoadFactor:    0,
					MaxCapacityGB: 8,
					AdminConfig: AdminConfig{
						ID:                   ID,
						AssignmentEnabled:    true,
						EncryptedCredentials: false,
						URI:                  "mongodb://localhost:27017",
						Labels:               []string{"l1", "l2"},
						SSL:                  true,
						WriteConcern:         "W2",
						UsedSize:             222,
					},
				},
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_GetOplog(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/oplog/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.BackupAdministrator.GetOplog(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetOplog returned error: %v", err)
	}

	expected := &Oplog{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_CreateOplog(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/oplog/mongoConfigs", func(w http.ResponseWriter, r *http.Request) {
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

	oplog := &Oplog{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}

	config, _, err := client.BackupAdministrator.CreateOplog(ctx, oplog)
	if err != nil {
		t.Fatalf("BackupAdministrator.CreateOplog returned error: %v", err)
	}

	expected := &Oplog{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_UpdateOplog(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/oplog/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

	oplog := &Oplog{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}

	config, _, err := client.BackupAdministrator.UpdateOplog(ctx, ID, oplog)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateOplog returned error: %v", err)
	}

	expected := &Oplog{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_DeleteOplog(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/oplog/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteOplog(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteOplog returned error: %v", err)
	}
}

func TestBackupAdministratorServiceOp_ListSyncs(t *testing.T) {
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

	config, _, err := client.BackupAdministrator.ListSyncs(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.ListSyncs returned error: %v", err)
	}

	expected := &Syncs{
		Results: []*Sync{
			{
				Blockstore: Blockstore{
					LoadFactor:    0,
					MaxCapacityGB: 8,
					AdminConfig: AdminConfig{
						ID:                   ID,
						AssignmentEnabled:    true,
						EncryptedCredentials: false,
						URI:                  "mongodb://localhost:27017",
						Labels:               []string{"l1", "l2"},
						SSL:                  true,
						WriteConcern:         "W2",
						UsedSize:             222,
					},
				},
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_GetSync(t *testing.T) {
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

	config, _, err := client.BackupAdministrator.GetSync(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetSync returned error: %v", err)
	}

	expected := &Sync{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_CreateSync(t *testing.T) {
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

	sync := &Sync{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}

	config, _, err := client.BackupAdministrator.CreateSync(ctx, sync)
	if err != nil {
		t.Fatalf("BackupAdministrator.CreateSync returned error: %v", err)
	}

	expected := &Sync{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_UpdateSync(t *testing.T) {
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

	sync := &Sync{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}

	config, _, err := client.BackupAdministrator.UpdateSync(ctx, ID, sync)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateSync returned error: %v", err)
	}

	expected := &Sync{
		Blockstore: Blockstore{
			LoadFactor:    0,
			MaxCapacityGB: 8,
			AdminConfig: AdminConfig{
				ID:                   ID,
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				WriteConcern:         "W2",
				UsedSize:             222,
			},
		},
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministratorServiceOp_DeleteSync(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/sync/mongoConfigs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteSync(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteSync returned error: %v", err)
	}
}

func TestBackupAdministratorServiceOp_ListDaemons(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/daemon/configs", func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.BackupAdministrator.ListDaemons(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.ListDaemons returned error: %v", err)
	}

	expected := &Daemons{
		Results: []*Daemon{
			{
				AdminConfig: AdminConfig{
					ID:                   ID,
					AssignmentEnabled:    true,
					EncryptedCredentials: false,
					Labels:               []string{"l1", "l2"},
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

func TestBackupAdministratorServiceOp_GetDaemon(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/daemon/configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.BackupAdministrator.GetDaemon(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetDaemon returned error: %v", err)
	}

	expected := &Daemon{
		AdminConfig: AdminConfig{
			ID:                   ID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			Labels:               []string{"l1", "l2"},
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

func TestBackupAdministratorServiceOp_UpdateDaemon(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/daemon/configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

	deamon := &Daemon{
		AdminConfig: AdminConfig{
			ID:                   ID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			Labels:               []string{"l1", "l2"},
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

	config, _, err := client.BackupAdministrator.UpdateDaemon(ctx, ID, deamon)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateDaemon returned error: %v", err)
	}

	expected := &Daemon{
		AdminConfig: AdminConfig{
			ID:                   ID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			Labels:               []string{"l1", "l2"},
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

func TestBackupAdministratorServiceOp_DeleteDaemn(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/daemon/configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteDaemon(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteDaemon returned error: %v", err)
	}
}

func TestBackupAdministratorServiceOp_ListProjectJobs(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/groups", func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.BackupAdministrator.ListProjectJobs(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.ListProjectJobs returned error: %v", err)
	}

	expected := &ProjectJobs{
		Results: []*ProjectJob{
			{
				AdminConfig: AdminConfig{
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

func TestBackupAdministratorServiceOp_GetProjectJob(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/groups/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.BackupAdministrator.GetProjectJob(ctx, ID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetProjectJob returned error: %v", err)
	}

	expected := &ProjectJob{

		AdminConfig: AdminConfig{
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

func TestBackupAdministratorServiceOp_UpdateProjectJob(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/groups/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

		AdminConfig: AdminConfig{
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

	config, _, err := client.BackupAdministrator.UpdateProjectJob(ctx, ID, projectJob)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateProjectJob returned error: %v", err)
	}

	expected := &ProjectJob{

		AdminConfig: AdminConfig{
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
