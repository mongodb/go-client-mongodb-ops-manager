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

const (
	blockstoreID = "1"
	fileSystemID = "1"
)

func TestBackupAdministrator_ListBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/mongoConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
			"assignmentEnabled" : true,
			"encryptedCredentials" : false,
			"id" : "1",
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
				ID:                   "1",
				AssignmentEnabled:    true,
				EncryptedCredentials: false,
				LoadFactor:           2,
				MaxCapacityGB:        8,
				URI:                  "mongodb://localhost:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  true,
				UsedSize:             222,
				WriteConcern:         "W2",
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministrator_GetBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/mongoConfigs/%s", blockstoreID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"encryptedCredentials" : false,
			"id" : "1",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 2,
			"maxCapacityGB" : 8,
			"uri" : "mongodb://localhost:27017",
			"ssl" : true,
			"usedSize" : 222,
			"writeConcern" : "W2"
}`)
	})

	config, _, err := client.BackupAdministrator.GetBlockstore(ctx, blockstoreID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetBlockstore returned error: %v", err)
	}

	expected := &Blockstore{
		ID:                   "1",
		AssignmentEnabled:    true,
		EncryptedCredentials: false,
		LoadFactor:           2,
		MaxCapacityGB:        8,
		URI:                  "mongodb://localhost:27017",
		Labels:               []string{"l1", "l2"},
		SSL:                  true,
		UsedSize:             222,
		WriteConcern:         "W2",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministrator_CreateBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/mongoConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"encryptedCredentials" : false,
			"id" : "1",
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
		ID:                   "1",
		AssignmentEnabled:    true,
		EncryptedCredentials: false,
		LoadFactor:           2,
		MaxCapacityGB:        8,
		URI:                  "mongodb://localhost:27017",
		Labels:               []string{"l1", "l2"},
		SSL:                  true,
		UsedSize:             222,
		WriteConcern:         "W2",
	}

	config, _, err := client.BackupAdministrator.CreateBlockstore(ctx, blockstore)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetBlockstore returned error: %v", err)
	}

	expected := &Blockstore{
		ID:                   "1",
		AssignmentEnabled:    true,
		EncryptedCredentials: false,
		LoadFactor:           2,
		MaxCapacityGB:        8,
		URI:                  "mongodb://localhost:27017",
		Labels:               []string{"l1", "l2"},
		SSL:                  true,
		UsedSize:             222,
		WriteConcern:         "W2",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministrator_UpdateBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/mongoConfigs/%s", blockstoreID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"encryptedCredentials" : false,
			"id" : "1",
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
		ID:                   "1",
		AssignmentEnabled:    true,
		EncryptedCredentials: false,
		LoadFactor:           2,
		MaxCapacityGB:        8,
		URI:                  "mongodb://localhost:27017",
		Labels:               []string{"l1", "l2"},
		SSL:                  true,
		UsedSize:             222,
		WriteConcern:         "W2",
	}

	config, _, err := client.BackupAdministrator.UpdateBlockstore(ctx, blockstoreID, blockstore)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateBlockstore returned error: %v", err)
	}

	expected := &Blockstore{
		ID:                   "1",
		AssignmentEnabled:    true,
		EncryptedCredentials: false,
		LoadFactor:           2,
		MaxCapacityGB:        8,
		URI:                  "mongodb://localhost:27017",
		Labels:               []string{"l1", "l2"},
		SSL:                  true,
		UsedSize:             222,
		WriteConcern:         "W2",
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestBackupAdministrator_DeleteBlockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/mongoConfigs/%s", blockstoreID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteBlockstore(ctx, blockstoreID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteBlockstore returned error: %v", err)
	}
}

func TestBackupAdministrator_ListFileSystemStoreConfigurations(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/fileSystemConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
			"assignmentEnabled" : true,
			"id" : "1",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 50,
			"mmapv1CompressionSetting" : "NONE",
			"storePath" : "/data/backup",
			"wtCompressionSetting" : "ZLIB"
		  } ],
		  "totalCount" : 1
}`)
	})

	config, _, err := client.BackupAdministrator.ListFileSystemStoreConfigurations(ctx, nil)
	if err != nil {
		t.Fatalf("BackupAdministrator.ListFileSystemStoreConfigurations returned error: %v", err)
	}

	expected := &FileSystemStoreConfigurations{
		Results: []*FileSystemStoreConfiguration{
			{
				ID:                       "1",
				Labels:                   []string{"l1", "l2"},
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

func TestBackupAdministrator_GetFileSystemStoreConfiguration(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/fileSystemConfigs/%s", fileSystemID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"id" : "1",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 50,
			"mmapv1CompressionSetting" : "NONE",
			"storePath" : "/data/backup",
			"wtCompressionSetting" : "ZLIB"
}`)
	})

	config, _, err := client.BackupAdministrator.GetFileSystemStoreConfiguration(ctx, fileSystemID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetFileSystemStoreConfigurations returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		ID:                       fileSystemID,
		Labels:                   []string{"l1", "l2"},
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

func TestBackupAdministrator_CreateFileSystemStoreConfiguration(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/fileSystemConfigs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"id" : "1",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 50,
			"mmapv1CompressionSetting" : "NONE",
			"storePath" : "/data/backup",
			"wtCompressionSetting" : "ZLIB"
}`)
	})

	fileSystem := &FileSystemStoreConfiguration{
		ID:                       fileSystemID,
		Labels:                   []string{"l1", "l2"},
		LoadFactor:               50,
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
		AssignmentEnabled:        true,
	}

	config, _, err := client.BackupAdministrator.CreateFileSystemStoreConfiguration(ctx, fileSystem)
	if err != nil {
		t.Fatalf("BackupAdministrator.CreateFileSystemStoreConfigurations returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		ID:                       fileSystemID,
		Labels:                   []string{"l1", "l2"},
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

func TestBackupAdministrator_UpdateFileSystemStoreConfiguration(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/fileSystemConfigs/%s", fileSystemID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			"assignmentEnabled" : true,
			"id" : "1",
			"labels" : [ "l1", "l2" ],
			"loadFactor" : 50,
			"mmapv1CompressionSetting" : "NONE",
			"storePath" : "/data/backup",
			"wtCompressionSetting" : "ZLIB"
}`)
	})

	fileSystem := &FileSystemStoreConfiguration{
		ID:                       fileSystemID,
		Labels:                   []string{"l1", "l2"},
		LoadFactor:               50,
		MMAPV1CompressionSetting: "NONE",
		StorePath:                "/data/backup",
		WTCompressionSetting:     "ZLIB",
		AssignmentEnabled:        true,
	}

	config, _, err := client.BackupAdministrator.UpdateFileSystemStoreConfiguration(ctx, fileSystemID, fileSystem)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateFileSystemStoreConfiguration returned error: %v", err)
	}

	expected := &FileSystemStoreConfiguration{
		ID:                       fileSystemID,
		Labels:                   []string{"l1", "l2"},
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

func TestBackupAdministrator_DeleteFileSystemStoreConfiguration(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/fileSystemConfigs/%s", fileSystemID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteFileSystemStoreConfiguration(ctx, blockstoreID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteFileSystemStoreConfiguration returned error: %v", err)
	}
}

func TestBackupAdministrator_ListS3Blockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/backup/snapshot/s3Configs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		  "results" : [ {
			 "assignmentEnabled" : true,
			  "acceptedTos": true,
			  "assignmentEnabled": true,
			  "awsAccessKey": "1",
			  "awsSecretKey": "1",
			  "encryptedCredentials": false,
			  "id": "1",
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
					ID:                   blockstoreID,
					AssignmentEnabled:    true,
					EncryptedCredentials: false,
					LoadFactor:           50,
					URI:                  "mongodb://127.0.0.1:27017",
					Labels:               []string{"l1", "l2"},
					SSL:                  false,
					WriteConcern:         "W2",
				},
				AWSAccessKey:           "1",
				AWSSecretKey:           "1",
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

func TestBackupAdministrator_GetS3Blockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/s3Configs/%s", blockstoreID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			 "assignmentEnabled" : true,
			  "acceptedTos": true,
			  "assignmentEnabled": true,
			  "awsAccessKey": "1",
			  "awsSecretKey": "1",
			  "encryptedCredentials": false,
			  "id": "1",
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

	config, _, err := client.BackupAdministrator.GetS3Blockstore(ctx, blockstoreID)
	if err != nil {
		t.Fatalf("BackupAdministrator.GetS3Blockstore returned error: %v", err)
	}

	expected := &S3Blockstore{
		Blockstore: Blockstore{
			ID:                   blockstoreID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			LoadFactor:           50,
			URI:                  "mongodb://127.0.0.1:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  false,
			WriteConcern:         "W2",
		},
		AWSAccessKey:           "1",
		AWSSecretKey:           "1",
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
			  "awsAccessKey": "1",
			  "awsSecretKey": "1",
			  "encryptedCredentials": false,
			  "id": "1",
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
			ID:                   blockstoreID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			LoadFactor:           50,
			URI:                  "mongodb://127.0.0.1:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  false,
			WriteConcern:         "W2",
		},
		AWSAccessKey:           "1",
		AWSSecretKey:           "1",
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
			ID:                   blockstoreID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			LoadFactor:           50,
			URI:                  "mongodb://127.0.0.1:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  false,
			WriteConcern:         "W2",
		},
		AWSAccessKey:           "1",
		AWSSecretKey:           "1",
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

func TestBackupAdministrator_UpdateS3Blockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/s3Configs/%s", blockstoreID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			 "assignmentEnabled" : true,
			  "acceptedTos": true,
			  "assignmentEnabled": true,
			  "awsAccessKey": "1",
			  "awsSecretKey": "1",
			  "encryptedCredentials": false,
			  "id": "1",
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
			ID:                   blockstoreID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			LoadFactor:           50,
			URI:                  "mongodb://127.0.0.1:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  false,
			WriteConcern:         "W2",
		},
		AWSAccessKey:           "1",
		AWSSecretKey:           "1",
		PathStyleAccessEnabled: false,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            true,
		SSEEnabled:             true,
	}

	config, _, err := client.BackupAdministrator.UpdateS3Blockstore(ctx, blockstoreID, blockstore)
	if err != nil {
		t.Fatalf("BackupAdministrator.UpdateS3Blockstore returned error: %v", err)
	}

	expected := &S3Blockstore{
		Blockstore: Blockstore{
			ID:                   blockstoreID,
			AssignmentEnabled:    true,
			EncryptedCredentials: false,
			LoadFactor:           50,
			URI:                  "mongodb://127.0.0.1:27017",
			Labels:               []string{"l1", "l2"},
			SSL:                  false,
			WriteConcern:         "W2",
		},
		AWSAccessKey:           "1",
		AWSSecretKey:           "1",
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

func TestBackupAdministrator_DeleteS3Blockstore(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/admin/backup/snapshot/s3Configs/%s", blockstoreID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.BackupAdministrator.DeleteS3Blockstore(ctx, blockstoreID)
	if err != nil {
		t.Fatalf("BackupAdministrator.DeleteS3Blockstore returned error: %v", err)
	}
}