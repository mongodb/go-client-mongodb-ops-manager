// Copyright 2019 MongoDB Inc
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

func TestAutomation_UpdateAgentVersion(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/api/public/v1.0/groups/%s/automationConfig/updateAgentVersions", projectID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w,
			`{
				  "automationAgentVersion": "10.2.7.5898",
                  "biConnectorVersion": "2.6.1"
				}`,
		)
	})

	agent, _, err := client.Automation.UpdateAgentVersion(ctx, projectID)
	if err != nil {
		t.Fatalf("Automation.UpdateAgentVersion returned error: %v", err)
	}

	expected := &AutomationConfigAgent{
		AutomationAgentVersion: "10.2.7.5898",
		BiConnectorVersion:     "2.6.1",
	}

	if diff := deep.Equal(agent, expected); diff != nil {
		t.Error(diff)
	}
}
