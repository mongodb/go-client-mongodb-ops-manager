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

func TestMeasurements_Disk(t *testing.T) {
	client, mux, teardown := setup()
	defer teardown()

	groups := "12345678"
	hostID := "1"
	disk := "disk"

	mux.HandleFunc(fmt.Sprintf("/groups/%s/hosts/%s/disks/%s/measurements", groups, hostID, disk), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
				  "end" : "2017-08-22T20:31:14Z",
				  "granularity" : "PT1M",
				  "groupId" : "12345678",
				  "hostId" : "1",
				  "links" : [ {
					"href" : "https://cloud.mongodb.com/api/public/v1.0/groups/12345678/hosts/1/disks/disk/measurements?granularity=PT1M&period=PT1M",
					"rel" : "self"
				  }, {
					"href" : "https://cloud.mongodb.com/api/public/v1.0/groups/12345678/hosts/1",
					"rel" : "http://mms.mongodb.com/hostID"
				  } ],
				  "measurements" : [ {
					"dataPoints" : [ {
					  "timestamp" : "2017-08-22T20:31:12Z",
					  "value" : null
					}, {
					  "timestamp" : "2017-08-22T20:31:14Z",
					  "value" : null
					} ],
					"name" : "DISK_PARTITION_IOPS_READ",
					"units" : "SCALAR_PER_SECOND"
				  }],
                  "partitionName":"xvdb",
				  "processId" : "shard-00-00.mongodb.net:27017",
				  "start" : "2017-08-22T20:30:45Z"
				}`)
	})

	opts := &atlas.ProcessMeasurementListOptions{
		Granularity: "PT1M",
		Period:      "PT1M",
	}

	measurements, _, err := client.Measurements.Disk(ctx, groups, hostID, disk, opts)
	if err != nil {
		t.Fatalf("Measurements.Disk returned error: %v", err)
	}

	expected := &atlas.ProcessDiskMeasurements{
		ProcessMeasurements: &atlas.ProcessMeasurements{
			End:         "2017-08-22T20:31:14Z",
			Granularity: "PT1M",
			GroupID:     "12345678",
			HostID:      "1",
			Links: []*atlas.Link{
				{
					Rel:  "self",
					Href: "https://cloud.mongodb.com/api/public/v1.0/groups/12345678/hosts/1/disks/disk/measurements?granularity=PT1M&period=PT1M",
				},
				{
					Href: "https://cloud.mongodb.com/api/public/v1.0/groups/12345678/hosts/1",
					Rel:  "http://mms.mongodb.com/hostID",
				},
			},
			Measurements: []*atlas.Measurements{
				{
					DataPoints: []*atlas.DataPoints{
						{
							Timestamp: "2017-08-22T20:31:12Z",
							Value:     nil,
						},
						{
							Timestamp: "2017-08-22T20:31:14Z",
							Value:     nil,
						},
					},
					Name:  "DISK_PARTITION_IOPS_READ",
					Units: "SCALAR_PER_SECOND",
				},
			},
			ProcessID: "shard-00-00.mongodb.net:27017",
			Start:     "2017-08-22T20:30:45Z",
		},
		PartitionName: "xvdb",
	}

	if diff := deep.Equal(measurements, expected); diff != nil {
		t.Error(diff)
	}
}
