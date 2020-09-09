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
)

const ID = "5628faffd4c606594adaa3b2"

func TestMaintenanceWindows_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/groups/%s/maintenanceWindows", groupID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			  "results": [
				{
					"alertTypeNames" : [ "BACKUP" ],
					"created" : "2015-10-22T15:04:31Z",
					"description" : "new description",
					"endDate" : "2015-10-23T23:30:00Z",
					"groupId" : "{PROJECT-ID}",
					"id" : "5628faffd4c606594adaa3b2",
					"startDate" : "2015-10-23T22:00:00Z",
					"updated" : "2015-10-22T15:04:31Z"
				  }, {
					"alertTypeNames" : [ "AGENT", "BACKUP" ],
					"created" : "2015-10-22T15:40:09Z",
					"endDate" : "2015-10-23T23:30:00Z",
					"groupId" : "{PROJECT-ID}",
					"id" : "56290359d4c606594adaafe8",
					"startDate" : "2015-10-23T22:00:00Z",
					"updated" : "2015-10-22T15:40:09Z"
				}
			  ],
			  "totalCount": 2
			}`)
	})

	maintenanceWindows, _, err := client.MaintenanceWindows.List(ctx, groupID)
	if err != nil {
		t.Fatalf("MaintenanceWindows.List returned error: %v", err)
	}

	expected := &MaintenanceWindows{
		Results: []*MaintenanceWindow{
			{
				ID:             "5628faffd4c606594adaa3b2",
				GroupID:        "{PROJECT-ID}",
				Created:        "2015-10-22T15:04:31Z",
				StartDate:      "2015-10-23T22:00:00Z",
				EndDate:        "2015-10-23T23:30:00Z",
				Updated:        "2015-10-22T15:04:31Z",
				AlertTypeNames: []string{"BACKUP"},
				Description:    "new description",
			},
			{
				ID:             "56290359d4c606594adaafe8",
				GroupID:        "{PROJECT-ID}",
				Created:        "2015-10-22T15:40:09Z",
				StartDate:      "2015-10-23T22:00:00Z",
				EndDate:        "2015-10-23T23:30:00Z",
				Updated:        "2015-10-22T15:40:09Z",
				AlertTypeNames: []string{"AGENT", "BACKUP"},
			},
		},
		TotalCount: 2,
	}

	if diff := deep.Equal(maintenanceWindows, expected); diff != nil {
		t.Error(diff)
	}
}

func TestMaintenanceWindows_Get(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()
	path := fmt.Sprintf("/groups/%s/maintenanceWindows/%s", groupID, ID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
					"alertTypeNames" : [ "BACKUP" ],
					"created" : "2015-10-22T15:04:31Z",
					"description" : "new description",
					"endDate" : "2015-10-23T23:30:00Z",
					"groupId" : "{PROJECT-ID}",
					"id" : "5628faffd4c606594adaa3b2",
					"startDate" : "2015-10-23T22:00:00Z",
					"updated" : "2015-10-22T15:04:31Z"
			}`)
	})

	maintenanceWindow, _, err := client.MaintenanceWindows.Get(ctx, groupID, ID)
	if err != nil {
		t.Fatalf("MaintenanceWindows.Get returned error: %v", err)
	}

	expected := &MaintenanceWindow{
		ID:             "5628faffd4c606594adaa3b2",
		GroupID:        "{PROJECT-ID}",
		Created:        "2015-10-22T15:04:31Z",
		StartDate:      "2015-10-23T22:00:00Z",
		EndDate:        "2015-10-23T23:30:00Z",
		Updated:        "2015-10-22T15:04:31Z",
		AlertTypeNames: []string{"BACKUP"},
		Description:    "new description",
	}

	if diff := deep.Equal(maintenanceWindow, expected); diff != nil {
		t.Error(diff)
	}
}

func TestMaintenanceWindows_Create(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()
	path := fmt.Sprintf("/groups/%s/maintenanceWindows", groupID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `{
					"alertTypeNames" : [ "BACKUP" ],
					"created" : "2015-10-22T15:04:31Z",
					"description" : "new description",
					"endDate" : "2015-10-23T23:30:00Z",
					"groupId" : "{PROJECT-ID}",
					"id" : "5628faffd4c606594adaa3b2",
					"startDate" : "2015-10-23T22:00:00Z",
					"updated" : "2015-10-22T15:04:31Z"
			}`)
	})

	maintenance := &MaintenanceWindow{
		EndDate:        "2015-10-23T23:30:00Z",
		Updated:        "2015-10-22T15:04:31Z",
		AlertTypeNames: []string{"BACKUP"},
		Description:    "new description",
	}
	maintenanceWindow, _, err := client.MaintenanceWindows.Create(ctx, groupID, maintenance)
	if err != nil {
		t.Fatalf("MaintenanceWindows.Create returned error: %v", err)
	}

	expected := &MaintenanceWindow{
		ID:             "5628faffd4c606594adaa3b2",
		GroupID:        "{PROJECT-ID}",
		Created:        "2015-10-22T15:04:31Z",
		StartDate:      "2015-10-23T22:00:00Z",
		EndDate:        "2015-10-23T23:30:00Z",
		Updated:        "2015-10-22T15:04:31Z",
		AlertTypeNames: []string{"BACKUP"},
		Description:    "new description",
	}

	if diff := deep.Equal(maintenanceWindow, expected); diff != nil {
		t.Error(diff)
	}
}

func TestMaintenanceWindows_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()
	path := fmt.Sprintf("/groups/%s/maintenanceWindows/%s", groupID, ID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprint(w, `{
					"alertTypeNames" : [ "BACKUP" ],
					"created" : "2015-10-22T15:04:31Z",
					"description" : "new description",
					"endDate" : "2015-10-23T23:30:00Z",
					"groupId" : "{PROJECT-ID}",
					"id" : "5628faffd4c606594adaa3b2",
					"startDate" : "2015-10-23T22:00:00Z",
					"updated" : "2015-10-22T15:04:31Z"
			}`)
	})

	maintenance := &MaintenanceWindow{
		EndDate:        "2015-10-23T23:30:00Z",
		Updated:        "2015-10-22T15:04:31Z",
		AlertTypeNames: []string{"BACKUP"},
		Description:    "new description",
	}
	maintenanceWindow, _, err := client.MaintenanceWindows.Update(ctx, groupID, ID, maintenance)
	if err != nil {
		t.Fatalf("MaintenanceWindows.Update returned error: %v", err)
	}

	expected := &MaintenanceWindow{
		ID:             "5628faffd4c606594adaa3b2",
		GroupID:        "{PROJECT-ID}",
		Created:        "2015-10-22T15:04:31Z",
		StartDate:      "2015-10-23T22:00:00Z",
		EndDate:        "2015-10-23T23:30:00Z",
		Updated:        "2015-10-22T15:04:31Z",
		AlertTypeNames: []string{"BACKUP"},
		Description:    "new description",
	}

	if diff := deep.Equal(maintenanceWindow, expected); diff != nil {
		t.Error(diff)
	}
}

func TestMaintenanceWindows_Delete(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/groups/%s/maintenanceWindows/%s", groupID, ID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.MaintenanceWindows.Delete(ctx, groupID, ID)
	if err != nil {
		t.Fatalf("MaintenanceWindows.Delete returned error: %v", err)
	}
}
