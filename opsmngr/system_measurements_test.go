package opsmngr

import (
	"fmt"
	"net/http"
	"testing"

	atlas "github.com/mongodb/go-client-mongodb-atlas/mongodbatlas"

	"github.com/go-test/deep"
)

func TestSystemMeasurements_List(t *testing.T) {
	setup()
	defer teardown()

	projectID := "6b8cd3c380eef5349ef77gf7"
	host := "host"

	path := fmt.Sprintf("/groups/%s/hosts/%s/measurements", projectID, host)

	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
					  "end": "2018-07-31T16:29:24Z",
					  "granularity": "P1T12H",
					  "groupId": "6b8cd3c380eef5349ef77gf7",
					  "hostId": "host",
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
		)
	})

	snapshots, _, err := client.SystemMeasurements.List(ctx, projectID, host, nil)
	if err != nil {
		t.Fatalf("Checkpoints.List returned error: %v", err)
	}

	var value float32 = 5.0

	expected := &atlas.ProcessMeasurements{
		End:         "2018-07-31T16:29:24Z",
		Granularity: "P1T12H",
		GroupID:     "6b8cd3c380eef5349ef77gf7",
		HostID:      "host",
		Links:       []*atlas.Link{},
		Measurements: []*atlas.Measurements{
			{
				DataPoints: []*atlas.DataPoints{
					{
						Timestamp: "2018-07-31T16:29:24Z",
						Value:     &value,
					},
				},
				Name:  "CONNECTIONS",
				Units: "SCALAR",
			},
			{
				DataPoints: []*atlas.DataPoints{
					{
						Timestamp: "2018-07-31T16:29:24Z",
						Value:     &value,
					},
				},
				Name:  "NETWORK_BYTES_IN",
				Units: "BYTES_PER_SECOND",
			},
			{
				DataPoints: []*atlas.DataPoints{
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
