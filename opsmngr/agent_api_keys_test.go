package opsmngr

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

func TestAgentAPIKeys_List(t *testing.T) {
	setup()

	defer teardown()
	projectID := "5e66185d917b220fbd8bb4d1"
	mux.HandleFunc(fmt.Sprintf("/groups/%s/agentapikeys", projectID), func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `
						[{
						  "_id" : "1",
						  "createdBy" : "PUBLIC_API",
						  "createdIpAddr" : "1",
						  "createdTime" : 1520458807291,
						  "createdUserId" : "21",
						  "desc" : "Agent API Key for this project",
						  "key" : "****************************8b87"
						}, {
						  "_id" : "2",
						  "createdBy" : "PROVISIONING",
						  "createdTime" : 1508871142864,
						  "desc" : "Generated by Provisioning",
						  "key" : "****************************39fe"
						}, {
						  "_id" : "3",
						  "createdBy" : "USER",
						  "createdIpAddr" : "1",
						  "createdTime" : 1507067499083,
						  "createdUserId" : "21",
						  "desc" : "Initial API Key",
						  "key" : "****************************70d7"
						}]
		`)
	})

	agentAPIKeys, _, err := client.AgentAPIKeys.List(ctx, projectID)
	if err != nil {
		t.Fatalf("client.AgentAPIKeys.List returned error: %v", err)
	}

	CreatedUserID := "21"
	CreatedIPAddr := "1"

	expected := &[]AgentAPIKey{
		{
			ID:            "1",
			Key:           "****************************8b87",
			Desc:          "Agent API Key for this project",
			CreatedTime:   1520458807291,
			CreatedUserID: &CreatedUserID,
			CreatedIPAddr: &CreatedIPAddr,
			CreatedBy:     "PUBLIC_API",
		},
		{
			ID:          "2",
			Key:         "****************************39fe",
			Desc:        "Generated by Provisioning",
			CreatedTime: 1508871142864,
			CreatedBy:   "PROVISIONING",
		},
		{
			ID:            "3",
			Key:           "****************************70d7",
			Desc:          "Initial API Key",
			CreatedTime:   1507067499083,
			CreatedUserID: &CreatedUserID,
			CreatedIPAddr: &CreatedIPAddr,
			CreatedBy:     "USER",
		},
	}

	if diff := deep.Equal(agentAPIKeys, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAgentAPIKeys_Create(t *testing.T) {
	setup()

	defer teardown()
	projectID := "5e66185d917b220fbd8bb4d1"
	mux.HandleFunc(fmt.Sprintf("/groups/%s/agentapikeys", projectID), func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
						  "_id" : "1",
						  "createdBy" : "PUBLIC_API",
						  "createdIpAddr" : "1",
						  "createdTime" : 1520458807291,
						  "createdUserId" : "21",
						  "desc" : "TEST",
						  "key" : "****************************8b87"
						}`)
	})

	agentRequest := AgentAPIKeysRequest{Desc: "TEST"}
	agentAPIKey, _, err := client.AgentAPIKeys.Create(ctx, projectID, agentRequest)

	if err != nil {
		t.Fatalf("client.AgentAPIKeys.Create returned error: %v", err)
	}

	CreatedUserID := "21"
	CreatedIPAddr := "1"

	expected := &AgentAPIKey{
		ID:            "1",
		Key:           "****************************8b87",
		Desc:          "TEST",
		CreatedTime:   1520458807291,
		CreatedUserID: &CreatedUserID,
		CreatedIPAddr: &CreatedIPAddr,
		CreatedBy:     "PUBLIC_API",
	}

	if diff := deep.Equal(agentAPIKey, expected); diff != nil {
		t.Error(diff)
	}
}
func TestAgentAPIKeys_Delete(t *testing.T) {
	setup()

	defer teardown()
	projectID := "5e66185d917b220fbd8bb4d1"
	agentAPIKey := "1"
	mux.HandleFunc(fmt.Sprintf("/groups/%s/agentapikeys/%s", projectID, agentAPIKey), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.AgentAPIKeys.Delete(ctx, projectID, agentAPIKey)
	if err != nil {
		t.Fatalf("client.AgentAPIKeys.Delete returned error: %v", err)
	}
}
