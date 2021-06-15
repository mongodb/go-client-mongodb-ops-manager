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
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestAgentsServiceOp_ListAgentLinks(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	if _, _, err := client.Agents.ListAgentLinks(ctx, ""); err == nil {
		t.Error("expected an error but got nil")
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/agents", projectID), func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
  			  "results" :[],
			  "links": [
					{
					  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents",
					  "rel":"self"
					},
					{
					  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1",
					  "rel":"http://mms.mongodb.com/group"
					},
					{
					  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents/MONITORING",
					  "rel":"http://mms.mongodb.com/monitoringAgents"
					},
					{
					  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents/BACKUP",
					  "rel":"http://mms.mongodb.com/backupAgents"
					},
					{
					  "href":"https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents/AUTOMATION",
					  "rel":"http://mms.mongodb.com/automationAgents"
					}
				],
			  "totalCount": 0
		}`)
	})

	links, _, err := client.Agents.ListAgentLinks(ctx, projectID)
	if err != nil {
		t.Fatalf("Agents.ListAgentLinks returned error: %v", err)
	}

	expected := &Agents{
		Links: []*atlas.Link{
			{
				Rel:  "self",
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents",
			},
			{
				Rel:  "http://mms.mongodb.com/group",
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1",
			},
			{
				Rel:  "http://mms.mongodb.com/monitoringAgents",
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents/MONITORING",
			},
			{
				Rel:  "http://mms.mongodb.com/backupAgents",
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents/BACKUP",
			},
			{
				Rel:  "http://mms.mongodb.com/automationAgents",
				Href: "https://cloud.mongodb.com/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents/AUTOMATION",
			},
		},
		Results:    []*Agent{},
		TotalCount: 0,
	}

	if diff := deep.Equal(links, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAgentsServiceOp_ListAgentsByType(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	const agentType = "MONITORING"

	if _, _, err := client.Agents.ListAgentsByType(ctx, "", agentType); err == nil {
		t.Error("expected an error but got nil")
	}

	if _, _, err := client.Agents.ListAgentsByType(ctx, projectID, ""); err == nil {
		t.Error("expected an error but got nil")
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/agents/%s", projectID, agentType), func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
						  "links" : [],
						  "results": [
							{
							  "confCount": 59,
							  "hostname": "example",
							  "isManaged": true,
							  "lastConf": "2015-06-18T14:21:42Z",
							  "lastPing": "2015-06-18T14:21:42Z",
							  "pingCount": 6,
							  "stateName": "ACTIVE",
							  "typeName": "MONITORING"
							}
						  ],
						  "totalCount": 1
						}`)
	})

	agent, _, err := client.Agents.ListAgentsByType(ctx, projectID, agentType)
	if err != nil {
		t.Fatalf("Agents.ListAgentsByType returned error: %v", err)
	}

	expected := &Agents{
		Links: []*atlas.Link{},
		Results: []*Agent{
			{
				TypeName:  "MONITORING",
				Hostname:  "example",
				ConfCount: 59,
				LastConf:  "2015-06-18T14:21:42Z",
				StateName: "ACTIVE",
				PingCount: 6,
				IsManaged: true,
				LastPing:  "2015-06-18T14:21:42Z",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(agent, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAgentsServiceOp_GlobalVersions(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	mux.HandleFunc("/api/public/v1.0/softwareComponents/versions", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
          "automationVersion": "10.14.0.6304",
		  "automationMinimumVersion": "10.2.17.5964",
		  "biConnectorVersion": "2.3.4",
		  "biConnectorMinimumVersion": "2.3.1",
		  "mongoDbToolsVersion": "100.0.1",
		  "links": [
			{
			  "href": "http://mms:9080/api/public/v1.0/softwareComponents/versions",
			  "rel": "self"
			}
		  ]
		}`)
	})

	agent, _, err := client.Agents.GlobalVersions(ctx)
	if err != nil {
		t.Fatalf("Agents.ListAgentsByType returned error: %v", err)
	}

	expected := &SoftwareVersions{
		AutomationVersion:         "10.14.0.6304",
		AutomationMinimumVersion:  "10.2.17.5964",
		BiConnectorVersion:        "2.3.4",
		BiConnectorMinimumVersion: "2.3.1",
		MongoDBToolsVersion:       "100.0.1",
		Links: []*atlas.Link{
			{
				Rel:  "self",
				Href: "http://mms:9080/api/public/v1.0/softwareComponents/versions",
			},
		},
	}

	if diff := deep.Equal(agent, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAgentsServiceOp_ProjectVersions(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	if _, _, err := client.Agents.ProjectVersions(ctx, ""); err == nil {
		t.Error("expected an error but got nil")
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/agents/versions", projectID), func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
		  "count": 0,
		  "entries": [],
		  "isAnyAgentNotManaged": false,
		  "isAnyAgentVersionDeprecated": false,
		  "isAnyAgentVersionOld": false,
		  "latestVersion": "10.14.0.6304",
		  "links": [{
			  "href": "http://mms:9080/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents/current",
			  "rel": "self"
			},
			{
			  "href": "http://mms:9080/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1",
			  "rel": "http://mms.mongodb.com/group"
			}
		  ],
		  "minimumAgentVersionDetected": "10.14.0.6304",
		  "minimumVersion": "5.0.0.309"
		}`)
	})

	agent, _, err := client.Agents.ProjectVersions(ctx, projectID)
	if err != nil {
		t.Fatalf("Agents.ListAgentsByType returned error: %v", err)
	}

	expected := &AgentVersions{
		Count:                       0,
		Entries:                     []*AgentVersion{},
		IsAnyAgentNotManaged:        false,
		IsAnyAgentVersionDeprecated: false,
		IsAnyAgentVersionOld:        false,
		Links: []*atlas.Link{
			{
				Href: "http://mms:9080/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1/agents/current",
				Rel:  "self",
			},
			{
				Href: "http://mms:9080/api/public/v1.0/groups/5e66185d917b220fbd8bb4d1",
				Rel:  "http://mms.mongodb.com/group",
			},
		},
		MinimumAgentVersionDetected: "10.14.0.6304",
		MinimumVersion:              "5.0.0.309",
	}

	if diff := deep.Equal(agent, expected); diff != nil {
		t.Error(diff)
	}
}
