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
	"github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestUnauth_CreateFirstUser(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &User{
		EmailAddress: "jane.doe@example.com",
		Password:     "password",
		FirstName:    "Jane",
		LastName:     "Doe",
	}

	mux.HandleFunc("/unauth/users", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, `{
			"apiKey": "1234abcd-ab12-cd34-ef56-1234abcd1234",
			"user": {
				"emailAddress": "jane.doe@example.com",			
				"id": "1234abcd-ab12-cd34-ef56-1234abcd1235",			
				"links": [
				  {
				   "href" : "https://cloud.mongodb.com/api/public/v1.0/users/1234abcd-ab12-cd34-ef56-1234abcd1235",
				   "rel" : "self"
				  }
				],
				"roles": [
				  {
					"roleName": "GLOBAL_OWNER"
				  }
				],
				"username": "jane.doe@example.com"
			}
		}`)
	})

	user, _, err := client.UnauthUsers.CreateFirstUser(ctx, createRequest, nil)
	if err != nil {
		t.Fatalf("Unauth.CreateFirstUser returned error: %v", err)
	}

	expected := &CreateUserResponse{
		APIKey: "1234abcd-ab12-cd34-ef56-1234abcd1234",
		User: &User{
			EmailAddress: "jane.doe@example.com",
			ID:           "1234abcd-ab12-cd34-ef56-1234abcd1235",
			Links: []*mongodbatlas.Link{
				{
					Href: "https://cloud.mongodb.com/api/public/v1.0/users/1234abcd-ab12-cd34-ef56-1234abcd1235",
					Rel:  "self",
				},
			},
			Roles: []*UserRole{
				{
					RoleName: "GLOBAL_OWNER",
				},
			},
			Username: "jane.doe@example.com",
		},
	}

	if diff := deep.Equal(user, expected); diff != nil {
		t.Error(diff)
	}
}
