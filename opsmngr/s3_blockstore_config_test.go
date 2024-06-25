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

func TestS3BlockstoreConfigServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/admin/backup/snapshot/s3Configs", func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.S3BlockstoreConfig.List(ctx, nil)
	if err != nil {
		t.Fatalf("S3BlockstoreConfig.List returned error: %v", err)
	}

	assignmentEnabled := true
	ssl := false
	encryptedCredentials := true
	pathStyleAccessEnabled := false
	acceptedTos := true
	sseEnabled := true
	loadFactor := int64(50)

	expected := &S3Blockstores{
		Results: []*S3Blockstore{
			{
				BackupStore: BackupStore{
					LoadFactor: &loadFactor,
					AdminBackupConfig: AdminBackupConfig{
						ID:                   ID,
						AssignmentEnabled:    &assignmentEnabled,
						EncryptedCredentials: &encryptedCredentials,
						URI:                  "mongodb://127.0.0.1:27017",
						Labels:               []string{"l1", "l2"},
						SSL:                  &ssl,
						WriteConcern:         "W2",
						UsedSize:             0,
					},
				},
				AWSAccessKey:           "5628faffd4c606594adaa3b2",
				AWSSecretKey:           "5628faffd4c606594adaa3b2",
				PathStyleAccessEnabled: &pathStyleAccessEnabled,
				S3AuthMethod:           "KEYS",
				S3BucketEndpoint:       "http://example.com/backupbucket",
				S3BucketName:           "bucketname",
				S3MaxConnections:       50,
				AcceptedTos:            &acceptedTos,
				SSEEnabled:             &sseEnabled,
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestS3BlockstoreConfigServiceOp_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/snapshot/s3Configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

	config, _, err := client.S3BlockstoreConfig.Get(ctx, ID)
	if err != nil {
		t.Fatalf("S3BlockstoreConfig.Get returned error: %v", err)
	}

	assignmentEnabled := true
	ssl := false
	encryptedCredentials := false
	pathStyleAccessEnabled := false
	acceptedTos := true
	sseEnabled := true
	loadFactor := int64(50)

	expected := &S3Blockstore{
		BackupStore: BackupStore{
			LoadFactor: &loadFactor,
			AdminBackupConfig: AdminBackupConfig{
				ID:                   ID,
				AssignmentEnabled:    &assignmentEnabled,
				EncryptedCredentials: &encryptedCredentials,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  &ssl,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: &pathStyleAccessEnabled,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            &acceptedTos,
		SSEEnabled:             &sseEnabled,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestS3BlockstoreConfigServiceOp_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/admin/backup/snapshot/s3Configs", func(w http.ResponseWriter, r *http.Request) {
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

	assignmentEnabled := true
	ssl := false
	encryptedCredentials := false
	pathStyleAccessEnabled := false
	acceptedTos := true
	sseEnabled := true
	loadFactor := int64(50)

	blockstore := &S3Blockstore{
		BackupStore: BackupStore{
			LoadFactor: &loadFactor,
			AdminBackupConfig: AdminBackupConfig{
				ID:                   ID,
				AssignmentEnabled:    &assignmentEnabled,
				EncryptedCredentials: &encryptedCredentials,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  &ssl,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: &pathStyleAccessEnabled,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            &acceptedTos,
		SSEEnabled:             &sseEnabled,
	}

	config, _, err := client.S3BlockstoreConfig.Create(ctx, blockstore)
	if err != nil {
		t.Fatalf("S3BlockstoreConfig.Create returned error: %v", err)
	}

	expected := &S3Blockstore{
		BackupStore: BackupStore{
			LoadFactor: &loadFactor,
			AdminBackupConfig: AdminBackupConfig{
				ID:                   ID,
				AssignmentEnabled:    &assignmentEnabled,
				EncryptedCredentials: &encryptedCredentials,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  &ssl,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: &pathStyleAccessEnabled,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            &acceptedTos,
		SSEEnabled:             &sseEnabled,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestS3BlockstoreConfigServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/snapshot/s3Configs/%s", ID), func(w http.ResponseWriter, r *http.Request) {
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

	assignmentEnabled := true
	ssl := false
	encryptedCredentials := false
	pathStyleAccessEnabled := false
	acceptedTos := true
	sseEnabled := true
	loadFactor := int64(50)

	blockstore := &S3Blockstore{
		BackupStore: BackupStore{
			LoadFactor: &loadFactor,
			AdminBackupConfig: AdminBackupConfig{
				ID:                   ID,
				AssignmentEnabled:    &assignmentEnabled,
				EncryptedCredentials: &encryptedCredentials,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  &ssl,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: &pathStyleAccessEnabled,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            &acceptedTos,
		SSEEnabled:             &sseEnabled,
	}

	config, _, err := client.S3BlockstoreConfig.Update(ctx, ID, blockstore)
	if err != nil {
		t.Fatalf("S3BlockstoreConfig.Update returned error: %v", err)
	}

	expected := &S3Blockstore{
		BackupStore: BackupStore{
			LoadFactor: &loadFactor,
			AdminBackupConfig: AdminBackupConfig{
				ID:                   ID,
				AssignmentEnabled:    &assignmentEnabled,
				EncryptedCredentials: &encryptedCredentials,
				URI:                  "mongodb://127.0.0.1:27017",
				Labels:               []string{"l1", "l2"},
				SSL:                  &ssl,
				WriteConcern:         "W2",
				UsedSize:             0,
			},
		},
		AWSAccessKey:           "5628faffd4c606594adaa3b2",
		AWSSecretKey:           "5628faffd4c606594adaa3b2",
		PathStyleAccessEnabled: &pathStyleAccessEnabled,
		S3AuthMethod:           "KEYS",
		S3BucketEndpoint:       "http://example.com/backupbucket",
		S3BucketName:           "bucketname",
		S3MaxConnections:       50,
		AcceptedTos:            &acceptedTos,
		SSEEnabled:             &sseEnabled,
	}

	if diff := deep.Equal(config, expected); diff != nil {
		t.Error(diff)
	}
}

func TestS3BlockstoreConfigServiceOp_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/backup/snapshot/s3Configs/%s", ID), func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.S3BlockstoreConfig.Delete(ctx, ID)
	if err != nil {
		t.Fatalf("S3BlockstoreConfig.Delete returned error: %v", err)
	}
}
