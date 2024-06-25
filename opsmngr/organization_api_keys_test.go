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
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

const jsonBlobOrg = `
		{
			"desc": "test-apikey",
			"id": "5c47503320eef5699e1cce8d",
			"privateKey": "********-****-****-db2c132ca78d",
			"publicKey": "ewmaqvdo",
			"roles": [
				{
					"groupId": "1",
					"roleName": "GROUP_OWNER"
				},
				{
					"orgId": "1",
					"roleName": "ORG_MEMBER"
				}
			]
		}
		`

func TestOrganizationAPIKeys_ListAPIKeys(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/orgs/1/apiKeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"results": [
				{
					"desc": "test-apikey",
					"id": "5c47503320eef5699e1cce8d",
					"privateKey": "********-****-****-db2c132ca78d",
					"publicKey": "ewmaqvdo",
					"roles": [
						{
							"groupId": "1",
							"roleName": "GROUP_OWNER"
						},
						{
							"orgId": "1",
							"roleName": "ORG_MEMBER"
						}
					]
				},
				{
					"desc": "test-apikey-2",
					"id": "5c47503320eef5699e1cce8f",
					"privateKey": "********-****-****-db2c132ca78f",
					"publicKey": "ewmaqvde",
					"roles": [
						{
							"groupId": "1",
							"roleName": "GROUP_OWNER"
						},
						{
							"orgId": "1",
							"roleName": "ORG_MEMBER"
						}
					]
				}
			],
			"totalCount": 2
		}`)
	})

	apiKeys, _, err := client.OrganizationAPIKeys.List(ctx, "1", nil)

	if err != nil {
		t.Fatalf("OrganizationAPIKeys.List returned error: %v", err)
	}

	expected := []APIKey{
		{
			ID:         "5c47503320eef5699e1cce8d",
			Desc:       testAPIKey,
			PrivateKey: "********-****-****-db2c132ca78d",
			PublicKey:  ewmaqvdo,
			Roles: []OMRole{
				{
					GroupID:  "1",
					RoleName: "GROUP_OWNER",
				},
				{
					OrgID:    "1",
					RoleName: "ORG_MEMBER",
				},
			},
		},
		{
			ID:         "5c47503320eef5699e1cce8f",
			Desc:       "test-apikey-2",
			PrivateKey: "********-****-****-db2c132ca78f",
			PublicKey:  "ewmaqvde",
			Roles: []OMRole{
				{
					GroupID:  "1",
					RoleName: "GROUP_OWNER",
				},
				{
					OrgID:    "1",
					RoleName: "ORG_MEMBER",
				},
			},
		},
	}
	if diff := deep.Equal(apiKeys, expected); diff != nil {
		t.Error(diff)
	}
}

func TestOrganizationAPIKeys_ListAPIKeysMultiplePages(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/orgs/1/apiKeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		dr := APIKeysResponse{
			Results: []APIKey{
				{
					ID:         "5c47503320eef5699e1cce8d",
					Desc:       testAPIKey,
					PrivateKey: "********-****-****-db2c132ca78d",
					PublicKey:  ewmaqvdo,
					Roles: []OMRole{
						{
							GroupID:  "1",
							RoleName: "GROUP_OWNER",
						},
						{
							OrgID:    "1",
							RoleName: "ORG_MEMBER",
						},
					},
				},
				{
					ID:         "5c47503320eef5699e1cce8f",
					Desc:       "test-apikey-2",
					PrivateKey: "********-****-****-db2c132ca78f",
					PublicKey:  "ewmaqvde",
					Roles: []OMRole{
						{
							GroupID:  "1",
							RoleName: "GROUP_OWNER",
						},
						{
							OrgID:    "1",
							RoleName: "ORG_MEMBER",
						},
					},
				},
			},
			Links: []*Link{
				{Href: "http://example.com/api/atlas/v1.0/orgs/1/apiKeys?pageNum=2&itemsPerPage=2", Rel: "self"},
				{Href: "http://example.com/api/atlas/v1.0/orgs/1/apiKeys?pageNum=2&itemsPerPage=2", Rel: "previous"},
			},
		}

		b, err := json.Marshal(dr)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Fprint(w, string(b))
	})

	_, resp, err := client.OrganizationAPIKeys.List(ctx, "1", nil)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestOrganizationAPIKeys_RetrievePageByNumber(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	jBlob := `
	{
		"links": [
			{
				"href": "http://example.com/api/atlas/v1.0/orgs/1/apikeys?pageNum=1&itemsPerPage=1",
				"rel": "previous"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/orgs/1/apikeys?pageNum=2&itemsPerPage=1",
				"rel": "self"
			},
			{
				"href": "http://example.com/api/atlas/v1.0/orgs/1/apikeys?itemsPerPage=3&pageNum=2",
				"rel": "next"
			}
		],
		"results": [
			{
				"desc": "test-apikey",
				"id": "5c47503320eef5699e1cce8d",
				"privateKey": "********-****-****-db2c132ca78d",
				"publicKey": "ewmaqvdo",
				"roles": [
					{
						"groupId": "1",
						"roleName": "GROUP_OWNER"
					},
					{
						"orgId": "1",
						"roleName": "ORG_MEMBER"
					}
				]
			}
		],
		"totalCount": 3
	}`

	mux.HandleFunc("/api/public/v1.0/orgs/1/apiKeys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, jBlob)
	})

	opt := &ListOptions{PageNum: 2}
	_, resp, err := client.OrganizationAPIKeys.List(ctx, "1", opt)
	if err != nil {
		t.Fatal(err)
	}

	checkCurrentPage(t, resp, 2)
}

func TestOrganizationAPIKeys_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	createRequest := &APIKeyInput{
		Desc:  "test-apiKey",
		Roles: []string{"GROUP_OWNER"},
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/orgs/%s/apiKeys", orgID), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"desc":  "test-apiKey",
			"roles": []interface{}{"GROUP_OWNER"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, jsonBlobOrg)
	})

	apiKey, _, err := client.OrganizationAPIKeys.Create(ctx, orgID, createRequest)
	if err != nil {
		t.Errorf("OrganizationAPIKeys.Create returned error: %v", err)
	}

	if desc := apiKey.Desc; desc != testAPIKey {
		t.Errorf("expected username '%s', received '%s'", "test-apikeye", desc)
	}

	if pk := apiKey.PublicKey; pk != ewmaqvdo {
		t.Errorf("expected publicKey '%s', received '%s'", orgID, pk)
	}
}

func TestOrganizationAPIKeys_GetAPIKey(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/orgs/1/apiKeys/5c47503320eef5699e1cce8d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"desc":"test-desc"}`)
	})

	apiKeys, _, err := client.OrganizationAPIKeys.Get(ctx, "1", "5c47503320eef5699e1cce8d")
	if err != nil {
		t.Errorf("OrganizationAPIKeys.Get returned error: %v", err)
	}

	expected := &APIKey{Desc: "test-desc"}

	if diff := deep.Equal(apiKeys, expected); diff != nil {
		t.Errorf("Clusters.Get = %v", diff)
	}
}

