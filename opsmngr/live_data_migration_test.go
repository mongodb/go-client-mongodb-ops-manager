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
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestLiveMigration_ConnectOrganizations(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/orgs/%s/liveExport/migrationLink", orgID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{"status": "SYNCED"}`)
	})

	linkToken := &atlas.LinkToken{LinkToken: "test"}
	response, _, err := client.LiveMigration.ConnectOrganizations(ctx, orgID, linkToken)
	if err != nil {
		t.Fatalf("LiveMigration.ConnectOrganizations returned error: %v", err)
	}

	expected := &ConnectionStatus{Status: "SYNCED"}

	if diff := deep.Equal(response, expected); diff != nil {
		t.Error(diff)
	}
}

func TestLiveMigration_ConnectionStatus(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/orgs/%s/liveExport/migrationLink/status", orgID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
				"status": "SYNCED"
			}`)
	})

	response, _, err := client.LiveMigration.ConnectionStatus(ctx, orgID)
	if err != nil {
		t.Fatalf("LiveMigration.ConnectionStatus returned error: %v", err)
	}

	expected := &ConnectionStatus{Status: "SYNCED"}

	if diff := deep.Equal(response, expected); diff != nil {
		t.Error(diff)
	}
}

func TestLiveMigration_DeleteConnection(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/orgs/%s/liveExport/migrationLink", orgID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.LiveMigration.DeleteConnection(ctx, orgID)
	if err != nil {
		t.Fatalf("LiveMigration.DeleteConnection returned error: %v", err)
	}
}
