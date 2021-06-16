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
)

const apiDesc = "test-apikeye"

func TestAPIKeys_ListAPIKeys(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/admin/apiKeys", func(w http.ResponseWriter, r *http.Request) {
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
							"roleName": "GLOBAL_OWNER"
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
							"roleName": "GLOBAL_READ_ONLY"
						}
					]
				}
			],
			"totalCount": 2
		}`)
	})

	apiKeys, _, err := client.GlobalAPIKeys.List(ctx, nil)

	if err != nil {
		t.Fatalf("APIKeys.List returned error: %v", err)
	}

	expected := []APIKey{
		{
			ID:         "5c47503320eef5699e1cce8d",
			Desc:       "test-apikey",
			PrivateKey: "********-****-****-db2c132ca78d",
			PublicKey:  "ewmaqvdo",
			Roles: []AtlasRole{
				{
					RoleName: "GLOBAL_OWNER",
				},
			},
		},
		{
			ID:         "5c47503320eef5699e1cce8f",
			Desc:       "test-apikey-2",
			PrivateKey: "********-****-****-db2c132ca78f",
			PublicKey:  "ewmaqvde",
			Roles: []AtlasRole{
				{
					RoleName: "GLOBAL_READ_ONLY",
				},
			},
		},
	}
	if diff := deep.Equal(apiKeys, expected); diff != nil {
		t.Error(diff)
	}
}

func TestAPIKeys_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	createRequest := &APIKeyInput{
		Desc:  "test-apiKey",
		Roles: []string{"GLOBAL_READ_ONLY"},
	}

	mux.HandleFunc("/api/public/v1.0/admin/apiKeys", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"desc":  "test-apiKey",
			"roles": []interface{}{"GLOBAL_READ_ONLY"},
		}

		jsonBlob := `
		{
			"desc": "test-apikeye",
			"id": "5c47503320eef5699e1cce8d",
			"privateKey": "********-****-****-db2c132ca78d",
			"publicKey": "ewmaqvdo",
			"roles": [
				{
					"roleName": "GLOBAL_READ_ONLY"
				}
			]
		}`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, jsonBlob)
	})

	apiKey, _, err := client.GlobalAPIKeys.Create(ctx, createRequest)
	if err != nil {
		t.Fatalf("APIKeys.Create returned error: %v", err)
	}

	if desc := apiKey.Desc; desc != apiDesc {
		t.Errorf("expected username '%s', received '%s'", apiDesc, desc)
	}

	if pk := apiKey.PublicKey; pk != "ewmaqvdo" {
		t.Errorf("expected publicKey '%s', received '%s'", orgID, pk)
	}
}

func TestAPIKeys_GetAPIKey(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/admin/apiKeys/5c47503320eef5699e1cce8d", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"desc":"test-desc"}`)
	})

	apiKeys, _, err := client.GlobalAPIKeys.Get(ctx, "5c47503320eef5699e1cce8d")
	if err != nil {
		t.Errorf("APIKey.Get returned error: %v", err)
	}

	expected := &APIKey{Desc: "test-desc"}

	if diff := deep.Equal(apiKeys, expected); diff != nil {
		t.Errorf("Clusters.Get = %v", diff)
	}
}

func TestAPIKeys_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	updateRequest := &APIKeyInput{
		Desc:  "test-apiKey",
		Roles: []string{"GLOBAL_READ_ONLY"},
	}

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/apiKeys/%s", "5c47503320eef5699e1cce8d"), func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"desc":  "test-apiKey",
			"roles": []interface{}{"GLOBAL_READ_ONLY"},
		}

		jsonBlob := `
		{
			"desc": "test-apikey",
			"id": "5c47503320eef5699e1cce8d",
			"privateKey": "********-****-****-db2c132ca78d",
			"publicKey": "ewmaqvdo",
			"roles": [
				{
					"roleName": "GLOBAL_READ_ONLY"
				}
			]
		}
		`

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

		fmt.Fprint(w, jsonBlob)
	})

	apiKey, _, err := client.GlobalAPIKeys.Update(ctx, "5c47503320eef5699e1cce8d", updateRequest)
	if err != nil {
		t.Fatalf("APIKeys.Create returned error: %v", err)
	}

	if desc := apiKey.Desc; desc != "test-apikey" {
		t.Errorf("expected username '%s', received '%s'", apiDesc, desc)
	}

	if pk := apiKey.PublicKey; pk != "ewmaqvdo" {
		t.Errorf("expected publicKey '%s', received '%s'", orgID, pk)
	}
}

func TestAPIKeys_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()
	apiKeyID := "5c47503320eef5699e1cce8d"

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/admin/apiKeys/%s", apiKeyID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.GlobalAPIKeys.Delete(ctx, apiKeyID)
	if err != nil {
		t.Errorf("APIKey.Delete returned error: %v", err)
	}
}
