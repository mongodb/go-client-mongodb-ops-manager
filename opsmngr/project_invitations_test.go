// Copyright 2021 MongoDB Inc
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

func TestProjects_Invitations(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/invites", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `[
               {
               "createdAt": "2021-02-18T21:05:40Z",
			   "expiresAt": "2021-03-20T21:05:40Z",
			   "id": "5a0a1e7e0f2912c554080adc",
			   "inviterUsername": "admin@example.com",
			   "groupId": "%[1]s",
			   "groupName": "jww-12-16",
			   "roles": [
				   "ORG_OWNER"
			   ],
			   "username": "wyatt.smith@example.com"},
			   {"createdAt": "2021-02-18T21:05:40Z",
			   "expiresAt": "2021-03-20T21:05:40Z",
			   "id": "5a0a1e7e0f2912c554080adc",
			   "inviterUsername": "admin@example.com",
			   "groupId": "%[1]s",
			   "groupName": "jww-12-16",
			   "roles": [
				   "ORG_OWNER"
			   ],
			   "teamIds": ["2"],
			   "username": "wyatt.smith@example.com"}]`,
			projectID,
		)
	})

	invitation, _, err := client.Projects.Invitations(ctx, projectID, nil)
	if err != nil {
		t.Fatalf("Projects.Invitations returned error: %v", err)
	}

	expected := []*Invitation{
		{
			ID:              "5a0a1e7e0f2912c554080adc",
			GroupID:         projectID,
			GroupName:       "jww-12-16",
			CreatedAt:       "2021-02-18T21:05:40Z",
			ExpiresAt:       "2021-03-20T21:05:40Z",
			InviterUsername: "admin@example.com",
			Username:        "wyatt.smith@example.com",
			Roles:           []string{"ORG_OWNER"},
		},
		{
			ID:              "5a0a1e7e0f2912c554080adc",
			GroupID:         projectID,
			GroupName:       "jww-12-16",
			CreatedAt:       "2021-02-18T21:05:40Z",
			ExpiresAt:       "2021-03-20T21:05:40Z",
			InviterUsername: "admin@example.com",
			Username:        "wyatt.smith@example.com",
			Roles:           []string{"ORG_OWNER"},
			TeamIDs:         []string{"2"},
		},
	}

	if diff := deep.Equal(invitation, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProjects_Invitation(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/invites/%s", projectID, invitationID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{
			   "createdAt": "2021-02-18T21:05:40Z",
			   "expiresAt": "2021-03-20T21:05:40Z",
			   "id": "5a0a1e7e0f2912c554080adc",
			   "inviterUsername": "admin@example.com",
			   "groupId": "%[1]s",
			   "groupName": "jww-12-16",
			   "roles": [
				   "ORG_OWNER"
			   ],
			   "username": "wyatt.smith@example.com"
		}`, projectID)
	})

	invitation, _, err := client.Projects.Invitation(ctx, projectID, invitationID)
	if err != nil {
		t.Fatalf("Projects.Invitation returned error: %v", err)
	}

	expected := &Invitation{
		ID:              "5a0a1e7e0f2912c554080adc",
		GroupID:         projectID,
		GroupName:       "jww-12-16",
		CreatedAt:       "2021-02-18T21:05:40Z",
		ExpiresAt:       "2021-03-20T21:05:40Z",
		InviterUsername: "admin@example.com",
		Username:        "wyatt.smith@example.com",
		Roles:           []string{"ORG_OWNER"},
	}

	if diff := deep.Equal(invitation, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProjects_InviteUser(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/invites", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `{
			   "createdAt": "2021-02-18T21:05:40Z",
			   "expiresAt": "2021-03-20T21:05:40Z",
			   "id": "5a0a1e7e0f2912c554080adc",
			   "inviterUsername": "admin@example.com",
			   "groupId": "%[1]s",
			   "groupName": "jww-12-16",
			   "roles": [
				   "ORG_OWNER"
			   ],
			   "username": "wyatt.smith@example.com"
		}`, projectID)
	})

	body := &Invitation{
		GroupID:         projectID,
		GroupName:       "jww-12-16",
		CreatedAt:       "2021-02-18T21:05:40Z",
		ExpiresAt:       "2021-03-20T21:05:40Z",
		InviterUsername: "admin@example.com",
		Username:        "wyatt.smith@example.com",
		Roles:           []string{"ORG_OWNER"},
	}

	invitation, _, err := client.Projects.InviteUser(ctx, projectID, body)
	if err != nil {
		t.Fatalf("Projects.InviteUser returned error: %v", err)
	}

	expected := &Invitation{
		ID:              "5a0a1e7e0f2912c554080adc",
		GroupID:         projectID,
		GroupName:       "jww-12-16",
		CreatedAt:       "2021-02-18T21:05:40Z",
		ExpiresAt:       "2021-03-20T21:05:40Z",
		InviterUsername: "admin@example.com",
		Username:        "wyatt.smith@example.com",
		Roles:           []string{"ORG_OWNER"},
	}

	if diff := deep.Equal(invitation, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProjects_UpdateInvitation(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/invites", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprintf(w, `{
			   "createdAt": "2021-02-18T21:05:40Z",
			   "expiresAt": "2021-03-20T21:05:40Z",
			   "id": "5a0a1e7e0f2912c554080adc",
			   "inviterUsername": "admin@example.com",
			   "groupId": "%[1]s",
			   "groupName": "jww-12-16",
			   "roles": [
				   "ORG_OWNER"
			   ],
			   "username": "wyatt.smith@example.com"
		}`, projectID)
	})

	body := &Invitation{
		GroupID:         projectID,
		GroupName:       "jww-12-16",
		CreatedAt:       "2021-02-18T21:05:40Z",
		ExpiresAt:       "2021-03-20T21:05:40Z",
		InviterUsername: "admin@example.com",
		Username:        "wyatt.smith@example.com",
		Roles:           []string{"ORG_OWNER"},
	}

	invitation, _, err := client.Projects.UpdateInvitation(ctx, projectID, body)
	if err != nil {
		t.Fatalf("Projects.UpdateInvitation returned error: %v", err)
	}

	expected := &Invitation{
		ID:              "5a0a1e7e0f2912c554080adc",
		GroupID:         projectID,
		GroupName:       "jww-12-16",
		CreatedAt:       "2021-02-18T21:05:40Z",
		ExpiresAt:       "2021-03-20T21:05:40Z",
		InviterUsername: "admin@example.com",
		Username:        "wyatt.smith@example.com",
		Roles:           []string{"ORG_OWNER"},
	}

	if diff := deep.Equal(invitation, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProjects_UpdateInvitationByID(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/invites/%s", projectID, invitationID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprintf(w, `{
			   "createdAt": "2021-02-18T21:05:40Z",
			   "expiresAt": "2021-03-20T21:05:40Z",
			   "id": "5a0a1e7e0f2912c554080adc",
			   "inviterUsername": "admin@example.com",
			   "groupId": "%[1]s",
			   "groupName": "jww-12-16",
			   "roles": [
				   "ORG_OWNER"
			   ],
			   "username": "wyatt.smith@example.com"
		}`, projectID)
	})

	body := &Invitation{
		ID:              invitationID,
		GroupID:         projectID,
		GroupName:       "jww-12-16",
		CreatedAt:       "2021-02-18T21:05:40Z",
		ExpiresAt:       "2021-03-20T21:05:40Z",
		InviterUsername: "admin@example.com",
		Username:        "wyatt.smith@example.com",
		Roles:           []string{"ORG_OWNER"},
	}

	invitation, _, err := client.Projects.UpdateInvitationByID(ctx, projectID, invitationID, body)
	if err != nil {
		t.Fatalf("Projects.UpdateInvitationByID returned error: %v", err)
	}

	expected := &Invitation{
		ID:              "5a0a1e7e0f2912c554080adc",
		GroupID:         projectID,
		GroupName:       "jww-12-16",
		CreatedAt:       "2021-02-18T21:05:40Z",
		ExpiresAt:       "2021-03-20T21:05:40Z",
		InviterUsername: "admin@example.com",
		Username:        "wyatt.smith@example.com",
		Roles:           []string{"ORG_OWNER"},
	}

	if diff := deep.Equal(invitation, expected); diff != nil {
		t.Error(diff)
	}
}

func TestProjects_DeleteInvitation(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/invites/%s", projectID, invitationID), func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Projects.DeleteInvitation(ctx, projectID, invitationID)
	if err != nil {
		t.Fatalf("Projects.DeleteInvitation returned error: %v", err)
	}
}
