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
	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"
)

func TestHost_List(t *testing.T) {
	setup()

	defer teardown()

	projectID := "6b8cd3c380eef5349ef77gf7"
	path := fmt.Sprintf("/groups/%s/hosts", projectID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
				  "totalCount" : 1,
				  "results" : [
					{
					  "alertsEnabled" : true,
					  "aliases": [ "{HOSTNAME}:26000", "{IP-ADDRESS}:26000" ],
					  "authMechanismName" : "SCRAM-SHA-1",
					  "clusterId" : "1",
					  "created" : "2014-04-22T19:56:50Z",
					  "deactivated" : false,
					  "groupId" : "6b8cd3c380eef5349ef77gf7",
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
					}
				  ]
			}`)
	})

	opts := HostListOptions{}

	alerts, _, err := client.Hosts.List(ctx, projectID, &opts)
	if err != nil {
		t.Fatalf("Hosts.List returned error: %v", err)
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
				MuninEnabled:       false,
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

func TestHost_Get(t *testing.T) {
	setup()

	defer teardown()

	projectID := "6b8cd3c380eef5349ef77gf7"
	hostID := "22"
	path := fmt.Sprintf("/groups/%s/hosts/%s", projectID, hostID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `
					{
					  "alertsEnabled" : true,
					  "aliases": [ "{HOSTNAME}:26000", "{IP-ADDRESS}:26000" ],
					  "authMechanismName" : "SCRAM-SHA-1",
					  "clusterId" : "1",
					  "created" : "2014-04-22T19:56:50Z",
					  "deactivated" : false,
					  "groupId" : "6b8cd3c380eef5349ef77gf7",
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
					}`)
	})

	alerts, _, err := client.Hosts.Get(ctx, projectID, hostID)
	if err != nil {
		t.Fatalf("Hosts.Get returned error: %v", err)
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
		MuninEnabled:       false,
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

func TestHost_GetByHostName(t *testing.T) {
	setup()

	defer teardown()

	projectID := "6b8cd3c380eef5349ef77gf7"
	hostName := "22"
	port := 26000
	path := fmt.Sprintf("/groups/%s/hosts/byName/%s:%d", projectID, hostName, port)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `
					{
					  "alertsEnabled" : true,
					  "aliases": [ "{HOSTNAME}:26000", "{IP-ADDRESS}:26000" ],
					  "authMechanismName" : "SCRAM-SHA-1",
					  "clusterId" : "1",
					  "created" : "2014-04-22T19:56:50Z",
					  "deactivated" : false,
					  "groupId" : "6b8cd3c380eef5349ef77gf7",
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
					}`)
	})

	alerts, _, err := client.Hosts.GetByHostname(ctx, projectID, hostName, port)
	if err != nil {
		t.Fatalf("Hosts.GetByHostname returned error: %v", err)
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
		MuninEnabled:       false,
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

func TestHost_Monitoring(t *testing.T) {
	setup()

	defer teardown()

	projectID := "6b8cd3c380eef5349ef77gf7"
	hostName := "server1.example.com"
	var port int32 = 27017
	path := fmt.Sprintf("/groups/%s/hosts", projectID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		_, _ = fmt.Fprint(w, `
					{
					  "alertsEnabled" : true,
					  "authMechanismName" : "SCRAM-SHA-1",
					  "created" : "2014-04-22T19:56:50Z",
					  "groupId" : "6b8cd3c380eef5349ef77gf7",
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
					}`)
	})

	host := &Host{
		Hostname: hostName,
		Port:     port,
	}

	alerts, _, err := client.Hosts.Monitoring(ctx, projectID, host)
	if err != nil {
		t.Fatalf("Hosts.Monitoring returned error: %v", err)
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
		Hostname:           hostName,
		ID:                 "22",
		JournalingEnabled:  false,
		Links:              []*atlas.Link{},
		LogsEnabled:        &logsEnabled,
		LowUlimit:          false,
		MuninEnabled:       false,
		Port:               port,
		ProfilerEnabled:    &profilerEnabled,
		SSLEnabled:         &sslEnabled,
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}

func TestHost_UpdateMonitoring(t *testing.T) {
	setup()

	defer teardown()

	hostID := "22"
	projectID := "6b8cd3c380eef5349ef77gf7"
	hostName := "server1.example.com"
	var port int32 = 27017
	path := fmt.Sprintf("/groups/%s/hosts/%s", projectID, hostID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		_, _ = fmt.Fprint(w, `
					{
					  "alertsEnabled" : true,
					  "authMechanismName" : "SCRAM-SHA-1",
					  "created" : "2014-04-22T19:56:50Z",
					  "groupId" : "6b8cd3c380eef5349ef77gf7",
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
					}`)
	})

	host := &Host{
		Hostname: hostName,
		Port:     port,
	}

	alerts, _, err := client.Hosts.UpdateMonitoring(ctx, projectID, hostID, host)
	if err != nil {
		t.Fatalf("Hosts.Monitoring returned error: %v", err)
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
		Hostname:           hostName,
		ID:                 "22",
		JournalingEnabled:  false,
		Links:              []*atlas.Link{},
		LogsEnabled:        &logsEnabled,
		LowUlimit:          false,
		MuninEnabled:       false,
		Port:               port,
		ProfilerEnabled:    &profilerEnabled,
		SSLEnabled:         &sslEnabled,
	}

	if diff := deep.Equal(alerts, expected); diff != nil {
		t.Error(diff)
	}
}

func TestHost_StopMonitoring(t *testing.T) {
	setup()

	defer teardown()

	hostID := "22"
	projectID := "6b8cd3c380eef5349ef77gf7"
	path := fmt.Sprintf("/groups/%s/hosts/%s", projectID, hostID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	_, err := client.Hosts.StopMonitoring(ctx, projectID, hostID)
	if err != nil {
		t.Fatalf("Hosts.Monitoring returned error: %v", err)
	}
}
