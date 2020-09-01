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
	"go.mongodb.org/atlas/mongodbatlas"
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

func TestProjects_ListUsers(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s/%s/users", projectBasePath, groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			"links": [{
				"href": "https://cloud.mongodb.com/api/public/v1.0/groups/users",
				"rel": "self"
			}],
			"results": [{
			   "emailAddress": "someone@example.com",
			   "firstName": "John",
			   "id": "59db8d1d87d9d6420df0613a",
			   "lastName": "Smith",
			   "links": [],
			   "roles": [{
				 "groupId": "59ea02e087d9d636b587a967",
				 "roleName": "GROUP_OWNER"
			   }, {
				 "groupId": "59db8d1d87d9d6420df70902",
				 "roleName": "GROUP_OWNER"
			   }, {
				 "orgId": "59db8d1d87d9d6420df0613f",
				 "roleName": "ORG_OWNER"
			   }],
			   "username": "someone@example.com"
			}, {
			   "emailAddress": "someone_else@example.com",
			   "firstName": "Jill",
			   "id": "59db8d1d87d9d6420df0613a",
			   "lastName": "Smith",
			   "links": [],
			   "roles": [{
				 "groupId": "59ea02e087d9d636b587a967",
				 "roleName": "GROUP_OWNER"
			   }, {
				 "groupId": "59db8d1d87d9d6420df70902",
				 "roleName": "GROUP_OWNER"
			   }, {
				 "orgId": "59db8d1d87d9d6420df0613f",
				 "roleName": "ORG_OWNER"
			   }],
			   "username": "someone@example.com"
			}]
		},`)
	})

	orgs, _, err := client.Projects.ListUsers(ctx, groupID, nil)
	if err != nil {
		t.Fatalf("Projects.ListUsers returned error: %v", err)
	}

	expected := []*User{
		{
			EmailAddress: "someone@example.com",
			FirstName:    "John",
			ID:           "59db8d1d87d9d6420df0613a",
			LastName:     "Smith",
			Links:        []*mongodbatlas.Link{},
			Roles: []*UserRole{
				{GroupID: "59ea02e087d9d636b587a967", RoleName: "GROUP_OWNER"},
				{GroupID: "59db8d1d87d9d6420df70902", RoleName: "GROUP_OWNER"},
				{OrgID: "59db8d1d87d9d6420df0613f", RoleName: "ORG_OWNER"},
			},
			Username: "someone@example.com",
		},
		{
			EmailAddress: "someone_else@example.com",
			FirstName:    "Jill",
			ID:           "59db8d1d87d9d6420df0613a",
			LastName:     "Smith",
			Links:        []*mongodbatlas.Link{},
			Roles: []*UserRole{
				{GroupID: "59ea02e087d9d636b587a967", RoleName: "GROUP_OWNER"},
				{GroupID: "59db8d1d87d9d6420df70902", RoleName: "GROUP_OWNER"},
				{OrgID: "59db8d1d87d9d6420df0613f", RoleName: "ORG_OWNER"},
			},
			Username: "someone@example.com",
		},
	}

	if diff := deep.Equal(orgs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProject_GetOneProject(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s/%s", projectBasePath, groupID), func(w http.ResponseWriter, r *http.Request) {
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
		"id": "5c8100bcf2a30b12ff88258f",
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

	projectResponse, _, err := client.Projects.Get(ctx, groupID)
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
		ID:              "5c8100bcf2a30b12ff88258f",
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
		OrgID: orgID,
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
			"orgId": "5a0a1e7e0f2912c554081adc",
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
		OrgID:            orgID,
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

	mux.HandleFunc(fmt.Sprintf("/groups/%s", groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Projects.Delete(ctx, groupID)
	if err != nil {
		t.Fatalf("Projects.Delete returned error: %v", err)
	}
}

func TestProject_RemoveUser(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	var str = fmt.Sprintf("/groups/%s/users/%s", groupID, userID)
	fmt.Println(str)

	mux.HandleFunc(fmt.Sprintf("/groups/%s/users/%s", groupID, userID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Projects.RemoveUser(ctx, groupID, userID)
	if err != nil {
		t.Fatalf("Projects.RemoveUser returned error: %v", err)
	}
}

func TestProject_AddTeamsToProject(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	projectID := "5a0a1e7e0f2912c554080adc"

	createRequest := []*mongodbatlas.ProjectTeam{{
		TeamID:    "{TEAM-ID}",
		RoleNames: []string{"GROUP_OWNER", "GROUP_READ_ONLY"},
	}}

	mux.HandleFunc(fmt.Sprintf("/groups/%s/teams", projectID), func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{
			"links": [{
				"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams",
				"rel": "self"
			}],
			"results": [{
				"links": [{
					"href": "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams/{TEAM-ID}",
					"rel": "self"
				}],
				"roleNames": ["GROUP_OWNER"],
				"teamId": "{TEAM-ID}"
			}],
			"totalCount": 1
		}`)
	})

	team, _, err := client.Projects.AddTeamsToProject(ctx, projectID, createRequest)
	if err != nil {
		t.Errorf("Projects.AddTeamsToProject returned error: %v", err)
	}

	expected := &mongodbatlas.TeamsAssigned{
		Links: []*mongodbatlas.Link{
			{
				Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams",
				Rel:  "self",
			},
		},
		Results: []*mongodbatlas.Result{
			{
				Links: []*mongodbatlas.Link{
					{
						Href: "https://cloud.mongodb.com/api/atlas/v1.0/groups/{GROUP-ID}/teams/{TEAM-ID}",
						Rel:  "self",
					},
				},
				RoleNames: []string{"GROUP_OWNER"},
				TeamID:    "{TEAM-ID}",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(team, expected); diff != nil {
		t.Error(diff)
	}
}
