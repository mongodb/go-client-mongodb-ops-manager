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
	"go.mongodb.org/atlas/mongodbatlas"
)

const userID = "56a10a80e4b0fd3b9a9bb0c2" //nolint:gosec // not a credential
const userName = "someone@example.com"

func TestUsers_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/users/%s", userID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
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
		}`)
	})

	projects, _, err := client.Users.Get(ctx, userID)
	if err != nil {
		t.Fatalf("Users.Get returned error: %v", err)
	}

	expected := &User{
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
	}

	if diff := deep.Equal(projects, expected); diff != nil {
		t.Error(diff)
	}
}

func TestUsers_GetByName(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/users/byName/%s", userName), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
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
		}`)
	})

	projects, _, err := client.Users.GetByName(ctx, userName)
	if err != nil {
		t.Fatalf("Users.GetByName returned error: %v", err)
	}

	expected := &User{
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
	}

	if diff := deep.Equal(projects, expected); diff != nil {
		t.Error(diff)
	}
}

func TestUsers_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/users", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
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
		}`)
	})

	createRequest := &User{
		Username:     "someone@example.com",
		Password:     "some_pass",
		FirstName:    "John",
		LastName:     "Smith",
		EmailAddress: "john.smith@mongodb.com",
		Links:        nil,
		Roles: []*UserRole{
			{RoleName: "ORG_OWNER", OrgID: "59db8d1d87d9d6420df0613f"},
		},
	}

	userResponse, _, err := client.Users.Create(ctx, createRequest)
	if err != nil {
		t.Fatalf("Users.Create returned error: %v", err)
	}

	expected := &User{
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
	}

	if diff := deep.Equal(userResponse, expected); diff != nil {
		t.Error(diff)
	}
}

func TestUsers_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/users/%s", userID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Users.Delete(ctx, userID)
	if err != nil {
		t.Fatalf("Users.Delete returned error: %v", err)
	}
}
