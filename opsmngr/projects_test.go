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
	"github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestProject_GetAllProjects(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"links": [{
				"href": "https://cloud.mongodb.com/api/public/v1.0/groups",
				"rel": "self"
			}],
			"results": [{
				"activeAgentCount": 0,
				"hostCounts": {
					"arbiter": 0,	
					"config": 0,
					"master": 0,
					"mongos": 0,
					"primary": 0,
					"secondary": 0,
					"slave": 0
				},
				"id": "56a10a80e4b0fd3b9a9bb0c2",
				"lastActiveAgent": "2016-03-09T18:19:37Z",
				"links": [{
					"href": "https://cloud.mongodb.com/api/public/v1.0/groups/56a10a80e4b0fd3b9a9bb0c2",
					"rel": "self"
				}],
				"name": "012i3091203jioawjioej",
				"orgId": "5980cfdf0b6d97029d82f86e",
				"publicApiEnabled": true,
				"replicaSetCount": 0,
				"shardCount": 0,
				"tags": []
			}, {
				"activeAgentCount": 0,
				"hostCounts": {
					"arbiter": 0,
					"config": 0,
					"master": 0,
					"mongos": 0,
					"primary": 0,
					"secondary": 0,
					"slave": 0
				},
				"id": "56aa691ce4b0a0e8c4be51f7",
				"lastActiveAgent": "2016-01-29T19:02:56Z",
				"links": [{
					"href": "https://cloud.mongodb.com/api/public/v1.0/groups/56aa691ce4b0a0e8c4be51f7",
					"rel": "self"
				}],
				"name": "1454008603036",
				"orgId": "5980d0040b6d97029d831798",
				"publicApiEnabled": true,
				"replicaSetCount": 0,
				"shardCount": 0,
				"tags": []
			}],
			"totalCount": 2
		}`)
	})

	projects, _, err := client.Projects.List(ctx, nil)
	if err != nil {
		t.Fatalf("Projects.List returned error: %v", err)
	}

	expected := &Projects{
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups",
				Rel:  "self",
			},
		},
		Results: []*Project{
			{
				ActiveAgentCount: 0,
				HostCounts: &HostCount{
					Arbiter:   0,
					Config:    0,
					Master:    0,
					Mongos:    0,
					Primary:   0,
					Secondary: 0,
					Slave:     0,
				},
				ID:              "56a10a80e4b0fd3b9a9bb0c2",
				LastActiveAgent: "2016-03-09T18:19:37Z",
				Links: []*mongodbatlas.Link{
					{
						Href: "https://cloud.mongodb.com/api/public/v1.0/groups/56a10a80e4b0fd3b9a9bb0c2",
						Rel:  "self",
					},
				},
				Name:             "012i3091203jioawjioej",
				OrgID:            "5980cfdf0b6d97029d82f86e",
				PublicAPIEnabled: true,
				ReplicaSetCount:  0,
				ShardCount:       0,
				Tags:             []*string{},
			},
			{
				ActiveAgentCount: 0,
				HostCounts: &HostCount{
					Arbiter:   0,
					Config:    0,
					Master:    0,
					Mongos:    0,
					Primary:   0,
					Secondary: 0,
					Slave:     0,
				},
				ID:              "56aa691ce4b0a0e8c4be51f7",
				LastActiveAgent: "2016-01-29T19:02:56Z",
				Links: []*mongodbatlas.Link{
					{
						Href: "https://cloud.mongodb.com/api/public/v1.0/groups/56aa691ce4b0a0e8c4be51f7",
						Rel:  "self",
					},
				},
				Name:             "1454008603036",
				OrgID:            "5980d0040b6d97029d831798",
				PublicAPIEnabled: true,
				ReplicaSetCount:  0,
				ShardCount:       0,
				Tags:             []*string{},
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(projects, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProject_GetOneProject(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/%s/%s", projectBasePath, projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		"activeAgentCount": 0,
		"hostCounts": {
			"arbiter": 0,
			"config": 0,
			"master": 0,
			"mongos": 0,
			"primary": 0,
			"secondary": 0,
			"slave": 0
		},
		"id": "56a10a80e4b0fd3b9a9bb0c2",
		"lastActiveAgent": "2016-03-09T18:19:37Z",
		"links": [{
			"href": "https://cloud.mongodb.com/api/public/v1.0/groups/56a10a80e4b0fd3b9a9bb0c2",
			"rel": "self"
		}],
		"name": "012i3091203jioawjioej",
		"orgId": "5980cfdf0b6d97029d82f86e",
		"publicApiEnabled": true,
		"replicaSetCount": 0,
		"shardCount": 0,
		"tags": []
	  }`)
	})

	projectResponse, _, err := client.Projects.Get(ctx, projectID)
	if err != nil {
		t.Fatalf("Projects.Get returned error: %v", err)
	}

	expected := &Project{
		ActiveAgentCount: 0,
		HostCounts: &HostCount{
			Arbiter:   0,
			Config:    0,
			Master:    0,
			Mongos:    0,
			Primary:   0,
			Secondary: 0,
			Slave:     0,
		},
		ID:              "56a10a80e4b0fd3b9a9bb0c2",
		LastActiveAgent: "2016-03-09T18:19:37Z",
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/56a10a80e4b0fd3b9a9bb0c2",
				Rel:  "self",
			},
		},
		Name:             "012i3091203jioawjioej",
		OrgID:            "5980cfdf0b6d97029d82f86e",
		PublicAPIEnabled: true,
		ReplicaSetCount:  0,
		ShardCount:       0,
		Tags:             []*string{},
	}

	if diff := deep.Equal(projectResponse, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProject_GetOneProjectByName(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	projectName := "012i3091203jioawjioej"

	mux.HandleFunc(fmt.Sprintf("/%s/byName/%s", projectBasePath, projectName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"activeAgentCount": 0,
			"hostCounts": {
				"arbiter": 0,
				"config": 0,
				"master": 0,
				"mongos": 0,
				"primary": 0,
				"secondary": 0,
				"slave": 0
			},
			"id": "56a10a80e4b0fd3b9a9bb0c2",
			"lastActiveAgent": "2016-03-09T18:19:37Z",
			"links": [{
				"href": "https://cloud.mongodb.com/api/public/v1.0/groups/56a10a80e4b0fd3b9a9bb0c2",
				"rel": "self"
			}],
			"name": "012i3091203jioawjioej",
			"orgId": "5980cfdf0b6d97029d82f86e",
			"publicApiEnabled": true,
			"replicaSetCount": 0,
			"shardCount": 0,
			"tags": []
		}`)
	})

	projectResponse, _, err := client.Projects.GetByName(ctx, projectName)
	if err != nil {
		t.Fatalf("Projects.Get returned error: %v", err)
	}

	expected := &Project{
		ActiveAgentCount: 0,
		HostCounts: &HostCount{
			Arbiter:   0,
			Config:    0,
			Master:    0,
			Mongos:    0,
			Primary:   0,
			Secondary: 0,
			Slave:     0,
		},
		ID:              "56a10a80e4b0fd3b9a9bb0c2",
		LastActiveAgent: "2016-03-09T18:19:37Z",
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/56a10a80e4b0fd3b9a9bb0c2",
				Rel:  "self",
			},
		},
		Name:             "012i3091203jioawjioej",
		OrgID:            "5980cfdf0b6d97029d82f86e",
		PublicAPIEnabled: true,
		ReplicaSetCount:  0,
		ShardCount:       0,
		Tags:             []*string{},
	}

	if diff := deep.Equal(projectResponse, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProject_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	createRequest := &Project{
		OrgID: "5a0a1e7e0f2912c554080adc",
		Name:  "ProjectFoobar",
	}

	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
			"activeAgentCount": 0,
			"hostCounts": {
				"arbiter": 0,
				"config": 0,
				"master": 0,
				"mongos": 0,
				"primary": 0,
				"secondary": 0,
				"slave": 0
			},
			"id": "56a10a80e4b0fd3b9a9bb0c2",
			"lastActiveAgent": "2016-03-09T18:19:37Z",
			"links": [{
				"href": "https://cloud.mongodb.com/api/public/v1.0/groups/56a10a80e4b0fd3b9a9bb0c2",
				"rel": "self"
			}],
			"name": "ProjectFoobar",
			"orgId": "5a0a1e7e0f2912c554080adc",
			"publicApiEnabled": true,
			"replicaSetCount": 0,
			"shardCount": 0,
			"tags": []
		}`)
	})

	project, _, err := client.Projects.Create(ctx, createRequest)
	if err != nil {
		t.Fatalf("Projects.Create returned error: %v", err)
	}

	expected := &Project{
		ActiveAgentCount: 0,
		HostCounts: &HostCount{
			Arbiter:   0,
			Config:    0,
			Master:    0,
			Mongos:    0,
			Primary:   0,
			Secondary: 0,
			Slave:     0,
		},
		ID:              "56a10a80e4b0fd3b9a9bb0c2",
		LastActiveAgent: "2016-03-09T18:19:37Z",
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/56a10a80e4b0fd3b9a9bb0c2",
				Rel:  "self",
			},
		},
		Name:             "ProjectFoobar",
		OrgID:            "5a0a1e7e0f2912c554080adc",
		PublicAPIEnabled: true,
		ReplicaSetCount:  0,
		ShardCount:       0,
		Tags:             []*string{},
	}

	if diff := deep.Equal(project, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProject_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/groups/%s", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Projects.Delete(ctx, projectID)
	if err != nil {
		t.Fatalf("Projects.Delete returned error: %v", err)
	}
}
