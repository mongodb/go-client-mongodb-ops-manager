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
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

const accessListID = "5f3cf81b89034c6b3c0a528e" //nolint:gosec // not a credential

func TestWhitelistAPIKeys_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/admin/whitelist", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
		  "results": [
			{
			  "id": "5f3cf81b89034c6b3c0a528e",
			  "cidrBlock": "172.20.0.1",
			  "created": "2020-08-19T13:17:01Z",
			  "description": "test",
			  "type": "GLOBAL_ROLE",
			  "updated": "2020-08-19T13:17:01Z"
			}
		  ],
		  "links": [
			{
			  "rel": "self",
			  "href": "http://mms:8080/api/public/v1.0/admin/whitelist?pageNum=1&itemsPerPage=100"
			}
		  ],
		  "totalCount": 1
		}`)
	})

	whitelistAPIKeys, _, err := client.GlobalAPIKeysWhitelist.List(ctx, nil)
	if err != nil {
		t.Fatalf("GlobalWhitelistAPIKeys.List returned error: %v", err)
	}

	expected := &GlobalWhitelistAPIKeys{
		Links: []*atlas.Link{
			{
				Href: "http://mms:8080/api/public/v1.0/admin/whitelist?pageNum=1&itemsPerPage=100",
				Rel:  "self",
			},
		},
		Results: []*GlobalWhitelistAPIKey{
			{
				ID:          accessListID,
				CidrBlock:   "172.20.0.1",
				Created:     "2020-08-19T13:17:01Z",
				Description: "test",
				Type:        "GLOBAL_ROLE",
				Updated:     "2020-08-19T13:17:01Z",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(whitelistAPIKeys, expected); diff != nil {
		t.Error(diff)
	}
}

func TestWhitelistAPIKeys_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/whitelist/%s", accessListID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"id": "5f3cf81b89034c6b3c0a528e",
			  "cidrBlock": "172.20.0.1",
			  "created": "2020-08-19T13:17:01Z",
			  "description": "test",
			  "type": "GLOBAL_ROLE",
			  "updated": "2020-08-19T13:17:01Z"
		}`)
	})

	whitelistAPIKey, _, err := client.GlobalAPIKeysWhitelist.Get(ctx, accessListID)
	if err != nil {
		t.Fatalf("GlobalWhitelistAPIKeys.Get returned error: %v", err)
	}

	expected := &GlobalWhitelistAPIKey{
		ID:          accessListID,
		CidrBlock:   "172.20.0.1",
		Created:     "2020-08-19T13:17:01Z",
		Description: "test",
		Type:        "GLOBAL_ROLE",
		Updated:     "2020-08-19T13:17:01Z",
	}

	if diff := deep.Equal(whitelistAPIKey, expected); diff != nil {
		t.Error(diff)
	}
}

func TestWhitelistAPIKeys_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	createRequest := &WhitelistAPIKeysReq{
		CidrBlock:   "77.54.32.11/32",
		Description: "test",
	}

	mux.HandleFunc("/api/public/v1.0/admin/whitelist", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"description": "test",
			"cidrBlock":   "77.54.32.11/32",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, `{
		  "id": "5f3cf81b89034c6b3c0a528e",
		  "cidrBlock": "172.20.0.1",
		  "created": "2020-08-19T13:17:01Z",
		  "description": "test",
		  "type": "GLOBAL_ROLE",
		  "updated": "2020-08-19T13:17:01Z"			
		}`)
	})

	whitelistAPIKey, _, err := client.GlobalAPIKeysWhitelist.Create(ctx, createRequest)
	if err != nil {
		t.Fatalf("GlobalWhitelistAPIKeys.Create returned error: %v", err)
	}

	expected := &GlobalWhitelistAPIKey{
		ID:          accessListID,
		CidrBlock:   "172.20.0.1",
		Created:     "2020-08-19T13:17:01Z",
		Description: "test",
		Type:        "GLOBAL_ROLE",
		Updated:     "2020-08-19T13:17:01Z",
	}

	if diff := deep.Equal(whitelistAPIKey, expected); diff != nil {
		t.Error(diff)
	}
}

func TestWhitelistAPIKeys_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/whitelist/%s", accessListID), func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.GlobalAPIKeysWhitelist.Delete(ctx, accessListID)
	if err != nil {
		t.Fatalf("GlobalWhitelistAPIKeys.Delete returned error: %v", err)
	}
}
