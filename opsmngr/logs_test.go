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
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

func TestLogs_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/groups/%s/logCollectionJobs", groupID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			  "results": [
				{
				  "childJobs": [
					{
					  "automationAgentId": "5c810cc4ff7a256345ff97bf",
					  "errorMessage": "null",
					  "finishDate": "2019-03-07T12:21:30Z",
					  "hostName": "server1.example.com",
					  "logCollectionType": "AUTOMATION_AGENT",
					  "path": "server1.example.com/automation_agent",
					  "startDate": "2019-03-07T12:21:24Z",
					  "status": "SUCCESS",
					  "uncompressedDiskSpaceBytes": 14686
					}
				  ],
				  "creationDate": "2019-03-07T12:21:24Z",
				  "downloadUrl": "https://127.0.0.1:8080/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/logCollectionJobs?verbose=true&pageNum=1&itemsPerPage=100",
				  "expirationDate": "2019-04-06T12:21:24Z",
				  "groupId": "5c8100bcf2a30b12ff88258f",
				  "id": "5c810cc4ff7a256345ff97b7",
				  "logTypes": [
					"AUTOMATION_AGENT",
					"MONGODB"
				  ],
				  "redacted": true,
				  "resourceName": "myReplicaSet",
				  "resourceType": "replicaset",
				  "rootResourceName": "myReplicaSet",
				  "rootResourceType": "replicaset",
				  "sizeRequestedPerFileBytes": 1000,
				  "status": "SUCCESS",
				  "uncompressedSizeTotalBytes": 63326,
				  "userId": "5c80f75fcf09a246878f67a4"
				},
				{
				  "childJobs": [
					{
					  "automationAgentId": "5c81086e014b76a3d85e1117",
					  "errorMessage": "null",
					  "finishDate": "2019-03-07T12:02:57Z",
					  "hostName": "server1.example.com:27027",
					  "logCollectionType": "MONGODB",
					  "path": "server1.example.com/27027/mongodb",
					  "startDate": "2019-03-07T12:02:54Z",
					  "status": "SUCCESS",
					  "uncompressedDiskSpaceBytes": 9292
					}
				  ],
				  "creationDate": "2019-03-07T12:02:54Z",
				  "downloadUrl": "https://127.0.0.1:8080/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/logCollectionJobs?verbose=true&pageNum=1&itemsPerPage=100",
				  "expirationDate": "2019-05-06T12:02:54Z",
				  "groupId": "5c8100bcf2a30b12ff88258f",
				  "id": "5c81086e014b76a3d85e1113",
				  "logTypes": [
					"MONGODB",
					"FTDC",
					"AUTOMATION_AGENT"
				  ],
				  "redacted": true,
				  "resourceName": "myReplicaSet",
				  "resourceType": "replicaset",
				  "rootResourceName": "myReplicaSet",
				  "rootResourceType": "replicaset",
				  "sizeRequestedPerFileBytes": 1000,
				  "status": "IN_PROGRESS",
				  "uncompressedSizeTotalBytes": 44518,
				  "userId": "5c80f75fcf09a246878f67a4"
				}
			  ],
			  "totalCount": 2
			}`)
	})

	logs, _, err := client.LogCollections.List(ctx, groupID, nil)
	if err != nil {
		t.Fatalf("LogCollectionJobs.List returned error: %v", err)
	}

	redacted := true

	expected := &LogCollectionJobs{
		Results: []*LogCollectionJob{
			{
				ID:               "5c810cc4ff7a256345ff97b7",
				GroupID:          groupID,
				UserID:           "5c80f75fcf09a246878f67a4",
				CreationDate:     "2019-03-07T12:21:24Z",
				ExpirationDate:   "2019-04-06T12:21:24Z",
				Status:           "SUCCESS",
				ResourceType:     "replicaset",
				ResourceName:     "myReplicaSet",
				RootResourceName: "myReplicaSet",
				RootResourceType: "replicaset",
				URL:              "https://127.0.0.1:8080/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/logCollectionJobs?verbose=true&pageNum=1&itemsPerPage=100",
				Redacted:         &redacted,
				LogTypes: []string{
					"AUTOMATION_AGENT",
					"MONGODB",
				},
				SizeRequestedPerFileBytes:  1000,
				UncompressedDiskSpaceBytes: 63326,
				ChildJobs: []*ChildJob{
					{
						AutomationAgentID:          "5c810cc4ff7a256345ff97bf",
						ErrorMessage:               "null",
						FinishDate:                 "2019-03-07T12:21:30Z",
						HostName:                   "server1.example.com",
						LogCollectionType:          "AUTOMATION_AGENT",
						Path:                       "server1.example.com/automation_agent",
						StartDate:                  "2019-03-07T12:21:24Z",
						Status:                     "SUCCESS",
						UncompressedDiskSpaceBytes: 14686,
					},
				},
			},
			{
				ID:               "5c81086e014b76a3d85e1113",
				GroupID:          groupID,
				UserID:           "5c80f75fcf09a246878f67a4",
				CreationDate:     "2019-03-07T12:02:54Z",
				ExpirationDate:   "2019-05-06T12:02:54Z",
				Status:           "IN_PROGRESS",
				ResourceType:     "replicaset",
				ResourceName:     "myReplicaSet",
				RootResourceName: "myReplicaSet",
				RootResourceType: "replicaset",
				URL:              "https://127.0.0.1:8080/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/logCollectionJobs?verbose=true&pageNum=1&itemsPerPage=100",
				Redacted:         &redacted,
				LogTypes: []string{
					"MONGODB",
					"FTDC",
					"AUTOMATION_AGENT",
				},
				SizeRequestedPerFileBytes:  1000,
				UncompressedDiskSpaceBytes: 44518,
				ChildJobs: []*ChildJob{
					{
						AutomationAgentID:          "5c81086e014b76a3d85e1117",
						ErrorMessage:               "null",
						FinishDate:                 "2019-03-07T12:02:57Z",
						HostName:                   "server1.example.com:27027",
						LogCollectionType:          "MONGODB",
						Path:                       "server1.example.com/27027/mongodb",
						StartDate:                  "2019-03-07T12:02:54Z",
						Status:                     "SUCCESS",
						UncompressedDiskSpaceBytes: 9292,
					},
				},
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(logs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestLogs_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	ID := "1"

	path := fmt.Sprintf("/groups/%s/logCollectionJobs/%s", groupID, ID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
				  "childJobs": [
					{
					  "automationAgentId": "5c810cc4ff7a256345ff97bf",
					  "errorMessage": "null",
					  "finishDate": "2019-03-07T12:21:30Z",
					  "hostName": "server1.example.com",
					  "logCollectionType": "AUTOMATION_AGENT",
					  "path": "server1.example.com/automation_agent",
					  "startDate": "2019-03-07T12:21:24Z",
					  "status": "SUCCESS",
					  "uncompressedDiskSpaceBytes": 14686
					}
				  ],
				  "creationDate": "2019-03-07T12:21:24Z",
				  "downloadUrl": "https://127.0.0.1:8080/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/logCollectionJobs?verbose=true&pageNum=1&itemsPerPage=100",
				  "expirationDate": "2019-04-06T12:21:24Z",
				  "groupId": "5c8100bcf2a30b12ff88258f",
				  "id": "5c810cc4ff7a256345ff97b7",
				  "logTypes": [
					"AUTOMATION_AGENT",
					"MONGODB"
				  ],
				  "redacted": true,
				  "resourceName": "myReplicaSet",
				  "resourceType": "replicaset",
				  "rootResourceName": "myReplicaSet",
				  "rootResourceType": "replicaset",
				  "sizeRequestedPerFileBytes": 1000,
				  "status": "SUCCESS",
				  "uncompressedSizeTotalBytes": 63326,
				  "userId": "5c80f75fcf09a246878f67a4",
                  "id" : "1"
				
			}`)
	})

	opts := LogListOptions{
		Verbose: true,
	}

	logs, _, err := client.LogCollections.Get(ctx, groupID, ID, &opts)
	if err != nil {
		t.Fatalf("LogCollectionJobs.Get returned error: %v", err)
	}

	redacted := true
	expected := &LogCollectionJob{
		ID:               ID,
		GroupID:          groupID,
		UserID:           "5c80f75fcf09a246878f67a4",
		CreationDate:     "2019-03-07T12:21:24Z",
		ExpirationDate:   "2019-04-06T12:21:24Z",
		Status:           "SUCCESS",
		ResourceType:     "replicaset",
		ResourceName:     "myReplicaSet",
		RootResourceName: "myReplicaSet",
		RootResourceType: "replicaset",
		URL:              "https://127.0.0.1:8080/api/public/v1.0/groups/5c8100bcf2a30b12ff88258f/logCollectionJobs?verbose=true&pageNum=1&itemsPerPage=100",
		Redacted:         &redacted,
		LogTypes: []string{
			"AUTOMATION_AGENT",
			"MONGODB",
		},
		SizeRequestedPerFileBytes:  1000,
		UncompressedDiskSpaceBytes: 63326,
		ChildJobs: []*ChildJob{
			{
				AutomationAgentID:          "5c810cc4ff7a256345ff97bf",
				ErrorMessage:               "null",
				FinishDate:                 "2019-03-07T12:21:30Z",
				HostName:                   "server1.example.com",
				LogCollectionType:          "AUTOMATION_AGENT",
				Path:                       "server1.example.com/automation_agent",
				StartDate:                  "2019-03-07T12:21:24Z",
				Status:                     "SUCCESS",
				UncompressedDiskSpaceBytes: 14686,
			},
		},
	}

	if diff := deep.Equal(logs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestLogs_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/groups/%s/logCollectionJobs", groupID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
			"id": "5c81086e014b76a3d85e1113" 
			}`)
	})

	redacted := true
	log := &LogCollectionJob{
		ResourceName: "my_deployment_1",
		ResourceType: "PROCESS",
		Redacted:     &redacted,
		LogTypes: []string{
			"FTDC",
			"MONGODB",
			"AUTOMATION_AGENT",
		},
		SizeRequestedPerFileBytes: 10000000,
	}

	logs, _, err := client.LogCollections.Create(ctx, groupID, log)
	if err != nil {
		t.Fatalf("LogCollectionJobs.Create returned error: %v", err)
	}

	expected := &LogCollectionJob{
		ID: "5c81086e014b76a3d85e1113",
	}

	if diff := deep.Equal(logs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestLogs_Extend(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	ID := "1"

	path := fmt.Sprintf("/groups/%s/logCollectionJobs/%s", groupID, ID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
	})

	log := &LogCollectionJob{
		ExpirationDate: "2019-04-06T12:02:54Z",
	}

	_, err := client.LogCollections.Extend(ctx, groupID, ID, log)
	if err != nil {
		t.Fatalf("LogCollectionJobs.Extend returned error: %v", err)
	}
}

func TestLogs_Retry(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	ID := "1"

	path := fmt.Sprintf("/groups/%s/logCollectionJobs/%s/retry", groupID, ID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
	})

	_, err := client.LogCollections.Retry(ctx, groupID, ID)
	if err != nil {
		t.Fatalf("LogCollectionJobs.Retry returned error: %v", err)
	}
}

func TestLogs_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	ID := "1"

	path := fmt.Sprintf("/groups/%s/logCollectionJobs/%s", groupID, ID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.LogCollections.Delete(ctx, groupID, ID)
	if err != nil {
		t.Fatalf("LogCollectionJobs.Delete returned error: %v", err)
	}
}

func TestLogs_Download(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	ID := "1"

	path := fmt.Sprintf("/groups/%s/logCollectionJobs/%s/download", groupID, ID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, "test")
	})

	buf := new(bytes.Buffer)
	_, err := client.Logs.Download(ctx, groupID, ID, buf)
	if err != nil {
		t.Fatalf("Logs.Download returned error: %v", err)
	}

	if buf.String() != "test" {
		t.Fatalf("Logs.Download returned error: %v", err)
	}
}
