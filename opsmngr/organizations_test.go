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

func TestOrganizations_GetAllOrganizations(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/orgs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"links": [{
				"href": "https://cloud.mongodb.com/api/public/v1.0/orgs",
				"rel": "self"
			}],
			"results": [{
				"id": "56a10a80e4b0fd3b9a9bb0c2",
				"links": [{
					"href": "https://cloud.mongodb.com/api/public/v1.0/orgs/56a10a80e4b0fd3b9a9bb0c2",
					"rel": "self"
				}],
				"name": "012i3091203jioawjioej"
			}, {
				"id": "56aa691ce4b0a0e8c4be51f7",
				"links": [{
					"href": "https://cloud.mongodb.com/api/public/v1.0/orgs/56aa691ce4b0a0e8c4be51f7",
					"rel": "self"
				}],
				"name": "1454008603036"
			}],
			"totalCount": 2
		}`)
	})

	orgs, _, err := client.Organizations.GetAllOrganizations(ctx)
	if err != nil {
		t.Fatalf("Organizations.GetAllOrganizations returned error: %v", err)
	}

	expected := &Organizations{
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/orgs",
				Rel:  "self",
			},
		},
		Results: []*Organization{
			{
				ID: "56a10a80e4b0fd3b9a9bb0c2",
				Links: []*mongodbatlas.Link{
					{
						Href: "https://cloud.mongodb.com/api/public/v1.0/orgs/56a10a80e4b0fd3b9a9bb0c2",
						Rel:  "self",
					},
				},
				Name: "012i3091203jioawjioej",
			},
			{
				ID: "56aa691ce4b0a0e8c4be51f7",
				Links: []*mongodbatlas.Link{
					{
						Href: "https://cloud.mongodb.com/api/public/v1.0/orgs/56aa691ce4b0a0e8c4be51f7",
						Rel:  "self",
					},
				},
				Name: "1454008603036",
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(orgs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestOrganizations_GetOneOrganization(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	ID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/%s/%s", orgsBasePath, ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
		"id": "56a10a80e4b0fd3b9a9bb0c2",
		"lastActiveAgent": "2016-03-09T18:19:37Z",
		"links": [{
			"href": "https://cloud.mongodb.com/api/public/v1.0/orgs/56a10a80e4b0fd3b9a9bb0c2",
			"rel": "self"
		}],
		"name": "012i3091203jioawjioej"
	  }`)
	})

	response, _, err := client.Organizations.GetOneOrganization(ctx, ID)
	if err != nil {
		t.Fatalf("Organizations.GetOneOrganization returned error: %v", err)
	}

	expected := &Organization{
		ID: "56a10a80e4b0fd3b9a9bb0c2",
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/orgs/56a10a80e4b0fd3b9a9bb0c2",
				Rel:  "self",
			},
		},
		Name: "012i3091203jioawjioej",
	}

	if diff := deep.Equal(response, expected); diff != nil {
		t.Error(diff)
	}
}

func TestOrganizations_GetProjects(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	ID := "5980cfdf0b6d97029d82f86e"

	mux.HandleFunc(fmt.Sprintf("/%s/%s/groups", orgsBasePath, ID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"links": [{
				"href": "https://cloud.mongodb.com/api/public/v1.0/orgs/5980cfdf0b6d97029d82f86e/groups",
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
			}],
			"totalCount": 1
		}`)
	})

	projects, _, err := client.Organizations.GetProjects(ctx, ID)
	if err != nil {
		t.Fatalf("Organizations.GetProjects returned error: %v", err)
	}

	expected := &Projects{
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/orgs/5980cfdf0b6d97029d82f86e/groups",
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
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(projects, expected); diff != nil {
		t.Error(diff)
	}
}

func TestOrganizations_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	createRequest := &Organization{
		Name: "OrgFoobar",
	}

	mux.HandleFunc("/orgs", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
			"id": "5a0a1e7e0f2912c554080adc",
			"links": [{
				"href": "https://cloud.mongodb.com/api/public/v1.0/orgs/5a0a1e7e0f2912c554080adc",
				"rel": "self"
			}],
			"name": "OrgFoobar"
		}`)
	})

	org, _, err := client.Organizations.Create(ctx, createRequest)
	if err != nil {
		t.Fatalf("Organizations.Create returned error: %v", err)
	}

	expected := &Organization{
		ID: "5a0a1e7e0f2912c554080adc",
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/public/v1.0/orgs/5a0a1e7e0f2912c554080adc",
				Rel:  "self",
			},
		},
		Name: "OrgFoobar",
	}

	if diff := deep.Equal(org, expected); diff != nil {
		t.Error(diff)
	}
}

func TestOrganizations_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	orgID := "5a0a1e7e0f2912c554080adc"

	mux.HandleFunc(fmt.Sprintf("/orgs/%s", orgID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Organizations.Delete(ctx, orgID)
	if err != nil {
		t.Fatalf("Organizations.Delete returned error: %v", err)
	}
}
