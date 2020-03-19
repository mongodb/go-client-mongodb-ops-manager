package opsmngr

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestAgents_ListLinks(t *testing.T) {
	setup()

	defer teardown()
	projectID := "5e66185d917b220fbd8bb4d1"
	mux.HandleFunc(fmt.Sprintf("/groups/%s/agents", projectID), func(w http.ResponseWriter, r *http.Request) {
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

	links, _, err := client.Agents.ListLinks(ctx, projectID)
	if err != nil {
		t.Fatalf("client.Agents.ListLinks returned error: %v", err)
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

func TestAgents_List(t *testing.T) {
	setup()

	defer teardown()
	projectID := "5e66185d917b220fbd8bb4d1"
	agentType := "MONITORING"
	mux.HandleFunc(fmt.Sprintf("/groups/%s/agents/%s", projectID, agentType), func(w http.ResponseWriter, r *http.Request) {
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

	agent, _, err := client.Agents.List(ctx, projectID, agentType)
	if err != nil {
		t.Fatalf("client.Agents.List returned error: %v", err)
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
