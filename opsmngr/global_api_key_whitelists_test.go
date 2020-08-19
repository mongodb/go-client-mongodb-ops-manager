package opsmngr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestWhitelistAPIKeys_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/admin/whitelist", func(w http.ResponseWriter, r *http.Request) {
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
				ID:          "5f3cf81b89034c6b3c0a528e",
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

	ipAddress := "5f3cf81b89034c6b3c0a528e"

	mux.HandleFunc(fmt.Sprintf("/admin/whitelist/%s", ipAddress), func(w http.ResponseWriter, r *http.Request) {
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

	whitelistAPIKey, _, err := client.GlobalAPIKeysWhitelist.Get(ctx, ipAddress)
	if err != nil {
		t.Fatalf("GlobalWhitelistAPIKeys.Get returned error: %v", err)
	}

	expected := &GlobalWhitelistAPIKey{
		ID:          "5f3cf81b89034c6b3c0a528e",
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

	createRequest := []*WhitelistAPIKeysReq{
		{
			CidrBlock:   "77.54.32.11/32",
			Description: "test",
		},
	}

	mux.HandleFunc("/admin/whitelist", func(w http.ResponseWriter, r *http.Request) {
		expected := []map[string]interface{}{
			{
				"description": "test",
				"cidrBlock":   "77.54.32.11/32",
			},
		}

		var v []map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("Decode json: %v", err)
		}

		if diff := deep.Equal(v, expected); diff != nil {
			t.Error(diff)
		}

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

	whitelistAPIKey, _, err := client.GlobalAPIKeysWhitelist.Create(ctx, createRequest)
	if err != nil {
		t.Fatalf("GlobalWhitelistAPIKeys.Create returned error: %v", err)
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
				ID:          "5f3cf81b89034c6b3c0a528e",
				CidrBlock:   "172.20.0.1",
				Created:     "2020-08-19T13:17:01Z",
				Description: "test",
				Type:        "GLOBAL_ROLE",
				Updated:     "2020-08-19T13:17:01Z",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(whitelistAPIKey, expected); diff != nil {
		t.Error(diff)
	}
}

func TestWhitelistAPIKeys_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	ipAddress := "5f3cf81b89034c6b3c0a528e"
	mux.HandleFunc(fmt.Sprintf("/admin/whitelist/%s", ipAddress), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.GlobalAPIKeysWhitelist.Delete(ctx, ipAddress)
	if err != nil {
		t.Fatalf("GlobalWhitelistAPIKeys.Delete returned error: %v", err)
	}
}
