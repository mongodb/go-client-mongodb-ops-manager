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
	atlas "go.mongodb.org/atlas/mongodbatlas"
)

func TestDeployments_ListHosts(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	path := fmt.Sprintf("/api/public/v1.0/groups/%s/hosts", projectID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{
				  "totalCount" : 1,
				  "results" : [
					{
					  "alertsEnabled" : true,
					  "aliases": [ "{HOSTNAME}:26000", "{IP-ADDRESS}:26000" ],
					  "authMechanismName" : "SCRAM-SHA-1",
					  "clusterId" : "1",
					  "created" : "2014-04-22T19:56:50Z",
					  "deactivated" : false,
					  "groupId" : "%[1]s",
					  "hasStartupWarnings" : false,
					  "hidden" : false,
					  "hostEnabled" : true,
					  "hostname" : "mongoHost",
					  "id" : "22",
					  "ipAddress": "127.0.0.1",
					  "journalingEnabled" : false,
					  "lastDataSizeBytes" : 633208918,
					  "lastIndexSizeBytes" : 101420524,
					  "lastPing" : "2016-08-18T11:23:41Z",
					  "links" : [  ],
					  "logsEnabled" : false,
					  "lowUlimit" : false,
					  "port" : 26000,
					  "profilerEnabled" : false,
					  "replicaSetName": "rs1",
					  "replicaStateName" : "PRIMARY",
					  "sslEnabled" : true,
					  "typeName": "REPLICA_PRIMARY",
					  "uptimeMsec": 1827300394,
					  "username" : "mongo",
					  "version" : "3.2.0"
					}
				  ]
			}`, projectID)
	})

	opts := HostListOptions{}

	alerts, _, err := client.Deployments.ListHosts(ctx, projectID, &opts)
	if err != nil {
		t.Fatalf("Deployments.ListHosts returned error: %v", err)
	}

	alertsEnabled := true
	logsEnabled := false
	profilerEnabled := false
	sslEnabled := true

	expected := &Hosts{
		Links: nil,
		Results: []*Host{
			{
				Aliases: []string{
					"{HOSTNAME}:26000", "{IP-ADDRESS}:26000",
				},
				AlertsEnabled:      &alertsEnabled,
				AuthMechanismName:  "SCRAM-SHA-1",
				ClusterID:          "1",
				Created:            "2014-04-22T19:56:50Z",
				Deactivated:        false,
				GroupID:            projectID,
				HasStartupWarnings: false,
				Hidden:             false,
				HostEnabled:        true,
				Hostname:           "mongoHost",
				ID:                 "22",
				IPAddress:          "127.0.0.1",
				JournalingEnabled:  false,
				LastDataSizeBytes:  633208918,
				LastIndexSizeBytes: 101420524,
				LastPing:           "2016-08-18T11:23:41Z",
				Links:              []*atlas.Link{},
				LogsEnabled:        &logsEnabled,
				LowUlimit:          false,
				Port:               26000,
				ProfilerEnabled:    &profilerEnabled,
				ReplicaSetName:     "rs1",
				ReplicaStateName:   "PRIMARY",
				SSLEnabled:         &sslEnabled,
				TypeName:           "REPLICA_PRIMARY",
				UptimeMsec:         1827300394,
				Version:            "3.2.0",
				Username:           "mongo",
			},
		},
		TotalCount: 1,
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDeployments_GetHost(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	hostID := "22"
	path := fmt.Sprintf("/api/public/v1.0/groups/%s/hosts/%s", projectID, hostID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `
					{
					  "alertsEnabled" : true,
					  "aliases": [ "{HOSTNAME}:26000", "{IP-ADDRESS}:26000" ],
					  "authMechanismName" : "SCRAM-SHA-1",
					  "clusterId" : "1",
					  "created" : "2014-04-22T19:56:50Z",
					  "deactivated" : false,
					  "groupId" : "%[1]s",
					  "hasStartupWarnings" : false,
					  "hidden" : false,
					  "hostEnabled" : true,
					  "hostname" : "mongoHost",
					  "id" : "22",
					  "ipAddress": "127.0.0.1",
					  "journalingEnabled" : false,
					  "lastDataSizeBytes" : 633208918,
					  "lastIndexSizeBytes" : 101420524,
					  "lastPing" : "2016-08-18T11:23:41Z",
					  "links" : [  ],
					  "logsEnabled" : false,
					  "lowUlimit" : false,
					  "port" : 26000,
					  "profilerEnabled" : false,
					  "replicaSetName": "rs1",
					  "replicaStateName" : "PRIMARY",
					  "sslEnabled" : true,
					  "typeName": "REPLICA_PRIMARY",
					  "uptimeMsec": 1827300394,
					  "username" : "mongo",
					  "version" : "3.2.0"
					}`, projectID)
	})

	alerts, _, err := client.Deployments.GetHost(ctx, projectID, hostID)
	if err != nil {
		t.Fatalf("Deployments.GetHost returned error: %v", err)
	}

	alertsEnabled := true
	logsEnabled := false
	profilerEnabled := false
	sslEnabled := true

	expected := &Host{
		Aliases: []string{
			"{HOSTNAME}:26000", "{IP-ADDRESS}:26000",
		},
		AlertsEnabled:      &alertsEnabled,
		AuthMechanismName:  "SCRAM-SHA-1",
		ClusterID:          "1",
		Created:            "2014-04-22T19:56:50Z",
		Deactivated:        false,
		GroupID:            projectID,
		HasStartupWarnings: false,
		Hidden:             false,
		HostEnabled:        true,
		Hostname:           "mongoHost",
		ID:                 "22",
		IPAddress:          "127.0.0.1",
		JournalingEnabled:  false,
		LastDataSizeBytes:  633208918,
		LastIndexSizeBytes: 101420524,
		LastPing:           "2016-08-18T11:23:41Z",
		Links:              []*atlas.Link{},
		LogsEnabled:        &logsEnabled,
		LowUlimit:          false,
		Port:               26000,
		ProfilerEnabled:    &profilerEnabled,
		ReplicaSetName:     "rs1",
		ReplicaStateName:   "PRIMARY",
		SSLEnabled:         &sslEnabled,
		TypeName:           "REPLICA_PRIMARY",
		UptimeMsec:         1827300394,
		Version:            "3.2.0",
		Username:           "mongo",
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDeployments_GetHostByHostname(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	hostName := "22"
	port := 26000
	path := fmt.Sprintf("/api/public/v1.0/groups/%s/hosts/byName/%s:%d", projectID, hostName, port)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `
					{
					  "alertsEnabled" : true,
					  "aliases": [ "{HOSTNAME}:26000", "{IP-ADDRESS}:26000" ],
					  "authMechanismName" : "SCRAM-SHA-1",
					  "clusterId" : "1",
					  "created" : "2014-04-22T19:56:50Z",
					  "deactivated" : false,
					  "groupId" : "%[1]s",
					  "hasStartupWarnings" : false,
					  "hidden" : false,
					  "hostEnabled" : true,
					  "hostname" : "mongoHost",
					  "id" : "22",
					  "ipAddress": "127.0.0.1",
					  "journalingEnabled" : false,
					  "lastDataSizeBytes" : 633208918,
					  "lastIndexSizeBytes" : 101420524,
					  "lastPing" : "2016-08-18T11:23:41Z",
					  "links" : [  ],
					  "logsEnabled" : false,
					  "lowUlimit" : false,
					  "muninEnabled" : false,
					  "port" : 26000,
					  "profilerEnabled" : false,
					  "replicaSetName": "rs1",
					  "replicaStateName" : "PRIMARY",
					  "sslEnabled" : true,
					  "typeName": "REPLICA_PRIMARY",
					  "uptimeMsec": 1827300394,
					  "username" : "mongo",
					  "version" : "3.2.0"
					}`, projectID)
	})

	alerts, _, err := client.Deployments.GetHostByHostname(ctx, projectID, hostName, port)
	if err != nil {
		t.Fatalf("Deployments.GetHostByHostname returned error: %v", err)
	}

	alertsEnabled := true
	logsEnabled := false
	profilerEnabled := false
	sslEnabled := true

	expected := &Host{
		Aliases: []string{
			"{HOSTNAME}:26000", "{IP-ADDRESS}:26000",
		},
		AlertsEnabled:      &alertsEnabled,
		AuthMechanismName:  "SCRAM-SHA-1",
		ClusterID:          "1",
		Created:            "2014-04-22T19:56:50Z",
		Deactivated:        false,
		GroupID:            projectID,
		HasStartupWarnings: false,
		Hidden:             false,
		HostEnabled:        true,
		Hostname:           "mongoHost",
		ID:                 "22",
		IPAddress:          "127.0.0.1",
		JournalingEnabled:  false,
		LastDataSizeBytes:  633208918,
		LastIndexSizeBytes: 101420524,
		LastPing:           "2016-08-18T11:23:41Z",
		Links:              []*atlas.Link{},
		LogsEnabled:        &logsEnabled,
		LowUlimit:          false,
		Port:               26000,
		ProfilerEnabled:    &profilerEnabled,
		ReplicaSetName:     "rs1",
		ReplicaStateName:   "PRIMARY",
		SSLEnabled:         &sslEnabled,
		TypeName:           "REPLICA_PRIMARY",
		UptimeMsec:         1827300394,
		Version:            "3.2.0",
		Username:           "mongo",
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}

const hostname = "server1.example.com"

func TestDeployments_StartMonitoring(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	var port int32 = 27017
	path := fmt.Sprintf("/api/public/v1.0/groups/%s/hosts", projectID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprintf(w, `
					{
					  "alertsEnabled" : true,
					  "authMechanismName" : "SCRAM-SHA-1",
					  "created" : "2014-04-22T19:56:50Z",
					  "groupId" : "%[1]s",
					  "hasStartupWarnings" : false,
					  "hidden" : false,
					  "hostEnabled" : true,
					  "hostname" : "server1.example.com",
					  "id" : "22",
					  "journalingEnabled" : false,
					  "links" : [  ],
					  "logsEnabled" : false,
					  "lowUlimit" : false,
					  "port" : 27017,
					  "profilerEnabled" : false,
					  "sslEnabled" : false
					}`, projectID)
	})

	host := &Host{
		Hostname: hostname,
		Port:     port,
	}

	alerts, _, err := client.Deployments.StartMonitoring(ctx, projectID, host)
	if err != nil {
		t.Fatalf("Deployments.StartMonitoring returned error: %v", err)
	}

	alertsEnabled := true
	logsEnabled := false
	profilerEnabled := false
	sslEnabled := false

	expected := &Host{
		AlertsEnabled:      &alertsEnabled,
		AuthMechanismName:  "SCRAM-SHA-1",
		Created:            "2014-04-22T19:56:50Z",
		GroupID:            projectID,
		HasStartupWarnings: false,
		Hidden:             false,
		HostEnabled:        true,
		Hostname:           hostname,
		ID:                 "22",
		JournalingEnabled:  false,
		Links:              []*atlas.Link{},
		LogsEnabled:        &logsEnabled,
		LowUlimit:          false,
		Port:               port,
		ProfilerEnabled:    &profilerEnabled,
		SSLEnabled:         &sslEnabled,
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDeployments_UpdateMonitoring(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	hostID := "22"
	var port int32 = 27017
	path := fmt.Sprintf("/api/public/v1.0/groups/%s/hosts/%s", projectID, hostID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprintf(w, `
					{
					  "alertsEnabled" : true,
					  "authMechanismName" : "SCRAM-SHA-1",
					  "created" : "2014-04-22T19:56:50Z",
					  "groupId" : "%[1]s",
					  "hasStartupWarnings" : false,
					  "hidden" : false,
					  "hostEnabled" : true,
					  "hostname" : "server1.example.com",
					  "id" : "22",
					  "journalingEnabled" : false,
					  "links" : [  ],
					  "logsEnabled" : false,
					  "lowUlimit" : false,
					  "muninEnabled" : false,
					  "port" : 27017,
					  "profilerEnabled" : false,
					  "sslEnabled" : false
					}`, projectID)
	})

	host := &Host{
		Hostname: hostname,
		Port:     port,
	}

	alerts, _, err := client.Deployments.UpdateMonitoring(ctx, projectID, hostID, host)
	if err != nil {
		t.Fatalf("Deployments.UpdateMonitoring returned error: %v", err)
	}

	alertsEnabled := true
	logsEnabled := false
	profilerEnabled := false
	sslEnabled := false

	expected := &Host{
		AlertsEnabled:      &alertsEnabled,
		AuthMechanismName:  "SCRAM-SHA-1",
		Created:            "2014-04-22T19:56:50Z",
		GroupID:            projectID,
		HasStartupWarnings: false,
		Hidden:             false,
		HostEnabled:        true,
		Hostname:           hostname,
		ID:                 "22",
		JournalingEnabled:  false,
		Links:              []*atlas.Link{},
		LogsEnabled:        &logsEnabled,
		LowUlimit:          false,
		Port:               port,
		ProfilerEnabled:    &profilerEnabled,
		SSLEnabled:         &sslEnabled,
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}

func TestDeployments_StopMonitoring(t *testing.T) {
	client, mux, teardown := setup()

	defer teardown()

	hostID := "22"
	path := fmt.Sprintf("/api/public/v1.0/groups/%s/hosts/%s", projectID, hostID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Deployments.StopMonitoring(ctx, projectID, hostID)
	if err != nil {
		t.Fatalf("Deployments.StopMonitoring returned error: %v", err)
	}
}
