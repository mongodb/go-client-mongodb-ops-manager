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
	"bytes"
	"fmt"
	"net/http"
	"testing"

	"github.com/go-test/deep"
)

func TestServerUsageServiceOp_GenerateDailyUsageSnapshot(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/usage/dailyCapture", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	_, err := client.ServerUsage.GenerateDailyUsageSnapshot(ctx)
	if err != nil {
		t.Fatalf("ServerUsage.GenerateDailyUsageSnapshot returned error: %v", err)
	}
}

func TestServerUsageServiceOp_UpdateProjectServerType(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/usage/groups/%s/defaultServerType", groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			   "serverType":{
				  "name":"RAM_POOL",
				  "label":"RAM Pool"
			   }
		}`)
	})
	_, err := client.ServerUsage.UpdateProjectServerType(ctx, groupID, nil)

	if err != nil {
		t.Fatalf("ServerUsage.UpdateProjectServerType returned error: %v", err)
	}
}

func TestServerUsageServiceOp_UpdateOrganizationServerType(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/usage/organizations/%s/defaultServerType", orgID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		_, _ = fmt.Fprint(w, `{
			   "serverType":{
				  "name":"RAM_POOL",
				  "label":"RAM Pool"
			   }
	}`)
	})
	_, err := client.ServerUsage.UpdateOrganizationServerType(ctx, orgID, nil)

	if err != nil {
		t.Fatalf("ServerUsage.UpdateOrganizationServerType returned error: %v", err)
	}
}

func TestServerUsageReportServiceOp_Download(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/usage/report", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, "testFile")
	})

	buf := new(bytes.Buffer)
	_, err := client.ServerUsageReport.Download(ctx, nil, buf)
	if err != nil {
		t.Fatalf("ServerUsageReport.Download returned error: %v", err)
	}

	if buf.String() != "testFile" {
		t.Fatalf("ServerUsageReport.Download returned error: %v", err)
	}
}

func TestServerUsageServiceOp_ListAllHostAssignment(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc("/usage/assignments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
					 "totalCount": 1,
					 "results": [{
					   "hostname": "virtual.host.lqhfcxlgzqtimcxf.internal.mongodb-2",
					   "processes": [{
						 "cluster": "sdivabux",
						 "groupName": "test",
						 "orgName": "5a0a1e7e0f2912c554081adc",
						 "groupId": "5c8100bcf2a30b12ff88258f",
						 "hasConflictingServerType": true,
						 "name": "replicaSecondary-0-proc1-run51839",
						 "processType": 8
					   }
					 ],
					   "serverType": {
						 "name": "RAM_POOL",
						 "label": "RAM Pool"
					   },
					   "isChargeable": true,
					   "memSizeMB": 178
					 }]
}`)
	})

	hostAssignments, _, err := client.ServerUsage.ListAllHostAssignment(ctx, nil)
	if err != nil {
		t.Fatalf("ServerUsage.ListAllHostAssignment returned error: %v", err)
	}

	isChargeable := true
	expected := &HostAssignments{
		Results: []*HostAssignment{
			{
				Hostname: "virtual.host.lqhfcxlgzqtimcxf.internal.mongodb-2",
				Processes: []*HostAssignmentProcess{
					{
						Cluster:                  "sdivabux",
						GroupName:                "test",
						OrgName:                  orgID,
						GroupID:                  groupID,
						Name:                     "replicaSecondary-0-proc1-run51839",
						HasConflictingServerType: true,
						ProcessType:              8,
					},
				},
				ServerType: &ServerType{
					Name:  "RAM_POOL",
					Label: "RAM Pool",
				},
				MemSizeMB:    178,
				IsChargeable: &isChargeable,
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(hostAssignments, expected); diff != nil {
		t.Error(diff)
	}
}

func TestServerUsageServiceOp_ProjectHostAssignments(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/usage/groups/%s/hosts", groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
					 "totalCount": 1,
					 "results": [{
					   "hostname": "virtual.host.lqhfcxlgzqtimcxf.internal.mongodb-2",
					   "processes": [{
						 "cluster": "sdivabux",
						 "groupName": "test",
						 "orgName": "5a0a1e7e0f2912c554081adc",
						 "groupId": "5c8100bcf2a30b12ff88258f",
						 "hasConflictingServerType": true,
						 "name": "replicaSecondary-0-proc1-run51839",
						 "processType": 8
					   }
					 ],
					   "serverType": {
						 "name": "RAM_POOL",
						 "label": "RAM Pool"
					   },
					   "isChargeable": true,
					   "memSizeMB": 178
					 }]
}`)
	})

	hostAssignments, _, err := client.ServerUsage.ProjectHostAssignments(ctx, groupID, nil)
	if err != nil {
		t.Fatalf("ServerUsage.ProjectHostAssignments returned error: %v", err)
	}

	isChargeable := true
	expected := &HostAssignments{
		Results: []*HostAssignment{
			{
				Hostname: "virtual.host.lqhfcxlgzqtimcxf.internal.mongodb-2",
				Processes: []*HostAssignmentProcess{
					{
						Cluster:                  "sdivabux",
						GroupName:                "test",
						OrgName:                  orgID,
						GroupID:                  groupID,
						Name:                     "replicaSecondary-0-proc1-run51839",
						HasConflictingServerType: true,
						ProcessType:              8,
					},
				},
				ServerType: &ServerType{
					Name:  "RAM_POOL",
					Label: "RAM Pool",
				},
				MemSizeMB:    178,
				IsChargeable: &isChargeable,
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(hostAssignments, expected); diff != nil {
		t.Error(diff)
	}
}

func TestServerUsageServiceOp_OrganizationHostAssignments(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/usage/organizations/%s/hosts", orgID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
					 "totalCount": 1,
					 "results": [{
					   "hostname": "virtual.host.lqhfcxlgzqtimcxf.internal.mongodb-2",
					   "processes": [{
						 "cluster": "sdivabux",
						 "groupName": "test",
						 "orgName": "5a0a1e7e0f2912c554081adc",
						 "groupId": "5c8100bcf2a30b12ff88258f",
						 "hasConflictingServerType": true,
						 "name": "replicaSecondary-0-proc1-run51839",
						 "processType": 8
					   }
					 ],
					   "serverType": {
						 "name": "RAM_POOL",
						 "label": "RAM Pool"
					   },
					   "isChargeable": true,
					   "memSizeMB": 178
					 }]
}`)
	})

	hostAssignments, _, err := client.ServerUsage.OrganizationHostAssignments(ctx, orgID, nil)
	if err != nil {
		t.Fatalf("ServerUsage.OrganizationHostAssignments returned error: %v", err)
	}

	isChargeable := true
	expected := &HostAssignments{
		Results: []*HostAssignment{
			{
				Hostname: "virtual.host.lqhfcxlgzqtimcxf.internal.mongodb-2",
				Processes: []*HostAssignmentProcess{
					{
						Cluster:                  "sdivabux",
						GroupName:                "test",
						OrgName:                  orgID,
						GroupID:                  groupID,
						Name:                     "replicaSecondary-0-proc1-run51839",
						HasConflictingServerType: true,
						ProcessType:              8,
					},
				},
				ServerType: &ServerType{
					Name:  "RAM_POOL",
					Label: "RAM Pool",
				},
				MemSizeMB:    178,
				IsChargeable: &isChargeable,
			},
		},
		TotalCount: 1,
	}
	if diff := deep.Equal(hostAssignments, expected); diff != nil {
		t.Error(diff)
	}
}

func TestServerUsageServiceOp_GetServerTypeOrganization(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/usage/organizations/%s/defaultServerType", orgID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
						 "name": "RAM_POOL",
						 "label": "RAM Pool"
}`)
	})

	serverType, _, err := client.ServerUsage.GetServerTypeOrganization(ctx, orgID)
	if err != nil {
		t.Fatalf("ServerUsage.GetServerTypeOrganization returned error: %v", err)
	}

	expected := &ServerType{
		Name:  "RAM_POOL",
		Label: "RAM Pool",
	}

	if diff := deep.Equal(serverType, expected); diff != nil {
		t.Error(diff)
	}
}

func TestServerUsageServiceOp_GetServerTypeProject(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/usage/groups/%s/defaultServerType", groupID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
						 "name": "RAM_POOL",
						 "label": "RAM Pool"
}`)
	})

	serverType, _, err := client.ServerUsage.GetServerTypeProject(ctx, groupID)
	if err != nil {
		t.Fatalf("ServerUsage.GetServerTypeProject returned error: %v", err)
	}

	expected := &ServerType{
		Name:  "RAM_POOL",
		Label: "RAM Pool",
	}

	if diff := deep.Equal(serverType, expected); diff != nil {
		t.Error(diff)
	}
}
