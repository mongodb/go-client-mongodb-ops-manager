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

func TestMeasurements_Host(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	hostID := "hostID"

	path := fmt.Sprintf("/api/public/v1.0/groups/%s/hosts/%s/measurements", projectID, hostID)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprintf(w, `{
					  "end": "2018-07-31T16:29:24Z",
					  "granularity": "P1T12H",
					  "groupId": "%[1]s",
					  "hostId": "hostID",
					  "links": [],
					  "measurements": [{
						"dataPoints": [{
						  "timestamp": "2018-07-31T16:29:24Z",
						  "value": 5.0
						}],
						"name": "CONNECTIONS",
						"units": "SCALAR"
					  }, {
						"dataPoints": [{
						  "timestamp": "2018-07-31T16:29:24Z",
						  "value": 5.0
						}],
						"name": "NETWORK_BYTES_IN",
						"units": "BYTES_PER_SECOND"
					  }, {
						"dataPoints": [{
						  "timestamp": "2018-07-31T16:29:24Z",
						  "value": 5.0
						}],
						"name": "NETWORK_BYTES_OUT",
						"units": "BYTES_PER_SECOND"
					  }],
					  "processId": "{MONGODB-PROCESS-FQDN}:{PORT}",
					  "start": "2018-07-31T16:29:24Z"
			}`,
			projectID,
		)
	})

	snapshots, _, err := client.Measurements.Host(ctx, projectID, hostID, nil)
	if err != nil {
		t.Fatalf("Measurements.Host returned error: %v", err)
	}

	var value float32 = 5.0

	expected := &ProcessMeasurements{
		End:         "2018-07-31T16:29:24Z",
		Granularity: "P1T12H",
		GroupID:     projectID,
		HostID:      "hostID",
		Links:       []*Link{},
		Measurements: []*Measurements{
			{
				DataPoints: []*DataPoints{
					{
						Timestamp: "2018-07-31T16:29:24Z",
						Value:     &value,
					},
				},
				Name:  "CONNECTIONS",
				Units: "SCALAR",
			},
			{
				DataPoints: []*DataPoints{
					{
						Timestamp: "2018-07-31T16:29:24Z",
						Value:     &value,
					},
				},
				Name:  "NETWORK_BYTES_IN",
				Units: "BYTES_PER_SECOND",
			},
			{
				DataPoints: []*DataPoints{
					{
						Timestamp: "2018-07-31T16:29:24Z",
						Value:     &value,
					},
				},
				Name:  "NETWORK_BYTES_OUT",
				Units: "BYTES_PER_SECOND",
			},
		},
		ProcessID: "{MONGODB-PROCESS-FQDN}:{PORT}",
		Start:     "2018-07-31T16:29:24Z",
	}

	if diff := deep.Equal(snapshots, expected); diff != nil {
		t.Error(diff)
	}
}
