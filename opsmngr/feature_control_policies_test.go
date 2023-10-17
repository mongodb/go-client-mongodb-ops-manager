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

func TestFeatureControlPoliciesServiceOp_List(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/groups/%s/controlledFeature", projectID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			  "created": "2019-08-29T15:03:24Z",
			  "updated": "2019-08-29T15:03:24Z",
			  "externalManagementSystem":
			  {
				"name": "Operator",
				"systemId": "6d6c139ae5528707b6e8e3b2",
				"version": "0.2.1"
			  },
			  "policies": [
				{"policy": "ExternallyManagedLock"},
				{"policy": "DisableUserManagement"},
				{"policy": "DisableAuthenticationMechanisms"},
				{"policy": "DisableSetMongodVersion"},
				{
				  "policy": "DisableSetMongodConfig",
				  "disabledParams": ["net.tls.CAFile"]
				}
			  ]
			}`)
	})

	logs, _, err := client.FeatureControlPolicies.List(ctx, projectID, nil)
	if err != nil {
		t.Fatalf("FeatureControlPolicies.List returned error: %v", err)
	}

	expected := &FeaturePolicy{
		Created: "2019-08-29T15:03:24Z",
		Updated: "2019-08-29T15:03:24Z",
		ExternalManagementSystem: &ExternalManagementSystem{
			Name:     "Operator",
			SystemID: "6d6c139ae5528707b6e8e3b2",
			Version:  "0.2.1",
		},
		Policies: []*Policy{
			{
				Policy: "ExternallyManagedLock",
			},
			{
				Policy: "DisableUserManagement",
			},
			{
				Policy: "DisableAuthenticationMechanisms",
			},
			{
				Policy: "DisableSetMongodVersion",
			},
			{
				Policy:         "DisableSetMongodConfig",
				DisabledParams: []string{"net.tls.CAFile"},
			},
		},
	}

	if diff := deep.Equal(logs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestFeatureControlPoliciesServiceOp_Update(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/groups/%s/controlledFeature", projectID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			  "created": "2019-08-29T15:03:24Z",
			  "updated": "2019-08-29T15:03:24Z",
			  "externalManagementSystem":
			  {
				"name": "Operator",
				"systemId": "6d6c139ae5528707b6e8e3b2",
				"version": "0.2.1"
			  },
			  "policies": [
				{"policy": "ExternallyManagedLock"},
				{"policy": "DisableUserManagement"},
				{"policy": "DisableAuthenticationMechanisms"},
				{"policy": "DisableSetMongodVersion"},
				{
				  "policy": "DisableSetMongodConfig",
				  "disabledParams": ["net.tls.CAFile"]
				}
			  ]
			}`)
	})

	policy := &FeaturePolicy{
		Created: "2019-08-29T15:03:24Z",
		Updated: "2019-08-29T15:03:24Z",
		ExternalManagementSystem: &ExternalManagementSystem{
			Name:     "Operator",
			SystemID: "6d6c139ae5528707b6e8e3b2",
			Version:  "0.2.1",
		},
		Policies: []*Policy{
			{
				Policy: "ExternallyManagedLock",
			},
			{
				Policy: "DisableUserManagement",
			},
			{
				Policy: "DisableAuthenticationMechanisms",
			},
			{
				Policy: "DisableSetMongodVersion",
			},
			{
				Policy:         "DisableSetMongodConfig",
				DisabledParams: []string{"net.tls.CAFile"},
			},
		},
	}

	logs, _, err := client.FeatureControlPolicies.Update(ctx, projectID, policy)
	if err != nil {
		t.Fatalf("FeatureControlPolicies.Update returned error: %v", err)
	}

	expected := &FeaturePolicy{
		Created: "2019-08-29T15:03:24Z",
		Updated: "2019-08-29T15:03:24Z",
		ExternalManagementSystem: &ExternalManagementSystem{
			Name:     "Operator",
			SystemID: "6d6c139ae5528707b6e8e3b2",
			Version:  "0.2.1",
		},
		Policies: []*Policy{
			{
				Policy: "ExternallyManagedLock",
			},
			{
				Policy: "DisableUserManagement",
			},
			{
				Policy: "DisableAuthenticationMechanisms",
			},
			{
				Policy: "DisableSetMongodVersion",
			},
			{
				Policy:         "DisableSetMongodConfig",
				DisabledParams: []string{"net.tls.CAFile"},
			},
		},
	}

	if diff := deep.Equal(logs, expected); diff != nil {
		t.Error(diff)
	}
}

func TestFeatureControlPoliciesServiceOp_ListAll(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/public/v1.0/groups/availablePolicies", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
			  "created": "2019-08-29T15:03:24Z",
			  "updated": "2019-08-29T15:03:24Z",
			  "externalManagementSystem":
			  {
				"name": "Operator",
				"systemId": "6d6c139ae5528707b6e8e3b2",
				"version": "0.2.1"
			  },
			  "policies": [
				{"policy": "ExternallyManagedLock"},
				{"policy": "DisableUserManagement"},
				{"policy": "DisableAuthenticationMechanisms"},
				{"policy": "DisableSetMongodVersion"},
				{
				  "policy": "DisableSetMongodConfig",
				  "disabledParams": ["net.tls.CAFile"]
				}
			  ]
			}`)
	})

	logs, _, err := client.FeatureControlPolicies.ListSupportedPolicies(ctx, nil)
	if err != nil {
		t.Fatalf("FeatureControlPolicies.ListSupportedPolicies returned error: %v", err)
	}

	expected := &FeaturePolicy{
		Created: "2019-08-29T15:03:24Z",
		Updated: "2019-08-29T15:03:24Z",
		ExternalManagementSystem: &ExternalManagementSystem{
			Name:     "Operator",
			SystemID: "6d6c139ae5528707b6e8e3b2",
			Version:  "0.2.1",
		},
		Policies: []*Policy{
			{
				Policy: "ExternallyManagedLock",
			},
			{
				Policy: "DisableUserManagement",
			},
			{
				Policy: "DisableAuthenticationMechanisms",
			},
			{
				Policy: "DisableSetMongodVersion",
			},
			{
				Policy:         "DisableSetMongodConfig",
				DisabledParams: []string{"net.tls.CAFile"},
			},
		},
	}

	if diff := deep.Equal(logs, expected); diff != nil {
		t.Error(diff)
	}
}