func TestOrganizationAPIKeys_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	updateRequest := &APIKeyInput{
		Desc:  "test-apiKey",
		Roles: []string{"GROUP_OWNER"},
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/orgs/%s/apiKeys/%s", orgID, "5c47503320eef5699e1cce8d"), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"desc":  "test-apiKey",
			"roles": []interface{}{"GROUP_OWNER"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, jsonBlobOrg)
	})

	apiKey, _, err := client.OrganizationAPIKeys.Update(ctx, orgID, "5c47503320eef5699e1cce8d", updateRequest)
	if err != nil {
		t.Fatalf("OrganizationAPIKeys.Create returned error: %v", err)
	}

	if desc := apiKey.Desc; desc != testAPIKey {
		t.Errorf("expected username '%s', received '%s'", "test-apikeye", desc)
	}

	if pk := apiKey.PublicKey; pk != ewmaqvdo {
		t.Errorf("expected publicKey '%s', received '%s'", orgID, pk)
	}
}

func TestOrganizationAPIKeys_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/orgs/%s/apiKeys/%s", orgID, apiKeyID), func(_ http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.OrganizationAPIKeys.Delete(ctx, orgID, apiKeyID)
	if err != nil {
		t.Errorf("OrganizationAPIKeys.Delete returned error: %v", err)
	}
}

func checkCurrentPage(t *testing.T, resp *Response, expectedPage int) {
	t.Helper()
	p, err := resp.CurrentPage()
	if err != nil {
		t.Fatal(err)
	}

	if p != expectedPage {
		t.Fatalf("expected current page to be '%d', was '%d'", expectedPage, p)
	}
}
